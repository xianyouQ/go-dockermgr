package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/astaxie/beego/validation"
    "errors"
	"github.com/astaxie/beego/logs"
	"encoding/json"
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
    idc.Status = IdcEnable
    err = checkIdc(idc)
    if err!=nil {
        return err
    }
    o := orm.NewOrm()
    pid,err = o.InsertOrUpdate(idc)  //返回值都是什么意思
    if err!=nil {
        return err
    }
    if pid != 0 {
        GlobalIdcConfList = append(GlobalIdcConfList,idc)
    } else {
        for _,IdcConfIter := range GlobalIdcConfList{
            if IdcConfIter.Id == idc.Id {
                IdcConfIter = idc
            }
        }
    }
    test,_:=json.Marshal(GlobalIdcConfList)
    logs.GetLogger("idcModel").Println(string(test))
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