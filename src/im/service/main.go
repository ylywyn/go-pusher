/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : main.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	"flag"
	log "im/common/log4go"
	"im/service/conf"
	"im/service/http"
	_ "im/service/http"
	"im/service/logic/models/utils"
	"im/service/logic/service/id"
	. "im/service/logic/service/send"
	"im/service/logic/service/session"
	"im/service/msg"
	"im/service/rpc"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//init step
	flag.Parse()
	if err := conf.InitConfig(); err != nil {
		panic(err.Error())
	}

	//log init
	log.LoadConfiguration(conf.Conf.LogConf)
	log.Info("[push service[version:%s]] start", Version)
	//

	//server init
	err := utils.InitDb(conf.Conf)
	if err != nil {
		log.Error("[push service] error:%s", err.Error())
		panic(err.Error())
	}

	//id service
	err = id.InitId(conf.Conf.Cluster, conf.Conf.IdService)
	if err != nil {
		log.Error("[push service] error:%s", err.Error())
		panic(err.Error())
	}

	//session
	session.NewMemSessions(conf.Conf.Cluster, runtime.NumCPU()*8)

	// processor
	p := msg.NewProcessor(conf.Conf)
	p.Start()

	// sender
	Send = NewSender(conf.Conf)
	Send.Start()

	// http
	http.StartHttpService()

	//rpc
	err = rpc.RunRpcService(conf.Conf.RpcService)
	if err != nil {
		log.Error("[RunRpcServer] error:%s", err.Error())
		panic(err.Error())
	}

	InitSignal()

	//exit step
	Send.Stop()

	p.Stop()

	utils.UninitDb()

	log.Close()
}
