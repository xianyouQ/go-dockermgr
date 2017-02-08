package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
)

type ServiceController struct {
	CommonController
}

func (c *ServiceController) AddOrUpdateService() {
	var err error
	newService := models.Service{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newService); err != nil {
		//handle error
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	err = models.AddOrUpdateService(o, &newService)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("serviceCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("serviceCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", newService)
}

func (c *ServiceController) DelService() {
	var err error
	oldService := models.Service{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &oldService); err != nil {
		//handle error
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	err = models.DelService(o, &oldService)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("serviceCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("serviceCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", nil)
}

func (c *ServiceController) GetSeparateCount() {
	count, err := beego.AppConfig.Int("service_separate_count")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", count)
}

func (c *ServiceController) GetService() {
	Services, err := models.GetServices()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", Services)
}
