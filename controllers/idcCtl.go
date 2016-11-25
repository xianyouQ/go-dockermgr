package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"
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
	IdcName := c.GetString("IdcName") 
	IdcCode := c.GetString("IdcCode")
	if IdcName == "" || IdcCode == "" {
		c.Rsp(false,"IdcName or IdcCode is empty",nil)
		return
	} 
	err := models.AddIdc(IdcName,IdcCode)
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		return
	}
	c.Rsp(true,"success",nil)
}
