package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"

	"encoding/json"
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
	err = models.AddOrUpdateRegistryConf(belongIdc.RegistryConf)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
	}
	err = models.AddOrUpdateIdc(&belongIdc)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
	}
	c.Rsp(true,"success",nil)
}

