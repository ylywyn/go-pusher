// automatically generated, do not modify

package msgbase

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type PayLoad struct {
	_tab flatbuffers.Table
}

func (rcv *PayLoad) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *PayLoad) General() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *PayLoad) Aps(obj *ApsPayLoad) *ApsPayLoad {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(ApsPayLoad)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *PayLoad) Extras() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func PayLoadStart(builder *flatbuffers.Builder) { builder.StartObject(3) }
func PayLoadAddGeneral(builder *flatbuffers.Builder, General flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(General), 0) }
func PayLoadAddAps(builder *flatbuffers.Builder, Aps flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(Aps), 0) }
func PayLoadAddExtras(builder *flatbuffers.Builder, Extras flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(Extras), 0) }
func PayLoadEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
