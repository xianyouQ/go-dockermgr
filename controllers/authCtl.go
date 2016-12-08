package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"

	"encoding/json"
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
		c.Rsp(false, err.Error(),nil)
		return
	}
	_,err = models.AddOrUpdateRole(&newRole)
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	//newRole.Id = int(id)
	c.Rsp(true, "success",newRole)
}

func (c *AuthController) DelRole() {
	var err error
	oldRole := models.Role{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &oldRole); err != nil {
		//handle error
		c.Rsp(false, err.Error(),nil)
		return
	}
	err = models.DelRole(&oldRole)
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return		
	}
	c.Rsp(true,"success",nil)
}



func (c *AuthController) GetRole() {
	Roles,err := models.GetRoleNodes()
	if err != nil {
		c.Rsp(false,err.Error(),nil)
		return	
	}
	c.Rsp(true,"success",Roles)
}

func (c *AuthController) GetNodes() {
	Nodes,err := models.GetNodes()
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
		return
	}
	c.Rsp(true,"success",Nodes)
}

func (c *AuthController) AddOrUpdateNode() {
	var err error
	//var id int64
	newNode := models.Node{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newNode); err != nil {
		//handle error
		c.Rsp(false, err.Error(),nil)
		return
	}
	_,err = models.AddOrUpdateNode(&newNode)
	if err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	//newRole.Id = int(id)
	c.Rsp(true, "success",newNode)
}