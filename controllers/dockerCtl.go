package controllers

import (
	m "github.com/xianyouQ/go-dockermgr/models"
    "github.com/xianyouQ/go-dockermgr/utils"
    outMarathon "github.com/xianyouQ/go-marathon"
    "encoding/json"
	"github.com/astaxie/beego/orm"
	"strings"
	"github.com/astaxie/beego/logs"
	"fmt"
)

type DockerController struct {
	CommonController
}



func (c *DockerController) DashBoard() {
    var err error
    var idcs []*m.IdcConf
	idcs,err = m.GetIdcs()
    if err != nil {
        c.Rsp(false,err.Error(),nil)
        return
    }
    for _,idc := range idcs {
        var client outMarathon.Marathon
        var mesosInfo *utils.MesosInfo
        idc.OthsData = make(map[string]interface{})
        if idc.MarathonSerConf == nil {
            idc.OthsData["status"] = false
            idc.OthsData["info"] = "Marathon Conf Not Set"
            continue
        }
       client,err = utils.NewMarathonClient(idc.MarathonSerConf.Server,idc.MarathonSerConf.HttpBasicAuthUser,idc.MarathonSerConf.HttpBasicPassword)
       if err != nil {
           idc.OthsData["status"] = false
           idc.OthsData["info"] = err.Error()
           continue
       }
       mesosInfo,err = utils.GetMesosInfo(client)
       if err != nil {
           idc.OthsData["status"] = false
           idc.OthsData["info"] = err.Error()
           continue
       }
       idc.OthsData["status"] = true
       idc.OthsData["mesos"] = mesosInfo
    }
    c.Rsp(true,"success",idcs)
}


type ServiceContainerForm struct {
    Service *m.Service
    Idc *m.IdcConf
    Scale int64
}

func (c *DockerController) GetContainers() {
	var err error
	mServiceContainerForm := ServiceContainerForm{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &mServiceContainerForm); err != nil {
		//handle error
		c.Rsp(false, err.Error(),nil)
		return
	}
    o := orm.NewOrm()
    var instances []*m.Ip
    instances,err = m.GetInstances(o,mServiceContainerForm.Service,mServiceContainerForm.Idc)
    if err != nil {
        c.Rsp(false,err.Error(),nil)
        return
    }
    var client outMarathon.Marathon
    client,err = utils.NewMarathonClient(mServiceContainerForm.Idc.MarathonSerConf.Server,mServiceContainerForm.Idc.MarathonSerConf.HttpBasicAuthUser,
    mServiceContainerForm.Idc.MarathonSerConf.HttpBasicPassword)
    if err != nil {
        c.Rsp(false,err.Error(),nil)
        return
    }
    var marathonApps []*outMarathon.Application
    marathonApps,err = utils.ListApplicationsFromGroup(mServiceContainerForm.Service.Name,client)
    if err != nil {
        c.Rsp(false,err.Error(),nil)
        return
    }
    for _,marathonApp := range marathonApps {
        for _,instance := range instances {
            marathonAppIpaddr := strings.Split(marathonApp.ID,"/")[1]
            if marathonAppIpaddr == instance.IpAddr {
                instance.MarathonData = marathonApp
            }
        }
    }
    c.Rsp(true,"success",instances)
}


func (c *DockerController) ScaleContainers() {
    var err error
    var containerCount int64
	mServiceContainerForm := ServiceContainerForm{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &mServiceContainerForm); err != nil {
		//handle error
		c.Rsp(false, err.Error(),nil)
		return
	}
    o := orm.NewOrm()
    
	if err != nil {
		c.Rsp(false, err.Error(),nil)
        return
	}
    containerCount,err = m.GetInstancesCount(o,mServiceContainerForm.Service,mServiceContainerForm.Idc)
    if err != nil {
        c.Rsp(false, err.Error(),nil)
        return
    }
    diff := mServiceContainerForm.Scale - containerCount
    var requestIp []*m.Ip
    var application *outMarathon.Application
    var client outMarathon.Marathon
    client,err = utils.NewMarathonClient(mServiceContainerForm.Idc.MarathonSerConf.Server,
    mServiceContainerForm.Idc.MarathonSerConf.HttpBasicAuthUser,mServiceContainerForm.Idc.MarathonSerConf.HttpBasicPassword)
    if err != nil {
        c.Rsp(false, err.Error(),nil)
        return
    }    
    application,err = utils.CreateMarathonAppFromJson(mServiceContainerForm.Service.MarathonConf)
    if err != nil {
        c.Rsp(false, err.Error(),nil)
        return
    }
    err = o.Begin()
    if diff > 0 {
        requestIp,err = m.RequestIp(o,mServiceContainerForm.Service,mServiceContainerForm.Idc,int(diff))
        if err != nil {
            c.Rsp(false, err.Error(),nil)
            err = o.Rollback()
            if err != nil {
                logs.GetLogger("AuthCtl").Printf("rollback error:%s",err.Error())
            }
            return
        }
        for _,ip := range requestIp {
            application.ID = fmt.Sprintf("/%s/%s",mServiceContainerForm.Service.Code,ip.IpAddr)
            application.Container.Docker.AddParameter("ip",ip.IpAddr)
            application.Container.Docker.AddParameter("mac-address",ip.MacAddr)
            //application.Container.Docker.AddParameter("hostname",XXXXX)
            _,err = utils.NewApplication(client,application)
            if err != nil {
                // 退回已经成功创建的容器
                c.Rsp(false, err.Error(),nil)
                err = o.Rollback()
                if err != nil {
                    logs.GetLogger("AuthCtl").Printf("rollback error:%s",err.Error())
                }
                return
            }
        }


    } else {
        requestIp,err = m.RecycleIp(o,mServiceContainerForm.Service,mServiceContainerForm.Idc,int(-diff))
        if err != nil {
            c.Rsp(false, err.Error(),nil)
            err = o.Rollback()
            if err != nil {
                logs.GetLogger("AuthCtl").Printf("rollback error:%s",err.Error())
            }
            return
        }


    }
    c.Rsp(true,"success",requestIp)

}


