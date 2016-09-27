// automatically generated, do not modify

package msgbase

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Msg struct {
	_tab flatbuffers.Table
}

func (rcv *Msg) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Msg) Type() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) To() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) Connid() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) ConnServerid() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) Gid() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) Text(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j * 1))
	}
	return 0
}

func (rcv *Msg) TextLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Msg) TextBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Msg) Time() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) Msgid() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) From() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) Appid() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) Platform() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Msg) Payload(obj *PayLoad) *PayLoad {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(PayLoad)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Msg) Option(obj *Options) *Options {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Options)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func MsgStart(builder *flatbuffers.Builder) { builder.StartObject(13) }
func MsgAddType(builder *flatbuffers.Builder, Type uint16) { builder.PrependUint16Slot(0, Type, 0) }
func MsgAddTo(builder *flatbuffers.Builder, To uint64) { builder.PrependUint64Slot(1, To, 0) }
func MsgAddConnid(builder *flatbuffers.Builder, Connid uint64) { builder.PrependUint64Slot(2, Connid, 0) }
func MsgAddConnServerid(builder *flatbuffers.Builder, ConnServerid uint16) { builder.PrependUint16Slot(3, ConnServerid, 0) }
func MsgAddGid(builder *flatbuffers.Builder, Gid uint64) { builder.PrependUint64Slot(4, Gid, 0) }
func MsgAddText(builder *flatbuffers.Builder, Text flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(Text), 0) }
func MsgStartTextVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(1, numElems, 1)
}
func MsgAddTime(builder *flatbuffers.Builder, Time uint64) { builder.PrependUint64Slot(6, Time, 0) }
func MsgAddMsgid(builder *flatbuffers.Builder, Msgid uint64) { builder.PrependUint64Slot(7, Msgid, 0) }
func MsgAddFrom(builder *flatbuffers.Builder, From uint64) { builder.PrependUint64Slot(8, From, 0) }
func MsgAddAppid(builder *flatbuffers.Builder, Appid uint16) { builder.PrependUint16Slot(9, Appid, 0) }
func MsgAddPlatform(builder *flatbuffers.Builder, Platform byte) { builder.PrependByteSlot(10, Platform, 0) }
func MsgAddPayload(builder *flatbuffers.Builder, Payload flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(11, flatbuffers.UOffsetT(Payload), 0) }
func MsgAddOption(builder *flatbuffers.Builder, Option flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(12, flatbuffers.UOffsetT(Option), 0) }
func MsgEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
