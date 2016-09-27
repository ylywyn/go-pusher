/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : redis.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package utils

import (
	"errors"
	"time"

	"github.com/PuerkitoBio/redisc"
	"github.com/garyburd/redigo/redis"
)

type RedisPool struct {
	cluster *redisc.Cluster
	pool    *redis.Pool
	addr    []string
	Name    string
}

func NewRedisPool(addr []string) (*RedisPool, error) {
	pool := &RedisPool{addr: addr}
	err := pool.InitRedis(addr)
	return pool, err
}

func (this *RedisPool) InitRedis(addr []string) error {

	if len(addr) > 1 {
		this.cluster = &redisc.Cluster{
			StartupNodes: addr,
			DialOptions:  []redis.DialOption{redis.DialConnectTimeout(6 * time.Second)},
			CreatePool:   this.createPool,
		}

		// initialize its mapping
		if err := this.cluster.Refresh(); err != nil {
			return err
		}
	} else if len(addr) == 1 {
		var err error
		var op redis.DialOption
		this.pool, err = this.createPool(addr[0], op)
		if err != nil {
			return err
		}
	}

	return nil
}

//func (this *)
func (this *RedisPool) createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {

	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	//this.addr = addr

	// initialize a new pool
	this.pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 300 * time.Second,
		Dial:        dialFunc,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	//test
	c := this.pool.Get()
	if c.Err() == nil {
		c.Close()
		return this.pool, nil
	} else {
		return nil, errors.New("Redis connect failed")
	}
}

func (this *RedisPool) Close() {
	if this.cluster != nil {
		this.cluster.Close()
	}
	if this.pool != nil {
		this.pool.Close()
	}
}

func (this *RedisPool) Get() redis.Conn {
	if this.cluster != nil {
		return this.cluster.Get()
	}
	if this.pool != nil {
		return this.pool.Get()
	}

	return nil
}
