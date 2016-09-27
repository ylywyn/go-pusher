/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : payload.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	"github.com/pquerna/ffjson/ffjson"
)

//官方文档 :
// The Remote Notification Payload
// https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/TheNotificationPayload.html#//apple_ref/doc/uid/TP40008194-CH107-SW1

type Aps struct {
	AlertContent     interface{} `json:"alert,omitempty"`
	Badge            interface{} `json:"badge,omitempty"`
	Category         string      `json:"category,omitempty"`
	ContentAvailable int         `json:"content-available,omitempty"`
	Sound            string      `json:"sound,omitempty"`
}

type Alert struct {
	Action       string   `json:"action,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	Body         string   `json:"body,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	Title        string   `json:"title,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
}

//string
//struct
func (aps *Aps) SetAlert(alert interface{}) *Aps {
	aps.AlertContent = alert
	return aps
}

//nil: 不改变旧直
//o  : 清零
//num: 新值
func (aps *Aps) SetBadge(b interface{}) *Aps {
	aps.Badge = b
	return aps
}

//string
func (aps *Aps) SetSound(sound string) *Aps {
	aps.Sound = sound
	return aps
}

func (aps *Aps) SetCategory(category string) *Aps {
	aps.Category = category
	return aps
}

//	{"aps":{"content-available":1}}
func (aps *Aps) SetContentAvailable() *Aps {
	aps.ContentAvailable = 1
	return aps
}

//Payload
type Payload struct {
	c map[string]interface{}
}

func NewPayload() *Payload {
	return &Payload{
		c: make(map[string]interface{}),
	}
}

func (p *Payload) SetAps(aps *Aps) *Payload {
	p.c["aps"] = aps
	return p
}

func (p *Payload) SetCustom(key string, val interface{}) *Payload {
	p.c[key] = val
	return p
}

// MarshalJSON returns the JSON encoded version of the Payload
func (p *Payload) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(p.c)
}
