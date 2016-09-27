/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : base.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package utils

import (
	"errors"
	log "im/common/log4go"
	"im/service/conf"
)

var Redis *RedisPool
var UserDBPool *MongoDBPool
var MsgDBPool *MongoDBPool

const (
	useDBName = "push_users_db"
	msgDBName = "push_msg_db"
)

func InitDb(conf *conf.Config) error {
	var err error
	if conf.Cluster {
		if len(conf.RedisAddr) == 0 || len(conf.RedisAddr[0]) == 0 {
			return errors.New("redis addr invalid")
		}

		log.Debug("[utils|InitDb|redis] conn redis: %s", conf.RedisAddr[0])

		Redis, err = NewRedisPool(conf.RedisAddr)
		if err != nil {
			return err
		}
	}

	log.Debug("[utils|InitDb|user MongoDB] conn MongoDB: %s", conf.MongodbUser)
	UserDBPool, err = NewMongoDBPool(conf.MongodbUser, useDBName)
	if err != nil {
		return err
	}

	log.Debug("[utils|InitDb|msg MongoDB] conn MongoDB: %s", conf.MongodbMsg)
	MsgDBPool, err = NewMongoDBPool(conf.MongodbMsg, msgDBName)
	if err != nil {
		return err
	}

	return nil
}

func UninitDb() {
	Redis.Close()
}
