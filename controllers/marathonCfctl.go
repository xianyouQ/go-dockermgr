package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"

	"encoding/json"
)

type IDCController struct {
	CommonController
}



func (c *IDCController) RequestIdcs() {
	Idcs,err := models.GetIdcs()
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		return
	}
	c.Rsp(true,"success",Idcs)
}

func (c *IDCController) AddIdc() {
	newIdc := models.IdcConf{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &newIdc);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	err := models.AddIdc(&newIdc)
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		return
	}
	c.Rsp(true,"success",newIdc)
}

func (c *IDCController) UpdateIdc() {
	newIdc := models.IdcConf{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &newIdc);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	err := models.UpdateIdc(&newIdc)
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		return
	}
	c.Rsp(true,"success",newIdc)
}