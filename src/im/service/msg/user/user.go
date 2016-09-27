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
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	"im/common/proto/entity/msg/login"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	"im/service/logic/service/send"
	"im/service/logic/service/session"
	"im/service/logic/service/user"
)

//长连接端登录
func Login(m *mb.MsgBase) error {

	//反序列化登陆信息
	lmsg := login.MsgLogin{}
	err := lmsg.UnSerialize([]byte(m.Text))
	if err != nil {
		return err
	}

	s := &session.Session{
		Uid:           lmsg.Uid,
		OsType:        lmsg.Platform,
		Appid:         m.Appid,
		ConnClusterId: uint32(m.ConnServerid),
		ConnSessionId: m.Connid,
		Key:           lmsg.Key,
		AckMsgId:      lmsg.LastMsgId,
	}

	retMsg := login.MsgLoginRet{
		Uid: lmsg.Uid,
		Ret: false,
	}

	//log.Debug("before user.Login:uid %d, appid:%d, passwd:%s", lmsg.Uid, uint32(m.Appid), lmsg.PassWd)
	_, err = user.LoginByUid(uint32(m.Appid), lmsg.Uid, lmsg.PassWd)
	if err != nil {
		//返回登录失败消息
		retMsg.Data = err.Error()
		sendLoginRetMsg(&retMsg, s)
		return err
	}

	log.Debug("[msg|process|Login] serverid:%d, connid:%d, uid:%d", uint32(m.ConnServerid), m.Connid, lmsg.Uid)

	old, err := user.UpdateSession(s)
	if err != nil {
		//返回登录失败消息
		retMsg.Data = err.Error()
		sendLoginRetMsg(&retMsg, s)
		return err
	}

	//	if len(old.Key) != 0 && s.Key != old.Key {
	//		//踢人动作
	//		if old.ConnSessionId != 0 {
	//			sendKickUserMsg(lmsg.Platform, old)
	//		}
	//	}

	//踢人动作
	if old.ConnSessionId != 0 && old.ConnSessionId != s.ConnSessionId {
		sendKickUserMsg(lmsg.Platform, old)
	}

	//返回登录成功消息
	retMsg.Ret = true
	retMsg.Data = "login ok"
	sendLoginRetMsg(&retMsg, s)
	return nil
}

//登录返回消息
func sendLoginRetMsg(ret *login.MsgLoginRet, s *session.Session) {
	r, _ := ret.Serialize()

	t := types.UserMsgTypeMT_USER_LOGIN_FAILED_REP
	if ret.Ret {
		t = types.UserMsgTypeMT_USER_LOGIN_REP
	}
	mbmsg := &mb.MsgBase{
		To:    s.ConnSessionId,
		Type:  uint16(t),
		Appid: s.Appid,
		Text:  string(r),
	}
	mbmsg.Connid = s.ConnSessionId
	data := mbmsg.Serialize()

	m := msg.MsgRaw{Body: data}
	m.Header.Ack = 0
	m.Header.Compress = 0
	m.Header.Length = uint16(len(data) - 4)
	m.Header.Repeat = 0
	m.Header.Type = mbmsg.Type
	m.Serialize()

	send.Send.SendMrToComet(&m, s.ConnClusterId)
}

//发送踢人消息
func sendKickUserMsg(newpl uint8, s *session.Session) {
	mbmsg := &mb.MsgBase{
		To:       s.ConnSessionId,
		Type:     types.UserMsgTypeMT_USER_KICKUSER,
		Appid:    s.Appid,
		Platform: newpl,
	}
	mbmsg.Connid = s.ConnSessionId
	data := mbmsg.Serialize()

	m := msg.MsgRaw{Body: data}
	m.Header.Ack = 0
	m.Header.Compress = 0
	m.Header.Length = uint16(len(data) - 4)
	m.Header.Repeat = 0
	m.Header.Type = mbmsg.Type
	m.Serialize()

	send.Send.SendMrToComet(&m, s.ConnClusterId)
}
