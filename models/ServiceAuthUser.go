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


func QueryUserAuthListByUser(username string) ([]*User,error) {
    o := orm.NewOrm()
	var err error
	var auths []*ServiceAuthUser
	users := make([]*User,0,5)
	if username == "" {
		_, err = o.QueryTable(beego.AppConfig.String("rbac_serviceauthuser_table")).Filter("ServiceAuth__Service__Id__isnull",true).RelatedSel().All(&auths)
	} else {
	_, err = o.QueryTable(beego.AppConfig.String("rbac_serviceauthuser_table")).Filter("ServiceAuth__Service__Id__isnull",true).Filter("User__Username__istartswith",username).RelatedSel().All(&auths)
	}
	skip := false
	for _,auth := range auths {
		skip = false
		for _,user := range users {
			if user.Id == auth.User.Id {
				user.ServiceAuths = append(user.ServiceAuths,auth.ServiceAuth)
			}
		}
		if skip == false {
			auth.User.ServiceAuths = append(auth.User.ServiceAuths,auth.ServiceAuth)
			auth.User.Password = ""
			users = append(users,auth.User)
		}
	}
	return  users,err
}