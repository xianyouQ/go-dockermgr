package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "time"
	"errors"
	"github.com/astaxie/beego/validation"
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
    Service *Service `orm:"rel(fk)" valid:"Required"`
    ReleaseTime time.Time `orm:"auto_now_add"`
    ImageTag string `orm:"size(20)" valid:"Required"`
    ReleaseUser *User `orm:"rel(fk)" valid:"Required"`
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
func checkReleaseTask(t *ReleaseTask) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&t)
	if !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}

func QueryRelease(o orm.Ormer,service *Service) ([]*ReleaseTask,error) {
    var ReleaseTaskList []*ReleaseTask
    _,err := o.QueryTable(beego.AppConfig.String("dockermgr_release_table")).Filter("Service__Id",service.Id).All(&ReleaseTaskList)
    return ReleaseTaskList,err
}

func NewRelease(o orm.Ormer,service *Service,releaseUser *User,imageTag string) (*ReleaseTask,error) {
    mReleaseTask := new(ReleaseTask)
    if imageTag == "" {
        return mReleaseTask,errors.New("not supported null imageTag")
    }
    mReleaseTask.Service = service
    mReleaseTask.ReleaseUser = releaseUser
    mReleaseTask.ImageTag = imageTag
    mReleaseTask.TaskStatus = NotReady
    _, err := o.Insert(mReleaseTask)
    return mReleaseTask,err
}

/*
func ReviewRelease(releaseTask *ReleaseTask,reviewUser *User) (*ReleaseTask,error) {

}

func OperationRelease(releaseTask *ReleaseTask,operationUser *User) (*ReleaseTask,error) {

}
func CancelRelease(releaseTask *ReleaseTask,cancelUser *User) (*ReleaseTask,error){

}
*/
