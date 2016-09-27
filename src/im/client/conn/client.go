/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description: golang im 简易客户端
 ******************************************************************/

package main

import (
	"bufio"
	"fmt"
	"im/common/proto/entity/msg"
	"im/common/proto/entity/msg/login"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	"im/conn-server/comet"
	"net"
	"strings"
	"time"
)

const (

	// PONG
	pongWait = 10 * time.Second

	// PING
	pingPeriod = (pongWait * 9) / 10
)

type TcpClient struct {
	Uid          uint64
	Pwd          string
	Conn         *net.TCPConn
	writeChannel chan string
	heartBeatMsg *msg.MsgRaw
	codec        comet.TcpLenCodec
	dataTest     []byte
	closed       bool
	count        int
}

func NewClient(uid uint64, pwd string, server *net.TCPAddr) (*TcpClient, error) {
	c := &TcpClient{
		Uid:          uid,
		Pwd:          pwd,
		writeChannel: make(chan string, 1024),
		closed:       false,
		count:        0,
	}
	var err error
	c.Conn, err = net.DialTCP("tcp", nil, server)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (this *TcpClient) Start() {
	go this.writeLoop()
	this.readLoop()
}

func (this *TcpClient) readLoop() {
	reader := bufio.NewReaderSize(this.Conn, 2048)
	for {
		m, err := this.codec.Read(reader)
		if nil != err {
			if !this.closed {
				fmt.Printf("[TcpClient|readLoop|codec.Read] %s\r\n", err.Error())
			}
			break
		}

		mb, err := this.codec.UnmarshalPacket(m)
		if nil != err {
			fmt.Printf("[TcpClient|readLoop|codec.Read] %s\r\n", err.Error())
			break
		}

		this.onMsgReveice(mb)

	}
}

func (this *TcpClient) writeLoop() {
	//发送登录消息
	this.writeBytes(this.makeLoginMsg())

	//心跳，定时写计时器
	var t time.Duration = time.Duration(Conf.HeartBeat)
	if Conf.WriteTestSpan > 0 {
		t = time.Duration(Conf.WriteTestSpan)
	}
	tickerHb := time.NewTicker(t * time.Second)

	defer func() {
		tickerHb.Stop()
		this.Conn.Close()
		if err := recover(); err != nil {
			fmt.Printf("[TcpClient|writeLoop|recover] %s\r\n", err)
		}
	}()

	var text string
	for {
		select {
		case text = <-this.writeChannel:
			{
				if text != "" {
					data := this.strToMsgBytes(text)
					if data != nil && len(data) > 10 {
						if err := this.writeBytes(data); err != nil {
							fmt.Printf("[TcpClient|writeLoop|writeBytes] %s\r\n", err.Error())
							goto end
						}
					}

				} else {
					goto end
				}
			}
		case <-tickerHb.C:
			if Conf.WriteTestSpan > 0 {
				if err := this.writeBytes(this.makeTestMsg()); err != nil {
					fmt.Printf("[TcpClient|writeLoop|ticker|MakeTestMsg] %s \r\n", err.Error())
					goto end
				}
			} else {
				if err := this.writeBytes(this.makeHeartBeat()); err != nil {
					fmt.Printf("[TcpClient|writeLoop|ticker|MakeHeartBeat] %s \r\n", err.Error())
					goto end
				}
			}

		}

	}
end:
}

func (this *TcpClient) Write(data string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[TcpClient|data|recover] %s\r\n", err)
		}
	}()

	select {
	case this.writeChannel <- data:
	default:
	}

	return nil
}

func (this *TcpClient) writeBytes(data []byte) error {
	l := 0
	ld := len(data)
	for {
		lw, err := this.Conn.Write(data[l:])
		if nil != err {
			fmt.Printf("[TcpClient|WriteN|] %s \r\n", err.Error())
			return err
		}

		l += lw
		if l >= ld {
			break
		}
	}

	return nil
}

func (this *TcpClient) makeHeartBeat() []byte {

	if this.heartBeatMsg == nil {
		this.heartBeatMsg = &msg.MsgRaw{Body: make([]byte, 4)}
		this.heartBeatMsg.Header.Ack = 0
		this.heartBeatMsg.Header.Compress = 0
		this.heartBeatMsg.Header.Length = 0
		this.heartBeatMsg.Header.Type = types.StatusMsgTypeMT_STATUS_HEARTBEAT
		this.heartBeatMsg.Serialize()
	}

	return this.heartBeatMsg.Body
}

func (this *TcpClient) makeLoginMsg() []byte {
	lgin := login.MsgLogin{}
	lgin.Uid = this.Uid
	lgin.PassWd = this.Pwd
	lgin.Platform = types.PlatformPF_DESKTOP

	d, _ := lgin.Serialize()

	m := &mb.MsgBase{
		Type:  uint16(types.UserMsgTypeMT_USER_LOGIN_REQ),
		From:  this.Uid,
		Time:  uint64(time.Now().UnixNano()),
		Appid: uint16(Conf.Appid),
		Text:  string(d),
	}

	data := m.Serialize()

	msgRaw := &msg.MsgRaw{Body: data}
	msgRaw.Header.Length = uint16(len(data) - msg.HEADER_LEN)
	msgRaw.Header.Type = m.Type
	msgRaw.Serialize()

	return msgRaw.Body
}

func (this *TcpClient) makeTestMsg() []byte {

	if this.dataTest == nil {
		var t string
		if len(Conf.WriteTestMsg) > 0 {
			t = Conf.WriteTestMsg
		} else {
			t = "hello, yangl"
		}

		this.dataTest = this.makeMsg(this.Uid, types.GeneralMsgTypeMT_GENERAL_MSG, t)
	}

	return this.dataTest
}

func (this *TcpClient) makeMsg(to uint64, types uint16, t string) []byte {
	m := &mb.MsgBase{
		Type:     types,
		From:     this.Uid,
		To:       to,
		Time:     uint64(time.Now().UnixNano()),
		Appid:    uint16(Conf.Appid),
		Text:     t,
		Platform: 1,
		Gid:      to,
	}

	data := m.Serialize()

	msgRaw := &msg.MsgRaw{Body: data}
	msgRaw.Header.Length = uint16(len(data) - msg.HEADER_LEN)
	msgRaw.Header.Type = m.Type
	msgRaw.Serialize()

	return msgRaw.Body
}

func (this *TcpClient) Close() {
	this.closed = true
	this.Conn.Close()
	close(this.writeChannel)
}

func (this *TcpClient) Closed() bool {
	return this.closed
}

//
func (this *TcpClient) strToMsgBytes(text string) []byte {

	t := types.GeneralMsgTypeMT_GENERAL_MSG
	if Conf.MsgTypes == "g" {
		t = types.GeneralMsgTypeMT_GENERAL_GROUP_MSG
	}

	return this.makeMsg(Conf.To, uint16(t), text)
}

func CheckInput(text string) bool {

	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return false
	}
	return true
}

func (this *TcpClient) onMsgReveice(m *mb.MsgBase) {
	switch m.Type {
	case types.GeneralMsgTypeMT_GENERAL_MSG,
		types.GeneralMsgTypeMT_GENERAL_GROUP_MSG:
		{
			this.count++
			if Conf.Num < 1001 {
				fmt.Printf("receive %d: from:%d, msg:%s \r\n", this.count, m.From, m.Text)
			}
		}
	case types.UserMsgTypeMT_USER_LOGIN_REP:
		fmt.Printf("用户 %d 登录成功！\r\n", Conf.Uid)
	case types.UserMsgTypeMT_USER_LOGIN_FAILED_REP:
		{
			fmt.Printf("用户 %d 登录失败！\r\n", Conf.Uid)
			this.Close()
		}
	case types.UserMsgTypeMT_USER_KICKUSER:
		{
			fmt.Printf("用户 %d 在其他平台登录！\r\n", Conf.Uid)
			this.Close()
		}
	case types.StatusMsgTypeMT_STATUS_CRITICAL_ERROR:
		{
			fmt.Printf("发生错误： %s！\r\n", m.Text)
			this.Close()
		}
	default:
	}
}
