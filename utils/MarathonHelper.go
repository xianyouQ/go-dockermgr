package utils

import (
    outMarathon "github.com/xianyouQ/go-marathon"
    "encoding/json"
    "net/http"
	"fmt"
	"io/ioutil"
    "strings"
)


type MesosInfo struct {
    CpuTotal float32 `json:"master/cpus_total"`
    CpuUserd float32 `json:"master/cpus_used"`
    CpuIdle float32 `json:"-"`
    CpuPercent float32 `json:"master/cpus_percent"`
    MemTotal float32 `json:"master/mem_total"`
    MemUserd float32 `json:"master/mem_used"`
    MemIdle float32 `json:"-"`
    MemPercents float32 `json:"master/mem_percent"`
    DiskTotal float32 `json:"master/disk_total"`
    DiskUserd float32 `json:"master/disk_used"`
    DiskIdle float32 `json:"-"`
    DiskPercent float32 `json:"master/disk_percent"`
    SlaveActive float32 `json:"master/slaves_active"`
    SlaveInActive float32 `json:"master/slaves_inactive"`
    TaskRunning float32 `json:"master/tasks_running"`
    TaskLost float32 `json:"master/tasks_lost"`
}




func CreateMarathonAppFromJson(conf string) (*outMarathon.Application,error) {
    MarathonApp := &outMarathon.Application{}
    err := json.Unmarshal([]byte(conf),&MarathonApp)
    if err != nil {
        return MarathonApp,err
    }

    return MarathonApp,nil
}

func  ListApplicationsFromGroup(name string,marathonClient outMarathon.Marathon) ([]*outMarathon.Application,error) {
    var Apps []*outMarathon.Application
    group,err := marathonClient.Group(name)
    if err != nil {
        return Apps,err
    }
    return group.Apps,nil
}

func  GetMesosInfo(marathonClient outMarathon.Marathon) (*MesosInfo,error) {
     mesosInfo := new(MesosInfo)
     marathonInfo,err := marathonClient.Info()
     if err !=nil {
         return mesosInfo,err
     }
     var api string
     if strings.HasSuffix(marathonInfo.MarathonConfig.MesosLeaderUrl,"/") {
         api = "metrics/snapshot"
     } else {
         api = "/metrics/snapshot"
     }
     mesosMetricsUrl := fmt.Sprintf("%s%s",marathonInfo.MarathonConfig.MesosLeaderUrl,api)
     resp,err := http.Get(mesosMetricsUrl)
     if err != nil {
         return mesosInfo,err
     }
     defer resp.Body.Close()
     body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return mesosInfo,err
    }
    err = json.Unmarshal(body, mesosInfo)
    if err != nil {
        return mesosInfo,err
    }
    mesosInfo.CpuIdle = mesosInfo.CpuTotal - mesosInfo.CpuUserd
    mesosInfo.MemIdle = mesosInfo.MemTotal - mesosInfo.MemUserd
    mesosInfo.DiskIdle = mesosInfo.DiskTotal - mesosInfo.DiskUserd
    return mesosInfo,nil
}

func NewMarathonClient(url string,user string,passwd string) (outMarathon.Marathon,error){
    config :=  outMarathon.NewDefaultConfig()
    config.URL = url
    config.HTTPBasicAuthUser = user
    config.HTTPBasicPassword = passwd
    client,err := outMarathon.NewClient(config)
    return client,err
}


func NewApplication(marathonClient outMarathon.Marathon,application *outMarathon.Application) (*outMarathon.Application,error){
    newApplication, err := marathonClient.CreateApplication(application); 
    return newApplication,err
    
}
