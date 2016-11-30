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
    IdcName string `orm:"size(30);unique" valid:"Required;MinSize(2)"`
    IdcCode string `orm:"size(30);unique" valid:"Required;MinSize(2)"`
    RegistryConf *RegistryConf `orm:"null;rel(one)"`
    MarathonSerConf *MarathonSerConf `orm:"null;rel(one)"`
    Cidrs []*Cidr `orm:"null;reverse(many)"`
}


var (
    GlobalIdcConfList []*IdcConf
)

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

func AddOrUpdateIdc(idc *IdcConf) error {
    var err error
    var pid int64
    err = checkIdc(idc)
    if err!=nil {
        return err
    }
    o := orm.NewOrm()
    if(idc.Id == 0) {
        idc.Status = IdcEnable
        pid,err = o.Insert(idc)
        if err != nil {
            return err
        }
    } else {
        _,err = o.Update(idc)
        if err != nil {
            return err
        }
    }
    if pid != 0 {
        GlobalIdcConfList = append(GlobalIdcConfList,idc)
    } else {
        for inx,IdcConfIter := range GlobalIdcConfList{
            if IdcConfIter.Id == idc.Id {
                GlobalIdcConfList[inx] = idc
            }
        }
    }
    return nil
}
func getIdcsfromOrm() ([]*IdcConf,error) {
    var tempCidrs []*Cidr
    var idcs []*IdcConf
    o := orm.NewOrm()
    _,err := o.QueryTable(beego.AppConfig.String("dockermgr_idc_table")).RelatedSel().All(&idcs)
    if err != nil {
        return idcs,err
    }
    tempCidrs,err = GetCidrFromOrm()
    if err!=nil {
         return idcs,err
    }
    for _,IdcConfIter := range idcs{
        for _,Cidriter := range tempCidrs {
            if Cidriter.BelongIdc.Id == IdcConfIter.Id {
                Cidriter.BelongIdc = nil
                IdcConfIter.Cidrs = append(IdcConfIter.Cidrs,Cidriter)
            }
        }
    }
    return idcs,nil
}

func GetIdcs() ([]*IdcConf,error) {
    var err error
    if GlobalIdcConfList == nil {
       GlobalIdcConfList,err = getIdcsfromOrm()
       if err!=nil {
           return GlobalIdcConfList,err
       }
    }
    return GlobalIdcConfList,nil
}



/*
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
*/

func init() {
    orm.RegisterModel(new(IdcConf))
}