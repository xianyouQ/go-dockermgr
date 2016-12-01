package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
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
		var accesslist []string
		var err error
		if requestUrl==rbac_auth_gateway || requestUrl==rbac_auth_signup {
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
						accesslist = listbysession.([]string)
					} else {
						accesslist,err = AccessList(uinfo.(m.User).Id)
						if err != nil {
							ctx.Output.SetStatus(403)
							ctx.Output.JSON(&map[string]interface{}{"status": false, "info": "获取权限信息失败"}, true, false)
						}
						ctx.Output.Session("accesslist",accesslist)
					}
				ret := AccessDecision(ctx.Request.RequestURI, accesslist)
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
func AccessDecision(url string,accesslist []string) bool {
	for _,access := range accesslist {
		if access == url {
			return true
		}
	}
	return false
}



func AccessList(uid int64) ([]string,error){
	var err error
	var auths []*m.ServiceAuth
	var roles []*m.Role
	nodes := make([]string,5,5)
	auths,err = m.GetAuthList(uid)
	if err != nil {
		return nodes,err
	}
	roles,err = m.GetRoleNodes()
	if err != nil {
		return nodes,err
	}
	for _,auth := range auths {
		for _,role := range roles {
			if auth.Role.Id == role.Id {
				if auth.Service == nil {
					for _,node := range role.Node {
						nodes = append(nodes,node.Name)
					}
				} else {
					for _,node := range role.Node {
						nodename := fmt.Sprintf("%s/%s",auth.Service.Name,node.Name)
						nodes = append(nodes,nodename)
					}	
					
				}
	
			} 
		}
	}
	return nodes,err
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
