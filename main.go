package main

import (
	"flag"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xianyouQ/go-dockermgr/models"
	_ "github.com/xianyouQ/go-dockermgr/routers"
)

func init() {
	logs.SetLogger("console")
	err := orm.RegisterDataBase("default", "mysql", "testfordjango:123456@/dockermgr?charset=utf8")
	if err != nil {
		panic(err.Error())
	} else {
		logs.GetLogger("init").Println("registry database successfully")
	}

}

func main() {

	init := flag.Bool("init", false, "init db and data")
	flag.Parse()
	if *init == true {
		initData()
		return
	}
	dbDebug, _ := beego.AppConfig.Bool("db_debug")
	orm.Debug = dbDebug
	beego.Run()
}

func initData() {
	var err error
	err = orm.RunSyncdb("default", true, true)
	if err != nil {
		panic(err.Error())
	}
	o := orm.NewOrm()

	adminRole := &models.Role{Name: "SYSTEM", Status: true}
	baseRole := &models.Role{Name: "BASE", Status: true}
	opsRole := &models.Role{Name: "OPS", Status: true, NeedAddAuth: true}
	qaRole := &models.Role{Name: "QA", Status: true, NeedAddAuth: true}
	devRole := &models.Role{Name: "DEV", Status: true, NeedAddAuth: true}
	err = models.AddOrUpdateRole(o, adminRole)
	if err != nil {
		panic(err.Error())
	}
	err = models.AddOrUpdateRole(o, baseRole)
	if err != nil {
		panic(err.Error())
	}
	err = models.AddOrUpdateRole(o, opsRole)
	if err != nil {
		panic(err.Error())
	}
	err = models.AddOrUpdateRole(o, qaRole)
	if err != nil {
		panic(err.Error())
	}
	err = models.AddOrUpdateRole(o, devRole)
	if err != nil {
		panic(err.Error())
	}
	Nodes := []models.Node{
		{Desc: "user managerment", Url: "/api/auth/user"},
		{Desc: "get users list", Url: "/api/auth/user/get"},
		{Desc: "passwd reset", Url: "/api/auth/passwd"},
		{Desc: "change passwd", Url: "/api/auth/passwd/change"},
		{Desc: "login/logout", Url: "/api/auth/sign"},
		{Desc: "get role list", Url: "/api/auth/get"},
		{Desc: "auth allocation for service", Url: "/api/auth/new"},
		{Desc: "add or update role", Url: "/api/auth/post"},
		{Desc: "get service's authlist", Url: "/api/auth/auths"},
		{Desc: "add/delete node to role", Url: "/api/authnode/post"},
		{Desc: "add or update node", Url: "/api/node/post"},
		{Desc: "get node list", Url: "/api/node/get"},
		{Desc: "add or update idc Conf", Url: "/api/idc"},
		{Desc: "get idc list", Url: "/api/idc/get"},
		{Desc: "add or update marathon conf", Url: "/api/marathon/conf"},
		{Desc: "add or update registry conf", Url: "/api/registry/conf"},
		{Desc: "add cidr for idc", Url: "/api/Cidr/Add"},
		{Desc: "add or update service", Url: "/api/service/Add"},
		{Desc: "delete service", Url: "/api/service/Delete"},
		{Desc: "get service split count", Url: "/api/service/count"},
		{Desc: "get service list", Url: "/api/service/get"},
		{Desc: "docker dashboard", Url: "/api/docker/dashboard"},
		{Desc: "container scale", Url: "/api/docker/scale"},
		{Desc: "get container list", Url: "/api/docker/list"},
		{Desc: "get release task list", Url: "/api/release/task"},
		{Desc: "new release task", Url: "/api/release/newtask"},
		{Desc: "release task review", Url: "/api/release/review"},
		{Desc: "operate task review", Url: "/api/release/operate"},
		{Desc: "create release task conf", Url: "/api/release/conf"},
		{Desc: "get release task conf", Url: "/api/release/getconf"},
		{Desc: "get release task status", Url: "/api/release/status"},
	}
	_, err = o.InsertMulti(10, Nodes)
	if err != nil {
		panic(err.Error())
	}
	adminUserName := beego.AppConfig.String("rbac_admin_user")
	if adminUserName == "" {
		panic("adminUser not Set in App.conf")
	}
	defaultPasswd := beego.AppConfig.String("rbac_auth_defaultpasswd")
	if defaultPasswd == "" {
		panic("defaultPasswd not Set in App.conf")
	}
	adminUser := &models.User{Username: adminUserName, Password: defaultPasswd, Repassword: defaultPasswd}
	_, err = models.AddUser(o, adminUser)
	if err != nil {
		panic(err.Error())
	}
	users := make([]*models.User, 0, 1)
	users = append(users, adminUser)
	err = models.AddUserAuth(o, users, adminRole, nil)
	if err != nil {
		panic(err.Error())
	}
}
