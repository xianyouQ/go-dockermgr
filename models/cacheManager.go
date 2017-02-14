package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"

	"github.com/astaxie/beego/logs"
)

var (
	bm cache.Cache
)

func init() {
	cacheType := beego.AppConfig.String("cache_type")
	cacheConfig := beego.AppConfig.String("cache_config")
	var err error
	bm, err = cache.NewCache(cacheType, cacheConfig)
	if err != nil {
		logs.Critical("initialize cache error occur:%s", err.Error())
		panic("initialize cache error")
	}
}

func GetIdcs() ([]*IdcConf, error) {
	var err error
	var idcs []*IdcConf
	if bm.IsExist("idcs") {
		return bm.Get("idcs").([]*IdcConf), nil
	} else {
		idcs, err = getIdcsfromOrm()
		if err != nil {
			return idcs, err
		}
		err = bm.Put("idcs", idcs, 600*time.Second)
		if err != nil {
			return idcs, err
		}
	}

	return idcs, nil
}

func DeleteCache(cache string) {
	bm.Delete(cache)
}
func GetRoleNodes() ([]*Role, error) {
	var err error
	var roles []*Role
	if bm.IsExist("roles") {
		return bm.Get("roles").([]*Role), nil
	} else {
		roles, err = GetRoleListFromOrm()
		if err != nil {
			return roles, err
		}
		err = bm.Put("roles", roles, 600*time.Second)
	}
	return roles, err
}

func GetServices() ([]*Service, error) {
	var err error
	var services []*Service
	if bm.IsExist("services") {
		return bm.Get("services").([]*Service), nil
	} else {
		services, err = GetServicesFromOrm()
		if err != nil {
			return services, err
		}
		err = bm.Put("services", services, 600*time.Second)
	}
	return services, err
}
