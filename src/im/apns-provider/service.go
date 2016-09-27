/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : service.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	"sync"
)

type Service struct {
	appid    int
	sendChan chan *Notification
	clients  []*Client
	stop     bool
	option   *Options
	wg       sync.WaitGroup
}

func NewService(op *Options) *Service {

	if op.sendChanSize < 1024 {
		op.sendChanSize = 1024
	}

	if op.connCount < 1 {
		op.connCount = 1
	}

	return &Service{
		appid:    op.id,
		sendChan: make(chan *Notification, op.sendChanSize),
		stop:     true,
		option:   op,
		clients:  make([]*Client, op.connCount),
	}
}

func (s *Service) Start() error {
	if s.stop {
		cert, pemErr := LoadPemFile(s.option.pemFile, s.option.pemPasswd)
		if pemErr != nil {
			log.Error("[Service|LoadPemFile]Cert Error:", pemErr)
			return pemErr
		}

		s.stop = false

		for i := 0; i < s.option.connCount; i++ {
			s.clients[i] = NewClient(cert)

			s.wg.Add(1)
			go s.sendLoop(s.clients[i])
		}
	}

	return nil
}

func (s *Service) sendLoop(c *Client) {
	defer func() {
		s.wg.Done()
		if err := recover(); err != nil {
			log.Error("[Service|sendLoop|recover]%s : %s", s.option.name, err)
		}
	}()

	for !s.stop {
		n, ok := <-s.sendChan
		if !ok {
			return
		}

		if n.Dev {
			c.Development()
		} else {
			c.Production()
		}

		//发送
		res, err := c.Push(n)
		log.Debug("[Service|sendLoop|Push] token:%s, return id:%s", n.DeviceToken, res.ApnsID)

		//失败是否重发？
		if err != nil {
			log.Error("[Service|sendLoop|Push]: %s", s.option.name, err.Error())
			return
		}

		if !res.Sent() {
			log.Error("[Service|sendLoop|Push Res]: %s", s.option.name, res.Reason)
		}
	}
}

func (s *Service) Stop() {
	s.stop = true
	close(s.sendChan)
	s.wg.Wait()
}

func (s *Service) Send(m *msg.ApnsTransMsg) error {
	defer func() {
		if err := recover(); nil != err {
			log.Error("[Service|Send|recover]%s", err)
		}
	}()

	for _, v := range m.Tokens {
		n := &Notification{
			DeviceToken: v,
			Topic:       s.option.bundleid,
			Payload:     m.Payload,
			Dev:         m.Dev,
		}
		select {
		case s.sendChan <- n:
		default:
			log.Error("[Service|Send] may be channel is full")
			//break
		}

	}
	return nil
}
