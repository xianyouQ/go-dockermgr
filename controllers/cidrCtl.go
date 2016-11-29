package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"
	"encoding/json"
)

type CidrController struct {
	CommonController
}





func (c *CidrController) AddCidr() {
	newCidr := models.Cidr{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &newCidr);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	err := models.AddCidr(&newCidr)
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		return
	}
	newCidr.BelongIdc = nil
	c.Rsp(true,"success",newCidr)
}

