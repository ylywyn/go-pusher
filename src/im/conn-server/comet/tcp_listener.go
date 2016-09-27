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
	"net"
)

var CONN_ERROR error = errors.New("STOP LISTENING")

//listener
type TcpListener struct {
	*net.TCPListener
	stop chan bool
}

//accept
func (self *TcpListener) Accept() (*net.TCPConn, error) {
	for {
		conn, err := self.AcceptTCP()
		select {
		case <-self.stop:
			return nil, CONN_ERROR
		default:
			// ok
		}

		if nil != err {
			return nil, err
		}

		return conn, err
	}
}
