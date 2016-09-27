/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : config.go
 *  Date   :
 *  Author : yangl
 *  Description: 提供集群唯一ID服务
 ******************************************************************/

package main

import (
	"fmt"
	. "im/common/id"
	"net"
	"sync"
	"time"
)

var SfsServs []*Sonyflake

type Option struct {
	SendChanCount uint32
	RecvChanCount uint32
}

type Server struct {
	config    *Option
	callback  ConnCallback
	protocol  Protocol
	exitChan  chan struct{}
	waitGroup *sync.WaitGroup
}

// NewServer creates a server
func NewServer(op *Option, callback ConnCallback, protocol Protocol) *Server {

	// Init Sonyflake
	SfsServs = make([]*Sonyflake, 10)
	for i := 0; i < 10; i++ {
		var st Settings
		SfsServs[i] = NewSonyflake(st)
		if SfsServs[i] == nil {
			panic("sonyflake not created")
		}
	}

	return &Server{
		config:    op,
		callback:  callback,
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

// Start starts service
func (s *Server) Start(listener *net.TCPListener, acceptTimeout time.Duration) {
	s.waitGroup.Add(1)
	defer func() {
		listener.Close()
		s.waitGroup.Done()
	}()

	for {
		select {
		case <-s.exitChan:
			return

		default:
		}

		//listener.SetDeadline(time.Now().Add(acceptTimeout))
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("conn err: %v", err)
			continue
		}
		NewConn(conn, s).Start()
	}
}

// Stop stops service
func (s *Server) Stop() {
	close(s.exitChan)
	s.waitGroup.Wait()
}
