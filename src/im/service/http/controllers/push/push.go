/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : push.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package push

import (
	"encoding/json"
	"errors"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	"im/service/http/controllers"
	"im/service/logic/service/id"
	. "im/service/logic/service/send"
)

const (
	MIN_JOSN_MSG = 64
	MAX_JSON_MSG = 4096
)

type PushController struct {
	controllers.BaseController
}

//新版 解析通用消息体
func (this *PushController) Push() {

	requestbody := this.Ctx.Input.RequestBody
	//check
	l := len(requestbody)
	if l < MIN_JOSN_MSG {
		this.SetError("json data invalid")
		return
	}
	if l > MAX_JSON_MSG {
		this.SetError("json data size too long")
		return
	}

	var m mb.MsgBase
	err := json.Unmarshal(requestbody, &m)
	if err != nil {
		this.SetError(err.Error())
		return
	}

	m.Msgid = id.GetId("webmsg")
	if m.Msgid == 0 {
		this.SetError("get msg id error")
		return
	}
	go sendAndSave(&m)
	//log.Debug("[http|Push] msgid:%d", m.Msgid)

	if err != nil {
		this.SetError(err.Error())
	} else {
		this.SetData(m.Msgid)
	}
}

func sendAndSave(m *mb.MsgBase) {

	var err error
	h := &msg.MsgHeader{
		Ack:      0,
		Compress: 0,
		Repeat:   0,
		Type:     m.Type,
	}

	switch m.Type {
	case types.GeneralMsgTypeMT_GENERAL_MSG,
		types.GeneralMsgTypeMT_GENERAL_NOTICE:
		err = Send.SendOne(h, m)
	case types.GeneralMsgTypeMT_GENERAL_GROUP_MSG,
		types.GeneralMsgTypeMT_GENERAL_GROUP_NOTICE:
		err = Send.SendToGroup(h, m)
	default:
		err = errors.New("不支持的消息类型")
	}

	if err != nil {
		log.Error("[http|Push|sendAndSave] error:%s", err.Error())
	}
}
