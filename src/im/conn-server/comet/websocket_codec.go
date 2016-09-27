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
	mb "im/common/proto/entity/msg/msgbase"
	//log "im/common/log4go"
	"io"
	"io/ioutil"

	"github.com/pquerna/ffjson/ffjson"
)

type WebSocketCodec struct {
}

//读取规定长度的数据
func (self *WebSocketCodec) Read(r io.Reader) (*mb.MsgBase, error) {

	p, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if len(p) == msg.HEADER_LEN {
		return nil, nil
	}

	return self.UnmarshalToMb(p)
}

////Josn bytes 转到 MsgBase
func (self *WebSocketCodec) UnmarshalToMb(data []byte) (*mb.MsgBase, error) {
	//log.Debug("[WebSocketCodec|UnmarshalToMb] %s", string(data))
	mb := &mb.MsgBase{}
	if err := ffjson.Unmarshal(data, mb); err != nil {
		return nil, err
	}

	return mb, nil
}

// MsgBase bytes to MsgRaw
func (self *WebSocketCodec) UnmarshalToMr(mb *mb.MsgBase) *msg.MsgRaw {

	data := mb.Serialize()
	m := &msg.MsgRaw{Body: data}

	m.Header.Ack = 0
	m.Header.Compress = 0
	m.Header.Length = uint16(len(data))
	m.Header.Repeat = 0
	m.Header.Type = mb.Type
	m.Serialize()

	return m
}

//MsgRaw To Json Bytes
func (self *WebSocketCodec) MarshalToJson(m *msg.MsgRaw) ([]byte, error) {
	mb := mb.MsgBase{}
	err := mb.UnSerialize(m.Body[msg.HEADER_LEN:])
	if err != nil {
		return nil, err
	}

	data, err := ffjson.Marshal(mb)
	if err != nil {
		return nil, err
	}

	return data, nil
}
