/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package comet

import (
	//"container/list"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/entity/msg/status"
	"im/common/pump"
	"im/common/store"
	"sync"
)

type Processor struct {
	pump        pump.IBasePump
	connManager *ConnManager
	store       *store.FileStore
	slock       sync.Mutex
}

var ProcessInstance *Processor

func NewProcessor(bp pump.IBasePump, cm *ConnManager) *Processor {
	p := &Processor{
		pump:        bp,
		connManager: cm,
	}
	return p
}

// 处理从连接处收到的 消息
func (this *Processor) ProcessRecv(m *msg.MsgRaw) {

	if m.Header.Type == 0 {
		log.Error("[Processor|ProcessRecv] bad msg type!")
		return
	}

	if this.pump.Pub(pump.MSG_CHANNEL_BACK, m.Body) {
		return
	}

	log.Warn("[Processor|ProcessRecv] too much msg, save to file")

	// 准备写入文件暂存
	this.PreparePersist(m)
}

func (this *Processor) PreparePersist(m *msg.MsgRaw) {

	//启动持久化机制
	this.slock.Lock()
	if this.store == nil {
		this.store = store.NewFileStore()
		go this.saving()
	}
	this.slock.Unlock()

	this.store.Write(m.Body)
}

func (this *Processor) saving() {
	l := this.pump.GetAllPubs()
	this.store.WriteList(l)
}

// 处理 从后端发往连接的消息
func (this *Processor) ProcessSend(data []byte) {
	connid, to := mb.UnSerializeConnId(data[msg.HEADER_LEN:])

	msg := &msg.MsgRaw{Body: data}

	c, err := this.connManager.GetConn(connid)
	if err == nil {
		c.Write(msg)
	} else {
		log.Debug("[Processor|ProcessSend] can't find id:%d", to)
		//清除session
		m := status.NewSessionClearMsg(to)
		this.pump.Pub(pump.MSG_CHANNEL_BACK, m.Body)
	}
}
