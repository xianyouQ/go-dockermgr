package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    . "github.com/xianyouQ/go-dockermgr/auth/models"
    "time"
)
const (
    NotReady = iota
    Ready
    Running
    Success
    Failed
    Cancel
)
type ReleaseTask struct {
    Id int `orm:"auto"`
    Service *Service `orm:"rel(fk)"`
    ReleaseTime time.Time `orm:"auto_now_add"`
    ImageTag string `orm:"size(20)"`
    ReleaseUser *User `orm:"rel(fk)"`
    OperationUser *User `orm:"rel(fk)"`
    ReviewUser *User `orm:"rel(fk)"`
    TaskStatus int
}

func ( this *ReleaseTask) TableName() string {
    return beego.AppConfig.String("dockermgr_release_table")
}
func init() {
    orm.RegisterModel(new(ReleaseTask))
}