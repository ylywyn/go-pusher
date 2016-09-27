/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description: 长连接服务
 ******************************************************************/

package main

import (
	"flag"
	log "im/common/log4go"
	"im/conn-server/comet"
	"im/conn-server/conf"
	"im/conn-server/web"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//conf init
	flag.Parse()
	if err := conf.InitConfig(); err != nil {
		panic(err.Error())
	}

	//log init
	log.LoadConfiguration(conf.Conf.LogConf)
	log.Info("[Globe|conn-server[version:%s]] start", Version)

	//server init
	comet.Server = comet.NewCometServer()
	comet.Server.Start()

	//http comet
	web.StartHttpMonitor()

	InitSignal()

	comet.Server.Stop()
	//exit
	log.Close()
}
