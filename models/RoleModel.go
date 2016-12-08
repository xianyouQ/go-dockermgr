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
		role.Id = id
		if role.NeedAddAuth == true {
			var services []*Service
			services,err = QueryService()
			if err != nil {
				return 0,err
			}
			for _,service := range services {
				_,err = NewServiceAuth(role,service)
				if err != nil {
					return 0,err
				}
			}
		}
	} else {
		_,err = o.Update(role)
		if err != nil {
			return 0,err
		}
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

func AddRoleNode(roleid int64, nodeid int64) (int64, error) {
	o := orm.NewOrm()
	role := Role{Id: roleid}
	node := Node{Id: nodeid}
	m2m := o.QueryM2M(&node, "Role")
	num, err := m2m.Add(&role)
	return num, err
}

/*
func DelUserRole(roleid int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user_roles").Filter("role_id", roleid).Delete()
	return err
}
func AddRoleUser(roleid int64, userid int64) (int64, error) {
	o := orm.NewOrm()
	role := Role{Id: roleid}
	user := User{Id: userid}
	m2m := o.QueryM2M(&user, "Role")
	num, err := m2m.Add(&role)
	return num, err
}

func GetUserByRoleId(roleid int64) (users []orm.Params, count int64) {
	o := orm.NewOrm()
	user := new(User)
	count, _ = o.QueryTable(user).Filter("Role__Role__Id", roleid).Values(&users)
	return users, count
}

func AccessList(uid int64) (list []orm.Params, err error) {
	var roles []orm.Params
	o := orm.NewOrm()
	role := new(Role)
	_, err = o.QueryTable(role).Filter("User__User__Id", uid).Values(&roles)
	if err != nil {
		return nil, err
	}
	var nodes []orm.Params
	node := new(Node)
	for _, r := range roles {
		_, err := o.QueryTable(node).Filter("Role__Role__Id", r["Id"]).Values(&nodes)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			list = append(list, n)
		}
	}
	return list, nil
}
*/