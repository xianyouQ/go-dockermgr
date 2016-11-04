package models

import (
    "github.com/astaxie/beego/orm"
)

type IdcConf struct {
    Id int `orm:"auto"`
    MgrConf *MgrConf `orm:ref(fk)`
    Cidrs []*Cidr `orm:ref()`
}