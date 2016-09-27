/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description: golang im 简易客户端， 测试使用
 ******************************************************************/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

var gClient *TcpClient

func TcpConnStart(uid uint64, pwd string, server *net.TCPAddr) {
	c, err := NewClient(uid, pwd, server)
	if err != nil {
		fmt.Println("TcpConnStart Error :%d", err.Error())
		return
	}
	gClient = c
	c.Start()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始化配置
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err.Error())
	}

	//校验配置
	if len(Conf.ServerAddr) == 0 {
		fmt.Println("服务器地址错误!")
		return
	}

	servAddr, err := net.ResolveTCPAddr("tcp4", Conf.ServerAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	n := Conf.Num
	if n <= 0 || n > 50000 {
		fmt.Println("连接个数错误!")
		return
	}

	if n == 1 {
		fmt.Printf("启动IM客户端，UID: %d \r\n", Conf.Uid)
		go TcpConnStart(Conf.Uid, Conf.Passwd, servAddr)
	} else {
		fmt.Printf("启动连接数测试模式，连接个数: %d \r\n", n)
		index := n + 1
		for i := 1; i < index; i++ {
			go TcpConnStart(uint64(Conf.Start+i), "1", servAddr)
			time.Sleep(4 * time.Millisecond)
		}
	}

	bio := bufio.NewReader(os.Stdin)
	for {

		if n == 1 {
			if gClient != nil && gClient.Closed() {
				os.Exit(0)
			}

			l, _, err := bio.ReadLine()
			if err != nil {
				continue
			}

			if len(l) == 0 || len(l) > 1024 {
				fmt.Println("")
				continue
			}

			line := string(l)
			if CheckInput(line) {
				if gClient != nil && !gClient.Closed() {
					gClient.Write(line)
				} else {
					os.Exit(0)
				}
			}
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
