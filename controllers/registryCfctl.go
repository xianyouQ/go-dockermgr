package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

type RegistryCfController struct {
	CommonController
}



func (c *RegistryCfController) AddOrUpdateRegistryConf() {
	var err error
	belongIdc := models.IdcConf{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &belongIdc);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	_,err = models.AddOrUpdateRegistryConf(o,belongIdc.RegistryConf)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("registryCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}

	err = models.AddOrUpdateIdc(o,&belongIdc)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("registryCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("registryCtl").Printf("commit error:%s",err.Error())
	}
	c.Rsp(true,"success",belongIdc.RegistryConf)
}

