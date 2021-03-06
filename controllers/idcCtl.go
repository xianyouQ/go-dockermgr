package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
)

type IDCController struct {
	CommonController
}

func (c *IDCController) RequestIdcs() {
	Idcs, err := models.GetIdcs()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}

	c.Rsp(true, "success", Idcs)
}

func (c *IDCController) AddOrUpdateIdc() {
	var err error
	newIdc := models.IdcConf{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newIdc); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
	}
	err = models.AddOrUpdateIdc(o, &newIdc)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("AuthCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("AuthCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", newIdc)
}
func (c *IDCController) DelIdc() {
	var err error
	var idcId int
	idcId, err = c.GetInt("idcId")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	delIdc := &models.IdcConf{}
	delIdc.Id = idcId
	o := orm.NewOrm()
	err = models.DelIdc(o, delIdc)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", nil)
}
