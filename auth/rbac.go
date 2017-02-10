package auth

import (
	"errors"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	m "github.com/xianyouQ/go-dockermgr/models"
	"github.com/xianyouQ/go-dockermgr/utils"
)

//check access and register user's nodes
func AccessRegister() {
	var Check = func(ctx *context.Context) {
		user_auth_type, _ := strconv.Atoi(beego.AppConfig.String("user_auth_type"))
		rbac_auth_gateway := beego.AppConfig.String("rbac_auth_gateway")
		rbac_auth_signup := beego.AppConfig.String("rbac_auth_signup")
		requestUrl := strings.ToLower(ctx.Request.RequestURI)
		var accesslist []*m.ServiceAuthUser
		var err error
		if requestUrl == rbac_auth_gateway || requestUrl == rbac_auth_signup {
			return
		}
		if user_auth_type > 0 {
			uinfo := ctx.Input.Session("userinfo")
			if uinfo == nil {
				//ctx.Redirect(302, rbac_auth_gateway)
				ctx.Output.SetStatus(401)
				ctx.Output.JSON(&map[string]interface{}{"status": false, "info": "未登录"}, true, false)
				return
			}
			//admin用户不用认证权限
			adminuser := beego.AppConfig.String("rbac_admin_user")
			if uinfo.(m.User).Username == adminuser {
				return
			}

			if user_auth_type == 1 {
				listbysession := ctx.Input.Session("accesslist")
				if listbysession != nil {
					accesslist = listbysession.([]*m.ServiceAuthUser)
				} else {
					accesslist, err = AccessList(uinfo.(m.User).Id)
					if err != nil {
						ctx.Output.SetStatus(403)
						ctx.Output.JSON(&map[string]interface{}{"status": false, "info": "获取权限信息失败"}, true, false)
					}
					ctx.Output.Session("accesslist", accesslist)
				}
				var ret bool
				ret, err = AccessDecision(ctx.Request.RequestURI, accesslist)
				if err != nil {
					logs.GetLogger("rbac").Printf("check auth fail,detail:%s", err.Error())
				}
				if !ret {
					ctx.Output.SetStatus(403)
					ctx.Output.JSON(&map[string]interface{}{"status": false, "info": "权限不足"}, true, false)
					return
				}
			}
		}
	}
	beego.InsertFilter("/api/*", beego.BeforeRouter, Check)
}

//To test whether permissions
func AccessDecision(url string, accesslist []*m.ServiceAuthUser) (bool, error) {
	baseUrl, params := ParseUrl(url)
	RoleNodes, err := m.GetRoleNodes()
	if err != nil {
		return false, err
	}
	for _, RoleNode := range RoleNodes {
		for _, Node := range RoleNode.Nodes {
			if Node.Url == baseUrl {
				for _, access := range accesslist {
					if RoleNode.Id == access.ServiceAuth.Role.Id {
						if RoleNode.NeedAddAuth == true {
							if serviceId, ok := params["serviceId"]; ok {
								serviceID, _ := strconv.Atoi(serviceId)
								if access.ServiceAuth.Service.Id == serviceID {
									return true, nil
								}
							}

						}
						if RoleNode.NeedAddAuth == false {
							return true, nil
						}
					}

				}
			}
		}
	}
	return false, nil

}

func ParseUrl(url string) (string, map[string]string) {
	urlSplits := strings.Split(url, "?")
	params := make(map[string]string)
	if len(urlSplits) == 1 {
		return urlSplits[0], params
	}
	paramSplit := strings.Split(urlSplits[1], "&")
	for _, param := range paramSplit {
		kv := strings.Split(param, "=")
		if len(kv) < 2 {
			continue
		}
		params[kv[0]] = kv[1]
	}
	return urlSplits[0], params
}

func AccessList(uid int64) ([]*m.ServiceAuthUser, error) {
	var err error
	var auths []*m.ServiceAuthUser
	auths, err = m.GetAuthList(uid)
	return auths, err
}

//check login
func CheckLogin(username string, password string) (user m.User, err error) {
	user = m.GetUserByUsername(username)
	if user.Id == 0 {
		return user, errors.New("用户不存在")
	}
	if user.Password != utils.Pwdhash(password) {
		return user, errors.New("密码错误")
	}
	return user, nil
}
