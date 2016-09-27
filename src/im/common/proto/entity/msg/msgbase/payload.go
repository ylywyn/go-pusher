/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package msgbase

import (
	"bytes"
	mb "im/common/proto/fbsgen/msg/msgbase"

	fb "github.com/google/flatbuffers/go"
)

// PayLoad
type PayLoad struct {
	Aps     *ApsPayLoad
	General string
	Extras  string
}

func (this *PayLoad) Serialize(builder *fb.Builder) fb.UOffsetT {

	var General_offset, Extrasoffset, Aps_offset fb.UOffsetT
	if len(this.General) != 0 {
		General_offset = builder.CreateString(this.General)
	}
	if len(this.Extras) != 0 {
		buf := bytes.NewBufferString(this.Extras)
		Extrasoffset = builder.CreateByteVector(buf.Bytes())
	}

	if this.Aps != nil {
		Aps_offset = this.Aps.Serialize(builder)
	}

	mb.PayLoadStart(builder)
	if General_offset > 0 {
		mb.PayLoadAddGeneral(builder, General_offset)
	}

	if Aps_offset > 0 {
		mb.PayLoadAddAps(builder, Aps_offset)
	}
	if Extrasoffset > 0 {
		mb.PayLoadAddExtras(builder, Extrasoffset)
	}

	return mb.PayLoadEnd(builder)
}
