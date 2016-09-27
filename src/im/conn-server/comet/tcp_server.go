/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : tcp_server.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package comet

import (
	//"encoding/json"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	"im/conn-server/conf"
	"net"
	"sync"
	"sync/atomic"
)

type TcpServer struct {
	options      *ServerOptions
	listener     *TcpListener
	shutDown     bool
	shutdownChan chan bool
	waitGroup    *sync.WaitGroup
	connManager  *ConnManager
	codec        *TcpLenCodec
	readCallBack func(c *TcpSession, p *msg.MsgRaw)

	Stat *CometStat //流量状态
}

func NewTcpServer(op *ServerOptions) *TcpServer {
	s := &TcpServer{
		options:      op,
		listener:     nil,
		shutDown:     true,
		shutdownChan: make(chan bool),
		waitGroup:    &sync.WaitGroup{},
		readCallBack: nil,
		codec:        &TcpLenCodec{},
		connManager:  op.Cm,
		Stat:         NewCometStat(conf.Conf.ServerId),
	}
	return s
}

func (s *TcpServer) ListenAndServer() error {
	addr, err := net.ResolveTCPAddr("tcp4", s.options.Addr)
	if err != nil {
		log.Error("[TcpServer|Listen] Resolve TCP Addr %s error: (%v)", s.options.Addr, err)
		return err
	}

	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Error("[TcpServer|Listen] Listen TCP error: (%v)", err)
		return err
	}

	s.shutDown = false
	log.Debug("[TcpServer|Listen] listen: \"%s\"", addr.String())

	s.listener = &TcpListener{l, s.shutdownChan}

	// split N core accept
	var i uint32 = 1
	rn := TcpAcceptRouts + 1
	for ; i < uint32(rn); i++ {
		go s.acceptTCP(i, s.listener)
	}
	return nil
}

func (s *TcpServer) acceptTCP(i uint32, l *TcpListener) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()

	var (
		conn *net.TCPConn
		err  error
		rawI = i
	)

	for {
		conn, err = l.Accept()
		if err != nil {
			log.Error("[TcpServer|acceptTCP] Listener error(%v)", err)
			continue
		}

		tcpconn := NewTcpSession(i, conn, s)
		go tcpconn.ServeTCP()

		i += TcpAcceptRouts
		if i >= MaxIndex {
			i = rawI
		}
	}
}

func (s *TcpServer) StopListen() {
	s.shutDown = true
	close(s.shutdownChan)
	s.listener.Close()

	s.waitGroup.Wait()
}

func (s *TcpServer) OnTcpConn(c *TcpSession, m *msg.MsgRaw) {
	if c.Id == 0 {
		c.Id = Ids.Get()
		s.connManager.PutConn(c.Id, c)
		atomic.AddUint32(&s.Stat.ConnCount, 1)
		log.Debug("[TcpServer|OnTcpConn]client id:%d", c.Id)
	}

	mb, err := s.codec.UnmarshalPacket(m)
	if nil != err {
		log.Error("[TcpServer|OnTcpConn|codec.UnmarshalPacket] %s", err.Error())
	}

	mb.Connid = c.Id
	mb.ConnServerid = uint16(conf.Conf.ServerId)

	data := mb.Serialize()
	for i := 0; i < msg.HEADER_LEN; i++ {
		data[i] = m.Body[i]
	}
	m.Body = data

	////////////////////
	//
	//	d, _ := json.Marshal(mb)
	//	if d != nil {
	//		log.Debug("[TcpServer|OnTcpConn] %s", string(d))
	//	}
}

func (s *TcpServer) OnTcpConnClose(c *TcpSession) {
	s.connManager.DelConn(c.Id)
	atomic.AddUint32(&s.Stat.ConnClosedCount, 1)
	log.Debug("[TcpServer|OnTcpConnClose] id:%d", c.Id)
}

func (s *TcpServer) SetReadCallBack(cb func(c *TcpSession, p *msg.MsgRaw)) {
	if s.shutDown {
		s.readCallBack = cb
	} else {
		panic("")
	}
}
