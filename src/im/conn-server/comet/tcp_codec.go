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
	"bufio"
	"errors"
	"im/common/proto/entity/msg"
	mb "im/common/proto/entity/msg/msgbase"
	tps "im/common/proto/fbsgen/msg/types"
	"io"

	log "im/common/log4go"
)

type TcpLenCodec struct {
}

//读取规定长度的数据
func (self *TcpLenCodec) Read(reader *bufio.Reader) (*msg.MsgRaw, error) {
	var d [msg.HEADER_LEN]byte
	n, err := io.ReadAtLeast(reader, d[:], msg.HEADER_LEN)

	//	if n != msg.HEADER_LEN {
	//		log.Warn("[LenCodec|Read: %d] read1 head n>HeaderLength", n)
	//	}

	if err != nil {
		return nil, err
	}

	h := msg.MsgHeader{}
	h.Parse(d)

	//心跳消息
	if h.Type == tps.StatusMsgTypeMT_STATUS_HEARTBEAT {
		//log.Debug("[LenCodec|Read] heart beat...")
		return nil, nil
	}

	if h.Type < 0 || h.Type > 4096 {
		// 将错误发送回去， 等客户端主动关闭连接 ??
		//c.writeErrorMsg("msg type error")
		return nil, errors.New("msg type error")
	}

	if h.Length < 32 || h.Length > msg.MAX_BODY_LEN {
		log.Error("[LenCodec|Read] %d", h.Length)
		return nil, errors.New("msg too big")
	}

	buff := make([]byte, int(h.Length)+msg.HEADER_LEN)
	n, err = io.ReadAtLeast(reader, buff[msg.HEADER_LEN:], int(h.Length))

	if n != int(h.Length) {
		log.Warn("[LenCodec|Read] read body n>BodyLength")
	}

	if err != nil {
		return nil, err
	}

	//赋值头部
	for i := 0; i < msg.HEADER_LEN; i++ {
		buff[i] = d[i]
	}

	return &msg.MsgRaw{h, buff}, nil
}

//反序列化
func (self *TcpLenCodec) UnmarshalPacket(m *msg.MsgRaw) (*mb.MsgBase, error) {
	mb := &mb.MsgBase{}
	err := mb.UnSerialize(m.Body[msg.HEADER_LEN:])
	return mb, err
}

//序列化
func (self *TcpLenCodec) MarshalPacket(packet *msg.MsgRaw) []byte {
	return packet.Body
}
