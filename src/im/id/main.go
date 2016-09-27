/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : config.go
 *  Date   :
 *  Author : yangl
 *  Description: 提供集群唯一ID服务
 ******************************************************************/

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//init step
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err.Error())
	}

	//fmt.Println(Conf.Addr)
	if len(Conf.Addr) == 0 {
		return
	}
	// 创建 tcp listener
	addr, err := net.ResolveTCPAddr("tcp4", Conf.Addr)
	checkError(err)

	listener, err := net.ListenTCP("tcp", addr)
	checkError(err)

	// 创建 serverR
	config := &Option{
		SendChanCount: 64,
		RecvChanCount: 64,
	}
	s := NewServer(config, &IdCallback{}, &IdProtocol{})

	// 启动服务
	go s.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chsig := make(chan os.Signal)
	signal.Notify(chsig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chsig)

	// 停止
	s.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
