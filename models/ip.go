package models

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	outMarathon "github.com/gambol99/go-marathon"
	"github.com/xianyouQ/go-dockermgr/utils"
)

const (
	IpUsed = iota
	IpUnUsed
)

type Cidr struct {
	Id        int      `orm:"auto"`
	Net       string   `orm:"size(20);unique"`
	StartIp   string   `orm:"size(20)"`
	EndIp     string   `orm:"size(20)"`
	IpList    []*Ip    `orm:"reverse(many)"`
	BelongIdc *IdcConf `orm:"rel(fk)"`
}

type Ip struct {
	Id            int    `orm:"auto"`
	BelongNet     *Cidr  `orm:"rel(fk)"`
	IpAddr        string `orm:"unique;size(20)"`
	MacAddr       string `orm:"unique;size(20)"`
	Status        int
	BelongService *Service                 `orm:"null;rel(fk)"`
	MarathonData  *outMarathon.Application `orm:"-"`
}

var (
	BaseMac = beego.AppConfig.String("basemacstring")
)

func (this *Ip) TableName() string {
	return beego.AppConfig.String("dockermgr_ip_table")
}

func (this *Cidr) TableName() string {
	return beego.AppConfig.String("dockermgr_cidr_table")
}

func GetCidrFromOrm() ([]*Cidr, error) {
	var Cidrs []*Cidr
	o := orm.NewOrm()
	cidr := new(Cidr)
	_, err := o.QueryTable(cidr).All(&Cidrs)
	if err != nil {
		return Cidrs, err
	}
	return Cidrs, nil
}

func AddCidr(o orm.Ormer, cidr *Cidr) error {
	var err error
	var newCidr utils.CidrHelper
	var idcs []*IdcConf
	newCidr, err = utils.NewCidrwithStartEnd(cidr.Net, cidr.StartIp, cidr.EndIp)
	if err != nil {
		return err
	}
	idcs, err = GetIdcs()
	if err != nil {
		return err
	}
	for _, Idciter := range idcs {
		for _, CidrIter := range Idciter.Cidrs {
			iterCidrHelper, _ := utils.NewCidrfromString(CidrIter.Net)
			if ok := iterCidrHelper.Overlaps(newCidr); ok {
				errorstring := fmt.Sprintf("new Cidr %s Overlaps with %s", cidr.Net, CidrIter)
				return errors.New(errorstring)
			}
		}
	}

	_, err = o.Insert(cidr)
	if err != nil {
		return err
	}
	for _, idcConfIter := range idcs {
		if idcConfIter.Id == cidr.BelongIdc.Id {
			idcConfIter.Cidrs = append(idcConfIter.Cidrs, cidr)
		}
	}
	IpList := make([]Ip, 0, 125)
	for _, iter := range newCidr.IpList() {
		newIp := new(Ip)
		newIp.BelongNet = cidr
		newIp.IpAddr = iter.String()
		newIp.MacAddr = utils.GetMacAddr(iter, BaseMac)
		newIp.Status = IpUnUsed
		IpList = append(IpList, *newIp)
	}
	_, err = o.InsertMulti(len(IpList), IpList)
	if err != nil {
		return err
	}
	return nil

}

func RequestIp(o orm.Ormer, service *Service, idc *IdcConf, num int) ([]*Ip, error) {
	var IpList []*Ip
	if len(idc.Cidrs) == 0 {
		return IpList, errors.New("no Cidr in this idc")
	}
	qnum, err := o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongNet__id__in", idc.Cidrs).Filter("Status", IpUnUsed).Limit(num).All(&IpList)
	if err != nil {
		return IpList, err
	}
	if qnum < int64(num) {
		return IpList, errors.New("No enough ip")
	}
	qnum, err = o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("Id__in", IpList).Update(orm.Params{"Status": IpUsed, "BelongService": service.Id})
	return IpList, nil
}

func RecycleIp(o orm.Ormer, service *Service, idc *IdcConf, num int) ([]*Ip, error) {
	var IpList []*Ip
	if len(idc.Cidrs) == 0 {
		return IpList, errors.New("no Cidr in this idc")
	}
	qnum, err := o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongNet__Id__in", idc.Cidrs).Filter("Status", IpUsed).Filter("BelongService", service).Limit(num).All(&IpList)
	if qnum < int64(num) {
		return IpList, errors.New("No enough ip to recycle")
	}
	qnum, err = o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("Id__in", IpList).Update(orm.Params{
		"Status": IpUnUsed, "BelongService": nil})
	return IpList, err
}

func GetCidrCount(o orm.Ormer, queryIdc *IdcConf) (int64, error) {
	qnum, err := o.QueryTable(beego.AppConfig.String("dockermgr_cidr_table")).Filter("BelongIdc__Id", queryIdc.Id).Count()
	return qnum, err
}

func DelCidr(o orm.Ormer, delCidr *Cidr) error {
	var qnum int64
	var err error
	qnum, err = o.QueryTable(beego.AppConfig.String("dockermgr_ip_table")).Filter("BelongNet__id__in", delCidr.Id).Filter("Status", IpUsed).Count()
	if err != nil {
		return err
	}
	if qnum > 0 {
		return errors.New("some Container still running in this Cidr")
	}
	_, err = o.Delete(delCidr)
	return err
}
func init() {
	orm.RegisterModel(new(Cidr), new(Ip))
}
