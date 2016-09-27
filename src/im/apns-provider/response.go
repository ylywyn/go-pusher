/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : response.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package main

import (
	"net/http"
	"strconv"
	"time"
)

const StatusSent = http.StatusOK

//文档
//https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/APNsProviderAPI.html#//apple_ref/doc/uid/TP40008194-CH101-SW1
//

// table 6-6
const (
	ReasonPayloadEmpty              = "PayloadEmpty"
	ReasonPayloadTooLarge           = "PayloadTooLarge"
	ReasonBadTopic                  = "BadTopic"
	ReasonTopicDisallowed           = "TopicDisallowed"
	ReasonBadMessageID              = "BadMessageId"
	ReasonBadExpirationDate         = "BadExpirationDate"
	ReasonBadPriority               = "BadPriority"
	ReasonMissingDeviceToken        = "MissingDeviceToken"
	ReasonBadDeviceToken            = "BadDeviceToken"
	ReasonDeviceTokenNotForTopic    = "DeviceTokenNotForTopic"
	ReasonUnregistered              = "Unregistered"
	ReasonDuplicateHeaders          = "DuplicateHeaders"
	ReasonBadCertificateEnvironment = "BadCertificateEnvironment"
	ReasonBadCertificate            = "BadCertificate"
	ReasonForbidden                 = "Forbidden"
	ReasonBadPath                   = "BadPath"
	ReasonMethodNotAllowed          = "MethodNotAllowed"
	ReasonTooManyRequests           = "TooManyRequests"
	ReasonIdleTimeout               = "IdleTimeout"
	ReasonShutdown                  = "Shutdown"
	ReasonInternalServerError       = "InternalServerError"
	ReasonServiceUnavailable        = "ServiceUnavailable"
	ReasonMissingTopic              = "MissingTopic"
)

type Response struct {

	// table 6-4
	StatusCode int
	Reason     string

	// Notification ID，UUID
	ApnsID    string
	Timestamp Time
}

func (c *Response) Sent() bool {
	return c.StatusCode == StatusSent
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	t.Time = time.Unix(int64(ts/1000), 0)
	return nil
}
