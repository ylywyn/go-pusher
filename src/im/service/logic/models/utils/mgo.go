/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : mgo.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package utils

import (
	"errors"
	log "im/common/log4go"

	"gopkg.in/mgo.v2"
)

var (
	ErrExist = errors.New("exist")
)

type MongoDBPool struct {
	session *mgo.Session
	name    string
	addr    string
}

func NewMongoDBPool(addr, name string) (*MongoDBPool, error) {
	mdb := &MongoDBPool{}
	if mdb.Init(addr, name) == nil {
		return mdb, nil
	} else {
		return nil, errors.New("con't connect to mongodb" + addr)
	}
}

func (this *MongoDBPool) Init(addr, name string) error {

	this.name = name
	this.addr = addr
	if this.session == nil {
		var err error
		this.session, err = mgo.Dial(addr)
		if err != nil {
			return err
		}
	}
	log.Debug("connect MongoDB: %s,%s, ok ", addr, name)
	this.session.SetMode(mgo.Monotonic, true)
	return nil
}

func (this *MongoDBPool) Session() (*mgo.Session, error) {
	if this.session != nil {
		return this.session.Clone(), nil
	} else {
		return nil, errors.New("session is nil")
	}

}

func (this *MongoDBPool) M(collection string, f func(*mgo.Collection) error) error {
	session, err := this.Session()

	if err != nil {
		return err
	}
	defer func() {
		session.Close()
		if err := recover(); err != nil {

			log.Error("db.M recover..")
		}
	}()

	c := session.DB(this.name).C(collection)
	return f(c)
}
