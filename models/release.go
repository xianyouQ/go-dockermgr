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
	Abandon
	Success
	Failed
	OutOfTolerant
	Cancel
)

type ReleaseTask struct {
	Id            int          `orm:"auto"`
	ReleaseConf   *ReleaseConf `orm:"rel(fk)" valid:"Required"`
	ReleaseTime   time.Time    `orm:"auto_now_add"`
	OperationTime time.Time    `orm:"null"`
	Service       *Service     `orm:"rel(fk)" valid:"Required"`
	ImageTag      string       `orm:"size(20)" valid:"Required"`
	ReleaseUser   *User        `orm:"null;rel(fk)" valid:"Required"`
	OperationUser *User        `orm:"null;rel(fk)"`
	ReviewUser    *User        `orm:"null;rel(fk)"`
	ReviewTime    time.Time    `orm:"null"`
	TaskStatus    int
	CancelUser    *User  `orm:"null;rel(fk)"`
	ReleaseDetail string `orm:"type(text)"`
	ReleaseResult string `orm:"type(text)"`
	ReleaseMsg    string `orm:"type(text)"`
}

type ReleaseConf struct {
	Id              int        `orm:"auto"`
	Service         *Service   `orm:"rel(fk)" valid:"Required"`
	FaultTolerant   int        `orm:"default(1)"`
	IdcParalle      int        `orm:"default(1)"`
	IdcInnerParalle int        `orm:"default(1)"`
	ReleaseIdc      []*IdcConf `orm:"rel(m2m);rel_through(github.com/xianyouQ/go-dockermgr/models.ReleaseConfIdc)"`
	TimeOut         int64
}

type ReleaseConfIdc struct {
	Id          int          `orm:"auto"`
	ReleaseConf *ReleaseConf `orm:"rel(fk)"`
	Idc         *IdcConf     `orm:"rel(fk)"`
}

func (this *ReleaseTask) TableName() string {
	return beego.AppConfig.String("dockermgr_release_table")
}

func (this *ReleaseConf) TableName() string {
	return beego.AppConfig.String("dockermgr_releaseconf_table")
}

func (this *ReleaseConfIdc) TableName() string {
	return beego.AppConfig.String("dockermgr_releaseconfidc_table")
}
func init() {
	orm.RegisterModel(new(ReleaseTask), new(ReleaseConf), new(ReleaseConfIdc))
}
func checkReleaseTask(t *ReleaseTask) (err error) {
	if t.ReleaseConf.Id < 1 {
		return errors.New("invaild releaseConf Id,Must greater 0")
	}
	if t.Service.Id < 1 {
		return errors.New("invaild Service Id,Must greater 0")
	}
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

func LoadReleaseConf(o orm.Ormer, releaseTask *ReleaseTask) error {
	var ReleaseConfIdcs []*ReleaseConfIdc
	var err error
	err = o.QueryTable(beego.AppConfig.String("dockermgr_release_table")).Filter("Id", releaseTask.Id).RelatedSel("ReleaseConf").One(releaseTask)
	if err != nil {
		return err
	}
	_, err = o.QueryTable(beego.AppConfig.String("dockermgr_releaseconfidc_table")).Filter("ReleaseConf__Id", releaseTask.ReleaseConf.Id).All(&ReleaseConfIdcs)
	if err != nil {
		return err
	}
	var idcs []*IdcConf
	idcs, err = GetIdcs()
	if err != nil {
		return err
	}
	for _, ReleaseConfIdc := range ReleaseConfIdcs {
		for _, idc := range idcs {
			if idc.Id == ReleaseConfIdc.Idc.Id {
				releaseTask.ReleaseConf.ReleaseIdc = append(releaseTask.ReleaseConf.ReleaseIdc, idc)

			}
		}

	}
	var services []*Service
	services, err = GetServices()
	if err != nil {
		return err
	}
	for _, service := range services {
		if service.Id == releaseTask.Service.Id {
			releaseTask.Service = service
			break
		}
	}
	return nil
}
func QueryRelease(o orm.Ormer, service *Service) ([]*ReleaseTask, error) {
	var ReleaseTaskList []*ReleaseTask
	_, err := o.QueryTable(beego.AppConfig.String("dockermgr_release_table")).Filter("Service__Id", service.Id).RelatedSel("ReviewUser", "ReleaseUser", "OperationUser", "CancelUser", "ReleaseConf").All(&ReleaseTaskList)
	return ReleaseTaskList, err
}

func QueryReleaseConf(o orm.Ormer, service *Service) (ReleaseConf, error) {
	var err error
	var releaseConf ReleaseConf
	err = o.QueryTable(beego.AppConfig.String("dockermgr_releaseconf_table")).Filter("Service__Id", service.Id).OrderBy("-Id").One(&releaseConf)
	if err != nil {
		return releaseConf, err
	}
	_, err = o.LoadRelated(&releaseConf, "ReleaseIdc")
	return releaseConf, err
}

func CreateOrUpdateRelease(o orm.Ormer, releaseTask *ReleaseTask, updatecols ...string) (int64, error) {
	var err error

	if releaseTask.Id == 0 {
		//releaseTask.TaskStatus = NotReady
		if err = checkReleaseTask(releaseTask); err != nil {
			return 0, err
		}
		_, err = o.Insert(releaseTask)
		return 0, err
	} else {
		if len(updatecols) == 0 {
			_, err = o.Update(releaseTask)
		} else {
			_, err = o.Update(releaseTask, updatecols...)
		}
		return 0, err
	}
}

func UpdateRelease(o orm.Ormer, releaseTask *ReleaseTask, checkStatus int, params orm.Params) (int64, error) {
	var err error
	var num int64
	if releaseTask.Id == 0 {
		err = errors.New("invalid pk value")
		return 0, err
	}
	num, err = o.QueryTable(beego.AppConfig.String("dockermgr_release_table")).Filter("Service__Id", releaseTask.Service.Id).Filter("Id", releaseTask.Id).Filter("TaskStatus", checkStatus).Update(params)
	return num, err
}

func CreateReleaseConf(o orm.Ormer, releaseConf *ReleaseConf) error {
	var err error
	if err = checkReleaseConf(releaseConf); err != nil {
		return err
	}
	_, err = o.Insert(releaseConf)
	if err != nil {
		return err
	}
	m2m := o.QueryM2M(releaseConf, "ReleaseIdc")
	_, err = m2m.Add(releaseConf.ReleaseIdc)
	if err != nil {
		return err
	}
	return nil

}
