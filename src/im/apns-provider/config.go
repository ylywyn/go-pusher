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
	"strings"
	//log "im/common/log4go"
)

var (
	gconf    *goconf.Config
	Conf     *Config
	confFile string
)

type Options struct {
	name         string
	id           int
	bundleid     string
	pemFile      string
	pemPasswd    string
	allowDevelop bool
	connCount    int
	sendChanSize int
}

type Config struct {
	pidFile    string
	logDir     string
	logConf    string
	pumpType   string
	pumpAddr   string
	subString  string
	servesOpts []*Options
}

func init() {
	flag.StringVar(&confFile, "c", "./apns-conf.conf", " apns_provider config file path")
}

func InitConfig() (err error) {

	Conf = &Config{
		pidFile:  "/tmp/im/apns.pid",
		logDir:   "./",
		logConf:  "./log.xml",
		pumpAddr: "",
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

	Conf.pidFile, _ = base.String("pidfile")
	Conf.logDir, _ = base.String("logdir")
	Conf.logConf, _ = base.String("logconf")
	Conf.subString, _ = base.String("substring")
	Conf.pumpType, _ = base.String("pump_type")
	Conf.pumpAddr, err = base.String("pump_addr")
	if err != nil {
		return err
	}

	sections := gconf.Sections()
	count := len(sections)
	if count == 1 {
		return errors.New("can't find  _service sections")
	}

	Conf.servesOpts = make([]*Options, 0, count)

	for _, k := range sections {
		if !strings.Contains(k, "_service") {
			continue
		}

		op := &Options{}
		ops := gconf.Get(k)

		num, _ := ops.Int("appid")
		op.id = int(num)

		op.bundleid, _ = ops.String("bundleid")
		op.name, _ = ops.String("name")
		op.pemFile, _ = ops.String("pem_file")
		op.pemPasswd, _ = ops.String("pem_passwd")
		if op.pemPasswd == "\"\"" {
			op.pemPasswd = ""
		}

		op.allowDevelop, _ = ops.Bool("allow_develop")

		num, _ = ops.Int("conn_count")
		op.connCount = int(num)

		num, _ = ops.Int("send_channel_size")
		op.sendChanSize = int(num)

		Conf.servesOpts = append(Conf.servesOpts, op)
	}
	//read service ops

	return nil
}
