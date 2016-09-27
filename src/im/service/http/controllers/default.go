/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package controllers

//"github.com/astaxie/beego"

type TestController struct {
	BaseController
}

func (c *TestController) GetTest() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	//c.TplName = "index.tpl"
	//c.SetError("hello")
	c.SetData("ok")
}
