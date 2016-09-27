/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package status

import (
	"im/common/proto/entity/msg"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	"time"
)

func NewStateMsg(tp int, err string) *msg.MsgRaw {
	m := &mb.MsgBase{
		Type: uint16(tp),
		Text: err,
		Time: uint64(time.Now().UnixNano()),
	}

	data := m.Serialize()

	msgRaw := &msg.MsgRaw{Body: data}
	msgRaw.Header.Length = uint16(len(data) - msg.HEADER_LEN)
	msgRaw.Header.Type = m.Type
	msgRaw.Serialize()
	return msgRaw
}

func NewSessionClearMsg(uid uint64) *msg.MsgRaw {
	m := &mb.MsgBase{
		Type: uint16(types.StatusMsgTypeMT_STATUS_SESSIONINVALIDATE),
		From: uid,
		Time: uint64(time.Now().UnixNano()),
	}

	data := m.Serialize()

	msgRaw := &msg.MsgRaw{Body: data}
	msgRaw.Header.Length = uint16(len(data) - msg.HEADER_LEN)
	msgRaw.Header.Type = m.Type
	msgRaw.Serialize()
	return msgRaw
}
