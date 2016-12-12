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


func NewServiceAuth(role *Role, service *Service) (int64,error){
    o := orm.NewOrm()
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




func AddUserAuth(uid int64,role *Role, service *Service) error {
    o := orm.NewOrm()
    mServiceAuth := ServiceAuth{Role:role,Service: service}
    m2m := o.QueryM2M(&mServiceAuth,"Users")
    addUser := User{Id:uid}
    _,err := m2m.Add(&addUser)
    return err
}


func DelUserAuth(uid int64,role *Role, service *Service) error {
    o := orm.NewOrm()
    mServiceAuth := ServiceAuth{Role:role,Service: service}
    m2m := o.QueryM2M(&mServiceAuth,"Users")
    addUser := User{Id:uid}
    _,err := m2m.Remove(&addUser)
    return err
}



func GetAuthList(uid int64) ([]*ServiceAuth,error) {
    o := orm.NewOrm()
	var mServiceAuth []*ServiceAuth
	_, err := o.QueryTable(beego.AppConfig.String("rbac_serviceauth_table")).Filter("Users__User__Id", uid).All(&mServiceAuth)
    return mServiceAuth,err
}

