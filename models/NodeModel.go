package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//节点表
type Node struct {
	Id     int64
	Title  string  `orm:"size(100)" form:"Title"  valid:"Required"`
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`
	
}

func (n *Node) TableName() string {
	return beego.AppConfig.String("rbac_node_table")
}

//验证用户信息
func checkNode(u *Node) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}

func init() {
	orm.RegisterModel(new(Node))
}

//get node list
func GetNodelist(page int64, page_size int64, sort string) ([]*Node,int64) {
	o := orm.NewOrm()
	var nodes []*Node
	var count int64
	qs := o.QueryTable(beego.AppConfig.String("rbac_node_table"))
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).All(&nodes)
	count, _ = qs.Count()
	return nodes, count
}

func ReadNode(nid int64) (Node, error) {
	o := orm.NewOrm()
	node := Node{Id: nid}
	err := o.Read(&node)
	if err != nil {
		return node, err
	}
	return node, nil
}

//添加用户
func AddNode(n *Node) (int64, error) {
	if err := checkNode(n); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	node := new(Node)
	node.Title = n.Title
	node.Name = n.Name
	id, err := o.Insert(node)
	return id, err
}

//更新用户
func UpdateNode(n *Node) (int64, error) {
	if err := checkNode(n); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	node := make(orm.Params)
	if len(n.Title) > 0 {
		node["Title"] = n.Title
	}
	if len(n.Name) > 0 {
		node["Name"] = n.Name
	}
	if len(node) == 0 {
		return 0, errors.New("update field is empty")
	}
	num, err := o.QueryTable(beego.AppConfig.String("rbac_node_table")).Filter("Id", n.Id).Update(node)
	return num, err
}

func DelNodeById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Node{Id: Id})
	return status, err
}


