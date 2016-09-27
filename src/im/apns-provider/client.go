/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : client.go
 *  Date   :
 *  Author : yangl
 *  Description: http2.0 client
 ******************************************************************/

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/net/http2"
)

const (
	HostDevelopment = "https://api.development.push.apple.com"
	HostProduction  = "https://api.push.apple.com"
)

var DefaultHost = HostDevelopment

type Client struct {
	HTTPClient  *http.Client
	Certificate tls.Certificate
	Host        string
}

//NewClient
func NewClient(cert tls.Certificate) *Client {
	tlsconf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	if len(cert.Certificate) > 0 {
		tlsconf.BuildNameToCertificate()
	}
	tr := &http2.Transport{
		TLSClientConfig: tlsconf,
	}
	return &Client{
		HTTPClient:  &http.Client{Transport: tr},
		Certificate: cert,
		Host:        DefaultHost,
	}
}

func (c *Client) Development() *Client {
	c.Host = HostDevelopment
	return c
}

func (c *Client) Production() *Client {
	c.Host = HostProduction
	return c
}

func (c *Client) Push(n *Notification) (*Response, error) {
	payload, err := ffjson.Marshal(n)

	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v/3/device/%v", c.Host, n.DeviceToken)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	setHeaders(req, n)
	httpRes, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	response := &Response{}
	response.StatusCode = httpRes.StatusCode
	response.ApnsID = httpRes.Header.Get("apns-id")

	decoder := json.NewDecoder(httpRes.Body)
	if err := decoder.Decode(&response); err != nil && err != io.EOF {
		return &Response{}, err
	}
	return response, nil
}

func setHeaders(r *http.Request, n *Notification) {
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	if n.Topic != "" {
		r.Header.Set("apns-topic", n.Topic)
	}
	if n.ApnsID != "" {
		r.Header.Set("apns-id", n.ApnsID)
	}
	if n.Priority > 0 {
		r.Header.Set("apns-priority", fmt.Sprintf("%v", n.Priority))
	}
	if !n.Expiration.IsZero() {
		r.Header.Set("apns-expiration", fmt.Sprintf("%v", n.Expiration.Unix()))
	}
}
