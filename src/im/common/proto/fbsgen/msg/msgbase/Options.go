// automatically generated, do not modify

package msgbase

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Options struct {
	_tab flatbuffers.Table
}

func (rcv *Options) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Options) TimeLive() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Options) StartTime() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Options) ApnsProduction() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Options) Command(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j * 1))
	}
	return 0
}

func (rcv *Options) CommandLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Options) CommandBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func OptionsStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func OptionsAddTimeLive(builder *flatbuffers.Builder, TimeLive uint32) { builder.PrependUint32Slot(0, TimeLive, 0) }
func OptionsAddStartTime(builder *flatbuffers.Builder, StartTime uint32) { builder.PrependUint32Slot(1, StartTime, 0) }
func OptionsAddApnsProduction(builder *flatbuffers.Builder, ApnsProduction byte) { builder.PrependByteSlot(2, ApnsProduction, 0) }
func OptionsAddCommand(builder *flatbuffers.Builder, Command flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(Command), 0) }
func OptionsStartCommandVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(1, numElems, 1)
}
func OptionsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
