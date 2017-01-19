package controllers

import (
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

type ReleaseController struct {
	CommonController
}


func (c *ReleaseController) NewReleaseTask () {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	err = models.CreateOrUpdateRelease(o,&releaseTask)
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s",err.Error())
	}
	c.Rsp(true,"success",releaseTask)
}

func (c *ReleaseController) ReviewReleaseTask () {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}

	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	releaseTask.TaskStatus = models.Ready
	err = models.CreateOrUpdateRelease(o,&releaseTask,"ReviewUser","TaskStatus")
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s",err.Error())
	}
	c.Rsp(true,"success",releaseTask)
}

func (c *ReleaseController) OperationReleaseTask () {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	releaseTask.TaskStatus = models.Running
	err = models.CreateOrUpdateRelease(o,&releaseTask,"OperationUser","TaskStatus")
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s",err.Error())
	}
	c.Rsp(true,"success",releaseTask)

}


func (c *ReleaseController) CheckReleaseTaskStatus () {

}

func (c *ReleaseController) CancelReleaseTask () {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}

	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	releaseTask.TaskStatus = models.Cancel
	err = models.CreateOrUpdateRelease(o,&releaseTask,"CancelUser","TaskStatus")
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s",err.Error())
	}
	c.Rsp(true,"success",releaseTask)
}

func (c *ReleaseController) CreateOrUpdateReleaseConf() {
	var err error
	releaseConf := models.ReleaseConf{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseConf);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	err = models.CreateOrUpdateReleaseConf(o,&releaseConf)
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s",err.Error())
	}
	c.Rsp(true,"success",releaseConf)
}

