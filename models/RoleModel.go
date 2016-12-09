package models

import (
	"errors"
	"log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//角色表
type Role struct {
	Id     int64
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`
	Status bool     `orm:"default(true)" form:"Status"`
	Nodes   []*Node `orm:"rel(m2m)"`
	NeedAddAuth bool `orm:"-" form:"NeedAddAuth" valid:"Required"`
	Services []*Service `orm:"rel(m2m);rel_through(github.com/xianyouQ/go-dockermgr/models.ServiceAuth)"`
}

func (r *Role) TableName() string {
	return beego.AppConfig.String("rbac_role_table")
}

func init() {
	orm.RegisterModel(new(Role))
}

func checkRole(g *Role) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&g)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

//get role list
func GetRoleListFromOrm() ([]*Role,error) {
	o := orm.NewOrm()
	var roles []*Role
	_,err := o.QueryTable(beego.AppConfig.String("rbac_role_table")).All(&roles)
	for _, role := range roles {
		_,err = o.LoadRelated(role, "Nodes")
		if err != nil {
			return	roles, err
		}
	}
	return roles, err
}


func AddOrUpdateRole(role *Role) (int64, error) {
	var id int64
	var err error
	if err = checkRole(role); err != nil {
		return 0, err
	}
	o := orm.NewOrm()

	if role.Id == 0 {
		id, err = o.Insert(role)
		if err != nil {
			return 0,err
		}
		if role.NeedAddAuth == true {
			var services []*Service
			services,err = QueryService()
			if err != nil {
				return 0,err
			}
			_,err = NewServiceAuths(role,services)
			if err != nil {
					return 0,err
			}
		} else {
			_,err = NewServiceAuth(role,nil)
			if err != nil {
					return 0,err
			}
		}
	} else {
		_,err = o.Update(role)
		if err != nil {
			return 0,err
		}
	}

	if id != 0 {
		UpdateRoleNodes(role,true)
	} else {
		UpdateRoleNodes(role,false)
	}
	return role.Id, err
}


func DelRole(role *Role)  error {
	o := orm.NewOrm()
	_, err := o.Delete(role)
	return err
}

func GetNodelistByRole(role *Role) (int64,error) {
	o := orm.NewOrm()
	count, err := o.QueryTable(beego.AppConfig.String("rbac_node_table")).Filter("Roles__Role__Id", role.Id).All(&role.Nodes)
	return  count,err
}

func AddRoleNode(role *Role, nodes []*Node) (int64, error) {
	o := orm.NewOrm()
	m2m := o.QueryM2M(role, "Nodes")
	num, err := m2m.Add(nodes)
	return num, err
}

func DelRoleNode(role *Role, nodes []*Node) (int64, error) {
	o := orm.NewOrm()
	m2m := o.QueryM2M(role, "Nodes")
	num, err := m2m.Remove(nodes)
	return num, err
}
