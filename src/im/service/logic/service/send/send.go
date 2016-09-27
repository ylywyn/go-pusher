/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package send

import (
	"github.com/pquerna/ffjson/ffjson"
	//"errors"
	"fmt"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	"im/common/pump"
	"im/service/conf"
	"im/service/logic/models/group"
	"im/service/logic/service/session"
)

type Sender struct {
	pump        pump.IBasePump
	pumpAddr    string
	apnsChannle string
}

var Send *Sender

func NewSender(c *conf.Config) *Sender {
	pump := pump.NewPump(c.PumpType)
	pump.EnablePub()

	s := &Sender{
		pump:     pump,
		pumpAddr: c.PumpAddr,
	}
	return s
}

func (this *Sender) Start() {
	this.pump.Connect(this.pumpAddr)
}

func (this *Sender) Stop() {
	this.pump.Stop()
}

//////////////////////////////////////////
//send
func (this *Sender) Send(h *msg.MsgHeader, msg *mb.MsgBase) error {
	switch h.Type {
	case types.GeneralMsgTypeMT_GENERAL_MSG:
		return this.SendOne(h, msg)
	case types.GeneralMsgTypeMT_GENERAL_GROUP_MSG:
		return this.SendToGroup(h, msg)
	default:

	}
	return nil
}

//////////////////////////////////////////
// send to one
func (this *Sender) SendOne(h *msg.MsgHeader, msg *mb.MsgBase) error {
	if msg.To != 0 {
		s := &session.Session{Uid: msg.To}
		err := s.Get()
		if err != nil {
			return err
		}

		return this.SendWithSession(h, msg, s)
	}
	return nil
}

func (this *Sender) SendWithSession(h *msg.MsgHeader, msg *mb.MsgBase, s *session.Session) error {
	if msg.Platform == types.PlatformPF_ALL {
		if s.OsType == types.PlatformPF_IOS {
			this.SendMbToApns(msg, s)
		} else {
			this.SendMbToComet(h, msg, s)
		}
	} else if msg.Platform == types.PlatformPF_IOS {
		this.SendMbToApns(msg, s)
	} else {
		this.SendMbToComet(h, msg, s)
	}
	return nil
}

func (this *Sender) SendMbToComet(h *msg.MsgHeader, mb *mb.MsgBase, s *session.Session) {
	if s.ConnSessionId < 1 || s.ConnClusterId < 1 {
		return
	}
	mb.Connid = s.ConnSessionId
	data := mb.Serialize()

	m := msg.MsgRaw{Body: data}
	m.Header = *h
	m.Header.Length = uint16(len(data) - 4)
	m.Serialize()

	this.SendMrToComet(&m, s.ConnClusterId)
}

func (this *Sender) SendMrToComet(m *msg.MsgRaw, cid uint32) {
	c := fmt.Sprintf(pump.MSG_CHANNEL_FRONT, cid)
	ret := this.pump.Pub(c, m.Body)
	if !ret {
		log.Error("[Sender|sendConnServer] pub error, may be channel is full")
	}
}

func (this *Sender) SendMbToApns(mb *mb.MsgBase, s *session.Session) {
	if s.ApnsToken == "" {
		return
	}
	//构造APNS格式
	var pl []byte
	if mb.Payload != nil {
		apt := msg.ApsPayloadTrans{
			Aps: msg.ApsSimple{
				Alert:             string(mb.Text),
				Sound:             mb.Payload.Aps.Sound,
				Badge:             mb.Payload.Aps.Badge,
				Category:          mb.Payload.Aps.Category,
				Content_available: mb.Payload.Aps.ContentavAilable,
			},
			Extras: msg.ApsExtras{
				Type:   mb.Type,
				From:   mb.From,
				To:     mb.To,
				Gid:    mb.Gid,
				Time:   mb.Time,
				Msgid:  mb.Msgid,
				Extras: mb.Payload.Extras,
			},
		}
		pl, _ = ffjson.Marshal(apt)
	} else {
		apt := msg.ApsPayloadTrans{
			Aps: msg.ApsSimple{
				Alert: string(mb.Text),
			},
			Extras: msg.ApsExtras{
				Type:  mb.Type,
				From:  mb.From,
				To:    mb.To,
				Gid:   mb.Gid,
				Time:  mb.Time,
				Msgid: mb.Msgid,
			},
		}
		pl, _ = ffjson.Marshal(apt)
	}

	tokens := []string{s.ApnsToken}

	atm := msg.ApnsTransMsg{
		Dev:        false,
		Appid:      mb.Appid,
		ExpireTime: 0,
		Payload:    string(pl),
		Tokens:     tokens,
	}

	data, _ := ffjson.Marshal(atm)
	c := fmt.Sprintf(pump.APNS_CHANNEL_FORMAT, mb.Appid)

	ret := this.pump.Pub(c, data)
	if !ret {
		log.Error("[Sender|sendApnsServer] pub error, may be channel is full")
	}
}

//////////////////////////////////////////
//send to group
func (this *Sender) SendToGroup(h *msg.MsgHeader, mbmsg *mb.MsgBase) error {
	if mbmsg.Gid != 0 {
		// 获取组成员
		gm := group.Group{Gid: mbmsg.Gid, Appid: uint32(mbmsg.Appid)}
		err := gm.GetMembers()
		if err != nil {
			return err
		}
		if gm.Members == nil {
			log.Info("[msg|GeneralGroupMsg] group:%d can't find members ", mbmsg.Gid)
			return nil
		}
		log.Info("[msg|GeneralGroupMsg] group:%v", gm.Members)
		//分别发送
		n := len(gm.Members)
		for i := 0; i < n; i++ {
			mbmsg.To = gm.Members[i]
			//mbmsg.To = gm.Members[i]
			err = this.SendOne(h, mbmsg)
			if err != nil {
				log.Info("[msg|GeneralGroupMsg] groupmsg . one :%d, %s", mbmsg.To, err.Error())
			}
		}
	}
	return nil
}
