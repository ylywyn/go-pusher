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

//为了减小锁的范围， 将连接分散在多个bucket中
type ConnManager struct {
	bucketCount int
	buckets     []*ConnBucket
}

func NewConnManager(count int) *ConnManager {
	cm := new(ConnManager)
	cm.Init(count)
	return cm
}

func (cm *ConnManager) Init(count int) {
	cm.bucketCount = count
	cm.buckets = make([]*ConnBucket, count)
	for i := 0; i < count; i++ {
		cm.buckets[i] = NewBucket()
	}
}

func (cm *ConnManager) GetBucketCount() int {
	return cm.bucketCount
}

func (cm *ConnManager) PutConn(id uint64, c Conn) error {
	idx := id % uint64(cm.bucketCount)
	b := cm.buckets[idx]
	err := b.Put(id, c)
	return err
}

func (cm *ConnManager) DelConn(id uint64) error {
	idx := id % uint64(cm.bucketCount)
	b := cm.buckets[idx]
	err := b.Del(id)
	return err
}

func (cm *ConnManager) GetConn(id uint64) (Conn, error) {
	idx := id % uint64(cm.bucketCount)
	b := cm.buckets[idx]
	c, err := b.Get(id)
	return c, err
}

func (cm *ConnManager) Clear() {
	for i := 0; i < cm.bucketCount; i++ {
		b := cm.buckets[i]
		b.Clear()
	}
}
