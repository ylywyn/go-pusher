/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : notification.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	"time"

	"github.com/pquerna/ffjson/ffjson"
)

const (
	PriorityLow = 5

	//立即触发
	PriorityHigh = 10
)

type Notification struct {
	ApnsID      string //赋给消息的UUID
	DeviceToken string
	Topic       string // 应用的bundle id
	Expiration  time.Time
	Priority    int
	Payload     interface{}
	Dev         bool
}

// MarshalJSON converts the notification payload to JSON.
func (n *Notification) MarshalJSON() ([]byte, error) {
	switch n.Payload.(type) {
	case string:
		return []byte(n.Payload.(string)), nil
	case []byte:
		return n.Payload.([]byte), nil
	default:
		return ffjson.Marshal(n.Payload)
	}
}
