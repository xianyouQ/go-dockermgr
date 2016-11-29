package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/xianyouQ/go-dockermgr/utils"
	"github.com/astaxie/beego/validation"
    "errors"
)

type MarathonSerConf struct {
    Id int `orm:"auto"`
    Server string `orm:"size(50);unique" valid:"Required"`
    HttpBasicAuthUser string `orm:"size(50)" valid:"Required"`
    HttpBasicPassword string `orm:"size(50)" valid:"Required"`
    MarathonConfTemplate string `orm:"type(text)" valid:"Required"`
    //BelongIdc *IdcConf `orm:"reverse(one)"` 防止解析json的时候循环解析
    
}

type RegistryConf struct {
    Id int `orm:"auto"`
    Server string `orm:"size(50);unique" valid:"Required"`
    UserName string `orm:"size(20)" valid:"Required"`
    Password string `orm:"size(20)" valid:"Required"`
    //BelongIdc *IdcConf `orm:"reverse(one)"`
}

func init() {
    orm.RegisterModel(new(RegistryConf),new(MarathonSerConf))
}


func ( this *MarathonSerConf) TableName() string {
    return beego.AppConfig.String("dockermgr_marathonconf_table")
}

func ( this *RegistryConf) TableName() string {
    return beego.AppConfig.String("dockermgr_registryconf_table")
}

func checkMarathonSerConf(conf *MarathonSerConf) error {
    valid := validation.Validation{}
	b, _ := valid.Valid(&conf)
	if !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
    return nil
}

func checkRegistryConf(conf *RegistryConf) error {
    valid := validation.Validation{}
	b, _ := valid.Valid(&conf)
	if !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
    return nil
}

func AddOrUpdateMarathonSerConf(newConf *MarathonSerConf) (int64,error) {
    var err error
    var id int64
    err = checkMarathonSerConf(newConf)
    if err != nil {
        return id,err
    }
    _,err = utils.CreateMarathonAppFromJson(newConf.MarathonConfTemplate)
    if err != nil {
        return id,err
    }
     o := orm.NewOrm()
    id,err = o.InsertOrUpdate(newConf)
    if err!=nil {
        return id,err
    }
    return id,nil
}


func AddOrUpdateRegistryConf(newConf *RegistryConf) (int64,error) {
    var err error
    var id int64
    err = checkRegistryConf(newConf)
    if err != nil {
        return id,err
    }
     o := orm.NewOrm()
    id,err = o.InsertOrUpdate(newConf)
    if err!=nil {
        return id,err
    }
    return id,nil
}
