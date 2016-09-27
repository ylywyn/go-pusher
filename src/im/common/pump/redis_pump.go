/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package pump

import (
	"container/list"
	log "im/common/log4go"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisPump struct {
	closed     bool
	connstr    string
	writeChan  chan pubItem
	pubEnable  bool
	subEnable  bool
	pool       *redis.Pool
	subProcess func(data []byte)
}

func (this *RedisPump) EnableSub() {
	this.subEnable = true
}

func (this *RedisPump) EnablePub() {
	this.pubEnable = true
}

//同步， 连接不上会一直等待
func (this *RedisPump) Connect(s string) {
	if !this.pubEnable && !this.subEnable {
		return
	}

	if len(s) > 0 {
		this.connstr = s
		this.writeChan = make(chan pubItem, MAX_PUBCHANNEL_SIZE)
		this.conn()
		this.closed = false
		go this.pub()
	}
}

func (this *RedisPump) conn() {

	// initialize a new pool
	this.pool = &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", this.connstr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

func (this *RedisPump) Stop() {

	if !this.closed {
		this.closed = true
		this.pool.Close()

		close(this.writeChan)
	}
}

func (this *RedisPump) Sub(c, group string) error {

	var err error
	go func() {
	begin:
		conn := this.pool.Get()
		defer conn.Close()

		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe(c)

		for {
			switch n := psc.Receive().(type) {
			case redis.Message:
				if this.subProcess != nil {
					this.subProcess(n.Data)
				}

			case redis.PMessage:

			case redis.Subscription:
				if n.Count == 0 {
					return
				}

			case error:
				conn.Close()
				time.Sleep(10 * time.Second)
				goto begin
			}
		}
	}()

	return err
}

func (this *RedisPump) Pub(channel string, m []byte) bool {

	//log.Debug("RedisPump Pub len(%d)", len(m))
	item := pubItem{c: channel, data: m}

	select {
	case this.writeChan <- item:
	default:
		log.Warn("RedisPump Pub Cache Full")
		return false
	}

	return true
}

func (this *RedisPump) pub() {
	for {
		item, ok := <-this.writeChan
		if ok {

			//err = this.connPub.Publish(item.c, item.data)
			conn := this.pool.Get()
			defer conn.Close()

			_, err := conn.Do("PUBLISH", item.c, item.data)

			if err != nil && (!this.closed) {
				log.Error("RedisPump pub routine error %s", err.Error())
			}
		} else {
			log.Debug("RedisPump pub routine quit")
			break
		}
	}
}

func (this *RedisPump) GetAllPubs() *list.List {
	l := list.New()

	for {
		select {
		case item, ok := <-this.writeChan:
			{
				if ok && l.Len() <= MAX_PUBCHANNEL_SIZE {
					l.PushBack(item)
				} else {
					return l
				}
			}
		default:
			return l
		}
	}
}

func (this *RedisPump) BindSubProcess(f func(data []byte)) {
	if f != nil {
		this.subProcess = f
	}
}
