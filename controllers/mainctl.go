package controllers

import (
	"github.com/xianyouQ/go-dockermgr/utils"
	"github.com/astaxie/beego/logs"
)

type MainController struct {
	CommonController
}



func (c *MainController) Get() {
	mesosInfo,err := utils.GetMesosInfo()
	if err != nil {
		logs.GetLogger("Main").Println(err)
	} 
	logs.GetLogger("Main").Println(mesosInfo.CpuTotal)
	c.TplName = "index.html"
}
