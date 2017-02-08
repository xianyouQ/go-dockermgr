package models

import (
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Service struct {
	Id           int            `orm:"auto"`
	Name         string         `orm:"size(20);unique"`
	Code         string         `orm:"size(20);unique"`
	Instances    []*Ip          `orm:"reverse(many)"`
	ReleaseVer   *ReleaseTask   `orm:"null;rel(fk)"`
	ReleaseTask  []*ReleaseTask `orm:"reverse(many)"`
	Roles        []*Role        `orm:"reverse(many)"`
	MarathonConf string         `orm:"type(text)"`
	//UsedReleaseConf *ReleaseConf `orm:"null;rel(fk)"`
}

func (this *Service) TableName() string {
	return beego.AppConfig.String("dockermgr_service_table")
}

func init() {
	orm.RegisterModel(new(Service))
}

func checkService(newService *Service) error {
	separateCount, err := beego.AppConfig.Int("service_separate_count")
	if err != nil {
		return err
	}
	valid := validation.Validation{}
	b, _ := valid.Valid(&newService)
	if !b {
		for _, err1 := range valid.Errors {
			return errors.New(err1.Message)
		}
	}
	codeSplits := strings.Split(newService.Code, "-")
	if len(codeSplits) != separateCount {
		return errors.New("invaild service code")
	}
	return nil

}
func AddOrUpdateService(o orm.Ormer, newService *Service, updatecols ...string) error {
	var err error
	var roles []*Role
	var pid int64
	err = checkService(newService)
	if err != nil {
		return err
	}
	if newService.Id == 0 {
		pid, err = o.Insert(newService)
		if err != nil {
			return err
		}

		roles, err = GetRoleNodes()
		if err != nil {
			return err
		}
		for _, role := range roles {
			if role.NeedAddAuth {
				_, err = NewServiceAuth(o, role, newService)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if len(updatecols) == 0 {
			_, err = o.Update(newService)
		} else {
			_, err = o.Update(newService, updatecols...)
		}
		if err != nil {
			return err
		}
	}
	if pid != 0 {
		UpdateServices(newService, true)
	}
	return err
}

func GetServicesFromOrm() ([]*Service, error) {
	var err error
	var Services []*Service
	o := orm.NewOrm()
	_, err = o.QueryTable(beego.AppConfig.String("dockermgr_service_table")).RelatedSel().All(&Services)
	if err != nil {
		return Services, err
	}
	return Services, nil
}

func DelService(o orm.Ormer, oldService *Service) error {
	var err error
	var count int64
	count, err = o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongService", oldService.Id).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Instance is not null")
	}
	_, err = o.Delete(oldService)
	if err != nil {
		return err
	}
	return nil
}

func GetInstances(o orm.Ormer, service *Service, idc *IdcConf) ([]*Ip, error) {
	var Ips []*Ip
	if len(idc.Cidrs) == 0 {
		return Ips, errors.New("no Cidr in this idc")
	}
	_, err := o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongService", service.Id).Filter("BelongNet__in", idc.Cidrs).All(&Ips)
	return Ips, err
}

func GetInstancesCount(o orm.Ormer, service *Service, idc *IdcConf) (int64, error) {
	if len(idc.Cidrs) == 0 {
		return 0, errors.New("no Cidr in this idc")
	}
	count, err := o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongService", service.Id).Filter("BelongNet__in", idc.Cidrs).Count()
	return count, err
}
