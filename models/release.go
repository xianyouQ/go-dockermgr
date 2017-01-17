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
    ReleaseDetail string `orm:"type(text)"`
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




func CreateOrUpdateRelease(o orm.Ormer,releaseTask *ReleaseTask,updatecols ...string) error {
    var err error
    if err = checkReleaseTask(releaseTask); err != nil {
		return err
	}
	if releaseTask.Id == 0 {
		_,err = o.Insert(releaseTask)
		return err
	} else {
		if len(updatecols) == 0 {
			_, err = o.Update(releaseTask)
		} else {
			_, err = o.Update(releaseTask,updatecols...)
		}
		return err
	}
}


