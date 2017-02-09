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
	serviceId, _ := c.GetInt("serviceId")
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseTask); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	if releaseTask.Service.Id != serviceId {
		c.Ctx.Output.SetStatus(403)
		c.Rsp(false, "permission denied", nil)
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
	_, err = models.CreateOrUpdateRelease(o, &releaseTask)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("releaseCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("releaseCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseTask)
}

func (c *ReleaseController) ReviewReleaseTask() {
	var err error
	var taskId, serviceId int
	var num int64
	releaseTask := models.ReleaseTask{}
	service := models.Service{}
	serviceId, _ = c.GetInt("serviceId")
	taskId, err = c.GetInt("taskId")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	service.Id = serviceId
	releaseTask.Id = taskId
	releaseTask.Service = &service
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	uinfo := c.Ctx.Input.Session("userinfo")
	reviewTime := time.Now()
	reviewUser := uinfo.(models.User)
	params := make(orm.Params)
	params["TaskStatus"] = models.Ready
	releaseTask.TaskStatus = models.Ready
	params["ReviewUser"] = reviewUser.Id
	releaseTask.ReviewUser = &reviewUser
	params["ReviewTime"] = reviewTime
	releaseTask.ReviewTime = reviewTime
	num, err = models.UpdateRelease(o, &releaseTask, models.NotReady, params)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("releaseCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}

	if num == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Rsp(false, "relaseTask info didn't match", nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("releaseCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("releaseCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseTask)
}

func (c *ReleaseController) OperationReleaseTask() {
	var err error
	var taskId, serviceId int
	releaseTask := models.ReleaseTask{}
	serviceId, _ = c.GetInt("serviceId")
	taskId, err = c.GetInt("taskId")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	releaseTask.Id = taskId

	o := orm.NewOrm()
	err = models.LoadReleaseConf(o, &releaseTask)
	if releaseTask.Service.Id != serviceId {
		c.Ctx.Output.SetStatus(403)
		c.Rsp(false, "permission denied", nil)
		return
	}
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
			logs.GetLogger("releaseCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	releaseTask.TaskStatus = models.Running
	uinfo := c.Ctx.Input.Session("userinfo")
	operationUser := uinfo.(models.User)
	operationTime := time.Now()
	releaseTask.OperationUser = &operationUser
	releaseTask.OperationTime = operationTime
	_, err = models.CreateOrUpdateRelease(o, &releaseTask, "OperationUser", "TaskStatus", "OperationTime")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("releaseCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}

	err = o.Commit()
	if err != nil {
		logs.GetLogger("releaseCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseTask)

}

func (c *ReleaseController) CheckReleaseTaskStatus() {
	var err error
	var mReleaseRoutine *ReleaseRoutine
	var taskId, serviceId int
	releaseTask := models.ReleaseTask{}
	serviceId, _ = c.GetInt("serviceId")
	taskId, err = c.GetInt("taskId")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	releaseTask.Id = taskId
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
	if releaseTask.Service.Id != serviceId {
		c.Ctx.Output.SetStatus(403)
		c.Rsp(false, "permission denied", nil)
		return
	}
	c.Rsp(true, "success", releaseTask)
}

func (c *ReleaseController) CreateReleaseConf() {
	var err error
	var serviceId int
	releaseConf := models.ReleaseConf{}
	serviceId, _ = c.GetInt("serviceId")
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &releaseConf); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	if releaseConf.Service.Id != serviceId {
		c.Ctx.Output.SetStatus(403)
		c.Rsp(false, "permission denied", nil)
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
			logs.GetLogger("releaseCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("releaseCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", releaseConf)
}

func (c *ReleaseController) QueryReleaseTasks() {
	var err error
	var tasks []*models.ReleaseTask
	serviceId, _ := c.GetInt("serviceId")
	querySerivce := models.Service{}
	querySerivce.Id = serviceId
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
	var serviceId int
	serviceId, _ = c.GetInt("serviceId")
	querySerivce := models.Service{}
	querySerivce.Id = serviceId
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
