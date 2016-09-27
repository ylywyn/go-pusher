// automatically generated, do not modify

package session

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Session struct {
	_tab flatbuffers.Table
}

func (rcv *Session) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Session) LoginType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Session) Admin() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Session) Appid() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Session) Uid() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Session) AckMsgId() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Session) Token() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Session) Key() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Session) ConnClusterId() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Session) ConnSessionId() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func SessionStart(builder *flatbuffers.Builder) { builder.StartObject(9) }
func SessionAddLoginType(builder *flatbuffers.Builder, LoginType byte) { builder.PrependByteSlot(0, LoginType, 0) }
func SessionAddAdmin(builder *flatbuffers.Builder, Admin byte) { builder.PrependByteSlot(1, Admin, 0) }
func SessionAddAppid(builder *flatbuffers.Builder, Appid uint16) { builder.PrependUint16Slot(2, Appid, 0) }
func SessionAddUid(builder *flatbuffers.Builder, Uid uint64) { builder.PrependUint64Slot(3, Uid, 0) }
func SessionAddAckMsgId(builder *flatbuffers.Builder, AckMsgId uint64) { builder.PrependUint64Slot(4, AckMsgId, 0) }
func SessionAddToken(builder *flatbuffers.Builder, Token flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(Token), 0) }
func SessionAddKey(builder *flatbuffers.Builder, Key flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(Key), 0) }
func SessionAddConnClusterId(builder *flatbuffers.Builder, ConnClusterId uint32) { builder.PrependUint32Slot(7, ConnClusterId, 0) }
func SessionAddConnSessionId(builder *flatbuffers.Builder, ConnSessionId uint64) { builder.PrependUint64Slot(8, ConnSessionId, 0) }
func SessionEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
