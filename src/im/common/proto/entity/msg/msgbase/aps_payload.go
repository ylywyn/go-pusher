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

// ApsPayLoad
type ApsPayLoad struct {
	Sound            string
	Badge            uint16
	ContentavAilable bool
	Category         string
}

func (this *ApsPayLoad) Serialize(builder *fb.Builder) fb.UOffsetT {
	var Sound_offset, Category_offset fb.UOffsetT
	if len(this.Sound) > 0 {
		Sound_offset = builder.CreateString(this.Sound)
	}
	if len(this.Category) > 0 {
		Category_offset = builder.CreateString(this.Category)
	}
	var ContentavAilable byte
	if this.ContentavAilable {
		ContentavAilable = 1
	} else {
		ContentavAilable = 0
	}

	mb.ApsPayLoadStart(builder)
	if Sound_offset > 0 {
		mb.ApsPayLoadAddSound(builder, Sound_offset)
	}
	mb.ApsPayLoadAddBadge(builder, this.Badge)
	mb.ApsPayLoadAddContentavAilable(builder, ContentavAilable)
	if Category_offset > 0 {
		mb.ApsPayLoadAddCategory(builder, Category_offset)
	}
	return mb.ApsPayLoadEnd(builder)
}
