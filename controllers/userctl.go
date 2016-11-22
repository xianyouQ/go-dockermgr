package controllers

import (
	m "github.com/xianyouQ/go-dockermgr/auth/models"
	"github.com/xianyouQ/go-dockermgr/auth"
	//"github.com/astaxie/beego/logs"
	"encoding/json"
)


type UserController struct {
	CommonController
}



func (this *UserController) AddUser() {
	u := m.User{}
	
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	id, err := m.AddUser(&u)
	
	if err == nil && id > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}

}

func (this *UserController) UpdateUser() {
	u := m.User{}
	if err := this.ParseForm(&u); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	id, err := m.UpdateUser(&u)
	if err == nil && id > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}

}

func (this *UserController) DelUser() {
	Id, _ := this.GetInt64("Id")
	status, err := m.DelUserById(Id)
	if err == nil && status > 0 {
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}

//登录
func (this *UserController) Login() {
	uinfo := this.Ctx.Input.Session("userinfo")
	if uinfo != nil {
		this.Rsp(false,"不可重复登陆");
		return
	}
	u := m.User{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}

	user, err := auth.CheckLogin(u.Username, u.Password)
	if err == nil {
		this.SetSession("userinfo", user)
		//accesslist, _ := auth.GetAccessList(user.Id)
		//this.SetSession("accesslist", accesslist)
		this.Rsp(true, u.Username)
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}

//退出
func (this *UserController) Logout() {
	this.DelSession("userinfo")
	this.Rsp(true, "退出成功")
}

//修改密码
func (this *UserController) Changepwd() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Rsp(false,"请先登录")
	}
	oldpassword := this.GetString("oldpassword")
	newpassword := this.GetString("newpassword")
	repeatpassword := this.GetString("repeatpassword")
	if newpassword != repeatpassword {
		this.Rsp(false, "两次输入密码不一致")
	}
	user, err := auth.CheckLogin(userinfo.(m.User).Username, oldpassword)
	if err == nil {
		var u m.User
		u.Id = user.Id
		u.Password = newpassword
		id, err := m.UpdateUser(&u)
		if err == nil && id > 0 {
			this.Rsp(true, "密码修改成功")
			return
		} else {
			this.Rsp(false, err.Error())
			return
		}
	}
	this.Rsp(false, "密码有误")

}