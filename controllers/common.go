package controllers

import (
	"github.com/astaxie/beego"


)

type CommonController struct {
	beego.Controller
}


func (this *CommonController) Rsp(status bool, info string,data interface{}) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": info,"data":data}
	this.ServeJSON()
}


