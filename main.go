package main

import (
	_ "github.com/xianyouQ/go-dockermgr/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
)

func init() {
	logs.SetLogger("console")
	err := orm.RegisterDataBase("default", "mysql", "testfordjango:123456@/dockermgr?charset=utf8")
	if err != nil {
		logs.GetLogger("init").Println(err)
	} else {
		logs.GetLogger("init").Println("registry database successfully")
	}
	//orm.RegisterDataBase("sqlite3","sqlite3","test.db")
	
}
func main() {
	//orm.RunCommand()
	orm.Debug = true
	beego.Run()
}

