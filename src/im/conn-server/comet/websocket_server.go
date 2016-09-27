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
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	"im/common/proto/entity/msg/msgbase"
	"im/conn-server/conf"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	options      *ServerOptions
	httpServer   *http.Server
	connManager  *ConnManager
	upgrader     *websocket.Upgrader
	codec        *WebSocketCodec
	readCallBack func(c *WebSocketSession, p *msg.MsgRaw)
	shutDown     bool
	Stat         *CometStat //流量状态
}

func NewWebSocketServer(op *ServerOptions) *WebSocketServer {
	var up = &websocket.Upgrader{
		ReadBufferSize:  op.ReadBufferSize,
		WriteBufferSize: op.WriteBufferSize,
	}
	up.CheckOrigin = func(r *http.Request) bool {
		// allow all connections by default
		return true
	}
	s := &WebSocketServer{
		options:      op,
		httpServer:   nil,
		shutDown:     true,
		readCallBack: nil,
		upgrader:     up,
		codec:        nil,
		connManager:  op.Cm,
		Stat:         NewCometStat(conf.Conf.ServerId)}
	return s
}

func (s *WebSocketServer) Start() {

	// handle
	hs := http.NewServeMux()
	hs.HandleFunc("/ws", s.serveWebSocket)

	s.httpServer = &http.Server{
		Handler:      hs,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go s.ListenAndServer()
}

func (s *WebSocketServer) ListenAndServer() error {

	s.httpServer.SetKeepAlivesEnabled(true)

	// Listener
	addr, err := net.ResolveTCPAddr("tcp4", s.options.Addr)
	if err != nil {
		log.Error("[WebSocketServer|ResolveTCPAddr] Resolve TCP Addr %s error: (%v)", s.options.Addr, err)
		return err
	}

	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Error("[WebSocketServer|ListenTCP] Listen TCP error: (%v)", err)
		return err
	}

	s.shutDown = false
	log.Debug("[WebSocketServer|Listen] listen: \"%s\"", addr.String())

	// Serve
	err = s.httpServer.Serve(l)
	if err != nil {
		log.Error("[WebSocketServer|Serve] error: (%v)", err)
		return err
	}

	return nil
}

func (s *WebSocketServer) StopListen() {
	s.shutDown = true
}

func (s *WebSocketServer) OnWebSocketConn(c *WebSocketSession, mb *msgbase.MsgBase) {
	s.connManager.PutConn(c.Id, c)
	atomic.AddUint32(&s.Stat.ConnCount, 1)
	log.Debug("[WebSocketServer|OnWebSocketConn]client id:%d", c.Id)

	if c.Id == 0 {
		c.Id = Ids.Get()
		s.connManager.PutConn(c.Id, c)
		atomic.AddUint32(&s.Stat.ConnCount, 1)
		log.Debug("[TcpServer|OnTcpConn]client id:%d", c.Id)
	}

	mb.Connid = c.Id
	mb.ConnServerid = uint16(conf.Conf.ServerId)
}

func (s *WebSocketServer) OnWebSocketClose(c *WebSocketSession) {
	s.connManager.DelConn(c.Id)
	atomic.AddUint32(&s.Stat.ConnClosedCount, 1)
	log.Debug("[WebSocketServer|OnWebSocketConn] id:%d", c.Id)
}

func (s *WebSocketServer) SetReadCallBack(cb func(c *WebSocketSession, p *msg.MsgRaw)) {
	if s.shutDown {
		s.readCallBack = cb
	} else {
		panic("")
	}
}

func (s *WebSocketServer) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("[WebSocketServer|serveWebSocket] %v", err)
		return
	}
	c := NewWebSocketSession(0, ws, s)
	c.ServeWebSocket()
}
