package utils


import (
    "net/http"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "fmt"
)

type ServerInfo struct {
    Host string
    Schema string
    Port int
    Certfile string
    Client *http.Client
}


func (self *ServerInfo) NewHttpClient() (*http.Client,error) {
    if self.Client != nil {
        return self.Client,nil
    }
    if self.Schema == "https" {
        pool := x509.NewCertPool()
        caCrt, err := ioutil.ReadFile(self.Certfile)
        if err != nil {
            return nil,err
        }
        pool.AppendCertsFromPEM(caCrt)
        tr := &http.Transport{
        TLSClientConfig: &tls.Config{RootCAs: pool},
        }
        self.Client = &http.Client{Transport: tr}
        return self.Client,nil
    } else if self.Schema == "http" {
        self.Client = &http.Client{}
        return self.Client,nil
    } else {
        return nil,nil
    }
   
} 
func (self ServerInfo) NewRequest(method string,url string) (*http.Request,error) {
    requestUrl := fmt.Sprintf("%s://%s%s",self.Schema,self.Host,url)
    request,err := http.NewRequest(method,requestUrl,nil)
    return request,err
}
