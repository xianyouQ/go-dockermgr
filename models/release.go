package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

const (
	NotReady = iota
	Ready
	Running
	Paused
	Success
	Failed
	Cancel
)

type ReleaseTask struct {
	Id            int          `orm:"auto"`
	ReleaseConf   *ReleaseConf `orm:"rel(fk)" valid:"Required"`
	ReleaseTime   time.Time    `orm:"auto_now_add"`
	Service       *Service     `orm:"rel(fk)" valid:"Required"`
	ImageTag      string       `orm:"size(20)" valid:"Required"`
	ReleaseUser   *User        `orm:"null;rel(fk)" valid:"Required"`
	OperationUser *User        `orm:"null;rel(fk)"`
	ReviewUser    *User        `orm:"null;rel(fk)"`
	TaskStatus    int
	CancelUser    *User  `orm:"null;rel(fk)"`
	ReleaseDetail string `orm:"type(text)"`
}

type ReleaseConf struct {
	Id              int        `orm:"auto"`
	Service         *Service   `orm:"rel(fk)" valid:"Required"`
	FaultTolerant   int        `orm:"default(1)"`
	IdcParalle      int        `orm:"default(1)"`
	IdcInnerParalle int        `orm:"default(1)"`
	ReleaseIdc      []*IdcConf `orm:"rel(m2m)"`
	TimeOut         int64
}

func (this *ReleaseTask) TableName() string {
	return beego.AppConfig.String("dockermgr_release_table")
}
func init() {
	orm.RegisterModel(new(ReleaseTask), new(ReleaseConf))
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

func checkReleaseConf(t *ReleaseConf) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&t)
	if !b {
		for _, err := range valid.Errors {
			return errors.New(err.Message)
		}
	}
	return nil
}

func QueryRelease(o orm.Ormer, service *Service) ([]*ReleaseTask, error) {
	var ReleaseTaskList []*ReleaseTask
	_, err := o.QueryTable(beego.AppConfig.String("dockermgr_release_table")).Filter("Service__Id", service.Id).All(&ReleaseTaskList)
	return ReleaseTaskList, err
}

func CreateOrUpdateRelease(o orm.Ormer, releaseTask *ReleaseTask, updatecols ...string) error {
	var err error
	if err = checkReleaseTask(releaseTask); err != nil {
		return err
	}
	if releaseTask.Id == 0 {
		_, err = o.Insert(releaseTask)
		return err
	} else {
		if len(updatecols) == 0 {
			_, err = o.Update(releaseTask)
		} else {
			_, err = o.Update(releaseTask, updatecols...)
		}
		return err
	}
}

func CreateOrUpdateReleaseConf(o orm.Ormer, releaseConf *ReleaseConf, updatecols ...string) error {
	var err error
	if err = checkReleaseConf(releaseConf); err != nil {
		return err
	}
	if releaseConf.Id == 0 {
		_, err = o.Insert(releaseConf)
		return err
	} else {
		if len(updatecols) == 0 {
			_, err = o.Update(releaseConf)
		} else {
			_, err = o.Update(releaseConf, updatecols...)
		}
		return err
	}
}
