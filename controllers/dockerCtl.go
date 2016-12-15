package controllers

import (
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

