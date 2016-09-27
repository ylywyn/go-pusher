/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : config.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package conf

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
	PidFile     string
	LogDir      string
	LogConf     string
	PumpType    string
	PumpAddr    string
	RedisAddr   []string
	MongodbUser string
	MongodbMsg  string
	IdService   string
	HttpService string
	RpcService  string
	Cluster     bool
}

func init() {
	flag.StringVar(&confFile, "conf", "./serv-conf.conf", " service config file path")
}

func InitConfig() (err error) {

	Conf = &Config{
		PidFile:  "/tmp/im/apns.pid",
		LogDir:   "./",
		LogConf:  "./log.xml",
		PumpAddr: "",
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

	Conf.PidFile, _ = base.String("pidfile")
	Conf.LogDir, _ = base.String("logdir")
	Conf.LogConf, _ = base.String("logconf")
	Conf.PumpType, _ = base.String("pump_type")
	Conf.PumpAddr, err = base.String("pump_addr")
	if err != nil {
		return err
	}

	Conf.RedisAddr, _ = base.Strings("redis", ",")
	Conf.MongodbUser, _ = base.String("mongodb_user")
	Conf.MongodbMsg, _ = base.String("mongodb_msg")
	Conf.IdService, _ = base.String("id_service")
	Conf.HttpService, _ = base.String("http_service")
	Conf.RpcService, _ = base.String("rpc_service")
	Conf.Cluster, _ = base.Bool("cluster")
	return nil
}
