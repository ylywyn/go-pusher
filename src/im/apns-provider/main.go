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
	"im/common/signal"
	"runtime"
)

func sendNotification() {
	cert, pemErr := LoadPemFile("dev_pie+.pem", "")
	if pemErr != nil {
		log.Debug("Cert Error:", pemErr)
		return
	}

	n := &Notification{}
	n.DeviceToken = "690f07a189532029fd79cb162d317a3ec0425b4ac257076678a3be3d1b6c991a"
	n.Topic = "com.htht.pieplus"
	n.Payload = []byte(`{
		  "aps" : {
			"alert" : "不错，可以的哦，哈哈哈!"
		  }
		}
	`)

	client := NewClient(cert).Development()
	res, err := client.Push(n)

	if err != nil {
		log.Debug("Error:", err)
		return
	}

	log.Debug("APNs ID: %s", res.ApnsID)
}

func main() {
	//test
	sendNotification()

	runtime.GOMAXPROCS(runtime.NumCPU())

	//conf init
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err.Error())
	}

	//log init
	log.LoadConfiguration(Conf.logConf)
	log.Info("[apns-provider version:%s] start", Version)

	//server init
	apns := NewApnsProvider()
	go apns.Start()

	signal.HandleSignal(signal.InitSignal())

	//exit
	apns.Stop()
	log.Close()
}
