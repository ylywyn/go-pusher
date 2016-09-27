/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package comet

import (
	"im/common/proto/entity/msg"
)

type Conn interface {
	Close()
	Write(m *msg.MsgRaw) error
	//Write(data []byte) error
}
