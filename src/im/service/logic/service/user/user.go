/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package user

import (
	"errors"
	"im/common/crypto/md5"
	//log "im/common/log4go"
	"im/service/logic/models/user"
	"im/service/logic/service/session"
	"time"
)

//登录
func Login(appid uint32, username string, passwd string) (*user.Users, error) {
	if appid <= 0 {
		return nil, errors.New("appid  invalidate")
	}

	if len(username) == 0 || len(passwd) == 0 {
		return nil, errors.New("username or password invalidate")
	}

	//1. 先从Users表中验证 账户密码
	password := md5.Md5([]byte(passwd))
	u := &user.Users{Username: username, Appid: appid}
	err := u.ReadByUserName()
	if err != nil || u.Password != password {
		return nil, errors.New("username or password error")
	}

	u.Lastlogintime = time.Now().Unix()
	u.Update()
	u.Password = ""
	return u, nil
}

func LoginByUid(appid uint32, uid uint64, passwd string) (*user.Users, error) {
	if appid <= 0 {
		return nil, errors.New("appid  invalidate")
	}

	if uid == 0 || len(passwd) == 0 {
		return nil, errors.New("uid or passwd  invalidate")
	}

	passwd = md5.Md5([]byte(passwd))
	u := &user.Users{Uid: uid, Appid: appid}
	err := u.Read()
	if err != nil || u.Password != passwd {
		return nil, errors.New("username or passwd error")
	}

	u.Lastlogintime = time.Now().Unix()
	u.Update()
	u.Password = ""
	return u, nil
}

//更新session
func UpdateSession(news *session.Session) (*session.Session, error) {
	//session
	old := &session.Session{Uid: news.Uid}
	old.Get()
	if old.AckMsgId != 0 {
		news.AckMsgId = old.AckMsgId
	}
	return old, news.Put()
}
