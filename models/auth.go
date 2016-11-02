package models

import (
    //"github.com/astaxie/beego"
    //"github.com/astaxie/beego/orm"
)

type User struct {
    UserName string `orm:"size(30)"`
}

type Permisson struct {
    
}