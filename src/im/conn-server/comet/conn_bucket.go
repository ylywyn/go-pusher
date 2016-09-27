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
	"sync"
)

var ErrNotFound = errors.New("not found")

type ConnBucket struct {
	lock     sync.RWMutex
	ConnsMap map[uint64]Conn // maps of connection
}

func NewBucket() (b *ConnBucket) {
	b = new(ConnBucket)
	b.ConnsMap = make(map[uint64]Conn, 512)
	return b
}

func (b *ConnBucket) Put(id uint64, c Conn) error {
	b.lock.Lock()
	if _, ok := b.ConnsMap[id]; !ok {
		b.ConnsMap[id] = c
		b.lock.Unlock()
		return nil
	}

	b.lock.Unlock()
	return ErrNotFound
}

func (b *ConnBucket) Del(id uint64) error {
	b.lock.Lock()
	if _, ok := b.ConnsMap[id]; ok {
		delete(b.ConnsMap, id)
		b.lock.Unlock()
		return nil
	}

	b.lock.Unlock()
	return ErrNotFound
}

func (b *ConnBucket) Get(id uint64) (Conn, error) {
	b.lock.Lock()
	if c, ok := b.ConnsMap[id]; ok {
		b.lock.Unlock()
		return c, nil
	}

	b.lock.Unlock()
	return nil, ErrNotFound
}

func (b *ConnBucket) Clear() {
	maps := make(map[uint64]Conn, len(b.ConnsMap))
	b.lock.Lock()
	for k, c := range b.ConnsMap {
		maps[k] = c
	}
	//清空
	b.ConnsMap = make(map[uint64]Conn)
	b.lock.Unlock()

	for _, c := range maps {
		c.Close()
	}
}
