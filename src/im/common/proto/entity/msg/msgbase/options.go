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
	mb "im/common/proto/fbsgen/msg/msgbase"

	fb "github.com/google/flatbuffers/go"
)

// Options
type Options struct {
	TimeLive       uint32
	StartTime      uint32
	ApnsProduction bool
	Command        []byte
}

func (this *Options) Serialize(builder *fb.Builder) fb.UOffsetT {
	var prod byte
	var Commandoffset fb.UOffsetT

	if this.ApnsProduction {
		prod = 1
	} else {
		prod = 0
	}

	if this.Command != nil {
		Commandoffset = builder.CreateByteVector(this.Command)
	}

	mb.OptionsStart(builder)
	mb.OptionsAddTimeLive(builder, this.TimeLive)
	mb.OptionsAddStartTime(builder, this.StartTime)
	mb.OptionsAddApnsProduction(builder, prod)
	if Commandoffset > 0 {
		mb.OptionsAddCommand(builder, Commandoffset)
	}
	return mb.OptionsEnd(builder)
}
