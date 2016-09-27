/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : status.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package status

import (
	log "im/common/log4go"
	mb "im/common/proto/entity/msg/msgbase"
	"im/service/logic/service/session"
)

//处理失效的session
func InvalidateSession(m *mb.MsgBase) error {
	log.Debug("[Status|InvalidateSession] uid:%d", m.From)
	s := &session.Session{Uid: m.From}
	return s.Del()
}

//处理客户端回复的ack消息
func AckMsg(m *mb.MsgBase) error {
	s := &session.Session{Uid: m.From}
	s.Get()

	if s.ConnSessionId != 0 {
		if m.Msgid > s.AckMsgId {
			s.AckMsgId = m.Msgid
			return s.Put()
		}
	}

	return nil
}
