package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/gambol99/go-marathon/marathon"
    "time"
)

type MgrConf struct{
    Id int `orm:"auto"`
    Version time.Time `orm:"auto_now"`
    Registrys *[]RegistryConf `orm:"rel(m2m)"`
    MarathonConfTemplate string `orm:"type(text)"`
}

type RegistryConf struct {
    Host string `orm:"size(20)"`
    Port int 
    Schema string `orm:"size(10)"`
    UserName string `orm:"size(20)"`
    Password string `orm:"size(20)"`
}


func (self *MgrConf) SetMarathonConf(conf string) error {
    if conf == self.MarathonConfTemplate {
        return nil
    }
    MarathonConf := &marathon.Application{}
    err := json.Unmarshal([]byte(conf),&stb)
    if err != nil {
        return err
    }
    self.MarathonConfTemplate = conf
    o := orm.NewOrm()
    if num,err := o.Update(self,"MarathonConfTemplate"); err !=nil {
         return  err
     }
} 