package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	m "github.com/xianyouQ/go-dockermgr/models"
	"github.com/xianyouQ/go-dockermgr/utils"
	outMarathon "github.com/xianyouQ/go-marathon"
)

type DockerController struct {
	CommonController
}

func (c *DockerController) DashBoard() {
	var err error
	var idcs []*m.IdcConf
	idcs, err = m.GetIdcs()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	for _, idc := range idcs {
		var client outMarathon.Marathon
		var mesosInfo *utils.MesosInfo
		idc.OthsData = make(map[string]interface{})
		if idc.MarathonSerConf == nil {
			idc.OthsData["status"] = false
			idc.OthsData["info"] = "Marathon Conf Not Set"
			continue
		}
		client, err = utils.NewMarathonClient(idc.MarathonSerConf.Server, idc.MarathonSerConf.HttpBasicAuthUser, idc.MarathonSerConf.HttpBasicPassword)
		if err != nil {
			idc.OthsData["status"] = false
			idc.OthsData["info"] = err.Error()
			continue
		}
		mesosInfo, err = utils.GetMesosInfo(client)
		if err != nil {
			idc.OthsData["status"] = false
			idc.OthsData["info"] = err.Error()
			continue
		}
		idc.OthsData["status"] = true
		idc.OthsData["mesos"] = mesosInfo
	}
	c.Rsp(true, "success", idcs)
}

func (c *DockerController) GetContainers() {
	var err error
	var serviceId, idcId int
	serviceId, _ = c.GetInt("serviceId")
	idcId, err = c.GetInt("idcId")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}

	var idcs []*m.IdcConf
	var queryIdc *m.IdcConf
	var services []*m.Service
	var queryService *m.Service
	idcs, err = m.GetIdcs()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	for _, idc := range idcs {
		if idc.Id == idcId {
			queryIdc = idc
		}
	}
	if queryIdc == nil {
		c.Rsp(false, "idc not found", nil)
		return
	}
	services, err = m.GetServices()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	for _, service := range services {
		if service.Id == serviceId {
			queryService = service
		}
	}
	if queryService == nil {
		c.Rsp(false, "service not found", nil)
		return
	}
	o := orm.NewOrm()
	var instances []*m.Ip
	instances, err = m.GetInstances(o, queryService, queryIdc)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	var client outMarathon.Marathon
	client, err = utils.NewMarathonClient(queryIdc.MarathonSerConf.Server, queryIdc.MarathonSerConf.HttpBasicAuthUser,
		queryIdc.MarathonSerConf.HttpBasicPassword)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	var marathonApps []*outMarathon.Application
	marathonApps, err = utils.ListApplicationsFromGroup(queryService.Code, client)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	for _, marathonApp := range marathonApps {
		for _, instance := range instances {
			marathonAppIpaddr := strings.Split(marathonApp.ID, "/")[2]
			if marathonAppIpaddr == instance.IpAddr {
				instance.MarathonData = marathonApp
			}
		}
	}
	c.Rsp(true, "success", instances)
}

func (c *DockerController) ScaleContainers() {
	var err error
	var scaleCount, serviceId, idcId int
	serviceId, _ = c.GetInt("serviceId")
	idcId, err = c.GetInt("idcId")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	scaleCount, err = c.GetInt("scaleCount")
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}

	var idcs []*m.IdcConf
	var queryIdc *m.IdcConf
	var services []*m.Service
	var queryService *m.Service
	idcs, err = m.GetIdcs()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	for _, idc := range idcs {
		if idc.Id == idcId {
			queryIdc = idc
		}
	}
	if queryIdc == nil {
		c.Rsp(false, "idc not found", nil)
		return
	}
	services, err = m.GetServices()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	for _, service := range services {
		if service.Id == serviceId {
			queryService = service
		}
	}
	if queryService == nil {
		c.Rsp(false, "service not found", nil)
		return
	}
	o := orm.NewOrm()
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	var containerCount int64
	containerCount, err = m.GetInstancesCount(o, queryService, queryIdc)
	if err != nil {
		c.Rsp(false, err.Error(), nil)

		return
	}
	diff := scaleCount - int(containerCount)
	if diff == 0 {
		c.Rsp(true, "success", nil)
	}
	var requestIp []*m.Ip
	var application *outMarathon.Application
	var client outMarathon.Marathon
	client, err = utils.NewMarathonClient(queryIdc.MarathonSerConf.Server,
		queryIdc.MarathonSerConf.HttpBasicAuthUser, queryIdc.MarathonSerConf.HttpBasicPassword)
	if err != nil {
		c.Rsp(false, err.Error(), nil)
		return
	}
	application, err = utils.CreateMarathonAppFromJson(queryService.MarathonConf)
	if err != nil {
		logs.GetLogger("dockerCtl").Printf("rollback error:%s", err.Error())
		c.Rsp(false, err.Error(), nil)
		return
	}

	err = o.Begin()
	if diff > 0 {
		requestIp, err = m.RequestIp(o, queryService, queryIdc, int(diff))
		if err != nil {
			c.Rsp(false, err.Error(), nil)
			err = o.Rollback()
			if err != nil {
				logs.GetLogger("dockerCtl").Printf("rollback error:%s", err.Error())
			}
			return
		}
		for idx, ip := range requestIp {
			application.ID = fmt.Sprintf("/%s/%s", queryService.Code, ip.IpAddr)
			//application.Container.Docker.EmptyParameters()
			if queryService.ReleaseVer == nil {

			} else {
				imageTag := fmt.Sprintf("%s:%s", queryService.Code, queryService.ReleaseVer.ImageTag)
				application.Container.Docker.Image = imageTag
			}
			application.Container.Docker.SetParameter("ip", ip.IpAddr)
			application.Container.Docker.SetParameter("mac-address", ip.MacAddr)
			//application.Container.Docker.AddParameter("hostname",XXXXX)
			_, err = client.CreateApplication(application)
			if err != nil {
				for iner := 0; iner <= idx; iner++ {
					applicationID := fmt.Sprintf("/%s/%s", queryService.Code, requestIp[iner].IpAddr)
					_, err := utils.DelApplication(client, applicationID)
					if err != nil {
						logs.GetLogger("dockerCtl").Printf("Stop Container failure:%s", applicationID)
					}
				}
				c.Rsp(false, err.Error(), nil)
				err = o.Rollback()
				if err != nil {
					logs.GetLogger("dockerCtl").Printf("rollback error:%s", err.Error())
				}
				return
			}
		}

	} else if diff < 0 {
		requestIp, err = m.RecycleIp(o, queryService, queryIdc, int(-diff))
		if err != nil {
			c.Rsp(false, err.Error(), nil)
			err = o.Rollback()
			if err != nil {
				logs.GetLogger("dockerCtl").Printf("rollback error:%s", err.Error())
			}
			return
		}
		for _, ip := range requestIp {
			applicationID := fmt.Sprintf("/%s/%s", queryService.Code, ip.IpAddr)
			_, err := utils.DelApplication(client, applicationID)
			if err != nil {
				logs.GetLogger("dockerCtl").Printf("Stop Container failure:%s", applicationID)
			}
		}

	}
	err = o.Commit()
	if err != nil {
		logs.GetLogger("dockerCtl").Printf("commit error:%s", err.Error())
	}
	c.Rsp(true, "success", requestIp)

}
