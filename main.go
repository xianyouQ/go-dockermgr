package main

import (
	_ "github.com/xianyouQ/go-dockermgr/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
	"flag"
	"github.com/xianyouQ/go-dockermgr/models"
)

func init() {
	logs.SetLogger("console")
	err := orm.RegisterDataBase("default", "mysql", "testfordjango:123456@/dockermgr?charset=utf8")
	if err != nil {
		logs.Critical(err.Error())
	} else {
		logs.GetLogger("init").Println("registry database successfully")
	}
	
}

func main() {
	var err error
	init := flag.Bool("init", false, "init db and data")
	flag.Parse()
	if *init == true {
		err = orm.RunSyncdb("default", true, true)
		if err != nil {
			logs.Critical(err.Error())
		}
		o := orm.NewOrm()
		adminUserName := beego.AppConfig.String("rbac_admin_user")
		if adminUserName == "" {
			logs.Critical("adminUser not Set in App.conf")
		}
		defaultPasswd := beego.AppConfig.String("rbac_auth_defaultpasswd")
		if defaultPasswd == "" {
			logs.Critical("defaultPasswd not Set in App.conf")
		}		
		adminUser := &models.User{Username:adminUserName,Password:defaultPasswd,Repassword:defaultPasswd}
		_,err = models.AddUser(o,adminUser)
		if err !=nil {
			logs.Critical(err.Error())
		}
		adminRole := &models.Role{Name:"SYSTEM",Status:true}
		baseRole := &models.Role{Name:"BASE",Status:true}
		_,err = models.AddOrUpdateRole(o,adminRole)
		if err !=nil {
			logs.Critical(err.Error())
		}
		_,err =models.AddOrUpdateRole(o,baseRole)
		if err !=nil {
			logs.Critical(err.Error())
		}
	}
	dbDebug,_ := beego.AppConfig.Bool("db_debug")
	orm.Debug = dbDebug
	beego.Run()
}

