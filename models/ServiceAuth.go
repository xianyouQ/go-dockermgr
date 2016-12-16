package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"errors"
)

type ServiceAuth struct {
	Id     int64
	Name   string  `orm:"null;size(100)"`
    Role  *Role  `orm:"rel(fk)"`
    Service *Service `orm:"null;rel(fk)"`
    Users   []*User `orm:"reverse(many)"`
}

func (s *ServiceAuth) TableName() string {
	return beego.AppConfig.String("rbac_serviceauth_table")
}

func init() {
	orm.RegisterModel(new(ServiceAuth))
}

func checkAuth(s *ServiceAuth) error {
	valid := validation.Validation{}
	b, _ := valid.Valid(&s)
	if !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}


func NewServiceAuth(o orm.Ormer,role *Role, service *Service) (int64,error){
    newServiceAuth := ServiceAuth{}
	if service == nil {
		newServiceAuth.Name = role.Name
	} else {
    	newServiceAuth.Name = fmt.Sprintf("%s.%s",service.Code,role.Name)
	    newServiceAuth.Service = service
	}

    newServiceAuth.Role = role

    id, err := o.Insert(&newServiceAuth)
	return id, err
}




func AddUserAuth(o orm.Ormer,user *User,role *Role, service *Service) error {
	var err error
	mServiceAuth := &ServiceAuth{}
    if service == nil {
		err = o.QueryTable(beego.AppConfig.String("rbac_serviceauth_table")).Filter("Role__Id",role.Id).One(mServiceAuth)
	} else {
		err = o.QueryTable(beego.AppConfig.String("rbac_serviceauth_table")).Filter("Role__Id",role.Id).Filter("Service__Id",service.Id).One(mServiceAuth)
	}
	if err != nil {
		return err
	}
    m2m := o.QueryM2M(mServiceAuth,"Users")
    _,err = m2m.Add(user)
	user.ServiceAuths = append(user.ServiceAuths,mServiceAuth)
    return err
}


func DelUserAuth(o orm.Ormer,user *User,role *Role, service *Service) error {
    mServiceAuth := ServiceAuth{Role:role,Service: service}
    m2m := o.QueryM2M(&mServiceAuth,"Users")
    _,err := m2m.Remove(&user)
    return err
}



func GetAuthList(uid int64) ([]*ServiceAuth,error) {
    o := orm.NewOrm()
	var mServiceAuth []*ServiceAuth
	_, err := o.QueryTable(beego.AppConfig.String("rbac_serviceauth_table")).Filter("Users__User__Id", uid).All(&mServiceAuth)
    return mServiceAuth,err
}

