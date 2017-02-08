package controllers

import (
	"encoding/json"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
)

type ReleaseController struct {
	CommonController
}

func (c *ReleaseController) NewReleaseTask() {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	releaseTask.TaskStatus = models.NotReady
	uinfo := c.Ctx.Input.Session("userinfo")
	releaseUser := uinfo.(models.User)
	releaseTask.ReleaseUser = &releaseUser
	err = models.CreateOrUpdateRelease(o, &releaseTask)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseTask)
}

func (c *ReleaseController) ReviewReleaseTask() {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}

	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	releaseTask.TaskStatus = models.Ready
	uinfo := c.Ctx.Input.Session("userinfo")
	reviewTime := time.Now()
	reviewUser := uinfo.(models.User)
	releaseTask.ReviewUser = &reviewUser
	releaseTask.ReviewTime = reviewTime
	err = models.CreateOrUpdateRelease(o, &releaseTask, "ReviewUser", "TaskStatus", "ReviewTime")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseTask)
}

func (c *ReleaseController) OperationReleaseTask() {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = models.LoadReleaseConf(o, releaseTask.ReleaseConf)
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}

	err = AddTask(&releaseTask)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	releaseTask.TaskStatus = models.Running
	uinfo := c.Ctx.Input.Session("userinfo")
	operationUser := uinfo.(models.User)
	operationTime := time.Now()
	releaseTask.OperationUser = &operationUser
	releaseTask.OperationTime = operationTime
	err = models.CreateOrUpdateRelease(o, &releaseTask, "OperationUser", "TaskStatus", "OperationTime")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}

	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseTask)

}

func (c *ReleaseController) CheckReleaseTaskStatus() {
	var err error
	var mReleaseRoutine *ReleaseRoutine
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	mReleaseRoutine, err = CheckTaskStatus(&releaseTask)
	if err != nil {
		o := orm.NewOrm()
		err = o.Read(&releaseTask)
		if err != nil {
			c.Rsp(false, err.Error(), nil)
			return
		}
	} else {
		var jsonByte []byte
		jsonByte, err = json.Marshal(mReleaseRoutine.ContainerResChan)
		if err != nil {
			c.Rsp(false, err.Error(), nil)
			return
		}
		releaseTask.ReleaseResult = string(jsonByte)
	}
	c.Rsp(true, "success", releaseTask)
}

func (c *ReleaseController) CancelReleaseTask() {
	var err error
	releaseTask := models.ReleaseTask{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}

	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	releaseTask.TaskStatus = models.Cancel
	err = models.CreateOrUpdateRelease(o, &releaseTask, "CancelUser", "TaskStatus")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseTask)
}

func (c *ReleaseController) CreateReleaseConf() {
	var err error
	releaseConf := models.ReleaseConf{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseConf); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	releaseConf.Id = 0
	err = models.CreateReleaseConf(o, &releaseConf)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("RegistryCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("RegistryCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseConf)
}

func (c *ReleaseController) QueryReleaseTasks() {
	var err error
	var tasks []*models.ReleaseTask
	querySerivce := models.Service{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &querySerivce); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	tasks, err = models.QueryRelease(o, &querySerivce)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", tasks)
}
func (c *ReleaseController) GetReleaseConf() {
	var err error
	var conf models.ReleaseConf
	querySerivce := models.Service{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &querySerivce); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	conf, err = models.QueryReleaseConf(o, &querySerivce)
	if err == orm.ErrNoRows {
		c.Rsp(false, "release conf not set", nil)
		return
	}
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", conf)
}
