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
    Server string `orm:"size(50);unique" valid:"Required;MinSize(5)"`
    HttpBasicAuthUser string `orm:"size(50)" valid:"Required;MinSize(2)"`
    HttpBasicPassword string `orm:"size(50)" valid:"Required;MinSize(2)"`
    MarathonConfTemplate string `orm:"type(text)" valid:"Required;MinSize(5)"`
    //BelongIdc *IdcConf `orm:"reverse(one)"` 防止解析json的时候循环解析
    
}

type RegistryConf struct {
    Id int `orm:"auto"`
    Server string `orm:"size(50);unique" valid:"Required;MinSize(5)"`
    UserName string `orm:"size(20)" valid:"Required;MinSize(2)"`
    Password string `orm:"size(20)" valid:"Required;MinSize(2)"`
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

func AddOrUpdateMarathonSerConf(o orm.Ormer,newConf *MarathonSerConf,updatecols ...string) error {
    var err error
    err = checkMarathonSerConf(newConf)
    if err != nil {
        return err
    }
    _,err = utils.CreateMarathonAppFromJson(newConf.MarathonConfTemplate)
    if err != nil {
        return err
    }
    if(newConf.Id == 0) {
        _,err = o.Insert(newConf)
        if err != nil {
            return err
        }
    } else {
        if len(updatecols) == 0 {
            _,err = o.Update(newConf)
        } else {
            _,err = o.Update(newConf,updatecols...)
        }
        if err != nil {
            return err
        }
    }
    return nil
}


func AddOrUpdateRegistryConf(o orm.Ormer,newConf *RegistryConf,updatecols ...string) error {
    var err error
    err = checkRegistryConf(newConf)
    if err != nil {
        return err
    }
    if(newConf.Id == 0) {
        _,err = o.Insert(newConf)
        if err != nil {
            return err
        }
    } else {
        if len(updatecols) == 0 {
            _,err = o.Update(newConf)
        } else {
            _,err = o.Update(newConf,updatecols...)
        }
        if err != nil {
            return err
        }
    }
    return nil
}
