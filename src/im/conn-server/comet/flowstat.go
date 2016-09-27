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
	"fmt"
	"sync/atomic"
)

//流量状态

type FlowStat struct {
	ReadCounts  int64
	ReadBytes   int64
	WriteCounts int64
	WriteBytes  int64
}

func (this *FlowStat) Stat() *FlowStat {
	stat := &FlowStat{
		ReadCounts:  this.ReadCounts,
		ReadBytes:   this.ReadBytes,
		WriteCounts: this.WriteCounts,
		WriteBytes:  this.WriteBytes,
	}

	this.Reset()
	return stat
}

func (this *FlowStat) IncrReadCounts() {
	atomic.AddInt64(&this.ReadCounts, int64(1))
}

func (this *FlowStat) IncrWriteCounts() {
	atomic.AddInt64(&this.WriteCounts, int64(1))
}

func (this *FlowStat) IncrReadBytes(n int32) {
	atomic.AddInt64(&this.ReadBytes, int64(n))
}

func (this *FlowStat) IncrWriteBytes(n int32) {
	atomic.AddInt64(&this.WriteBytes, int64(n))
}

func (this *FlowStat) Reset() {
	atomic.StoreInt64(&this.ReadCounts, 0)
	atomic.StoreInt64(&this.ReadBytes, 0)
	atomic.StoreInt64(&this.WriteCounts, 0)
	atomic.StoreInt64(&this.WriteBytes, 0)
}

func (this *FlowStat) String() string {
	return fmt.Sprintf("read:%d/%d\twrite:%d/%d", this.ReadBytes, this.ReadCounts,
		this.WriteBytes, this.WriteCounts)
}
