package models

import (
    //"github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/xianyouQ/go-dockermgr/utils"
    "time"
)

type MgrConf struct {
    Id int `orm:"auto"`
    Version time.Time `orm:"auto_now"`
    Registrys *RegistryConf `orm:"rel(fk)"`
    MarathonConfTemplate string `orm:"type(text)"`
    MarathonSerConf MarathonSerConf `orm:"rel(fk)"`
}

type MarathonSerConf struct {
    Server string `orm:"size(50)"`
    HTTPBasicAuthUser string `orm:"size(50)"`
    HTTPBasicPassword string `orm:"size(50)"`
    PollingWaitTime int `orm:"size(50)"`
}

type RegistryConf struct {
    Host string `orm:"size(20)"`
    Port int 
    Schema string `orm:"size(10)"`
    UserName string `orm:"size(20)"`
    Password string `orm:"size(20)"`
}

func init() {
    orm.RegisterModel(new(MgrConf),new(RegistryConf))
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
    if num,err := o.Update(self,"MarathonConfTemplate"); err !=nil {
         return  err
     }
     return nil
} 

//func (slf *RegistryConf) 