package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/gambol99/go-marathon/marathon"
    "encoding/json"
)

type Service struct {
    Id int `orm:"auto"`
    Name string `orm:"size(20);unique"`
    //Image string `orm:"size(100)"`
    //Tag string `orm:"size(100)"`
    Instances *[]Ip `orm:"reverse(many)"`
    ReleaseTask *[]ReleaseTask `orm:"reverse(many)"`
    MarathonConf string `orm:"type(text)"`
}

func init() {
    orm.RegisterModel(new(Service))
}

func AddService(serivcename string) error {

}

func DelService(servicename string) error {

}


func (self *Service) SetMarathonConf(conf string) error {
    if conf == self.MarathonConf {
        return nil
    }
    MarathonConf := &marathon.Application{}
    err := json.Unmarshal([]byte(conf),&stb)
    if err != nil {
        return err
    }
    o := orm.NewOrm()
    self.MarathonConf = conf
    if num,err := o.Update(self,"MarathonConf"); err !=nil {
         return  err
     }

    return nil

} 

func (self Service) GetInstances() ([]*Ip,error) {
    o := orm.NewOrm()
    var Ips []*Ip
    num,err := o.QueryTable("Ip").Filter("BelongService",self.Id).RelatedSel().All(&Ips)
    if err != nil {
        return Ips,err
    }
    return Ips,nil
}

func (self Service) GetInstancesWithCidr(cidr Cidr) ([]*Ip,error) {
    o := orm.NewOrm()
    var Ips []*Ip
    num,err := o.QueryTable("Ip").Filter("BelongService",self.Id).Filter("BelongNet",cidr.Id).RelatedSel().All(&Ips)
    if err != nil {
        return Ips,err
    }
    return Ips,nil
}