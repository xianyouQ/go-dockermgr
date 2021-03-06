package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
)

type MarathonCfController struct {
	CommonController
}

func (c *MarathonCfController) AddOrUpdateMarathonConf() {
	var err error
	belongIdc := models.IdcConf{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &belongIdc); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
	}
	err = models.AddOrUpdateMarathonSerConf(o, belongIdc.MarathonSerConf)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("AuthCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = models.AddOrUpdateIdc(o, &belongIdc, "MarathonSerConf")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		if err != nil {
			logs.GetLogger("AuthCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("AuthCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", belongIdc.MarathonSerConf)
}
