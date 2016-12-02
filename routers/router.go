package routers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/auth/user",&controllers.UserController{},"post:AddUser;put:UpdateUser;delete:DelUser")
	beego.Router("/api/auth/sign",&controllers.UserController{},"post:Login;get:Logout")
	beego.Router("/api/idc",&controllers.IDCController{},"post:AddOrUpdateIdc;get:RequestIdcs")
	beego.Router("/api/marathon/conf",&controllers.MarathonCfController{},"post:AddOrUpdateMarathonConf")
	beego.Router("/api/registry/conf",&controllers.RegistryCfController{},"post:AddOrUpdateRegistryConf")
	beego.Router("/api/Cidr/Add",&controllers.CidrController{},"post:AddCidr")
	beego.Router("/api/service/Add",&controllers.ServiceController{},"post:AddService")
	beego.Router("/api/service/Delete",&controllers.ServiceController{},"post:DelService")
	beego.Router("/api/service/count",&controllers.ServiceController{},"get:GetSeparateCount")
	beego.Router("/api/service/get",&controllers.ServiceController{},"get:GetService")
	beego.Router("/api/docker/dashboard",&controllers.DockerController{},"get:DashBoard")
}
