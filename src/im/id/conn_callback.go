/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : config.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	"fmt"
	"strconv"
	"time"
)

var (
	reqCmd    = []byte("req")
	pingCmd   = []byte("ping")
	quitCmd   = []byte("quit")
	pongCmd   = []byte(":pong")
	unKnowCmd = []byte(":unknow cmd")
)

type IdCallback struct {
}

func (this *IdCallback) OnConnect(c *Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println(time.Now(), " ,OnConnect:", addr)
	return true
}

func (this *IdCallback) OnMessage(c *Conn, p Packet) bool {
	packet := p.(*IdPacket)
	command := packet.GetData()

	if command == nil || len(command) < 3 {
		return false
	}
	//fmt.Println("OnMessage:", string(command))
	switch command[0] {
	case reqCmd[0]:
		if len(command) == 4 {
			return this.processReq(c, int(command[3])%10)
		} else {
			return false
		}
	case pingCmd[0]:
		c.AsyncWritePacket(NewIdPacket("", pongCmd), 0)
	case quitCmd[0]:
		return false
	default:
		c.AsyncWritePacket(NewIdPacket("", unKnowCmd), 0)
	}

	return true
}

func (this *IdCallback) OnClose(c *Conn) {
	fmt.Println(time.Now(), " ,OnClose:", c.GetExtraData())
}

func (this *IdCallback) processReq(c *Conn, n int) bool {
	id, err := SfsServs[n].NextID()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	//id = 1907281294393530
	//id = 1907281294393531
	//fmt.Println(id)
	//
	//	m := Decompose(id)
	//	fmt.Println(m["id"])
	//	fmt.Println(m["time"])
	//	fmt.Println(m["sequence"])
	//	fmt.Println(m["machine-id"])
	//
	data := strconv.FormatUint(id, 10)
	err = c.AsyncWritePacket(NewIdPacket("", []byte(data)), 0)
	if err != nil {
		return false
	}
	return true
}
