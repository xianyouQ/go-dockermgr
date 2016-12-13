package models

import (
	"github.com/astaxie/beego/orm"
    "github.com/astaxie/beego"
)

type ServiceAuthUser struct {
    Id    int64
    ServiceAuth *ServiceAuth `orm:"rel(fk)"`
    User *User `orm:"rel(fk)"`
}


func init() {
	orm.RegisterModel(new(ServiceAuthUser))
}

func (s *ServiceAuthUser) TableName() string {
	return beego.AppConfig.String("rbac_serviceauthuser_table")
}



func QueryUserAuthList() ([]*ServiceAuthUser,error) {
    o := orm.NewOrm()
	var auths []*ServiceAuthUser
	_, err := o.QueryTable(beego.AppConfig.String("rbac_serviceauthuser_table")).RelatedSel().All(&auths)
	return  auths,err
}
