package controllers

import (
	m "github.com/xianyouQ/go-dockermgr/models"
    "github.com/xianyouQ/go-dockermgr/utils"
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
        utils.NewMarathonClient(idc)
    }
}

