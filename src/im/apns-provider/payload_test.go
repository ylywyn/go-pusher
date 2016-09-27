/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : payload_test.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main_test

import (
	apns "im/apns-provider"
	"testing"
)

func TestPayloadEmpty(t *testing.T) {
	pl := apns.NewPayload()
	b, err := pl.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON Error:%s", err.Error())
	}
	if string(b) != `{}` {
		t.Error("error result")
	}
}

func TestPayloadKV(t *testing.T) {
	pl := apns.NewPayload()
	pl.SetCustom("key", "value")
	b, err := pl.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON Error:%s", err.Error())
	}
	if string(b) != `{"key":"value"}` {
		t.Error("error result")
	}
}

func TestPayloadAps(t *testing.T) {
	pl := apns.NewPayload()

	aps := apns.Aps{}
	aps.SetAlert("hello")
	aps.SetBadge(nil)
	aps.SetSound("1.mp3")
	aps.SetContentAvailable()
	aps.SetCategory("category")

	pl.SetAps(&aps)

	b, err := pl.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON Error:%s", err.Error())
	}
	if string(b) != `{"aps":{"alert":"hello","category":"category","content-available":1,"sound":"1.mp3"}}` {
		t.Error("error result")
	}
}

func TestPayloadAlert(t *testing.T) {

	alert := apns.Alert{}
	alert.Title = "title"
	alert.Body = "body"

	aps := apns.Aps{}
	aps.SetAlert(alert)
	aps.SetBadge(nil)
	aps.SetSound("1.mp3")

	pl := apns.NewPayload()
	pl.SetAps(&aps)

	b, err := pl.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON Error:%s", err.Error())
	}
	if string(b) != `{"aps":{"alert":{"body":"body","title":"title"},"sound":"1.mp3"}}` {
		t.Error("error result")
	}
}

func TestPayloadApsKV(t *testing.T) {
	pl := apns.NewPayload()

	aps := apns.Aps{}
	aps.SetAlert("hello")
	aps.SetBadge(nil)
	aps.SetSound("1.mp3")
	aps.SetContentAvailable()
	aps.SetCategory("category")

	pl.SetAps(&aps)
	pl.SetCustom("key", "value")

	b, err := pl.MarshalJSON()
	if err != nil {
		t.Error("MarshalJSON Error:%s", err.Error())
	}
	if string(b) != `{"aps":{"alert":"hello","category":"category","content-available":1,"sound":"1.mp3"},"key":"value"}` {
		t.Error("error result")
	}
}
