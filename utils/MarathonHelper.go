package utils

import (
    marathon "github.com/gambol99/go-marathon"
    "github.com/astaxie/beego"
    "encoding/json"
    "net"
    "net/http"
    "crypto/tls"
    "time"
	"github.com/rogpeppe/godef/vendor/9fans.net/go/plan9/client"
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

func (slf *MarathonClient) ListApplications() error {
    applications,err := slf.Client.Applications()
    if err != nil {
        return err
    }
    //for _, application := range applications.Apps {
    //
    //}

}

func (slf *MarathonClient) ListGroups() {

}