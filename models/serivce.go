package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/xianyouQ/go-dockermgr/utils"
	"github.com/astaxie/beego/validation"
    "strings"
    "errors"

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


func checkService(newService *Service) error {
    separator := beego.AppConfig.String("service_separator")
    separateCount,err := beego.AppConfig.Int("service_separate_count")
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
    codeSplits := strings.Split(newService.Code,separator)
    if len(codeSplits) != separateCount {
        return errors.New("invaild service code")
    }
	return nil

}
func AddService(newService *Service) (int64,error) {
    var err error
    var id int64
    err = checkService(newService)
    if err != nil {
        return id,err
    }
    o := orm.NewOrm()
    id,err = o.Insert(newService)
    if err!=nil {
        return id,err
    }
    return id,nil
}


func QueryService() ([]*Service,error) {
    var err error
    var Services []*Service
    o := orm.NewOrm()
    _,err = o.QueryTable(beego.AppConfig.String("dockermgr_service_table")).All(&Services)
    if err !=nil {
        return Services,err
    }
    return  Services,nil
}

func DelService(oldService *Service) error {
    var err error
    var count int64
    o := orm.NewOrm()
    count,err = o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongService",oldService.Id).Count()
    if err != nil {
        return err
    }
    if count > 0 {
        return errors.New("Instance is not null")
    } 
    _,err = o.Delete(oldService)
    if err != nil {
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


/*

func (self Service) GetInstancesWithIdc(idc IdcConf) ([]*Ip,error) {
    o := orm.NewOrm()
    var Ips []*Ip
    _,err := o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongService",self.Id).Filter("BelongNet__in",IdcConf.Cidr).RelatedSel().All(&Ips)
    if err != nil {
        return Ips,err
    }
    return Ips,nil
}
*/