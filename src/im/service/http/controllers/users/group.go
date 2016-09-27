/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : group.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/
package users

import (
	"encoding/json"
	log "im/common/log4go"
	"im/service/http/controllers"
	"im/service/logic/models/group"
	"im/service/logic/models/user"
	sgroup "im/service/logic/service/group"
	"im/service/logic/service/session"
	"time"
)

type GroupController struct {
	controllers.BaseController
}

// 创建群组 Post Josn
func (this *GroupController) CreateGroup() {
	//1.验证session
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	//2. 读取body
	var g group.Group
	body := this.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &g)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	l := len(g.Name)
	if l <= 0 || l > 21 {
		this.SetError("请设置群组名字, 不超过7个汉字")
		return
	}

	//3.创建群组
	g.Ownerid = s.Uid
	g.Creattime = time.Now().Unix()
	g.Appid = uint32(s.Appid)
	gid, err := sgroup.GroupCreate(s.Uid, &g)
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData(gid)
	}
}

//删除群组 GET
func (this *GroupController) DeleteGroup() {

	//1.检查session
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	//2. 获取get来的群组id
	gid, err := this.GetInt64("gid")
	if err != nil {
		this.SetError("gid不正确")
		return
	}
	notice, _ := this.GetBool("notice")

	g := group.Group{Gid: uint64(gid), Appid: uint32(s.Appid)}
	err = sgroup.GroupDelete(s.Uid, &g, notice)
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData("删除成功")
	}
}

// 获取群组详情 Get
func (this *GroupController) GetGroupInfo() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	gid, err := this.GetInt64("gid")
	if err != nil {
		this.SetError(err.Error())
		return
	}

	g := group.Group{Gid: uint64(gid), Appid: uint32(s.Appid)}
	err = g.Read()
	if err != nil || g.Name == "" {
		this.SetError(err.Error())
	} else {
		this.SetData(g)
	}
}

//获取自己创建的所有群组 GET
func (this *GroupController) GetCreateGroups() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	uext := &user.UserExt{Uid: uint64(s.Uid), Appid: uint32(s.Appid)}
	err = uext.GetGroupInfo()

	if err != nil {
		this.SetError(err.Error())
		return
	} else {
		l := len(uext.CreateGroups)
		groups := make([]group.Group, l)

		for i := 0; i < l; i++ {
			g := group.Group{Gid: uext.CreateGroups[i], Appid: uint32(s.Appid)}
			err = g.Read()
			if err != nil {
				log.Error("[GroupController|GetCreateGroups] %s", err.Error())
			} else {
				groups[i] = g
			}
		}

		this.SetData(groups)
	}
}

//获取加入的群组 GET
func (this *GroupController) GetJoinGroups() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	uext := &user.UserExt{Uid: s.Uid, Appid: uint32(s.Appid)}
	err = uext.GetGroupInfo()

	if err != nil {
		this.SetError(err.Error())
		return
	} else {

		l := len(uext.CreateGroups)
		groups := make([]group.Group, l)

		for i := 0; i < l; i++ {
			g := group.Group{Gid: uext.Groups[i], Appid: uint32(s.Appid)}
			err = g.Read()
			if err != nil {
				log.Error("[GroupController|GetJoinGroups] %s", err.Error())
			} else {
				groups[i] = g
			}
		}

		this.SetData(groups)
	}
}

//添加群组成员 Get
func (this *GroupController) AddGroupMember() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	////验证被添加ID是否有效
	uidaim64, _ := this.GetInt64("uid")
	uidaim := uint64(uidaim64)

	user := &user.Users{Uid: uidaim, Appid: uint32(s.Appid)}
	err = user.Read()

	if err != nil || user.Appid != uint32(s.Appid) {
		this.SetError("未找到要添加的成员")
		return
	}

	//2. 读取参数
	gid, _ := this.GetInt64("gid")
	notice, _ := this.GetBool("notice")
	g := group.Group{Gid: uint64(gid), Appid: uint32(s.Appid)}

	//3.添加
	err = sgroup.GroupAddMember(s.Uid, uidaim, &g, notice)
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData("添加成功")
	}
}

//添加群组成员
func (this *GroupController) AddGroupMembers() {

	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	type addMembers struct {
		Uids []uint64
		Gid  uint64
	}
	var mbs addMembers
	body := this.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &mbs)

	if err != nil {
		this.SetError(err.Error())
		return
	}

	length := len(mbs.Uids)
	for i := 0; i < length; i++ {
		uidaim := mbs.Uids[i]
		if s.Uid == uidaim {
			continue
		}

		user := &user.Users{Uid: uidaim, Appid: uint32(s.Appid)}
		err = user.Read()

		if err != nil {
			log.Error("[GroupController|AddGroupMembers|check users] %s", err.Error())
			this.SetError("未找到要添加的成员")
			return
		}

		gp := group.Group{Gid: mbs.Gid, Appid: uint32(s.Appid)}
		err = sgroup.GroupAddMember(s.Uid, uidaim, &gp, true)
		if err != nil {
			this.SetError(err.Error())
			return
		}
	}
	this.SetData("添加成功")

}

//删除群组成员
func (this *GroupController) RemoveGroupMembers() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	type DelMembers struct {
		Uids []uint64
		Gid  uint64
	}
	var delMembers DelMembers
	body := this.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &delMembers)

	if err != nil {
		this.SetError(err.Error())
		return
	}

	//2. 读取参数 校验是否有权限删除成员
	if delMembers.Gid <= 0 {
		this.SetError("参数错误，请确认添加群组id.")
		return
	}

	if len(delMembers.Uids) == 0 {
		this.SetError("参数错误:要删除成员列表错误。")
		return
	}

	//3.删除
	for _, uidaim := range delMembers.Uids {
		///校验群组无法删除自己 需要解散群
		if s.Uid == uidaim {
			this.SetError("参数错误:群主不能移除自己，需要解散群。")
			return
		}
		gp := group.Group{Gid: delMembers.Gid, Appid: uint32(s.Appid)}
		err = sgroup.GroupDelMember(s.Uid, uidaim, &gp, true)
		if err != nil {
			this.SetError(err.Error())
		} else {
			this.SetData("删除成功")
		}
	}
}

//获取群组成员列表
func (this *GroupController) GetGroupMemberList() {

	//1.检查session
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	gid, err := this.GetInt64("gid")
	if err != nil {
		this.SetError(err.Error())
		return
	}

	//2. 获取成员列表
	g := group.Group{Gid: uint64(gid), Appid: uint32(s.Appid)}
	err = g.GetMembers()
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData(g.Members)
	}
}

//群组成员退群
func (this *GroupController) QuitGroup() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	//2. 读取参数
	gid, _ := this.GetInt64("gid")
	notice, _ := this.GetBool("notice")

	//3.退群
	gp := group.Group{Gid: uint64(gid), Appid: uint32(s.Appid)}
	err = sgroup.GroupQuit(s.Uid, &gp, notice)
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData("退群成功")
	}
}

//获取群组版本
func (this *GroupController) GetGroupVersion() {
	sessionid := this.GetString("sessionid")
	if len(sessionid) == 0 {
		this.SetError("please login")
		return
	}

	s := session.Session{}
	err := s.GetBySessionId(sessionid)
	if err != nil {
		this.SetError("GetBySessionId failed:" + err.Error())
		return
	}

	// 获取get来的群组id
	gid, err := this.GetInt64("gid")
	if err != nil {
		this.SetError("gid不正确")
		return
	}

	g := group.Group{Gid: uint64(gid), Appid: uint32(s.Appid)}
	g.Read()
	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData(g.Version)
	}
}

//修改群组信息
//func (this *GroupController) UpdateGroupInfo() {

//	sessionid := this.GetString("sessionid")
//	if len(sessionid) == 0 {
//		this.SetError("please login")
//		return
//	}

//	s := session.Session{}
//	err := s.GetBySessionId(sessionid)
//	if err != nil {
//		this.SetError("GetBySessionId failed:" + err.Error())
//		return
//	}

//	var gup group.Group
//	body := this.Ctx.Input.RequestBody
//	err = json.Unmarshal(body, &gup)
//	if err != nil {
//		this.SetError(err.Error())
//		return
//	}

//	l := len(gup.Name)
//	if l <= 0 || l > 24 {
//		this.SetError("请设置群组名字, 不超过8个汉字")
//		return
//	}

//	g := group.Group{Gid: gup.Gid}
//	err = g.Read()
//	if err != nil {
//		this.SetError(err.Error())
//		return
//	}

//	if g.Ownerid != s.Uid {
//		this.SetError("没有权限修改群信息")
//		return
//	}

//	if len(gup.Name) > 0 {
//		g.Name = gup.Name
//	}

//	if len(gup.Desc) > 0 {
//		g.Desc = gup.Desc
//	}

//	g.Version = g.Version + 1
//	err = g.Update()
//	if err != nil {
//		this.SetError(err.Error())
//	} else {
//		this.SetData("修改成功")
//	}
//}
