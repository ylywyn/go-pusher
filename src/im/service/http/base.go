/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package http

import (
	log "im/common/log4go"
	"im/service/conf"
	_ "im/service/http/routers"

	"github.com/astaxie/beego"
)

func StartHttpService() {
	//beego.BConfig.RunMode = "dev"
	beego.BConfig.RunMode = "prod"
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.AutoRender = false

	log.Debug("[httpservice|init|Run] run: %s", conf.Conf.HttpService)
	go beego.Run(conf.Conf.HttpService)
}
