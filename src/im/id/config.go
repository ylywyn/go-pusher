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
)

var (
	Conf     *Config
	confFile string
	gconf    *goconf.Config
)

type Config struct {
	Addr string
}

func init() {
	flag.StringVar(&confFile, "conf", "./id-conf.conf", " service config file path")
}

func InitConfig() (err error) {

	Conf = &Config{
		Addr: ":9999",
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

	Conf.Addr, err = base.String("addr")
	return err
}
