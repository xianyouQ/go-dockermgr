package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/astaxie/beego/validation"
    "errors"
)

const (
    IdcEnable =iota
    IdcDisable
)
type IdcConf struct {
    Id int `orm:"pk;auto"`
    Status int
    IdcName string `orm:"size(30);unique" valid:"Required"`
    IdcCode string `orm:"size(30);unique" valid:"Required"`
    MgrConf *MgrConf `orm:"null;rel(one)"`
    Cidrs []*Cidr `orm:"null;reverse(many)"`
}


func ( this *IdcConf) TableName() string {
    return beego.AppConfig.String("dockermgr_idc_table")
}

func checkIdc(idc *IdcConf) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&idc)
	if !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}

func AddIdc(idc *IdcConf) error {
    var err error
    idc.Status = IdcEnable
    err = checkIdc(idc)
    o := orm.NewOrm()
    _,err = o.Insert(idc)
    if err!=nil {
        return err
    }
    return nil
}
func GetIdcs() ([]*IdcConf,error) {
    var idcs []*IdcConf
    o := orm.NewOrm()
    _,err := o.QueryTable(beego.AppConfig.String("dockermgr_idc_table")).RelatedSel().All(&idcs)
    if err != nil {
        return idcs,err
    }
    return idcs,nil
}

func GetMgrConf(code string) (*IdcConf,error) {
    var idc *IdcConf
    o := orm.NewOrm()
    _,err := o.QueryTable(beego.AppConfig.String("dockermgr_idc_table")).All(&idc,"MgrConf")
    if err != nil {
        return idc,err
    }
    return idc,nil
}

func UpdateIdc(idc *IdcConf) error {
     var err error
     err = checkIdc(idc)
     if err!=nil {
        return err
     }
     o := orm.NewOrm()
    _,err = o.Update(idc)
    if err!=nil {
        return err
    }
    return nil
}

func toggleStatus(code string,status int) error {
    idc := IdcConf{IdcCode:code,Status:status}
    o := orm.NewOrm()
    if _, err := o.Update(&idc,"Status"); err != nil {
        return err
    }
    return nil
}

func EnableIdc (code string) error {
    return toggleStatus(code,IdcEnable)
}

func DisableIdc (code string) error {
    return toggleStatus(code,IdcDisable)
}

func init() {
    orm.RegisterModel(new(IdcConf))
}