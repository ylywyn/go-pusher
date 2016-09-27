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

//1.流量,
//2.连接个数
//

type CometStat struct {
	Id              int
	ConnClosedCount uint32
	ConnCount       uint32
	FlowStat        FlowStat
}

func NewCometStat(id int) *CometStat {
	return &CometStat{
		Id:              id,
		ConnClosedCount: 0,
		ConnCount:       0,
	}
}

type CometStats struct {
	TcpStats       CometStat
	WebSocketStats CometStat
}
