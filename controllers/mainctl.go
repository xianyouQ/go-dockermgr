package controllers

import (
	"github.com/xianyouQ/go-dockermgr/auth"
	"github.com/astaxie/beego"
)

type MainController struct {
	CommonController
}



func (c *MainController) Get() {
	c.TplName = "index.html"
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