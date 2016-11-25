package routers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/auth/user",&controllers.UserController{},"post:AddUser;put:UpdateUser;delete:DelUser")
	beego.Router("/api/auth/sign",&controllers.UserController{},"post:Login;get:Logout")
	beego.Router("/api/idc",&controllers.IDCController{},"post:AddIdc;get:RequestIdcs")

}
