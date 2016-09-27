/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package status

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
	case types.StatusMsgTypeMT_STATUS_ACK_FROMCLIENT:
		err = AckMsg(mb)
	case types.StatusMsgTypeMT_STATUS_SESSIONINVALIDATE:
		err = InvalidateSession(mb)
	default:
		return errors.New("status: type undefined")
	}

	return err
}
