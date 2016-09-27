/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package conf

import (
	"flag"

	"im/common/conf"
)

var (
	gconf    *goconf.Config
	confFile string
	Conf     *Config
)

func init() {
	flag.StringVar(&confFile, "c", "./conn-conf.conf", " conn_server config file path")
}

type Config struct {
	// base
	PidFile   string `goconf:"base:pidfile"`
	LogDir    string `goconf:"base:logdir"`
	LogConf   string `goconf:"base:logconf"`
	MaxConn   uint32 `goconf:"base:maxconn"`
	ServerId  int    `goconf:"base:server.id"`
	PumpType  string `goconf:"base:pump_type"`
	PumpAddr  string `goconf:"base:pump_addr"`
	HeartBeat bool   `goconf:"base:heartbeat"`
	TimeOut   int    `goconf:"base:timeout"`

	// tcp
	TCPBind        string `goconf:"tcp:bind"`
	TCPReadbuf     int    `goconf:"tcp:readbuf"`
	TCPWritebuf    int    `goconf:"tcp:writebuf"`
	TCPSendChannel int    `goconf:"tcp:sendchannel"`
	TCPRecvChannel int    `goconf:"tcp:recvchannel"`

	// websocket
	WebsocketBind        string `goconf:"websocket:bind"`
	WebsocketReadbuf     int    `goconf:"websocket:readbuf"`
	WebsocketWritebuf    int    `goconf:"websocket:writebuf"`
	WebsocketSendChannel int    `goconf:"websocket:sendchannel"`
	WebsocketRecvChannel int    `goconf:"websocket:recvchannel"`

	// http
	HTTPBind string `goconf:"http:bind"`
}

func NewConfig() *Config {
	return &Config{
		// base section
		PidFile:   "/tmp/im/comet.pid",
		LogDir:    "./",
		LogConf:   "./log/xml",
		MaxConn:   800000,
		PumpType:  "nats",
		PumpAddr:  "nats://localhost:4222",
		ServerId:  1000,
		HeartBeat: true,
		TimeOut:   300,

		// tcp
		TCPBind:        "localhost:9088",
		TCPReadbuf:     2048,
		TCPWritebuf:    2048,
		TCPSendChannel: 512,
		TCPRecvChannel: 512,

		// websocket
		WebsocketBind:        "localhost:9078",
		WebsocketReadbuf:     2048,
		WebsocketWritebuf:    2048,
		WebsocketSendChannel: 512,
		WebsocketRecvChannel: 512,

		// http
		HTTPBind: "localhost:9058",
	}
}

// InitConfig init the global config.
func InitConfig() (err error) {
	Conf = NewConfig()
	gconf = goconf.New()
	if err = gconf.Parse(confFile); err != nil {
		return err
	}
	if err := gconf.Unmarshal(Conf); err != nil {
		return err
	}
	return nil
}

func ReloadConfig() (*Config, error) {
	conf := NewConfig()
	ngconf, err := gconf.Reload()
	if err != nil {
		return nil, err
	}
	if err := ngconf.Unmarshal(conf); err != nil {
		return nil, err
	}
	gconf = ngconf
	return conf, nil
}
