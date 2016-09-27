/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : comet.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package comet

import (
	"fmt"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	"im/common/pump"
	"im/conn-server/conf"
	"runtime"
	"time"
)

var (
	Server *CometServer
)

type CometServer struct {
	name        string
	closed      bool
	pump        pump.IBasePump
	connManager *ConnManager
	tcpServer   *TcpServer
	wsServer    *WebSocketServer
}

func CreateTcpServer(cm *ConnManager) *TcpServer {
	op := &ServerOptions{}
	op.Addr = conf.Conf.TCPBind
	op.HeartBeat = conf.Conf.HeartBeat
	op.IdleTime = time.Duration(conf.Conf.TimeOut) * time.Second
	op.ReadBufferSize = conf.Conf.TCPReadbuf
	op.WriteBufferSize = conf.Conf.TCPWritebuf
	op.ReadChannelSize = conf.Conf.TCPRecvChannel
	op.WriteChannelSize = conf.Conf.TCPSendChannel
	op.Cm = cm
	tcp := NewTcpServer(op)

	return tcp
}

func CreateWebSocketServer(cm *ConnManager) *WebSocketServer {
	op := &ServerOptions{}
	op.Addr = conf.Conf.WebsocketBind
	op.HeartBeat = conf.Conf.HeartBeat
	op.IdleTime = time.Duration(conf.Conf.TimeOut) * time.Second
	op.ReadBufferSize = conf.Conf.WebsocketReadbuf
	op.WriteBufferSize = conf.Conf.WebsocketWritebuf
	op.ReadChannelSize = conf.Conf.WebsocketRecvChannel
	op.WriteChannelSize = conf.Conf.WebsocketSendChannel
	op.Cm = cm
	ws := NewWebSocketServer(op)

	return ws
}

func NewCometServer() *CometServer {

	//消息泵
	msgpump := pump.NewPump(conf.Conf.PumpType)

	//已验证连接管理
	connmanager := NewConnManager(runtime.NumCPU() * 64)

	//Tcp 服务器
	tcpserver := CreateTcpServer(connmanager)

	//WebSocketServer
	wsserver := CreateWebSocketServer(connmanager)

	ims := &CometServer{
		closed:      true,
		connManager: connmanager,
		tcpServer:   tcpserver,
		wsServer:    wsserver,
		pump:        msgpump,
	}

	//消息处理器
	ProcessInstance = NewProcessor(msgpump, ims.connManager)
	return ims
}

func (im *CometServer) Start() {
	if im.closed {
		// msg pump
		im.pump.EnablePub()
		im.pump.EnableSub()
		im.pump.BindSubProcess(ProcessInstance.ProcessSend)
		im.pump.Connect(conf.Conf.PumpAddr)
		c := fmt.Sprintf(pump.MSG_CHANNEL_FRONT, conf.Conf.ServerId)
		im.pump.Sub(c, c)

		// tcp server
		im.tcpServer.SetReadCallBack(im.tcpProcessRecvMsg)
		err := im.tcpServer.ListenAndServer()
		if err != nil {
			log.Error("[ImServer|Start|tcpListenAndServer] %s", err.Error())
			panic(err.Error())
		}

		// web socket
		im.wsServer.SetReadCallBack(im.webSocktProcessRecv)
		im.wsServer.Start()

		im.closed = false
	}
}

func (im *CometServer) Stop() {
	if !im.closed {
		im.closed = true
		im.pump.Stop()
		im.tcpServer.StopListen()
		im.wsServer.StopListen()
		im.connManager.Clear()
	}
}

func (im *CometServer) tcpProcessRecvMsg(s *TcpSession, m *msg.MsgRaw) {
	ProcessInstance.ProcessRecv(m)
}

func (im *CometServer) webSocktProcessRecv(c *WebSocketSession, m *msg.MsgRaw) {
	ProcessInstance.ProcessRecv(m)
}

func (im *CometServer) Stat() *CometStats {
	st := &CometStats{}
	st.TcpStats = *(im.tcpServer.Stat)
	return st
}
