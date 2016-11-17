package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/auth"
)

type CommonController struct {
	beego.Controller
}


func (this *CommonController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJSON()
}

func init() {
	
	beego.SetStaticPath("/css","static/css")
	beego.SetStaticPath("/img","static/img")
	beego.SetStaticPath("/js","static/js")
	beego.SetStaticPath("/vendor","static/vendor")
	beego.SetStaticPath("/fonts","static/fonts")
	beego.SetStaticPath("/tpl","views")
	beego.SetStaticPath("/l10n","static/i10n")
	auth.AccessRegister()
}
