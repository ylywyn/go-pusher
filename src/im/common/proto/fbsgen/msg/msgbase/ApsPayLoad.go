// automatically generated, do not modify

package msgbase

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ApsPayLoad struct {
	_tab flatbuffers.Table
}

func (rcv *ApsPayLoad) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ApsPayLoad) Sound() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ApsPayLoad) Badge() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ApsPayLoad) ContentavAilable() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ApsPayLoad) Category() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ApsPayLoadStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func ApsPayLoadAddSound(builder *flatbuffers.Builder, Sound flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(Sound), 0) }
func ApsPayLoadAddBadge(builder *flatbuffers.Builder, Badge uint16) { builder.PrependUint16Slot(1, Badge, 0) }
func ApsPayLoadAddContentavAilable(builder *flatbuffers.Builder, ContentavAilable byte) { builder.PrependByteSlot(2, ContentavAilable, 0) }
func ApsPayLoadAddCategory(builder *flatbuffers.Builder, Category flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(Category), 0) }
func ApsPayLoadEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
