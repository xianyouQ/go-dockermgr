package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/utils"
	//"fmt"
)

type MainController struct {
	beego.Controller
}
var registryClient utils.RegistryClient

func init() {
	beego.SetStaticPath("/css","static/css")
	beego.SetStaticPath("/img","static/img")
	beego.SetStaticPath("/js","static/js")
	beego.SetStaticPath("/vendor","static/vendor")
	beego.SetStaticPath("/fonts","static/fonts")
	beego.SetStaticPath("/tpl","views")
	beego.SetStaticPath("/l10n","static/i10n")
	var registryserver,registryauthserver utils.ServerInfo
	registryPort,_ := beego.AppConfig.Int("registryserver.port")
	registryAuthPort,_ := beego.AppConfig.Int("registryauthserver.port")
	registryserver = utils.ServerInfo{Host: beego.AppConfig.String("registryserver.host"),Schema: beego.AppConfig.String("registryserver.schema"),Port: registryPort}
	registryauthserver = utils.ServerInfo{Host: beego.AppConfig.String("registryauthserver.host"),Schema: beego.AppConfig.String("registryauthserver.schema"),
		Port: registryAuthPort}
	registryClient = utils.RegistryClient{Server: registryserver,TokenAuthServer: registryauthserver,TokenAuthService: beego.AppConfig.String("registryserver.tokenauthservice"),
		TokenMap: map[string]string{},Username: beego.AppConfig.String("registryserver.username"),Password: beego.AppConfig.String("registryserver.password")}
}
func (c *MainController) Get() {
	c.TplName = "index.html"
}
