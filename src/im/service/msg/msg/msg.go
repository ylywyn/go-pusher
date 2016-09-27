/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package msg

import (
	"errors"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	"im/service/logic/service/id"
	. "im/service/logic/service/send"
	"im/service/logic/service/session"
	"strconv"
)

func GeneralMsg(h msg.MsgHeader, m *mb.MsgBase) error {
	// 获取一个msg id
	msgId := id.GetId("msg")
	if msgId == 0 {
		return errors.New("get msg id error")
	}

	// 回复 ack
	if h.Ack == 1 && m.Msgid != 0 && m.From != 0 {
		sendAckMsg(m.From, m.Msgid, strconv.FormatUint(msgId, 10))
	}

	m.Msgid = msgId
	return Send.SendOne(&h, m)
}

func GeneralGroupMsg(h msg.MsgHeader, m *mb.MsgBase) error {

	// 获取一个msg id
	msgId := id.GetId("group")
	if msgId == 0 {
		return errors.New("get msg id error")
	}

	//回复 ack
	if h.Ack == 1 && m.Msgid != 0 && m.From != 0 {
		sendAckMsg(m.From, m.Msgid, strconv.FormatUint(msgId, 10))
	}

	m.Msgid = msgId
	Send.SendToGroup(&h, m)

	return nil
}

// 服务端已经接收，回复ack
func sendAckMsg(uid, msgid uint64, newid string) {
	s := &session.Session{Uid: uid}
	s.Get()

	if s.ConnSessionId == 0 || s.ConnClusterId == 0 {
		log.Debug("[msg|sendAckMsg|Session] get session faild")
		return
	}

	mbmsg := &mb.MsgBase{
		To:     s.Uid,
		Connid: s.ConnSessionId,
		Type:   types.StatusMsgTypeMT_STATUS_ACK_FROMSERVER,
		Appid:  s.Appid,
		Msgid:  msgid,
		Text:   newid,
	}

	data := mbmsg.Serialize()

	m := msg.MsgRaw{Body: data}
	m.Header.Ack = 0
	m.Header.Compress = 0
	m.Header.Length = uint16(len(data) - 4)
	m.Header.Repeat = 0
	m.Header.Type = mbmsg.Type
	m.Serialize()

	Send.SendMrToComet(&m, s.ConnClusterId)
}
