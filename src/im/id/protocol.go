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
	"bytes"
	"net"
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
}

var (
	endTag = []byte("\r\n")
)

// Packet
type IdPacket struct {
	pType string
	pData []byte
}

func (p *IdPacket) Serialize() []byte {
	p.pData = append(p.pData, endTag...)
	return p.pData
}

func (p *IdPacket) GetType() string {
	return p.pType
}

func (p *IdPacket) GetData() []byte {
	return p.pData
}

func NewIdPacket(pType string, pData []byte) *IdPacket {
	return &IdPacket{
		pType: pType,
		pData: pData,
	}
}

type IdProtocol struct {
}

func (this *IdProtocol) ReadPacket(conn *net.TCPConn) (Packet, error) {
	fullBuf := bytes.NewBuffer([]byte{})
	for {
		data := make([]byte, 16)
		readLengh, err := conn.Read(data)

		if err != nil { //EOF, or worse
			return nil, err
		}

		if readLengh == 0 {
			return nil, ErrConnClosing
		} else {
			fullBuf.Write(data[:readLengh])

			index := bytes.Index(fullBuf.Bytes(), endTag)
			if index > -1 {
				command := fullBuf.Next(index)
				return NewIdPacket("", command), nil
			}
		}
	}
}
