/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package login

import (
	//"encoding/json"

	"github.com/pquerna/ffjson/ffjson"
)

type MsgLogin struct {
	Uid       uint64
	User      string
	PassWd    string
	Key       string
	Platform  uint8
	LastMsgId uint64
}

func (msg *MsgLogin) UnSerialize(buf []byte) error {
	err := ffjson.Unmarshal(buf, msg)
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgLogin) Serialize() ([]byte, error) {

	return ffjson.Marshal(msg)
}

type MsgLoginRet struct {
	Uid  uint64
	Ret  bool
	Data string
}

func (msg *MsgLoginRet) Serialize() ([]byte, error) {
	return ffjson.Marshal(msg)
}

func (msg *MsgLoginRet) UnSerialize(buf []byte) error {
	err := ffjson.Unmarshal(buf, msg)
	if err != nil {
		return err
	}
	return nil
}
