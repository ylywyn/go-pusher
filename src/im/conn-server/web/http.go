/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author :
 *  Description:
 ******************************************************************/

package web

import (
	"encoding/json"
	"fmt"
	log "im/common/log4go"
	"im/conn-server/conf"
	"net"
	"net/http"
	"time"
)

const (
	RET          = "ret"
	DATA         = "data"
	OK           = "true"
	CONTENT_TYPE = "Content-Type"
	JSON_TYPE    = "application/json;charset=utf-8"
	JSONP_TYPE   = "application/javascript;charset=utf-8"
)

func StartHttpMonitor() {
	httpSerMux := http.NewServeMux()
	httpSerMux.HandleFunc("/stats/cpu/", getCpuStat)
	httpSerMux.HandleFunc("/stats/mem/", getMemStat)
	httpSerMux.HandleFunc("/stats/network/", getNetWorkStat)
	httpSerMux.HandleFunc("/stats/all/", getAllStat)

	log.Debug("[Web | Http listen ] %s", conf.Conf.HTTPBind)
	go httpListen(httpSerMux, conf.Conf.HTTPBind)
}

func httpListen(mux *http.ServeMux, addr string) error {
	hs := http.Server{
		Handler:      mux,
		ReadTimeout:  15,
		WriteTimeout: 15,
	}

	hs.SetKeepAlivesEnabled(true)

	addrTcp, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Error("[Web|ResolveTCPAddr] Resolve TCP Addr %s error: (%v)", addr, err)
		return err
	}

	l, err := net.ListenTCP("tcp4", addrTcp)
	if err != nil {
		log.Error("[Web|httpListen] Listen TCP error: (%v)", err)
		return err
	}

	err = hs.Serve(l)
	if err != nil {
		log.Error("[Web|Serve] error: (%v)", err)
		return err
	}

	return nil
}

func getCpuStat(w http.ResponseWriter, r *http.Request) {

	p := r.URL.Query()
	cb := p.Get("cb")
	retMap := map[string]interface{}{RET: OK, DATA: nil}
	defer writeResponse(w, r, retMap, cb, time.Now())

	if r.Method != "GET" {
		retMap[RET] = false
		retMap[DATA] = "Please Use Get Method"
		return
	}

	CpuStat(retMap)
	return
}

func getMemStat(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()
	cb := p.Get("cb")
	retMap := map[string]interface{}{RET: OK, DATA: nil}
	defer writeResponse(w, r, retMap, cb, time.Now())

	if r.Method != "GET" {
		retMap[RET] = false
		retMap[DATA] = "Please Use Get Method"
		return
	}

	MemStat(retMap)
	return
}

func getNetWorkStat(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()
	cb := p.Get("cb")
	retMap := map[string]interface{}{RET: OK, DATA: nil}
	defer writeResponse(w, r, retMap, cb, time.Now())

	if r.Method != "GET" {
		retMap[RET] = false
		retMap[DATA] = "Please Use Get Method"
		return
	}

	NetWorkStat(retMap)
	return
}

func getAllStat(w http.ResponseWriter, r *http.Request) {
}

func writeResponse(w http.ResponseWriter, r *http.Request, ret map[string]interface{}, cb string, start time.Time) {
	data, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, "Json Marshal Error", 500)
		log.Error("[Http|writeResponse]json.Marshal(\"%v\") error(%v)", ret, err)
		return
	}

	retStr := ""
	if cb == "" {
		w.Header().Set(CONTENT_TYPE, JSON_TYPE)
		retStr = string(data)
	} else {
		w.Header().Set(CONTENT_TYPE, JSONP_TYPE)
		retStr = fmt.Sprintf("%s(%s)", cb, string(data))
	}

	if _, err := w.Write([]byte(retStr)); err != nil {
		log.Error("[Http|writeResponse]w.Write(\"%s\") error(%v)", retStr, err)
	}
}
