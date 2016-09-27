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
	//log "im/common/log4go"
	"im/common/proto/entity/msg"
	"im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
)

func ProcessMsg(mr *msg.MsgRaw) error {

	// 这里可以根据需要来处理Header
	//

	mb := &msgbase.MsgBase{}
	err := mb.UnSerialize(mr.Body[msg.HEADER_LEN:])
	if err != nil {
		return err
	}

	switch mb.Type {
	case types.UserMsgTypeMT_USER_LOGIN_REQ:
		err = Login(mb)
	default:
		return errors.New("user: type undefined")
	}

	return err
}
