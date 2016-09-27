/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package msg

const HEADER_LEN = 4
const MAX_BODY_LEN = 2048 + 200

type MsgHeader struct {
	Ack      byte
	Compress byte
	Repeat   byte
	Type     uint16
	Length   uint16
}

func (header *MsgHeader) Parse(h [HEADER_LEN]byte) {
	header.Length = uint16(h[0])<<8 + uint16(h[1])
	header.Type = (uint16(h[3])&0xf0)<<4 + uint16(h[2])
	header.Repeat = (h[3] >> 2) & 0x03
	header.Compress = (h[3] >> 1) & 0x01
	header.Ack = h[3] & 0x01
}

type MsgRaw struct {
	Header MsgHeader
	Body   []byte
}

func NewMsgRaw(h *MsgHeader) *MsgRaw {
	return nil
}

func (m *MsgRaw) Len() int {
	return len(m.Body)
}

func (m *MsgRaw) Serialize() {
	if m.Body != nil {
		m.Body[0] = byte(m.Header.Length >> 8)
		m.Body[1] = uint8(m.Header.Length)
		m.Body[2] = byte(m.Header.Type)
		m.Body[3] = byte((m.Header.Type >> 4) & 0xf0)
	}
}

func (m *MsgRaw) GetMsgType() uint16 {
	return (uint16(m.Body[3])&0xf0)<<4 + uint16(m.Body[2])
}
