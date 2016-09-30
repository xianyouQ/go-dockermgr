package routers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
