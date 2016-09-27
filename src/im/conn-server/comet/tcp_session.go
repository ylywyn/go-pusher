/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : tcp_session.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package comet

import (
	"bufio"
	"errors"
	"fmt"
	"im/common/proto/entity/msg"
	"im/common/proto/entity/msg/status"
	"im/common/proto/fbsgen/msg/types"
	"io"
	"net"
	"sync/atomic"
	"time"

	log "im/common/log4go"
)

var CONN_CLOSED_ERROR error = errors.New("CONN CLOSED ERROR")

//ReadChannel可去掉， wait fix
type TcpSession struct {
	Id         uint64
	Authed     bool
	status     int32 //此状态并没有加锁保护
	remoteAddr string
	errMsg     string
	Conn       *net.TCPConn
	br         *bufio.Reader
	bw         *bufio.Writer
	//ReadChannel  chan *msg.MsgRaw
	WriteChannel chan *msg.MsgRaw

	server *TcpServer
}

func NewTcpSession(i uint32, c *net.TCPConn, s *TcpServer) *TcpSession {
	op := s.options

	//c.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
	if op.HeartBeat {
		//客户端主动发送心跳包
		c.SetReadDeadline(time.Now().Add(op.IdleTime))
	} else {
		//tcp的KeepAlive， 并缩短KeepAlive时长
		c.SetKeepAlive(true)
		c.SetKeepAlivePeriod(op.IdleTime)
	}

	// 如果回复心跳这样的小包可能有用
	// 这里可能用处不大，因为一个包大于100字节了
	c.SetNoDelay(true)

	// Tcp/Ip协议栈的套接字缓冲，和具体的实现相关
	if op.ReadBufferSize > 512 {
		c.SetReadBuffer(op.ReadBufferSize)
	}

	if op.WriteBufferSize > 512 {
		c.SetWriteBuffer(op.WriteBufferSize)
	}

	session := &TcpSession{
		Id:         0,
		Authed:     false,
		Conn:       c,
		status:     StatusRunning,
		remoteAddr: c.RemoteAddr().String(),
		br:         bufio.NewReaderSize(c, op.ReadBufferSize),
		bw:         bufio.NewWriterSize(c, op.WriteBufferSize),
		//ReadChannel:  make(chan *msg.MsgRaw, op.ReadChannelSize),
		WriteChannel: make(chan *msg.MsgRaw, op.WriteChannelSize),
		server:       s,
	}

	return session
}

func (this *TcpSession) LocalAddr() string {
	return this.Conn.LocalAddr().String()
}

func (this *TcpSession) RemoteAddr() string {
	return this.Conn.RemoteAddr().String()
}

//启动，开始读写操作
func (this *TcpSession) ServeTCP() {
	this.server.waitGroup.Add(1)

	//写读数据
	go this.writeLoop()
	this.readLoop()

	s := atomic.LoadInt32(&this.status)
	if s > 0 {
		this.Close()
	}

	this.server.waitGroup.Done()
}

//read
func (this *TcpSession) readLoop() {

	defer func() {
		if err := recover(); err != nil {
			log.Error("[TcpSession|readLoop|recover|%s] %s", this.remoteAddr, err)
		}
	}()

	codec := this.server.codec
	op := this.server.options
	stat := this.server.Stat

	for this.status == StatusRunning {

		// 读一个完整的包
		m, err := codec.Read(this.br)
		if nil != err {
			//this.errMsg = err.Error()
			log.Error("[TcpSession|readLoop|codec.Read|%s] %s", this.remoteAddr, err.Error())
			break
		}

		//重置允许的IdleTime
		if op.HeartBeat {
			if this.Authed {
				//没有登录成功之前，不允许心跳
				this.Conn.SetReadDeadline(time.Now().Add(op.IdleTime))
			}
			if m == nil {
				continue
			}
		}

		if this.Id == 0 {
			// 第一个消息不是登陆消息，则关闭连接
			if m.Header.Type == types.UserMsgTypeMT_USER_LOGIN_REQ {
				this.server.OnTcpConn(this, m)
			} else {
				this.errMsg = "must first auth conn, you type is "
				log.Error("[TcpSession|readLoop|m.Header.Type|%s] %d", this.remoteAddr, m.Header.Type)
				break
			}
		} else if m.Header.Type == types.UserMsgTypeMT_USER_LOGIN_REQ {
			this.server.OnTcpConn(this, m)
		}

		//调试
		//log.Debug("[TcpSession|Read]ConnId:%d, ConnStatus:%d", this.Id, this.status)

		//dispatchmsg
		//this.ReadChannel <- p
		this.server.readCallBack(this, m)

		//读流量统计

		stat.FlowStat.IncrReadCounts()
		stat.FlowStat.IncrReadBytes(int32(m.Len()))
	}
	//log.Debug("[TcpSession|ServeTCP] readLoop over~")
	this.Close()
}

//write
func (this *TcpSession) writeLoop() {

	this.server.waitGroup.Add(1)
	defer func() {
		//应该让客户端去Close(), 单客户端容易断网，容易导致文件描述符泄露
		this.Conn.Close()
		this.server.waitGroup.Done()
		if err := recover(); err != nil {
			log.Error("[TcpSession|writeLoop|recover|%s] %s", this.remoteAddr, err)
		}
	}()

	var p *msg.MsgRaw
	var err error
	for this.status == StatusRunning {
		p = <-this.WriteChannel
		if p != nil {
			err = this.write(p)
			if err != nil {
				this.errMsg = err.Error()
				break
			}
		} else {
			//channel已经关闭，说明已经关闭
			break
		}
	}

	// 在这里处理 发送最后的错误消息，以及关闭连接
	// 线程安全，在这里发送错误信息
	if this.errMsg != "" {
		log.Debug("[TcpSession|writeLoop|WriteErrorMsg] %s", this.errMsg)
		m := status.NewStateMsg(types.StatusMsgTypeMT_STATUS_CRITICAL_ERROR, this.errMsg)
		err := this.WriteN(m)
		if err != nil {
			log.Error("[TcpSession|writeLoop] %s", err.Error())
		}
	}

	//this.Conn.
	//log.Debug("[TcpSession|ServeTCP] writeLoop over~")
	this.Conn.CloseWrite()

	//不延迟的话， 最后一个消息发不出去就close了
	t := time.NewTimer(time.Second * 3)
	<-t.C
}

func (this *TcpSession) Write(p *msg.MsgRaw) error {
	defer func() {
		if err := recover(); nil != err {
			log.Error("[TcpSession|Write|recover|%s]%s", this.remoteAddr, err)
		}
	}()

	if this.Authed == false {
		if p.GetMsgType() == types.UserMsgTypeMT_USER_LOGIN_REP {
			this.Authed = true
		} else {
			//登录失败，强制关闭连接？
		}
	}

	//log.Debug("[TcpSession|Write]ConnId:%d, ConnStatus:%d", this.Id, this.status)
	if this.status == StatusRunning {
		select {
		case this.WriteChannel <- p:
			return nil
		default:
			return errors.New(fmt.Sprintf("[TcpSession|Write|%s] channel full", this.remoteAddr))
		}
	}
	return CONN_CLOSED_ERROR
}

//这需要带缓冲吗？ 反正最后也要 flush wait fix
//带缓冲的写
func (this *TcpSession) write(p *msg.MsgRaw) error {

	data := p.Body
	if nil == data || len(data) == 0 {
		log.Error("[TcpSession|write|%s] EMPTY PACKET", this.remoteAddr)
		return nil
	}

	this.Conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))

	l := 0
	ld := len(data)
	for {
		lw, err := this.bw.Write(data[l:])
		if nil != err {
			log.Error("[TcpSession|write|%s] %s", this.remoteAddr, err.Error())
			//链接是关闭的
			if err != io.ErrShortWrite {
				return err
			}

			//如果没有写够则再写一次
			if err == io.ErrShortWrite {
				this.bw.Reset(this.Conn)
			}
		}

		l += lw
		//write finish
		if l == ld {
			break
		}
	}

	this.bw.Flush()

	//读流量统计
	this.server.Stat.FlowStat.IncrWriteCounts()
	this.server.Stat.FlowStat.IncrWriteBytes(int32(ld))

	return nil
}

//is close
func (this *TcpSession) Closed() bool {
	return atomic.LoadInt32(&this.status) == StatusClosed
}

//
func (this *TcpSession) Close() {

	s := atomic.LoadInt32(&this.status)
	if atomic.CompareAndSwapInt32(&this.status, s, StatusClosed) {

		// 连接管理清除

		if this.Id != 0 {
			this.server.OnTcpConnClose(this)
		}

		//close(this.ReadChannel)

		//处理WriteChannel里剩余的消息
		//		var p *msg.MsgRaw
		//		dispmsg := this.server.packetDispatcher
		//		for {
		//			p = nil
		//			select {
		//			case p = <-this.WriteChannel:
		//			default:
		//				break
		//			}

		//			if p == nil {
		//				break
		//			} else {
		//				dispmsg(this, p)
		//			}
		//		}

		close(this.WriteChannel)

		log.Debug("[TcpSession|Close]")
	}
}

func (this *TcpSession) WriteN(p *msg.MsgRaw) error {

	data := p.Body
	if nil == data || len(data) == 0 {
		log.Error("[TcpSession|WriteN|%s] EMPTY PACKET", this.remoteAddr)
		return nil
	}

	l := 0
	ld := len(data)
	//log.Debug("[TcpSession|write] read write: %d", ld)
	for {
		lw, err := this.Conn.Write(data[l:])
		if nil != err {
			log.Error("[TcpSession|WriteN|%s] %s", this.remoteAddr, err.Error())
			return err
		}

		l += lw
		//write finish
		if l >= ld {
			break
		}
	}

	//读流量统计
	this.server.Stat.FlowStat.IncrWriteCounts()
	this.server.Stat.FlowStat.IncrWriteBytes(int32(ld))

	return nil
}
