package routers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/auth/user",&controllers.UserController{},"post:AddUser;put:UpdateUser;delete:DelUser;get:GetUserList")
	beego.Router("/api/auth/passwd",&controllers.UserController{},"post:ResetPwd")
	beego.Router("/api/auth/sign",&controllers.UserController{},"post:Login;get:Logout")
	beego.Router("/api/auth/get",&controllers.AuthController{},"get:GetRole")
	beego.Router("/api/auth/new",&controllers.AuthController{},"post:AddUserAuth")
	beego.Router("/api/auth/post",&controllers.AuthController{},"post:AddOrUpdateRole")
	beego.Router("/api/auth/auths",&controllers.AuthController{},"get:GetUserAuthList")
	beego.Router("/api/authnode/post",&controllers.AuthController{},"post:UpdateRoleNode")
	beego.Router("/api/node/post",&controllers.AuthController{},"post:AddOrUpdateNode")
	beego.Router("/api/node/get",&controllers.AuthController{},"get:GetNodes")
	beego.Router("/api/idc",&controllers.IDCController{},"post:AddOrUpdateIdc;get:RequestIdcs")
	beego.Router("/api/marathon/conf",&controllers.MarathonCfController{},"post:AddOrUpdateMarathonConf")
	beego.Router("/api/registry/conf",&controllers.RegistryCfController{},"post:AddOrUpdateRegistryConf")
	beego.Router("/api/Cidr/Add",&controllers.CidrController{},"post:AddCidr")
	beego.Router("/api/service/Add",&controllers.ServiceController{},"post:AddOrUpdateService")
	beego.Router("/api/service/Delete",&controllers.ServiceController{},"post:DelService")
	beego.Router("/api/service/count",&controllers.ServiceController{},"get:GetSeparateCount")
	beego.Router("/api/service/get",&controllers.ServiceController{},"get:GetService")
	beego.Router("/api/docker/dashboard",&controllers.DockerController{},"get:DashBoard")
	beego.Router("/api/docker/scale",&controllers.DockerController{},"post:ScaleContainers")
}
