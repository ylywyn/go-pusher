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
	"errors"
	"fmt"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	//"im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	//"io"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (

	// PONG
	pongWait = 180 * time.Second

	// PING
	pingPeriod = (pongWait * 9) / 10
)

//ReadChannel可去掉， wait fix
type WebSocketSession struct {
	Id           uint64
	Authed       bool
	status       int32 //此状态并没有加锁保护
	remoteAddr   string
	errMsg       string
	conn         *websocket.Conn
	WriteChannel chan *msg.MsgRaw

	server *WebSocketServer
}

func NewWebSocketSession(i uint32, c *websocket.Conn, s *WebSocketServer) *WebSocketSession {
	op := s.options

	if op.HeartBeat {
		c.SetReadDeadline(time.Now().Add(op.IdleTime))
		c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(op.IdleTime)); return nil })
	}

	c.SetReadLimit(msg.MAX_BODY_LEN)

	session := &WebSocketSession{
		Id:           0,
		Authed:       false,
		status:       StatusRunning,
		conn:         c,
		remoteAddr:   c.RemoteAddr().String(),
		WriteChannel: make(chan *msg.MsgRaw, op.WriteChannelSize),
		server:       s,
	}

	return session
}

func (this *WebSocketSession) LocalAddr() string {
	return this.conn.LocalAddr().String()
}

func (this *WebSocketSession) RemoteAddr() string {
	return this.conn.RemoteAddr().String()
}

//启动，开始读写操作
func (this *WebSocketSession) ServeWebSocket() {

	//写读数据
	go this.writeLoop()
	this.readLoop()

	s := atomic.LoadInt32(&this.status)
	if s > 0 {
		this.Close()
	}
}

//read
func (this *WebSocketSession) readLoop() {

	defer func() {
		if err := recover(); err != nil {
			log.Error("[WebSocketSession|readLoop|recover|%s] %s", this.remoteAddr, err)
		}
	}()

	codec := this.server.codec
	op := this.server.options
	stat := this.server.Stat

	for this.status == StatusRunning {

		_, r, err := this.conn.NextReader()
		if err != nil {
			//this.errMsg = err.Error()
			log.Error("[WebSocketSession|NextReader] %s", this.remoteAddr, err.Error())
			break
		}

		// 读一个完整的包
		mb, err := codec.Read(r)
		if nil != err {
			this.errMsg = err.Error()
			log.Error("[WebSocketSession|readLoop|codec.Read|%s] %s", this.remoteAddr, err.Error())
			break
		}

		//重置允许的IdleTime
		if op.HeartBeat {
			//没有心跳，不用
			//this.conn.SetReadDeadline(time.Now().Add(op.IdleTime))
			if this.Authed {
				//没有登录成功之前，不允许心跳
				//this.Conn.SetReadDeadline(time.Now().Add(op.IdleTime))
			}
			if mb == nil {
				continue
			}
		}

		//Josn转到MsgBase
		//mb, err := codec.UnmarshalToMb(m)

		if this.Id == 0 {
			// 第一个消息不是登陆消息，则关闭连接
			if mb.Type == types.UserMsgTypeMT_USER_LOGIN_REQ {
				this.server.OnWebSocketConn(this, mb)
			} else {
				this.errMsg = "must first auth conn, you type is "
				log.Error("[TcpSession|readLoop|m.Header.Type|%s] %d", this.remoteAddr, mb.Type)
				break
			}
		} else if mb.Type == types.UserMsgTypeMT_USER_LOGIN_REQ {
			this.server.OnWebSocketConn(this, mb)
		}

		//调试
		log.Debug("[WebSocketSession|Read]ConnId:%d, ConnStatus:%d", this.Id, this.status)

		m := codec.UnmarshalToMr(mb)
		m.Header.Type = mb.Type
		this.server.readCallBack(this, m)

		//读流量统计
		stat.FlowStat.IncrReadCounts()
		stat.FlowStat.IncrReadBytes(int32(m.Len()))
	}
	log.Debug("[WebSocketSession|writeLoop] readLoop over~")
	this.Close()
}

//write
func (this *WebSocketSession) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		this.conn.Close()
		if err := recover(); err != nil {
			log.Error("[WebSocketSession|writeLoop|recover|%s] %s", this.remoteAddr, err)
		}
	}()

	var p *msg.MsgRaw
	for this.status == StatusRunning {
		select {
		case p = <-this.WriteChannel:
			{
				if p != nil {
					//这里转成 Json bytes
					data, err := this.server.codec.MarshalToJson(p)
					if err != nil {
						log.Error("[WebSocketSession|writeLoop|codec.MarshalPacket|%s] %s", this.remoteAddr, err.Error())
						continue
					}
					if err := this.writeBytes(websocket.TextMessage, data); err != nil {
						log.Error("[WebSocketSession|writeLoop|writeBytes] %s", err.Error())
						this.errMsg = err.Error()
						break
					}
				} else {
					//channel已经关闭，说明已经关闭
					break
				}
			}
		case <-ticker.C:
			if err := this.writeBytes(websocket.PingMessage, []byte{}); err != nil {
				break
			}

		}

	}

	this.writeBytes(websocket.CloseMessage, []byte{})
	//this.Conn.
	log.Debug("[WebSocketSession|writeLoop] writeLoop over~")
}

func (this *WebSocketSession) Write(p *msg.MsgRaw) error {
	//防止写入已经关闭的WriteChannel引起的panic
	defer func() {
		if err := recover(); nil != err {
			log.Error("[WebSocketSession|Write|recover|%s]%s", this.remoteAddr, err)
		}
	}()

	if this.Authed == false {
		if p.GetMsgType() == types.UserMsgTypeMT_USER_LOGIN_REP {
			this.Authed = true
		} else {
			//登录失败，强制关闭连接？
		}
	}

	log.Debug("[WebSocketSession|Write]ConnId:%d, ConnStatus:%d", this.Id, this.status)
	if this.status == StatusRunning {
		select {
		case this.WriteChannel <- p:
			return nil
		default:
			return errors.New(fmt.Sprintf("[WebSocketSession|Write|%s] channel full", this.remoteAddr))
		}
	}
	return CONN_CLOSED_ERROR
}

//real write
func (this *WebSocketSession) writeBytes(mt int, data []byte) error {
	this.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
	err := this.conn.WriteMessage(mt, data)
	//this.conn.SetWriteDeadline(time.Time{})
	if err != nil {
		return err
	}
	this.server.Stat.FlowStat.IncrWriteCounts()
	this.server.Stat.FlowStat.IncrWriteBytes(int32(len(data)))
	return nil
}

//is close
func (this *WebSocketSession) Closed() bool {
	return atomic.LoadInt32(&this.status) == StatusClosed
}

//
func (this *WebSocketSession) Close() {

	s := atomic.LoadInt32(&this.status)
	if atomic.CompareAndSwapInt32(&this.status, s, StatusClosed) {

		// 连接管理清除
		if this.Id != 0 {
			this.server.OnWebSocketClose(this)
		}

		close(this.WriteChannel)

		log.Debug("[WebSocketSession|Close]")
	}
}
