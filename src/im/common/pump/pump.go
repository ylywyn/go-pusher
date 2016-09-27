/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package pump

import "container/list"

// channel定义
const (
	//%d: appid
	APNS_CHANNEL_FORMAT = "push_ioschannel_%d"
	MSG_CHANNEL_BACK    = "push_msgchannel_back"

	//%d: connserver的cluster id
	MSG_CHANNEL_FRONT   = "push_msgchannel_%d"
	MAX_PUBCHANNEL_SIZE = 20000
)

type IBasePump interface {
	Connect(s string)
	Stop()
	Sub(c, group string) error
	Pub(c string, m []byte) bool
	EnableSub()
	EnablePub()
	BindSubProcess(fun func(data []byte))
	GetAllPubs() *list.List
}

type pubItem struct {
	c    string
	data []byte
}

func NewPump(t string) IBasePump {
	var pump IBasePump
	switch t {
	case "nats":
		pump = &NatsPump{}
	case "redis":
		pump = &RedisPump{}
	default:
		pump = nil
	}

	return pump
}
