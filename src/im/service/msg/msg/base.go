/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : base.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/
package msg

import (
	"errors"
	"im/common/proto/entity/msg"
	"im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
)

func ProcessMsg(mr *msg.MsgRaw) error {

	// 这里可以根据需要来处理Header
	mb := &msgbase.MsgBase{}
	err := mb.UnSerialize(mr.Body[msg.HEADER_LEN:])
	if err != nil {
		return err
	}

	switch mb.Type {
	case types.GeneralMsgTypeMT_GENERAL_MSG:
		err = GeneralMsg(mr.Header, mb)
	case types.GeneralMsgTypeMT_GENERAL_GROUP_MSG:
		err = GeneralGroupMsg(mr.Header, mb)
	default:
		return errors.New("msg: type undefined")
	}

	return err
}
