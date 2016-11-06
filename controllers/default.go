package controllers

import (
	"github.com/astaxie/beego"
	//"fmt"
)

type MainController struct {
	beego.Controller
}


func init() {
	beego.SetStaticPath("/css","static/css")
	beego.SetStaticPath("/img","static/img")
	beego.SetStaticPath("/js","static/js")
	beego.SetStaticPath("/vendor","static/vendor")
	beego.SetStaticPath("/fonts","static/fonts")
	beego.SetStaticPath("/tpl","views")
	beego.SetStaticPath("/l10n","static/i10n")

}
func (c *MainController) Get() {
	c.TplName = "index.html"
}
