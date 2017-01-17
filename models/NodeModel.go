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
	Desc  string  `orm:"size(100)" form:"Title"  valid:"Required"`
	Url   string  `orm:"size(100)" form:"Name"  valid:"Required"`
	Roles   []*Role `orm:"reverse(many)"`
	Active bool  `orm:"-"`
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
func GetNodes() ([]*Node,error) {
	o := orm.NewOrm()
	var Nodes []*Node
	_,err := o.QueryTable(beego.AppConfig.String("rbac_node_table")).All(&Nodes)
	return Nodes, err
}




//更新用户
func AddOrUpdateNode(o orm.Ormer,node *Node,updatecols ...string) (error) {
	var err error
	if err = checkNode(node); err != nil {
		return err
	}
	if node.Id == 0 {
		_,err = o.Insert(node)
		return err
	} else {
		if len(updatecols) == 0 {
			_, err = o.Update(node)
		} else {
			_, err = o.Update(node,updatecols...)
		}
		return err
	}

}

func DelNodeById(o orm.Ormer,Id int64) (int64, error) {
	status, err := o.Delete(&Node{Id: Id})
	return status, err
}


