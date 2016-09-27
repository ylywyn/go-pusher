/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description: ID服务，集群从网络获取，非集群本机产生
 ******************************************************************/

package id

import (
	cid "im/common/id"
	log "im/common/log4go"
	"im/id/client"
)

var clusterID bool

func InitId(cluster bool, addr string) error {
	if cluster {
		id.IdPool = &id.IdsPool{}
		err := id.IdPool.Init(addr)
		if err != nil {
			id.IdPool = nil
			return err
		}
		clusterID = true
	} else {
		initIdBuilder()
	}
	return nil
}

func GetId(key string) uint64 {
	if clusterID {
		n, err := id.IdPool.Id(key)
		if err != nil {
			n = 0
			log.Error("[id|GetId] error:%s", err.Error())
		}
		return n
	} else {
		return ids.get()
	}
}

//本集群，本机唯一ID
var ids *IdBuilder

type IdBuilder struct {
	start  chan bool
	idchan chan uint64
	sf     *cid.Sonyflake
}

func initIdBuilder() {
	var st cid.Settings

	ids = &IdBuilder{
		start:  make(chan bool),
		idchan: make(chan uint64, 4096),
		sf:     cid.NewSonyflake(st),
	}

	go ids.create()
}

func (this *IdBuilder) create() {
	this.fill()
	for {
		<-this.start
		this.fill()
	}
}

func (this *IdBuilder) get() uint64 {
	var id uint64 = 0

	for {
		select {
		case id = <-this.idchan:
			//log.Debug("[IdBuilder|Get] get id:%d", id)
			return id
		default:
			this.start <- true
		}
	}
}

func (this *IdBuilder) fill() {
	for {
		id, err := this.sf.NextID()
		if err != nil {
			log.Error("[IdBuilder|fill] error:%s", err.Error())
			continue
		}

		select {
		case this.idchan <- id:
			continue
		default:
			return
		}
	}
}
