package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
)

type AuthController struct {
	CommonController
}

func (c *AuthController) AddOrUpdateRole() {
	var err error
	//var id int64
	newRole := models.Role{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newRole); err != nil {
		//handle error
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	err = models.AddOrUpdateRole(o, &newRole)
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
	c.Rsp(true, "success", newRole)
}

func (c *AuthController) DelRole() {
	var err error
	oldRole := models.Role{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &oldRole); err != nil {
		//handle error
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	err = models.DelRole(o, &oldRole)
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
	c.Rsp(true, "success", nil)
}

func (c *AuthController) GetRole() {
	Roles, err := models.GetRoleNodes()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", Roles)
}

func (c *AuthController) GetNodes() {
	Nodes, err := models.GetNodes()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", Nodes)
}

func (c *AuthController) AddOrUpdateNode() {
	var err error
	//var id int64
	newNode := models.Node{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newNode); err != nil {
		//handle error
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	err = models.AddOrUpdateNode(o, &newNode)
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
	c.Rsp(true, "success", newNode)
}

func (c *AuthController) UpdateRoleNode() {
	var err error
	oldRole := models.Role{}
	activeNodes := make([]*models.Node, 0, 5)
	inActiveNodes := make([]*models.Node, 0, 5)
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &oldRole); err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	for _, node := range oldRole.Nodes {
		if node.Active == true {
			activeNodes = append(activeNodes, node)
		} else {
			inActiveNodes = append(inActiveNodes, node)
		}
	}
	_, err = models.AddRoleNode(o, &oldRole, activeNodes)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("AuthCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}

	_, err = models.DelRoleNode(o, &oldRole, inActiveNodes)
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
	c.Rsp(true, "success", nil)
}
func (c *AuthController) GetUserAuthList() {
	serviceId, _ := c.GetInt("serviceId")
	newService := &models.Service{Id: serviceId}
	auths, err := models.QueryUserAuthList(newService)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", auths)

}

type UserAuthForm struct {
	Users   []*models.User
	Service *models.Service
	Role    *models.Role
}

func (c *AuthController) AddUserAuth() {
	var err error
	serviceId, _ := c.GetInt("serviceId")
	mUserAuthForm := UserAuthForm{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &mUserAuthForm); err != nil {
		//handle error
		c.Rsp(false, err.Error(), nil)
		return
	}
	if serviceId != mUserAuthForm.Service.Id {
		c.Ctx.Output.SetStatus(403)
		c.Rsp(false, "permission denied", nil)
		return
	}
	o := orm.NewOrm()
	err = models.AddUserAuth(o, mUserAuthForm.Users, mUserAuthForm.Role, mUserAuthForm.Service)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	c.Rsp(true, "success", mUserAuthForm.Users)
}
