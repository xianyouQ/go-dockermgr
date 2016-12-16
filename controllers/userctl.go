package controllers

import (
	m "github.com/xianyouQ/go-dockermgr/models"
	"github.com/xianyouQ/go-dockermgr/auth"
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
)


type UserController struct {
	CommonController
}



func (this *UserController) AddUser() {
	u := m.User{}
	var err error
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error(),nil)
		return
	}
	u.Password = beego.AppConfig.String("rbac_auth_defaultpasswd")
	u.Repassword = beego.AppConfig.String("rbac_auth_defaultpasswd")
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		this.Rsp(false, err.Error(),nil)
		return
	}
	_, err = m.AddUser(o,&u)
	if err == nil {
		this.Rsp(true, "Success",u)
		err = o.Commit()
		if err != nil {
			logs.GetLogger("userCtl").Printf("commit error:%s",err.Error())
		}
		return
	} else {
		this.Rsp(false, err.Error(),nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("userCtl").Printf("rollback error:%s",err.Error())
		}
		return
	}

}

func (this *UserController) UpdateUser() {
	u := m.User{}
	if err := this.ParseForm(&u); err != nil {
		//handle error
		this.Rsp(false, err.Error(),nil)
		return
	}
	o := orm.NewOrm()
	id, err := m.UpdateUser(o,&u)
	if err == nil && id > 0 {
		this.Rsp(true, "Success",nil)
		return
	} else {
		this.Rsp(false, err.Error(),nil)
		return
	}

}

func (this *UserController) DelUser() {
	Id, _ := this.GetInt64("Id")
	o := orm.NewOrm()
	status, err := m.DelUserById(o,Id)
	if err == nil && status > 0 {
		this.Rsp(true, "Success",nil)
		return
	} else {
		this.Rsp(false, err.Error(),nil)
		return
	}
}

//登录
func (this *UserController) Login() {
	uinfo := this.Ctx.Input.Session("userinfo")
	data := make(map[string]string)
	if uinfo != nil {
		data["Username"] = uinfo.(m.User).Username
		//datajson,_ := json.Marshal(data)
		this.Rsp(false,"不可重复登陆", data);
		return
	}
	u := m.User{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error(),nil)
		return
	}

	user, err := auth.CheckLogin(u.Username, u.Password)
	if err == nil {
		this.SetSession("userinfo", user)
		//accesslist, _ := auth.GetAccessList(user.Id)
		//this.SetSession("accesslist", accesslist)
		data["Username"] = user.Username
		//datajson,_ := json.Marshal(data)
		this.Rsp(true, "登陆成功",data)
		return
	} else {
		this.Rsp(false, err.Error(),nil)
		return
	}
}

//退出
func (this *UserController) Logout() {
	this.DelSession("userinfo")
	this.Rsp(true, "退出成功",nil)
}

//修改密码
func (this *UserController) Changepwd() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Rsp(false,"请先登录",nil)
	}
	oldpassword := this.GetString("oldpassword")
	newpassword := this.GetString("newpassword")
	repeatpassword := this.GetString("repeatpassword")
	if newpassword != repeatpassword {
		this.Rsp(false, "两次输入密码不一致",nil)
	}
	user, err := auth.CheckLogin(userinfo.(m.User).Username, oldpassword)
	if err == nil {
		var u m.User
		u.Id = user.Id
		u.Password = newpassword
		o := orm.NewOrm()
		id, err := m.UpdateUser(o,&u)
		if err == nil && id > 0 {
			this.Rsp(true, "密码修改成功",nil)
			return
		} else {
			this.Rsp(false, err.Error(),nil)
			return
		}
	}
	this.Rsp(false, "密码有误",nil)

}

/*
func (this *UserController) GetUserList() {
	pageId,err := this.GetInt("pageId")
	UserName := this.GetString("username")
	users,err := m.Getuserlist(UserName,int64(pageId),10, "username")
	if err != nil {
		this.Rsp(false,err.Error(),nil)
		return
	}
	this.Rsp(true,"success",users)
}
*/
func (this *UserController) GetUserList() {
	//pageId,err := this.GetInt("pageId")
	UserName := this.GetString("username")
	users,err := m.QueryUserAuthListByUser(UserName)
	if err != nil {
		this.Rsp(false,err.Error(),nil)
		return
	}
	this.Rsp(true,"success",users)
}