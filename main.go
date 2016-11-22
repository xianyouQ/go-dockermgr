package main

import (
	_ "github.com/xianyouQ/go-dockermgr/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"

)

func init() {
	orm.RegisterDataBase("default", "mysql", "testfordjango:123456@/dockermgr?charset=utf8")
	logs.SetLogger("console")
}
func main() {
	orm.RunCommand()
	beego.Run()
}

