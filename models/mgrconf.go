package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/xianyouQ/go-dockermgr/utils"
    "time"
)

type MgrConf struct {
    Id int `orm:"auto"`
    Version time.Time `orm:"auto_now"`
    Registrys *RegistryConf `orm:"rel(one)"`
    MarathonConfTemplate string `orm:"type(text)"`
    MarathonSerConf *MarathonSerConf `orm:"rel(one)"`
    BelongIdc *IdcConf `orm:"reverse(one)"`
}

type MarathonSerConf struct {
    Id int `orm:"auto"`
    Server string `orm:"size(50)"`
    HttpBasicAuthUser string `orm:"size(50)"`
    HttpBasicPassword string `orm:"size(50)"`
    PollingWaitTime int `orm:"size(50)"`
    BelongMgrConf *MgrConf `orm:"reverse(one)"`
}

type RegistryConf struct {
    Id int `orm:"auto"`
    Server string `orm:"size(50)"`
    UserName string `orm:"size(20)"`
    Password string `orm:"size(20)"`
    BelongMgrConf *MgrConf `orm:"reverse(one)"`
}

func init() {
    orm.RegisterModel(new(MgrConf),new(RegistryConf),new(MarathonSerConf))
}


func ( this *MgrConf) TableName() string {
    return beego.AppConfig.String("dockermgr_mgrconf_table")
}

func ( this *MarathonSerConf) TableName() string {
    return beego.AppConfig.String("dockermgr_marathonconf_table")
}

func ( this *RegistryConf) TableName() string {
    return beego.AppConfig.String("dockermgr_registryconf_table")
}

func (self *MgrConf) SetMarathonConf(conf string) error {
    if conf == self.MarathonConfTemplate {
        return nil
    }
    _,err := utils.CreateMarathonAppFromJson(conf)
    if err != nil {
        return err
    }
    self.MarathonConfTemplate = conf
    o := orm.NewOrm()
    if _,err := o.Update(self,"MarathonConfTemplate"); err !=nil {
         return  err
     }
     return nil
} 

//func (slf *RegistryConf) 