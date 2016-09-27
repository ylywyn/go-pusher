/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : user.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/
package users

import (
	"encoding/json"
	"im/common/crypto/md5"
	//log "im/common/log4go"
	"im/common/proto/fbsgen/msg/types"
	"im/service/http/controllers"
	"im/service/logic/models/user"
	"im/service/logic/service/id"
	"im/service/logic/service/session"
	suser "im/service/logic/service/user"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	JSONSTR = ".json"
)

type UserController struct {
	controllers.BaseController
}

type Login struct {
	Username string
	Password string
	Appid    uint32
	Token    string
}

type LoginUid struct {
	Uid      uint64
	Password string
	Appid    uint32
	Token    string
}

type LoginResult struct {
	User      *user.Users
	Sessionid string `json:"sessionid"`
}

//注册 post
func (this *UserController) SignUp() {

	//1.获取post的json消息
	var u user.Users
	jsonmsg := this.Ctx.Input.RequestBody
	err := json.Unmarshal(jsonmsg, &u)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	if len(u.Username) < 2 || len(u.Password) < 2 || u.Appid < 1 {
		this.SetError("username or password or appid unsuitable")
		return
	}

	//2. 检查username是否可用
	r, err := user.CheckUserName(u.Username, u.Appid)
	if r || (err == mgo.ErrNotFound) {
		//获取 uid
		uid := id.GetId("SignUp")
		if uid == 0 {
			this.SetError("get is error")
			return
		}
		u.Uid = uid
		u.Password = md5.Md5([]byte(u.Password))
		u.Registertime = time.Now().Unix()
		err = u.Insert()
		if err != nil {
			this.SetError(err.Error())
		} else {
			this.SetData(uid)
		}

		//插入用户扩展信息表
		uext := &user.UserExt{Uid: uid, Appid: u.Appid}
		uext.Insert()

	} else {
		this.SetError(err.Error())
	}
}

//IOS需要设置token
func (this *UserController) Login() {
	//1. 解析json
	var login Login
	jsonmsg := this.Ctx.Input.RequestBody
	err := json.Unmarshal(jsonmsg, &login)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	//2. 登录
	u, err := suser.Login(login.Appid, login.Username, login.Password)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	//3. session (踢人动作,仅当当前登陆是IOS才会发生）
	s := session.Session{Uid: u.Uid}
	s.Get()
	if s.Appid != 0 {
		if len(login.Token) > 0 && login.Token != s.ApnsToken {
			//踢人[上次是其它平台或者是其它ApnsToken]
			s.OsType = types.PlatformPF_IOS
			s.ApnsToken = login.Token
			s.Put()
		}
	} else {
		s.OsType = 0
		if len(login.Token) > 0 {
			s.OsType = types.PlatformPF_IOS
			s.ApnsToken = login.Token
		}
		s.Admin = 0
		if u.Admin {
			s.Admin = 1
		}
		s.Appid = uint16(u.Appid)
		s.Put()
	}

	sessionid := session.NewSessionId(login.Appid, u.Uid, "")
	this.SetData(&LoginResult{u, sessionid})
}

//通过UID登录，IOS需要设置token
func (this *UserController) LoginByUid() {
	//1. 解析json
	var login LoginUid
	jsonmsg := this.Ctx.Input.RequestBody
	err := json.Unmarshal(jsonmsg, &login)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	//2. 登录
	u, err := suser.LoginByUid(login.Appid, login.Uid, login.Password)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	//3. session (踢人动作,仅当当前登陆是IOS才会发生）
	s := session.Session{Uid: u.Uid}
	s.Get()
	if s.Appid != 0 {
		if len(login.Token) > 0 && login.Token != s.ApnsToken {
			//踢人[上次是其它平台或者是其它ApnsToken]
			s.OsType = types.PlatformPF_IOS
			s.ApnsToken = login.Token
			s.Put()
		}
	} else {
		s.OsType = 0
		if len(login.Token) > 0 {
			s.OsType = types.PlatformPF_IOS
			s.ApnsToken = login.Token
		}
		s.Admin = 0
		if u.Admin {
			s.Admin = 1
		}
		s.Appid = uint16(u.Appid)
		s.Put()
	}

	sessionid := session.NewSessionId(login.Appid, u.Uid, "")
	this.SetData(&LoginResult{u, sessionid})
}

//仅限IOS系统调用，设置token
func (this *UserController) Logout() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("Logout failed:" + err.Error())
		return
	}
	if s.OsType == types.PlatformPF_IOS {
		s.Del()
	} else {
		if s.ConnClusterId == 0 {
			s.Del()
		}
	}

	this.SetData("ok")
}

//从数据库中查看用户名是否可用（是否存在）
func (this *UserController) CheckUserName() {
	username := this.GetString("username")
	appid, err := this.GetInt32("appid")
	if err != nil {
		this.SetError(err.Error())
		return
	}
	brst, err := user.CheckUserName(username, uint32(appid))

	if brst {
		this.SetData("用户名可用")
	} else {
		this.SetError(err.Error())
	}
}

// GET 获取用户详情
func (this *UserController) GetUserInfo() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	uid, err := this.GetInt64("uid")
	if err != nil {
		this.SetError("uid error")
		return
	}

	u := &user.Users{Uid: uint64(uid), Appid: uint32(s.Appid)}
	err = u.Read()
	if err != nil {
		this.SetError("can't find user:" + err.Error())
	} else {
		u.Password = ""
		this.SetData(u)
	}
}

//GET 获取所有用户
func (this *UserController) GetUsersInfo() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	if s.Admin == 0 {
		this.SetError("非管理员不能获取所有用户信息")
		return
	}

	depid, err := this.GetInt("depid")
	if err != nil {
		depid = -2
	}

	users, err := user.GetUsersInfo(uint32(s.Appid), uint32(depid))
	if err != nil {
		this.SetError(err.Error())
	} else {
		for i := 0; i < len(users); i++ {
			users[i].Password = ""
		}
		this.SetData(users)
	}
}

//删除用户 GET
func (this *UserController) RemoveUser() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	deluid, err := this.GetInt64("uid")
	deluid64 := uint64(deluid)
	if err != nil {
		this.SetError("uid不正确")
		return
	}

	if s.Admin != 1 && s.Uid != deluid64 {
		this.SetError("非管理员用户无法删除其他用户")
		return
	}

	u := &user.Users{Uid: deluid64, Appid: uint32(s.Appid)}
	err = u.Delete()

	// 删除用户扩展信息表
	uext := &user.UserExt{Uid: deluid64, Appid: uint32(s.Appid)}
	uext.Delete()

	if s.Uid == deluid64 {
		s.Del()
	}
	//
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData("删除成功")
	}
}

func (this *UserController) AddFriend() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("sessionid is nil")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	uidaim64, _ := this.GetInt64("uid")
	uidaim := uint64(uidaim64)

	//验证被添加ID是否有效
	if s.Uid == uidaim {
		this.SetError("参数错误：可能不应该添加自己为好友")
		return
	}

	user := &user.Users{Uid: uidaim, Appid: uint32(s.Appid)}
	err = user.Read()
	if err != nil {
		this.SetError("该用户不存在")
		return
	}

	err = suser.AddFriend(s.Uid, uidaim, uint32(s.Appid), true)
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData("ok")
	}
}

func (this *UserController) RemoveFriend() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	uidaim64, _ := this.GetInt64("uid")
	uidaim := uint64(uidaim64)

	//验证被添加ID是否有效
	if s.Uid == uidaim {
		this.SetError("参数错误：不能删除自己")
		return
	}

	user := &user.Users{Uid: uidaim, Appid: uint32(s.Appid)}
	err = user.Read()
	if err != nil {
		this.SetError("该用户不存在")
		return
	}

	err = suser.DelFriend(s.Uid, uidaim, uint32(s.Appid), true)
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData("ok")
	}
}

func (this *UserController) GetFriendsList() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	uext := user.UserExt{Uid: s.Uid, Appid: uint32(s.Appid)}
	err = uext.GetFriends()
	if err != nil {
		this.SetError(err.Error())
		return
	}
	this.SetData(uext.Friends)
}

func (this *UserController) GetFriendVersion() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	uext := user.UserExt{Uid: s.Uid, Appid: uint32(s.Appid)}
	err = uext.GetFriendsVersion()
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData(uext.FriendsVersion)
	}
}

// 更改密码
func (this *UserController) UpdatePassword() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("please login")
		return
	}

	type PwdModify struct {
		pwd    string
		modify string
	}
	var mp PwdModify
	jsondata := this.Ctx.Input.RequestBody
	err = json.Unmarshal(jsondata, &mp)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	if len(mp.pwd) < 1 || len(mp.modify) < 3 {
		this.SetError("错误的密码长度")
		return
	}

	//读取user信息，获取数据库中的密码，检查是否正确
	userinfo := &user.Users{Uid: s.Uid, Appid: uint32(s.Appid)}
	err = userinfo.Read()
	if err != nil || userinfo.Username == "" {
		this.SetError("sessionid 未能找到用户")
		return
	}

	mp.pwd = md5.Md5([]byte(mp.pwd))
	mp.modify = md5.Md5([]byte(mp.modify))

	if userinfo.Password != mp.pwd {
		this.SetError("密码不正确，请重新输入")
		return
	}
	userinfo.Password = mp.modify
	err = userinfo.Update()
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData("修改成功")
	}
}
