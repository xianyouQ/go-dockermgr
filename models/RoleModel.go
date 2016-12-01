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
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"`
	Node   []*Node `orm:"rel(m2m)"`
	//User   []*User `orm:"reverse(many)"`
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
	_,err := o.QueryTable(beego.AppConfig.String("rbac_role_table")).RelatedSel("Node").All(&roles)
	
	return roles, err
}


func AddRole(r *Role) (int64, error) {
	if err := checkRole(r); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	role := new(Role)
	role.Name = r.Name
	role.Status = r.Status

	id, err := o.Insert(role)
	return id, err
}

func UpdateRole(r *Role) (int64, error) {
	if err := checkRole(r); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	role := make(orm.Params)
	if len(r.Name) > 0 {
		role["Name"] = r.Name
	}
	if r.Status != 0 {
		role["Status"] = r.Status
	}
	if len(role) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Role
	num, err := o.QueryTable(table).Filter("Id", r.Id).Update(role)
	return num, err
}

func DelRoleById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Role{Id: Id})
	return status, err
}

func GetNodelistByRoleId(Id int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Role__Role__Id", Id).Values(&nodes)
	return nodes, count
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