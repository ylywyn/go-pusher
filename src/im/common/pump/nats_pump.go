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
	"math/rand"
	"strings"
	"time"

	log "im/common/log4go"

	"github.com/nats-io/nats"
)

type NatsPump struct {
	closed     bool
	connstr    string
	writeChan  chan pubItem
	pubEnable  bool
	subEnable  bool
	connPub    *nats.Conn
	connSub    *nats.Conn
	subProcess func(data []byte)
}

func (this *NatsPump) EnableSub() {
	this.subEnable = true
}

func (this *NatsPump) EnablePub() {
	this.pubEnable = true
}

//同步， 连接不上会一直等待
func (this *NatsPump) Connect(s string) {
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

func (this *NatsPump) conn() {

	urls := strings.Split(this.connstr, ",")
	l := len(urls)

	if l > 1 {
		//随机连接
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if this.subEnable {
			this.subConn(urls[r.Intn(l)])
		}
		if this.pubEnable {
			this.pubConn(urls[r.Intn(l)])
		}

	} else {
		if this.subEnable {
			this.subConn(this.connstr)
		}

		if this.pubEnable {
			this.pubConn(this.connstr)
		}
	}
}

func (this *NatsPump) subConn(url string) {

	var err error
	for {
		this.subClose()

		this.connSub, err = nats.Connect(url, nats.ReconnectWait(15*time.Second),
			nats.DisconnectHandler(func(nc *nats.Conn) {
				log.Error("[NatsPump|Sub] Disconnected!")
			}),
			nats.ReconnectHandler(func(_ *nats.Conn) {
				log.Info("[NatsPump|Sub] Reconnected %v ", this.connSub.ConnectedUrl())
			}))

		//错误回调
		if err == nil {
			this.connSub.SetErrorHandler(func(*nats.Conn, *nats.Subscription, error) {
				log.Error("[NatsPump|Sub] Conn Error %v", this.connSub.ConnectedUrl())
				this.subConn(url)
			})
			break
		} else {
			log.Error(err.Error())
			time.Sleep(15 * time.Second)
		}
	}

	log.Debug("[NatsPump|Sub] Connected %s", url)
}

func (this *NatsPump) pubConn(url string) {
	var err error
	for {
		this.pubClose()

		this.connPub, err = nats.Connect(url, nats.ReconnectWait(15*time.Second),
			nats.DisconnectHandler(func(nc *nats.Conn) {
				log.Error("[NatsPump|Pub] Disconnected!")
			}),
			nats.ReconnectHandler(func(_ *nats.Conn) {
				log.Info("[NatsPump|Pub]  Reconnected %v ", this.connPub.ConnectedUrl())
			}))

		//错误回调
		if err == nil {
			this.connPub.SetErrorHandler(func(*nats.Conn, *nats.Subscription, error) {
				log.Error("[NatsPump|Pub]  Conn Error %v", this.connSub.ConnectedUrl())
				this.pubConn(url)
			})
			break
		} else {
			log.Error(err.Error())
			time.Sleep(15 * time.Second)
		}
	}

	log.Debug("[NatsPump|Pub] Connected %s", url)
}

func (this *NatsPump) subClose() {
	if this.connSub != nil && !this.connSub.IsClosed() {
		this.connSub.Close()
	}
}

func (this *NatsPump) pubClose() {
	if this.connPub != nil && !this.connPub.IsClosed() {
		this.connPub.Close()
	}
}

func (this *NatsPump) Stop() {

	if !this.closed {
		this.closed = true
		this.subClose()
		this.pubClose()

		close(this.writeChan)
	}
}

func (this *NatsPump) Sub(c, group string) error {

	var err error
	if len(group) == 0 {
		_, err = this.connSub.Subscribe(c, this.subCallBack)
	} else {
		this.connSub.QueueSubscribe(c, group, this.subCallBack)
	}

	return err
}

func (this *NatsPump) Pub(channel string, m []byte) bool {

	//log.Debug("NatsPump Pub len(%d)", len(m))
	item := pubItem{c: channel, data: m}

	select {
	case this.writeChan <- item:
	default:
		log.Warn("NatsPump Pub Cache Full")
		return false
	}

	return true
}

func (this *NatsPump) pub() {

	var err error
	for {
		item, ok := <-this.writeChan
		if ok {

			err = this.connPub.Publish(item.c, item.data)

			if err != nil && (!this.closed) {
				this.pubConn(this.connstr)
				log.Error("NatsPump pub routine error %s", err.Error())
			}
		} else {
			log.Debug("NatsPump pub routine quit")
			break
		}
	}
}

func (this *NatsPump) subCallBack(m *nats.Msg) {
	if m != nil {
		if this.subProcess != nil {
			this.subProcess(m.Data)
		}
	}
}

func (this *NatsPump) GetAllPubs() *list.List {
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

func (this *NatsPump) BindSubProcess(f func(data []byte)) {
	if f != nil {
		this.subProcess = f
	}
}
