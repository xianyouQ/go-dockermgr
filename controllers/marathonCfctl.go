package controllers

import (
	"github.com/xianyouQ/go-dockermgr/models"

	"encoding/json"
)

type MarathonCfController struct {
	CommonController
}



func (c *MarathonCfController) AddOrUpdateMarathonConf() {
	var err error
	var id int64
	belongIdc := models.IdcConf{}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &belongIdc);err != nil {
		c.Rsp(false, err.Error(),nil)
		return
	}
	id,err = models.AddOrUpdateMarathonSerConf(belongIdc.MarathonSerConf)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
	}
	if id != 0  {
		belongIdc.MarathonSerConf.Id = int(id)
	}
	err = models.AddOrUpdateIdc(&belongIdc)
	if err !=nil {
		c.Rsp(false,err.Error(),nil)
	}
	c.Rsp(true,"success",belongIdc.MarathonSerConf)
}

/*
{
  "id": "/youxianqin",
  "apps": [
  {
    "id": "/foo",
    "cmd": "while [ true ] ; do echo 'Hello Marathon' ; sleep 5 ; done",
    "args": null,
    "user": null,
    "instances": 1,
    "cpus": 1,
    "mem": 128,
    "disk": 1024,
    "executor": "",
    "constraints": [],
    "uris": [],
    "fetch": [],
    "storeUrls": [],
	"requirePorts": false,
    "backoffSeconds": 1,
    "backoffFactor": 1.15,
    "maxLaunchDelaySeconds": 3600,
    "container": {
      "type": "DOCKER",
      "volumes": [],
      "docker": {
        "image": "centos:net",
        "network": "BRIDGE",
        "privileged": false,
        "parameters": [
          {
            "key": "ip",
            "value": "10.208.177.152"
          },
          {
            "key": "net",
            "value": "dockerbr0"
          },
          {
            "key": "name",
            "value": "test-marathon"
          }
        ],
        "forcePullImage": false
      }
    },
    "healthChecks": [],
    "readinessChecks": [],
    "dependencies": [],
    "upgradeStrategy": {
      "minimumHealthCapacity": 1,
      "maximumOverCapacity": 1
    },
    "labels": {},
    "acceptedResourceRoles": null,
    "ipAddress": null,
    "residency": null
	}
  ],
  "groups": [],
  "dependencies": []
}

*/