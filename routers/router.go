package routers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/auth/user", &controllers.UserController{}, "post:AddUser;delete:DelUser;get:GetUserList") //system;get:ops
	beego.Router("/api/auth/passwd", &controllers.UserController{}, "post:ResetPwd;put:Changepwd")               //post:system;put:base
	//beego.Router("/api/auth/passwd", &controllers.UserController{}, "put:Changepwd")
	beego.Router("/api/auth/sign", &controllers.UserController{}, "post:Login;get:Logout")                  //base
	beego.Router("/api/auth/get", &controllers.AuthController{}, "get:GetRole")                             //base
	beego.Router("/api/auth/new", &controllers.AuthController{}, "post:AddUserAuth")                        //ops
	beego.Router("/api/auth/post", &controllers.AuthController{}, "post:AddOrUpdateRole")                   //system
	beego.Router("/api/auth/auths", &controllers.AuthController{}, "get:GetUserAuthList")                   //ops
	beego.Router("/api/authnode/post", &controllers.AuthController{}, "post:UpdateRoleNode")                //system
	beego.Router("/api/node/post", &controllers.AuthController{}, "post:AddOrUpdateNode")                   //system
	beego.Router("/api/node/get", &controllers.AuthController{}, "get:GetNodes")                            //system
	beego.Router("/api/idc", &controllers.IDCController{}, "post:AddOrUpdateIdc;get:RequestIdcs")           //post:system;get:base
	beego.Router("/api/marathon/conf", &controllers.MarathonCfController{}, "post:AddOrUpdateMarathonConf") //system
	beego.Router("/api/registry/conf", &controllers.RegistryCfController{}, "post:AddOrUpdateRegistryConf") //system
	beego.Router("/api/Cidr/Add", &controllers.CidrController{}, "post:AddCidr")                            //system
	beego.Router("/api/service/Add", &controllers.ServiceController{}, "post:AddOrUpdateService")           //system
	beego.Router("/api/service/Delete", &controllers.ServiceController{}, "post:DelService")                //system
	beego.Router("/api/service/count", &controllers.ServiceController{}, "get:GetSeparateCount")            //base
	beego.Router("/api/service/get", &controllers.ServiceController{}, "get:GetService")                    //base
	beego.Router("/api/docker/dashboard", &controllers.DockerController{}, "get:DashBoard")                 //system
	beego.Router("/api/docker/scale", &controllers.DockerController{}, "post:ScaleContainers")              //ops
	beego.Router("/api/docker/list", &controllers.DockerController{}, "post:GetContainers")                 //ops
	beego.Router("/api/release/task", &controllers.ReleaseController{}, "post:QueryReleaseTasks")           //ops dev qa
	beego.Router("/api/release/newtask", &controllers.ReleaseController{}, "post:NewReleaseTask")           //dev
	beego.Router("/api/release/review", &controllers.ReleaseController{}, "post:ReviewReleaseTask")         //qa
	beego.Router("/api/release/operate", &controllers.ReleaseController{}, "post:OperationReleaseTask")     //ops
	beego.Router("/api/release/cancel", &controllers.ReleaseController{}, "post:CancelReleaseTask")         //ops dev qa
	beego.Router("/api/release/conf", &controllers.ReleaseController{}, "post:CreateReleaseConf")           //ops
	beego.Router("/api/release/getconf", &controllers.ReleaseController{}, "post:GetReleaseConf")           //ops dev (qa)
	beego.Router("/api/release/status", &controllers.ReleaseController{}, "post:CheckReleaseTaskStatus")    //ops dev qa

}
