package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/xianyouQ/go-dockermgr/utils"
)

type Service struct {
    Id int `orm:"auto"`
    Name string `orm:"size(20);unique"`
    Code string `orm:"size(20);unique"`
    Instances []*Ip `orm:"reverse(many)"`
    ReleaseTask []*ReleaseTask `orm:"reverse(many)"`
    MarathonConf string `orm:"type(text)"`
}

func ( this *Service) TableName() string {
    return beego.AppConfig.String("dockermgr_service_table")
}

func init() {
    orm.RegisterModel(new(Service))
}

func AddService(name string,code string) error {
    service := Service{Name:name,Code:code}
    o := orm.NewOrm()
    _,err := o.Insert(service)
    if err!=nil {
        return err
    }
    return nil
}


func (self *Service) SetMarathonConf(conf string) error {
    if conf == self.MarathonConf {
        return nil
    }
    _,err := utils.CreateMarathonAppFromJson(conf)
    if err != nil {
        return err
    }
    o := orm.NewOrm()
    self.MarathonConf = conf
    if _,err := o.Update(self,"MarathonConf"); err !=nil {
         return  err
     }

    return nil

} 

func (self Service) GetInstances() ([]*Ip,error) {
    o := orm.NewOrm()
    var Ips []*Ip
    _,err := o.QueryTable("Ip").Filter("BelongService",self.Id).RelatedSel().All(&Ips)
    if err != nil {
        return Ips,err
    }
    return Ips,nil
}

func (self Service) GetInstancesWithCidr(cidr Cidr) ([]*Ip,error) {
    o := orm.NewOrm()
    var Ips []*Ip
    _,err := o.QueryTable("Ip").Filter("BelongService",self.Id).Filter("BelongNet",cidr.Id).RelatedSel().All(&Ips)
    if err != nil {
        return Ips,err
    }
    return Ips,nil
}