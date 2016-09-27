/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : process.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package msg

import (
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	"im/common/pump"
	"im/common/store"
	"im/service/conf"
	smsg "im/service/msg/msg"
	"im/service/msg/status"
	"im/service/msg/user"
	"sync"
)

const MAX_ROUTINS = 50000

type Processor struct {
	pump            pump.IBasePump
	pumpAddr        string
	maxRoutines     int
	maxRoutinesChan chan bool
	store           *store.FileStore
	slock           sync.Mutex
}

func NewProcessor(c *conf.Config) *Processor {
	pump := pump.NewPump(c.PumpType)
	pump.EnableSub()

	p := &Processor{
		pump:            pump,
		pumpAddr:        c.PumpAddr,
		maxRoutines:     MAX_ROUTINS,
		maxRoutinesChan: make(chan bool, MAX_ROUTINS),
	}

	pump.BindSubProcess(p.processSub)
	return p
}

func (this *Processor) Start() {
	this.pump.Connect(this.pumpAddr)
	this.pump.Sub(pump.MSG_CHANNEL_BACK, pump.MSG_CHANNEL_BACK)
}

func (this *Processor) Stop() {
	this.pump.Stop()
}

// 处理从连接处收到的 消息
func (this *Processor) processSub(d []byte) {
	m := &msg.MsgRaw{Body: d}

	var h [msg.HEADER_LEN]byte
	for i := 0; i < msg.HEADER_LEN; i++ {
		h[i] = d[i]
	}
	m.Header.Parse(h)

	this.maxRoutinesChan <- true
	go this.processMsg(m)
}

//使用指定数量的 Routines来处理消息
func (this *Processor) processMsg(m *msg.MsgRaw) {

	defer func() {
		<-this.maxRoutinesChan
		if err := recover(); err != nil {
			log.Error("[Processor|processMsg|recover] %s", err)
		}
	}()

	//log.Debug("[Processor|processMsg] type:%d", m.Header.Type)

	t := m.Header.Type / 100
	var err error
	switch t {
	case 0:
		err = smsg.ProcessMsg(m)
	case 1:
		err = user.ProcessMsg(m)
	case 4:
		err = status.ProcessMsg(m)
	default:
		log.Error("[Processor|processMsg] unknow service type: %d", m.Header.Type)
	}

	if err != nil {
		log.Error("[Processor|processMsg] error: %s", err.Error())
	}
}
