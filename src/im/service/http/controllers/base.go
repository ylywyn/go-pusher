/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   : base.go
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/
package controllers

import (
	"bytes"
	log "im/common/log4go"

	"github.com/astaxie/beego"
)

type JsonRet struct {
	Result bool        `json:"result"`
	Err    string      `json:"error"`
	Data   interface{} `json:"data"`
}
type BaseController struct {
	beego.Controller
	NotAutoJson bool
	jsonp       bool
	result      JsonRet
}

const (
	ErrJson = `{"result":false,"error":"","data":null}`
	OkJson  = `{"result":true,"error":"","data":`
)

func (this *BaseController) Prepare() {
	callback := this.GetString("callback")
	if callback != "" {
		this.jsonp = true
	}
}

func (this *BaseController) Finish() {
	if this.NotAutoJson {
		if this.result.Result {
			str := this.result.Data.(string)
			buf := bytes.Buffer{}
			buf.Grow(len(str) + 128)
			buf.WriteString(OkJson)
			buf.WriteString(str)
			buf.WriteString("}")
			this.Ctx.Output.ContentType(".json")
			this.Ctx.WriteString(buf.String())

		} else {
			this.Ctx.Output.ContentType(".json")
			this.Ctx.WriteString(ErrJson)
		}
	} else {
		if this.jsonp {
			this.Data["jsonp"] = this.result
			this.ServeJSONP()
		} else {
			this.Data["json"] = this.result
			this.ServeJSON()
		}

	}
}

func (this *BaseController) SetError(err string) {
	this.result.Err = err
	this.result.Result = false
}

func (this *BaseController) SetData(data interface{}) {
	this.NotAutoJson = false
	this.result.Data = data
	this.result.Result = true
}

func (this *BaseController) CheckErr(err error) bool {
	if err != nil {
		log.Debug(err.Error())
		return false
	} else {
		return true
	}
}
