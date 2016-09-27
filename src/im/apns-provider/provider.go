/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : provider.go
 *  Date   :
 *  Author : yangl
 *  Description: receive msg and send
 ******************************************************************/

package main

import (
	"fmt"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	"im/common/pump"

	"github.com/pquerna/ffjson/ffjson"
)

type ApnsProvider struct {
	closed   bool
	pump     pump.IBasePump
	servMaps map[int]*Service
}

func NewApnsProvider() *ApnsProvider {
	// 消息泵
	msgPump := pump.NewPump(Conf.pumpType)

	return &ApnsProvider{
		closed:   true,
		pump:     msgPump,
		servMaps: nil,
	}
}

func (p *ApnsProvider) Start() {
	if p.closed {
		// msg pump
		p.pump.BindSubProcess(p.Send)
		p.pump.EnableSub()
		p.pump.Connect(Conf.pumpAddr)

		//service
		count := len(Conf.servesOpts)
		if count == 0 {
			log.Error("[ApnsProvider|Start] services options is 0")
			return
		}

		p.servMaps = make(map[int]*Service, count)
		for _, v := range Conf.servesOpts {
			s := NewService(v)
			p.servMaps[v.id] = s

			err := s.Start()
			if err != nil {
				log.Error("[ApnsProvider|Start] service start error id %d : %s ", v.id, err.Error())
				continue
			} else {
				log.Debug("[ApnsProvider|Start] service start id : %d", v.id)
			}

			substr := fmt.Sprintf(Conf.subString, v.id)
			err = p.pump.Sub(substr, substr)
			if err != nil {
				log.Error("[ApnsProvider|Start] sub error : %s ", err.Error())
				continue
			}
			log.Debug("[ApnsProvider|Start] sub : %s", substr)
		}
		p.closed = false
	}
}

func (p *ApnsProvider) Stop() {
	if !p.closed {
		p.closed = true
		p.pump.Stop()

		for _, s := range p.servMaps {
			s.Stop()
		}
	}
}

//处理接到的消息, data为ApnsTransMsg 的json格式
func (p *ApnsProvider) Send(data []byte) {

	m := &msg.ApnsTransMsg{}
	err := ffjson.Unmarshal(data, m)
	if err != nil {
		log.Error("[ApnsProvider|Send|Unmarshal] error:%s", err.Error())
		return
	}

	if s, ok := p.servMaps[int(m.Appid)]; ok {
		err = s.Send(m)
		if err != nil {
			log.Error("[ApnsProvider|Send|ServiceSend] error: %s", err.Error())
			return
		}
	} else {
		log.Error("[ApnsProvider|Send] error: bad appid")
		return
	}
}
