package models

import (
    "github.com/astaxie/beego/cache"
	"github.com/astaxie/beego"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
    bm cache.Cache
)

func init() {
    cacheType := beego.AppConfig.String("cache_type")
    cacheConfig := beego.AppConfig.String("cache_config")
    var err error
    bm,err = cache.NewCache(cacheType,cacheConfig)
    if err != nil {
        logs.Critical("initialize cache error occur:%s",err.Error())
    }
}


func GetIdcs() ([]*IdcConf,error) {
    var err error
    var idcs []*IdcConf
    if bm.IsExist("idcs") {
        return bm.Get("idcs").([]*IdcConf),nil
    } else {
        idcs,err = getIdcsfromOrm()
       if err!=nil {
           return idcs,err
       }
       err = bm.Put("idcs",idcs,600*time.Second)
       if err != nil {
           return idcs,err
       }
    }

    return idcs,nil
}


func GetRoleNodes() ([]*Role,error){
    var err error
    var roles []*Role
    if bm.IsExist("roles") {
        return bm.Get("roles").([]*Role),nil
    } else {
        roles,err = GetRoleListFromOrm()
        if err != nil {
            return roles,err
        }
        err = bm.Put("roles",roles,600*time.Second)
    }
    return roles,err
}
func UpdateRoleNodes(Updaterole *Role,isnew bool) {
    if !bm.IsExist("roles") {
        return
    }
    roles := bm.Get("roles").([]*Role)
    if isnew {
        roles = append(roles,Updaterole)
    } else {
        for inx,role := range roles{
            if role.Id == Updaterole.Id {
                roles[inx] = Updaterole
            }
        }
    }
    bm.Put("roles",roles,3600*time.Second)
}

func UpdateIdcs(UpdateIdc *IdcConf,isnew bool) {
    if !bm.IsExist("idcs") {
        return
    }
    idcs := bm.Get("idcs").([]*IdcConf)
    if isnew {
        idcs = append(idcs,UpdateIdc)
    } else {
        for inx,idc := range idcs{
            if idc.Id == UpdateIdc.Id {
                idcs[inx] = UpdateIdc
            }
        }
    }
    bm.Put("idcs",idcs,3600*time.Second)
}