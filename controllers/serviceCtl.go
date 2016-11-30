package controllers

import (
	"github.com/astaxie/beego"
	"github.com/xianyouQ/go-dockermgr/models"

	"encoding/json"
)

type ServiceController struct {
	CommonController
}

func (c *ServiceController) AddService() {
	var err error
	var id int64
	newService := models.Service{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newService); err != nil {
		//handle error
		c.Rsp(false, err.Error(),nil)
		return
	}
	id,err = models.AddService(&newService)
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	newService.Id = int(id)
	c.Rsp(true, "success",newService)
}

func (c *ServiceController) DelService() {
	var err error
	oldService := models.Service{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &oldService); err != nil {
		//handle error
		c.Rsp(false, err.Error(),nil)
		return
	}
	err = models.DelService(&oldService)
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return		
	}
	c.Rsp(true,"success",nil)
}

func (c *ServiceController) GetSeparateCount() {
	count,err := beego.AppConfig.Int("service_separate_count")
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	c.Rsp(true,"success",count)
}

func (c *ServiceController) GetService() {
	Services,err := models.QueryService()
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		return
	}
	c.Rsp(true,"success",Services)
}