package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
)

type Service struct {
    Id int `orm:"auto"`
    Name string `orm:"size(20);unique"`
    Image string `orm:"size(100)"`
    Tag string `orm:"size(100)"`
    Instances *[]Ip `orm:"reverse(many)"`
}

type ContainerConf struct {

    
}
func AddService() {

}

func DelService() {

}
