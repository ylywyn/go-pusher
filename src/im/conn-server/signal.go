/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : signal.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	log "im/common/log4go"
	"im/conn-server/conf"
	"os"
	"os/signal"
	"syscall"
)

// InitSignal register signals handler.
func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("comet[%s] get a signal %s", Version, s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			reload()
		default:
			return
		}
	}
}

func reload() {
	newConf, err := conf.ReloadConfig()
	if err != nil {
		log.Error("ReloadConfig() error(%v)", err)
		return
	}
	conf.Conf = newConf
}
