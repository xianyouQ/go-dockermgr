package models

import (
    "github.com/astaxie/beego/orm"
)

const (
    IdcEnable =iota
    IdcDisable
)
type IdcConf struct {
    Id int `orm:"auto"`
    Status int
    IdcName string `orm:"size(30);unique"`
    IdcCode string `orm:"size(30);unique"`
    MgrConf *MgrConf `orm:"ref(one)"`
    Cidrs []*Cidr `orm:"reverse(many)"`
}

func AddIdc(name string,code string) error {
    idc := IdcConf{IdcName:name,IdcCode:code,Status:IdcEnable}
    o := orm.NewOrm()
    _,err := o.Insert(idc)
    if err!=nil {
        return err
    }
    return nil
}

func toggleStatus(name string,status int) error {
    idc := IdcConf{IdcName:name,Status:status}
    o := orm.NewOrm()
    if _, err := o.Update(&idc,"Status"); err != nil {
        return err
    }
    return nil
}

func EnableIdc (name string) error {
    return toggleStatus(name,IdcEnable)
}

func DisableIdc (name string) error {
    return toggleStatus(name,IdcDisable)
}

func init() {
    orm.RegisterModel(new(IdcConf))
}