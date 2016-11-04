package utils

import (
    marathon "github.com/gambol99/go-marathon"
    "github.com/astaxie/beego"
    "encoding/json"
    "net"
    "net/http"
    "crypto/tls"
    "time"
)

type MarathonClient struct {
    Client marathon.Marathon 
}

var (
    marathonClient MarathonClient = MarathonClient{}
)
func init() {
    marathonURL := beego.AppConfig.String("marathonUrl")
    config :=  marathon.NewDefaultConfig()
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

    client,err := marathon.NewClient(config)
    if err != nil {
        //log.Fatalf("Failed to create a client for marathon, error: %s", err)
    }
    marathonClient.Client = client
}

func CreateMarathonAppFromJson(conf string) (*marathon.Application,error) {
    MarathonApp := &marathon.Application{}
    err := json.Unmarshal([]byte(conf),&MarathonApp)
    if err != nil {
        return MarathonApp,err
    }

    return MarathonApp,nil
}

func (slf *MarathonClient) ListApplicationsFromGroup(name string) ([]*marathon.Application,error) {
    var Apps []*marathon.Application
    group,err := slf.Client.Group(name)
    if err != nil {
        return Apps,err
    }
    return group.Apps,nil
}
