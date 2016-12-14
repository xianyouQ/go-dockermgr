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



func QueryUserAuthList(service *Service) ([]*ServiceAuthUser,error) {
    o := orm.NewOrm()
	var auths []*ServiceAuthUser
	_, err := o.QueryTable(beego.AppConfig.String("rbac_serviceauthuser_table")).Filter("ServiceAuth__Service__Id",service.Id).RelatedSel().All(&auths)
	return  auths,err
}


func QueryUserList(username string) ([]*ServiceAuthUser,error) {
    o := orm.NewOrm()
	var auths []*ServiceAuthUser
	_, err := o.QueryTable(beego.AppConfig.String("rbac_serviceauthuser_table")).Filter("ServiceAuth__Service__Id__isnull",true).Filter("User__Username__istartswith",username).RelatedSel("User").All(&auths)
	return  auths,err
}