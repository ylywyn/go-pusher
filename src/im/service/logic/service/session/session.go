/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : session.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package session

import (
	"errors"
	"fmt"
	"im/common/crypto/md5"
	"im/common/proto/fbsgen/session"
	"im/service/logic/models/utils"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	fb "github.com/google/flatbuffers/go"
)

const KEY_FORMAT = "PUSH_USER_%d"

var (
	ErrNotFound     = errors.New("not found")
	ErrBadSessionId = errors.New("Bad SessionId")
)

type Session struct {
	OsType byte
	Admin  byte
	Appid  uint16
	Uid    uint64
	//单点登录。多点的话，应该保存多个ConnServerId， ConnId
	ConnClusterId uint32
	ConnSessionId uint64
	AckMsgId      uint64
	ApnsToken     string
	Key           string
}

func (this *Session) Get() error {
	//先从内存查找
	s := memSessions.GetSession(this.Uid)
	if s != nil {
		*this = *s
		return nil
	} else {
		if !clusterSession {
			return ErrNotFound
		}
	}

	//redis查找
	key := fmt.Sprintf(KEY_FORMAT, this.Uid)

	c := utils.Redis.Get()
	data, err := redis.Bytes(c.Do("GET", key))
	c.Close()

	if err != nil || data == nil {
		return err
	}

	return this.unSerialize(data)
}

//sid : appid_uid_md5(key)
func (this *Session) GetBySessionId(sid string) error {

	sids := strings.Split(sid, "_")
	if len(sids) != 3 {
		return ErrBadSessionId
	}

	var err error
	this.Uid, err = strconv.ParseUint(sids[1], 10, 0)
	if err != nil {
		return err
	}

	//	key := fmt.Sprintf(KEY_FORMAT, this.Uid)

	//	c := utils.Redis.Get()
	//	data, err := redis.Bytes(c.Do("GET", key))
	//	c.Close()

	//	if err != nil || data == nil {
	//		return err
	//	}
	//	this.unSerialize(data)
	err = this.Get()
	if err != nil {
		return err
	}

	flag := fmt.Sprintf("%s_%s", strconv.Itoa(int(this.Appid)), strconv.FormatUint(this.Uid, 10))
	if md5.Md5([]byte(flag)) == sids[2] {
		return nil
	}
	return ErrBadSessionId
}

//
func (this *Session) Put() error {
	memSessions.PutSession(this.Uid, this)
	if !clusterSession {
		return nil
	}

	data := this.serialize()

	key := fmt.Sprintf(KEY_FORMAT, this.Uid)

	c := utils.Redis.Get()
	_, err := c.Do("SET", key, data)
	c.Close()

	return err
}

func (this *Session) Del() error {
	memSessions.DelSession(this.Uid)
	if !clusterSession {
		return nil
	}

	key := fmt.Sprintf(KEY_FORMAT, this.Uid)

	c := utils.Redis.Get()
	_, err := c.Do("DEL", key)
	c.Close()

	return err
}

//序列化
func (this *Session) serialize() []byte {

	builder := fb.NewBuilder(300)

	var keyof, tokenof fb.UOffsetT
	if len(this.Key) > 0 {
		keyof = builder.CreateString(this.Key)
	}

	if len(this.ApnsToken) > 0 {
		tokenof = builder.CreateString(this.ApnsToken)
	}

	session.SessionStart(builder)
	session.SessionAddLoginType(builder, this.OsType)
	session.SessionAddAdmin(builder, this.Admin)
	session.SessionAddAppid(builder, this.Appid)
	session.SessionAddUid(builder, this.Uid)
	session.SessionAddConnClusterId(builder, this.ConnClusterId)
	session.SessionAddConnSessionId(builder, this.ConnSessionId)
	session.SessionAddAckMsgId(builder, this.AckMsgId)

	if keyof > 0 {
		session.SessionAddKey(builder, keyof)
	}

	if tokenof > 0 {
		session.SessionAddToken(builder, tokenof)
	}

	u := session.SessionEnd(builder)
	builder.Finish(u)

	return builder.Bytes[builder.Head():]
}

//反序列化
func (this *Session) unSerialize(buf []byte) error {
	n := fb.GetUOffsetT(buf)
	m := &session.Session{}
	m.Init(buf, n)

	this.OsType = m.LoginType()
	this.Admin = m.Admin()
	this.Appid = m.Appid()
	this.Uid = m.Uid()
	this.ConnClusterId = m.ConnClusterId()
	this.ConnSessionId = m.ConnSessionId()
	this.Key = string(m.Key())
	this.ApnsToken = string(m.Token())

	return nil
}

func NewSessionId(appid uint32, uid uint64, key string) string {
	flag := fmt.Sprintf("%s_%s", strconv.Itoa(int(appid)), strconv.FormatUint(uid, 10))
	return fmt.Sprintf("%s_%s", flag, md5.Md5([]byte(flag)))
}
