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
	Id          int64
	Name        string     `orm:"size(100)" form:"Name"  valid:"Required"`
	Status      bool       `orm:"default(true)" form:"Status"`
	Nodes       []*Node    `orm:"rel(m2m)"`
	NeedAddAuth bool       `orm:"-" form:"NeedAddAuth" valid:"Required"`
	Services    []*Service `orm:"rel(m2m);rel_through(github.com/xianyouQ/go-dockermgr/models.ServiceAuth)"`
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
func GetRoleListFromOrm() ([]*Role, error) {
	o := orm.NewOrm()
	var roles []*Role
	_, err := o.QueryTable(beego.AppConfig.String("rbac_role_table")).All(&roles)
	for _, role := range roles {
		_, err = o.LoadRelated(role, "Nodes")
		if err != nil {
			return roles, err
		}
	}
	return roles, err
}

func AddOrUpdateRole(o orm.Ormer, role *Role, updatecols ...string) error {
	var id int64
	var err error
	if err = checkRole(role); err != nil {
		return err
	}
	if role.Id == 0 {
		id, err = o.Insert(role)
		if err != nil {
			return err
		}
		if role.NeedAddAuth == true {
			var services []*Service
			services, err = GetServices()
			if err != nil {
				return err
			}
			for _, service := range services {
				_, err = NewServiceAuth(o, role, service)
				if err != nil {
					return err
				}
			}
		} else {
			_, err = NewServiceAuth(o, role, nil)
			if err != nil {
				return err
			}
		}
	} else {
		if len(updatecols) == 0 {
			_, err = o.Update(role)
		} else {
			_, err = o.Update(role, updatecols...)
		}
		if err != nil {
			return err
		}
	}

	if id != 0 {
		UpdateRoleNodes(role, true)
	} else {
		UpdateRoleNodes(role, false)
	}
	return err
}

func DelRole(o orm.Ormer, role *Role) error {
	_, err := o.Delete(role)
	return err
}

func GetNodelistByRole(role *Role) (int64, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable(beego.AppConfig.String("rbac_node_table")).Filter("Roles__Role__Id", role.Id).All(&role.Nodes)
	return count, err
}

func AddRoleNode(o orm.Ormer, role *Role, nodes []*Node) (int64, error) {
	if len(nodes) <= 0 {
		return 0, nil
	}
	m2m := o.QueryM2M(role, "Nodes")
	num, err := m2m.Add(nodes)
	return num, err
}

func DelRoleNode(o orm.Ormer, role *Role, nodes []*Node) (int64, error) {
	if len(nodes) <= 0 {
		return 0, nil
	}
	m2m := o.QueryM2M(role, "Nodes")
	num, err := m2m.Remove(nodes)
	return num, err
}

func QueryRole(name string) (*Role, error) {
	roles, err := GetRoleNodes()
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.Name == name {
			return role, nil
		}
	}
	err = errors.New("role not found")
	return nil, err
}
