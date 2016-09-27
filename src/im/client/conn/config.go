/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : config.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	"errors"
	"flag"
	"im/common/conf"
	"strconv"
)

var (
	Conf     *Config
	confFile string
	gconf    *goconf.Config
)

type Config struct {
	ServerAddr string
	Num        int
	Start      int
	Uid        uint64
	Appid      int
	Passwd     string
	HeartBeat  int
	To         uint64
	MsgTypes   string

	WriteTestSpan int
	WriteTestMsg  string
}

func init() {
	flag.StringVar(&confFile, "conf", "./conf.conf", " service config file path")
}

func InitConfig() (err error) {

	Conf = &Config{
		ServerAddr: ":9088",
	}

	gconf = goconf.New()
	if err = gconf.Parse(confFile); err != nil {
		return err
	}

	// read base section
	base := gconf.Get("base")
	if base == nil {
		return errors.New("can't find base section")
	}

	Conf.ServerAddr, err = base.String("server")

	appid, err := base.String("appid")
	Conf.Appid, err = strconv.Atoi(appid)

	hb, err := base.String("heartbeat")
	Conf.HeartBeat, err = strconv.Atoi(hb)

	uid, err := base.String("uid")
	Conf.Uid, err = strconv.ParseUint(uid, 10, 0)

	Conf.Passwd, err = base.String("pwd")

	to, err := base.String("msg_to")
	Conf.To, err = strconv.ParseUint(to, 10, 0)

	Conf.MsgTypes, err = base.String("msg_type")

	//测试模式参数
	cnum, err := base.String("clientnum")
	Conf.Num, err = strconv.Atoi(cnum)

	start, err := base.String("clientidstart")
	Conf.Start, err = strconv.Atoi(start)

	span, err := base.String("write_test_span")
	Conf.WriteTestSpan, err = strconv.Atoi(span)

	Conf.WriteTestMsg, err = base.String("write_test_msg")

	return err
}
