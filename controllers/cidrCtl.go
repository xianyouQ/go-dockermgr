package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
)

type CidrController struct {
	CommonController
}

func (c *CidrController) AddCidr() {
	newCidr := models.Cidr{}
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newCidr); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
	}
	err = models.AddCidr(o, &newCidr)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("AuthCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	newCidr.BelongIdc = nil
	err = o.Commit()
	if err != nil {
		logs.GetLogger("AuthCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", newCidr)
}
func (c *CidrController) DelCidr() {
	var err error
	var cidrId int
	cidrId, err = c.GetInt("cidrId")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	delCidr := &models.Cidr{}
	delCidr.Id = cidrId
	o := orm.NewOrm()
	err = models.DelCidr(o, delCidr)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", nil)
}
