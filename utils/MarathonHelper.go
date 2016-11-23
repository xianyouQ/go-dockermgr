package utils

import (
    "github.com/astaxie/beego"
    outMarathon "github.com/xianyouQ/go-dockermgr/3rd/github.com/gambol99/go-marathon"
    "encoding/json"
    "net"
    "net/http"
    "crypto/tls"
    "time"
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
}

var (
    marathonClient outMarathon.Marathon 
    marathonURL string
)
func init() {
    marathonURL = beego.AppConfig.String("marathonUrl")
    config :=  outMarathon.NewDefaultConfig()
    config.URL = marathonURL 
    config.HTTPClient = &http.Client{
    Timeout: (time.Duration(10) * time.Second),
    Transport: &http.Transport{
        Dial: (&net.Dialer{
            Timeout:   10 * time.Second,
            KeepAlive: 10 * time.Second,
        }).Dial,
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
            },
        },
    }

    client,err := outMarathon.NewClient(config)
    if err != nil {
        //log.Fatalf("Failed to create a client for marathon, error: %s", err)
    }
    marathonClient = client
}

func CreateMarathonAppFromJson(conf string) (*outMarathon.Application,error) {
    MarathonApp := &outMarathon.Application{}
    err := json.Unmarshal([]byte(conf),&MarathonApp)
    if err != nil {
        return MarathonApp,err
    }

    return MarathonApp,nil
}

func  ListApplicationsFromGroup(name string) ([]*outMarathon.Application,error) {
    var Apps []*outMarathon.Application
    group,err := marathonClient.Group(name)
    if err != nil {
        return Apps,err
    }
    return group.Apps,nil
}

func  GetMesosInfo() (*MesosInfo,error) {
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
         return mesosInfo,nil
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