/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description: 提供集群唯一ID服务
 ******************************************************************/

package id

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// conn is the low-level implementation of Conn
type conn struct {

	// Shared
	mu   sync.Mutex
	err  error
	conn net.Conn

	// Read
	readTimeout time.Duration
	br          *bufio.Reader

	// Write
	writeTimeout time.Duration
	bw           *bufio.Writer
}

func Dial(network, address string) (Conn, error) {
	dialer := xDialer{}
	return dialer.Dial(network, address)
}

func DialTimeout(network, address string, connectTimeout, readTimeout, writeTimeout time.Duration) (Conn, error) {
	netDialer := net.Dialer{Timeout: connectTimeout}
	dialer := xDialer{
		NetDial:      netDialer.Dial,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	return dialer.Dial(network, address)
}

type xDialer struct {
	NetDial func(network, addr string) (net.Conn, error)

	ReadTimeout time.Duration

	WriteTimeout time.Duration
}

// Dial connects to the Redis server at address on the named network.
func (d *xDialer) Dial(network, address string) (Conn, error) {
	dial := d.NetDial
	if dial == nil {
		dial = net.Dial
	}
	netConn, err := dial(network, address)
	if err != nil {
		return nil, err
	}
	return &conn{
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  d.ReadTimeout,
		writeTimeout: d.WriteTimeout,
	}, nil
}

// NewConn returns a new Redigo connection for the given net connection.
func NewConn(netConn net.Conn, readTimeout, writeTimeout time.Duration) Conn {
	return &conn{
		conn:         netConn,
		bw:           bufio.NewWriter(netConn),
		br:           bufio.NewReader(netConn),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (c *conn) Close() error {
	c.mu.Lock()
	err := c.err
	if c.err == nil {
		c.err = errors.New("client: closed")
		err = c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) fatal(err error) error {
	c.mu.Lock()
	if c.err == nil {
		c.err = err
		c.conn.Close()
	}
	c.mu.Unlock()
	return err
}

func (c *conn) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *conn) writeString(s string) error {
	_, err := c.bw.WriteString(s + "\r\n")
	return err
}

func (c *conn) writeBytes(p []byte) error {
	c.bw.Write(p)
	_, err := c.bw.WriteString("\r\n")
	return err
}

type protocolError string

func (pe protocolError) Error() string {
	return fmt.Sprintf("client: %s (possible server error or unsupported concurrent read by application)", string(pe))
}

func (c *conn) readLine() ([]byte, error) {
	p, err := c.br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	i := len(p) - 2
	if i < 0 || p[i] != '\r' {
		return nil, protocolError("bad response line terminator123")
	}
	return p[:i], nil
}

// parseInt parses an integer reply.
func parseInt(p []byte) (uint64, error) {
	l := len(p)
	if l == 0 || l > 32 {
		return 0, protocolError("malformed integer")
	}

	return strconv.ParseUint(string(p), 10, 0)
}

func parseString(p []byte) (string, error) {
	l := len(p)
	if l < 2 {
		return "", protocolError("malformed integer")
	}

	return string(p), nil
}

func (c *conn) readReply() (*Reply, error) {

	line, err := c.readLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, protocolError("short response line")
	}
	r := &Reply{}
	switch line[0] {
	case ':':
		r.str = string(line[1:])
		return r, nil
	default:
		r.value, err = parseInt(line)
		return r, err
	}

	return nil, protocolError("unexpected response line")
}

func (c *conn) Send(cmd string) error {
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if err := c.writeString(cmd); err != nil {
		return c.fatal(err)
	}
	return nil
}

func (c *conn) Flush() error {
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	if err := c.bw.Flush(); err != nil {
		return c.fatal(err)
	}
	return nil
}

func (c *conn) Receive() (reply *Reply, err error) {
	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	if reply, err = c.readReply(); err != nil {
		return nil, c.fatal(err)
	}

	return reply, nil
}

func (c *conn) Do(cmd string) (*Reply, error) {

	if cmd == "" {
		return nil, nil
	}

	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	if err := c.writeString(cmd); err != nil {
		return nil, c.fatal(err)
	}

	if err := c.bw.Flush(); err != nil {
		return nil, c.fatal(err)
	}

	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}

	reply, e := c.readReply()
	if e != nil {
		return nil, c.fatal(e)
	}

	return reply, nil
}
