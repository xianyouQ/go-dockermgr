package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"

	"encoding/json"
)

type MarathonCfController struct {
	CommonController
}



func (c *MarathonCfController) AddOrUpdateMarathonConf() {
	var err error
	var id int64
	belongIdc := models.IdcConf{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &belongIdc);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	id,err = models.AddOrUpdateMarathonSerConf(belongIdc.MarathonSerConf)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
    return
	}
	if id != 0  {
		belongIdc.MarathonSerConf.Id = int(id)
	}
	err = models.AddOrUpdateIdc(&belongIdc)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
    return
	}
	c.Rsp(true,"success",belongIdc.MarathonSerConf)
}
