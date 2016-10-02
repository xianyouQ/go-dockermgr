package utils


/*
deprecated,use github.com/heroku/docker-registry-client instead
*/
import (
    "net/http"
    "io/ioutil"
    "fmt"
    "errors"
    "encoding/base64"
)

type RegistryClient struct {
    Server ServerInfo
    TokenAuthServer ServerInfo
    TokenAuthService string
    TokenMap map[string]string
    Username string
    Password string
}

func (self RegistryClient) ErrorHandler(err error) {

}
func (self RegistryClient) GeneralBase64AuthString() string {
    authentication := fmt.Sprintf("%s:%s",self.Username,self.Password)
    return base64.URLEncoding.EncodeToString([]byte(authentication))
}
func (self *RegistryClient) GetToken(scope string,regeneral bool) (string,error) {
    var tokenstr string
    if Token,ok := self.TokenMap[scope]; ok && !regeneral {
        return Token,nil
    } else {
        httpClient,_ := self.TokenAuthServer.NewHttpClient()
        queryUrl := fmt.Sprintf("/auth?account=%s&scope=%s&service=%s",self.Username,scope,self.TokenAuthService)
        httpRequest,_ := self.Server.NewRequest("GET",queryUrl)
        base64String := fmt.Sprintf("Basic %s",self.GeneralBase64AuthString())
        httpRequest.Header.Add("Authorization",base64String)
        response,err := httpClient.Do(httpRequest)
        self.ErrorHandler(err)
        defer response.Body.Close()
        body,err := ioutil.ReadAll(response.Body)
        self.ErrorHandler(err)
        bodymap,err := Json2Map(body)
        self.ErrorHandler(err)
        if token,ok := bodymap["token"]; ok {
            tokenstr = token.(string)
            self.TokenMap[scope] = tokenstr
            return tokenstr,nil
        } else {
            err = errors.New("token header no found")
            return tokenstr,err
        }

    }

}
func (self *RegistryClient) NewRequest(method string,url string,scope string,headers map[string]string) (*http.Response,error) {
    var response *http.Response
    token,_ := self.GetToken(scope,false)
    Authorization := fmt.Sprintf("Bearer %s",token)
    httpClient,_ := self.TokenAuthServer.NewHttpClient()
    httpRequest,_ := self.Server.NewRequest("GET",url)
    httpRequest.Header.Add("Authorization",Authorization)
    for k,v := range headers {
        httpRequest.Header.Add(k,v)
    }
    response,_ = httpClient.Do(httpRequest)
    if response.StatusCode == http.StatusUnauthorized {
        token,_ := self.GetToken(scope,true)
        Authorization := fmt.Sprintf("Bearer %s",token)
        httpRequest.Header.Set("Authorization",Authorization)
        response,_ = httpClient.Do(httpRequest)
        if response.StatusCode == http.StatusUnauthorized {
            return response,errors.New("token error")
        }
    }
    return response,nil
}
func (self RegistryClient) GetCatalog() []byte{
    url := "/v2/_catalog"
    response,_ := self.NewRequest("GET",url,"registry:catalog:*",map[string]string{})
    defer response.Body.Close()
    body,_ := ioutil.ReadAll(response.Body)
    return  body
}