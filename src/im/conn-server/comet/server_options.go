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
	"time"
)

//网络层参数
type ServerOptions struct {
	Addr string

	MaxDispatcherNum chan int //最大分发处理协程数
	ReadBufferSize   int      //读取缓冲大小
	WriteBufferSize  int      //写入缓冲大小
	WriteChannelSize int      //写异步channel长度
	ReadChannelSize  int      //读异步channel长度

	HeartBeat bool          //客户端是否发送心跳包
	IdleTime  time.Duration //连接空闲时间
	Cm        *ConnManager
}
