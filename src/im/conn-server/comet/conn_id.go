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
	"im/common/id"
	log "im/common/log4go"
)

//本机唯一即可
var Ids *IdBuilder

type IdBuilder struct {
	start  chan bool
	idchan chan uint64
	sf     *id.Sonyflake
}

func init() {
	var st id.Settings

	Ids = &IdBuilder{
		start:  make(chan bool),
		idchan: make(chan uint64, 4096),
		sf:     id.NewSonyflake(st),
	}

	go Ids.create()
}

func (this *IdBuilder) create() {
	this.fill()
	for {
		<-this.start
		this.fill()
	}
}

func (this *IdBuilder) Get() uint64 {
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
