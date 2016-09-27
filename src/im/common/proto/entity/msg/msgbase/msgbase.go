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
	"im/common/proto/fbsgen/msg/types"
	//log "code.google.com/p/log4go"
	fb "github.com/google/flatbuffers/go"
)

// 通讯的最终Msg
type MsgBase struct {
	Type         uint16
	Appid        uint16
	From         uint64
	To           uint64
	Connid       uint64
	ConnServerid uint16
	Gid          uint64
	Text         string
	Time         uint64
	Msgid        uint64
	Platform     uint8
	Payload      *PayLoad
	Options      *Options
}

//让出了4个字节 留给头部
func (msg *MsgBase) Serialize() []byte {

	builder := fb.NewBuilder(256)

	var pl, op, text fb.UOffsetT

	// Payload
	if msg.Payload != nil {
		pl = msg.Payload.Serialize(builder)
	}

	// Options
	if msg.Options != nil {
		op = msg.Options.Serialize(builder)
	}

	if len(msg.Text) > 0 {
		text = builder.CreateByteVector([]byte(msg.Text))
	}

	//msg start
	mb.MsgStart(builder)

	mb.MsgAddType(builder, msg.Type)
	mb.MsgAddTo(builder, msg.To)
	mb.MsgAddTime(builder, msg.Time)
	mb.MsgAddMsgid(builder, msg.Msgid)
	mb.MsgAddFrom(builder, msg.From)
	mb.MsgAddAppid(builder, msg.Appid)
	mb.MsgAddPlatform(builder, byte(msg.Platform))
	mb.MsgAddConnid(builder, msg.Connid)
	mb.MsgAddConnServerid(builder, msg.ConnServerid)
	mb.MsgAddGid(builder, msg.Gid)

	if text > 0 {
		mb.MsgAddText(builder, text)
	}

	if pl > 0 {
		mb.MsgAddPayload(builder, pl)
	}

	if op > 0 {
		mb.MsgAddOption(builder, op)
	}

	u := mb.MsgEnd(builder)
	builder.Finish(u)
	return builder.Bytes[builder.Head()-4:]
}

func (msg *MsgBase) UnSerialize(buf []byte) error {

	n := fb.GetUOffsetT(buf)
	m := mb.Msg{}
	m.Init(buf, n)

	msg.Type = m.Type()
	msg.Appid = m.Appid()
	msg.To = m.To()
	msg.From = m.From()
	msg.Platform = m.Platform()
	msg.Msgid = m.Msgid()
	msg.Time = m.Time()
	msg.Gid = m.Gid()
	msg.Connid = m.Connid()
	msg.ConnServerid = m.ConnServerid()

	l := m.TextLength()
	if l != 0 {
		msg.Text = string(m.TextBytes())
	}

	// option
	if op := m.Option(nil); op != nil {
		var prod bool = false
		if op.ApnsProduction() == 1 { //flatbuffer 的bool转成了byte，不知道有没有强转的好方法
			prod = true
		}
		msg.Options = &Options{op.TimeLive(), op.StartTime(), prod, op.CommandBytes()}
	}

	// payload
	if pl := m.Payload(nil); pl != nil {

		msg.Payload = &PayLoad{
			General: string(pl.General()),
			Extras:  string(pl.Extras())}

		// ios apns
		if (msg.Platform&types.PlatformPF_ALL) == types.PlatformPF_ALL ||
			(msg.Platform&types.PlatformPF_IOS) == types.PlatformPF_IOS {
			if aps := pl.Aps(nil); aps != nil {
				var ca bool = false
				if aps.ContentavAilable() == 1 {
					ca = true
				}
				msg.Payload.Aps = &ApsPayLoad{string(aps.Sound()), aps.Badge(), ca, string(aps.Category())}
			}
		}
	}

	return nil
}

func UnSerializeConnId(buf []byte) (uint64, uint64) {
	n := fb.GetUOffsetT(buf)
	m := mb.Msg{}
	m.Init(buf, n)

	return m.Connid(), m.To()
}
