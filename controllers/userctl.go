package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/auth"
	m "github.com/xianyouQ/go-dockermgr/models"
	"github.com/xianyouQ/go-dockermgr/utils"
)

type UserController struct {
	CommonController
}

func (this *UserController) AddUser() {
	u := m.User{}
	var err error
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error(), nil)
		return
	}
	u.Password = beego.AppConfig.String("rbac_auth_defaultpasswd")
	u.Repassword = beego.AppConfig.String("rbac_auth_defaultpasswd")
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		this.Rsp(false, err.Error(), nil)
		return
	}
	_, err = m.AddUser(o, &u)
	if err == nil {
		this.Rsp(true, "Success", u)
		err = o.Commit()
		if err != nil {
			logs.GetLogger("userCtl").Printf("commit error:%s", err.Error())
		}
		return
	} else {
		this.Rsp(false, err.Error(), nil)
		err = o.Rollback()
		if err != nil {
			logs.GetLogger("userCtl").Printf("rollback error:%s", err.Error())
		}
		return
	}

}

/*
func (this *UserController) UpdateUser() {
	u := m.User{}
	if err := this.ParseForm(&u); err != nil {
		//handle error
		this.Rsp(false, err.Error(), nil)
		return
	}
	o := orm.NewOrm()
	id, err := m.UpdateUser(o, &u)
	if err == nil && id > 0 {
		this.Rsp(true, "Success", nil)
		return
	} else {
		this.Rsp(false, err.Error(), nil)
		return
	}

}
*/

func (this *UserController) DelUser() {
	Id, _ := this.GetInt64("Id")
	o := orm.NewOrm()
	status, err := m.DelUserById(o, Id)
	if err == nil && status > 0 {
		this.Rsp(true, "Success", nil)
		return
	} else {
		this.Rsp(false, err.Error(), nil)
		return
	}
}

//登录
func (this *UserController) Login() {
	uinfo := this.GetSession("userinfo")
	data := make(map[string]interface{})
	if uinfo != nil {
		data["Username"] = uinfo.(m.User).Username
		data["auth"] = this.GetSession("accesslist").([]*m.ServiceAuthUser)
		this.Rsp(false, "不可重复登陆", data)
		return
	}
	u := m.User{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error(), nil)
		return
	}

	user, err := auth.CheckLogin(u.Username, u.Password)
	if err != nil {
		this.Rsp(false, err.Error(), nil)
		return
	}
	this.SetSession("userinfo", user)
	data["Username"] = user.Username
	var accesslist []*m.ServiceAuthUser
	accesslist, err = auth.AccessList(user.Id)
	if err != nil {
		logs.GetLogger("userCtl").Printf("get auth fail,detail:%s", err.Error())
	} else {
		data["auth"] = accesslist
		this.SetSession("accesslist", accesslist)
	}
	this.Rsp(true, "登陆成功", data)
	return

}

//退出
func (this *UserController) Logout() {
	this.DelSession("userinfo")
	this.Rsp(true, "退出成功", nil)
}

//修改密码
func (this *UserController) Changepwd() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Rsp(false, "请先登录", nil)
	}
	u := m.User{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error(), nil)
		return
	}

	if u.Password != u.Repassword {
		this.Rsp(false, "两次输入密码不一致", nil)
	}
	user, err := auth.CheckLogin(userinfo.(m.User).Username, u.OldPassword)
	if err == nil {
		u.Id = user.Id
		u.Password = utils.Pwdhash(u.Password)
		o := orm.NewOrm()
		id, err := m.UpdateUser(o, &u, "Password")
		if err == nil && id > 0 {
			this.Rsp(true, "密码修改成功", nil)
			return
		} else {
			this.Rsp(false, err.Error(), nil)
			return
		}
	}
	this.Rsp(false, "密码有误", nil)

}

func (this *UserController) ResetPwd() {
	var err error
	u := m.User{}
	if err = json.Unmarshal(this.Ctx.Input.RequestBody, &u); err != nil {
		//handle error
		this.Rsp(false, err.Error(), nil)
		return
	}
	u.Password = beego.AppConfig.String("rbac_auth_defaultpasswd")
	o := orm.NewOrm()
	newU := m.User{}
	newU.Id = u.Id
	newU.Password = u.Password
	_, err = m.UpdateUser(o, &newU)
	if err == nil {
		this.Rsp(true, "密码重置成功", nil)
		return
	} else {
		this.Rsp(false, err.Error(), nil)
		return
	}
}

func (this *UserController) GetUserList() {
	//pageId,err := this.GetInt("pageId")
	UserName := this.GetString("username")
	users, err := m.QueryUserAuthListByUser(UserName)
	if err != nil {
		this.Rsp(false, err.Error(), nil)
		return
	}
	this.Rsp(true, "success", users)
}
