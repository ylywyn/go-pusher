/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package id

import (
	"errors"
	"time"
)

var IdPool *IdsPool

type IdsPool struct {
	pool *Pool //私有
}

func (this *IdsPool) Init(addr string) error {
	dialFunc := func() (c Conn, err error) {
		c, err = Dial("tcp", addr)

		return c, nil
	}
	//this.addr = addr_
	// initialize a new pool
	this.pool = &Pool{
		MaxIdle:     64,
		IdleTimeout: 240 * time.Second,
		Dial:        dialFunc,
	}

	//test
	c := this.Get()
	if c != nil {
		return nil
	} else {
		return errors.New("id connect failed")
	}
	c.Close()

	return nil
}

func (this *IdsPool) Get() Conn {
	if this.pool != nil {
		return this.pool.Get()
	} else {
		return nil
	}
}

func (this *IdsPool) Id(key string) (uint64, error) {
	c := this.pool.Get()
	if c.Err() != nil {
		return 0, errors.New("con't conn")
	}
	r, err := c.Do("req" + string(key[0]))
	c.Close()
	if err != nil {
		return 0, err
	} else {
		return r.Uint64(), nil
	}
}
