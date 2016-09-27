/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description: 提供集群唯一ID服务
 ******************************************************************/

package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Error type
var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

// Conn
type Conn struct {
	srv         *Server
	conn        *net.TCPConn
	extraData   interface{}
	closeOnce   sync.Once
	closeFlag   int32
	closeChan   chan struct{}
	sendChan    chan Packet
	receiveChan chan Packet
}

// connection callbacks
type ConnCallback interface {
	OnConnect(*Conn) bool
	OnMessage(*Conn, Packet) bool
	OnClose(*Conn)
}

// newConn returns a wrapper of raw conn
func NewConn(conn *net.TCPConn, srv *Server) *Conn {
	return &Conn{
		srv:         srv,
		conn:        conn,
		closeChan:   make(chan struct{}),
		sendChan:    make(chan Packet, srv.config.SendChanCount),
		receiveChan: make(chan Packet, srv.config.RecvChanCount),
	}
}

func (c *Conn) GetExtraData() interface{} {
	return c.extraData
}

func (c *Conn) PutExtraData(data interface{}) {
	c.extraData = data
}

func (c *Conn) GetRawConn() *net.TCPConn {
	return c.conn
}

func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		atomic.StoreInt32(&c.closeFlag, 1)
		close(c.closeChan)
		c.conn.Close()
		c.srv.callback.OnClose(c)
	})
}

func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

func (c *Conn) AsyncReadPacket(timeout time.Duration) (Packet, error) {
	if c.IsClosed() {
		return nil, ErrConnClosing
	}

	if timeout == 0 {
		select {
		case p := <-c.receiveChan:
			return p, nil

		default:
			return nil, ErrReadBlocking
		}

	} else {
		select {
		case p := <-c.receiveChan:
			return p, nil

		case <-c.closeChan:
			return nil, ErrConnClosing

		case <-time.After(timeout):
			return nil, ErrReadBlocking
		}
	}
}

func (c *Conn) AsyncWritePacket(p Packet, timeout time.Duration) error {
	if c.IsClosed() {
		return ErrConnClosing
	}

	if timeout == 0 {
		select {
		case c.sendChan <- p:
			return nil

		default:
			return ErrWriteBlocking
		}

	} else {
		select {
		case c.sendChan <- p:
			return nil

		case <-c.closeChan:
			return ErrConnClosing

		case <-time.After(timeout):
			return ErrWriteBlocking
		}
	}
}

// Start
func (c *Conn) Start() {
	if !c.srv.callback.OnConnect(c) {
		return
	}

	//go c.handleLoop()
	go c.readLoop()
	go c.writeLoop()
}

func (c *Conn) readLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("readLoop exit", err)
		}
		c.Close()
		c.srv.waitGroup.Done()
	}()

	for {
		select {
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		default:
		}

		p, err := c.srv.protocol.ReadPacket(c.conn)
		if err != nil {
			fmt.Println("readLoop ReadPacket err:", err)
			return
		}

		//c.receiveChan <- p
		if !c.srv.callback.OnMessage(c, p) {
			fmt.Println("handleLoop OnMessage error")
			return
		}
	}
}

func (c *Conn) writeLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("writeLoop exit", err)
		}
		c.Close()
		c.srv.waitGroup.Done()
	}()

	for {
		select {
		case p := <-c.sendChan:
			if _, err := c.conn.Write(p.Serialize()); err != nil {
				return
			}

		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return

		}
	}
}

func (c *Conn) handleLoop() {
	c.srv.waitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("handleLoop exit", err)
		}
		c.Close()
		c.srv.waitGroup.Done()
	}()

	for {
		select {
		case p := <-c.receiveChan:
			if !c.srv.callback.OnMessage(c, p) {
				fmt.Println("handleLoop OnMessage error")
				return
			}
		case <-c.srv.exitChan:
			return

		case <-c.closeChan:
			return
		}
	}
}
