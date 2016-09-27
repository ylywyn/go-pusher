/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : mem.go
 *  Date   :
 *  Author : yangl
 *  Description: 内存缓存session
 ******************************************************************/

package session

import (
	"sync"
)

var clusterSession bool
var memSessions *MemSessionManager

type MemSessionManager struct {
	bucketCount int
	buckets     []*MemSessionBucket
}

func NewMemSessions(cluster bool, count int) *MemSessionManager {
	clusterSession = cluster
	memSessions = new(MemSessionManager)
	memSessions.Init(count)
	return memSessions
}

func (cm *MemSessionManager) Init(count int) {
	cm.bucketCount = count
	cm.buckets = make([]*MemSessionBucket, count)
	for i := 0; i < count; i++ {
		cm.buckets[i] = NewMemSessionBucket()
	}
}

func (cm *MemSessionManager) GetBucketCount() int {
	return cm.bucketCount
}

func (cm *MemSessionManager) PutSession(id uint64, c *Session) error {
	idx := id % uint64(cm.bucketCount)
	b := cm.buckets[idx]
	err := b.Put(id, c)
	return err
}

func (cm *MemSessionManager) DelSession(id uint64) error {
	idx := id % uint64(cm.bucketCount)
	b := cm.buckets[idx]
	err := b.Del(id)
	return err
}

func (cm *MemSessionManager) GetSession(id uint64) *Session {
	idx := id % uint64(cm.bucketCount)
	b := cm.buckets[idx]
	return b.Get(id)
}

func (cm *MemSessionManager) Clear() {
	for i := 0; i < cm.bucketCount; i++ {
		b := cm.buckets[i]
		b.Clear()
	}
}

//////////////////////
type MemSessionBucket struct {
	lock     sync.RWMutex
	sessions map[uint64]*Session
}

func NewMemSessionBucket() *MemSessionBucket {
	sb := new(MemSessionBucket)
	sb.sessions = make(map[uint64]*Session, 512)
	return sb
}

func (b *MemSessionBucket) Put(id uint64, s *Session) error {
	b.lock.Lock()
	b.sessions[id] = s
	b.lock.Unlock()
	return nil
}

func (b *MemSessionBucket) Del(id uint64) error {
	b.lock.Lock()
	delete(b.sessions, id)
	b.lock.Unlock()
	return nil
}

func (b *MemSessionBucket) Get(id uint64) *Session {
	b.lock.Lock()
	if c, ok := b.sessions[id]; ok {
		b.lock.Unlock()
		return c
	}

	b.lock.Unlock()
	return nil
}

func (b *MemSessionBucket) Clear() {
	b.lock.Lock()
	b.sessions = make(map[uint64]*Session)
	b.lock.Unlock()
}
