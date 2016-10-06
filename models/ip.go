package models

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "github.com/xianyouQ/go-dockermgr/utils"
    "errors"
    "fmt"
    "net"
)

const (
    IpUsed = iota
    IpUnUsed
    IpAllocated
)
type Cidr struct {
    Id int `orm:"auto"`
    Net string `orm:"size(20);unique"`
    StartIp string `orm:"size(20)"`
    EndIp string `orm:"size(20)"`
    IpList []*Ip `orm:"reverse(many)"`

}

type Ip struct {
    Id int `orm:"auto"`
    BelongNet *Cidr `orm:"rel(fk)"`
    IpAddr string `orm:"unique;size(20)"`
    MacAddr string `orm:"unique;size(20)"`
    Status int 
    BelongService *Service `orm:"rel(fk)"`
}


//generate no ieee mac for container
func GetMacAddr(Ip net.IP,base string) string {
    Ip = Ip.To4()
    result := fmt.Sprintf("%s:%x:%x:%x",base,self.Ip[1],self.Ip[2],self.Ip[3])
    return result
}
//

func GetCidrFromOrm() []utils.CidrHelper {
    CidrList := make([]utils.CidrHelper,0,5)
    var Cidrs []Cidr
    o := orm.NewOrm()
    cidr := new(Cidr)
    qs := o.QueryTable(cidr).All(&Cidrs)
    for _,iter range Cidrs {
        mCidrHelper := utils.NewCidrfromString(iter.Net)
        CidrList = append(CidrList,mCidrHelper)
    }
    return CidrList
}

var (
    GlobalCidrList = GetCidrFromOrm()
    BaseMac = beego.AppConfig.String("basemacstring")
)


func AddCidr(net string,start string,end string) error {
    newCidr,err := utils.NewCidrwithStartEnd(net,start,end)
    if err != nil {

    }
    for _,iter range GlobalCidrList {
        if ok := iter.Overlaps(newCidr); ok {
            errorstring := fmt.Sprintf("new Cidr %s Overlaps with %s",net,iter)
            return errors.New(errorstring)
        }
    }
    GlobalCidrList = append(GlobalCidrList,newCidr)
    o := orm.NewOrm()
    mcidr = new(Cidr)
    mcidr.Net = newCidr.Net.String()
    mcidr.StartIp = newCidr.StartIp.String()
    mcidr.EndIp = newCidr.EndIp.String()
    id,err := o.Insert(&mcidr)
    if err != nil {

    }
    IpList := make([]Ip,0,125)
    for iter range newCidr.IpList() {
        newIp = new(Ip)
        newIp.BelongNet = mcidr
        newIp.IpAddr = iter.String()
        newIp.MacAddr = GetMacAddr(newIp,BaseMac)
        newIp.Status = IpUnUsed
        IpList = append(IpList,newIp)
    }
    num,err := o.InsertMulti(len(IpList),IpList)
    if err != nil {

    }

}


func RequestIp(service Service,cidr Cidr,num int) ([]Ip,error){
    var IpList []Ip
    ip := Ip{}
    o := orm.NewOrm()
    _:=o.Begin()
    qnum,err := o.QueryTable(ip).Filter("BelongNet__id",cidr.Id).Filter("Status",IpUnUsed).Limit(num).Update(orm.Params{
        "Status":IpAllocated
    })
    if err != nil {
        _ := o.Rollback()
        return IpList,err
    }
    if qnum < num {
        _ := o.Rollback()
        return IpList,errors.New("No enough ip")
    }
    o.Commit()
    return IpList,nil
}

func RecycleIp (IpList []Ip) error {
    
}




func init() {
    orm.RegisterModel(new(Cidr),new(Ip))
}