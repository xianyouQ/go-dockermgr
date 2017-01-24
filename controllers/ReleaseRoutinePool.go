package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/xianyouQ/go-dockermgr/models"
	"github.com/xianyouQ/go-dockermgr/utils"
	outMarathon "github.com/xianyouQ/go-marathon"
)

const (
	InstanceReady = iota
	InstanceStopping
	InstanceStarting
)

type ContainerRes struct {
	Instance   *models.Ip
	ReleaseMsg string
}

type ReleaseRoutine struct {
	ReleaseTask        *models.ReleaseTask
	ContainerResChan   []*ContainerRes
	Abandon            chan struct{}
	Done               bool
	FaultOutOfTolerant chan struct{}
	ErrorMsg           string
}
type IdcReleaseInstance struct {
	idcCode          string
	client           outMarathon.Marathon
	runningInstances []*ReleaseInstace
	done             bool
	instances        chan *models.Ip
}
type ReleaseInstace struct {
	InstanceIp *models.Ip
	Status     int
	LastCheck  int64
}

var (
	TaskChannel         chan *ReleaseRoutine
	RoutinePoolSize     int
	Stop                chan struct{}
	ReleaseRoutineTasks []ReleaseRoutine
)

func init() {
	TaskChannelSize, err := beego.AppConfig.Int("task_channel_maxsize")
	if err != nil {
		panic("task_channel_maxsize must be integer")

	}
	TaskChannel = make(chan *ReleaseRoutine, TaskChannelSize)
	ReleaseRoutineTasks = make([]ReleaseRoutine, TaskChannelSize)
	Stop = make(chan struct{})
	RoutinePoolSize, err = beego.AppConfig.Int("routinepool_size")
	if err != nil {
		panic("task_channel_maxsize must be  integer")
	}
	StartReleaseRoutinePool()
}

func idcTaskMonitor(clients <-chan *IdcReleaseInstance, task *ReleaseRoutine) {
	var err error
	var errCount int
	IdcReleaseInstances := make([]*IdcReleaseInstance, 3, 3)
	for {
		if len(IdcReleaseInstances) < task.ReleaseTask.ReleaseConf.IdcParalle {
			select {
			case newIdcRelease := <-clients:
				IdcReleaseInstances = append(IdcReleaseInstances, newIdcRelease)
			case <-time.After(time.Second):
				if task.Done == true && len(IdcReleaseInstances) == 0 {
					return
				}
			}
		}

		for outindex := 0; outindex < len(IdcReleaseInstances); {
			idc := IdcReleaseInstances[outindex]
			for index := 0; index < len(idc.runningInstances); {
				instance := idc.runningInstances[index]
				applicationString := fmt.Sprintf("/%s/%s", task.ReleaseTask.ReleaseConf.Service.Code, instance.InstanceIp.IpAddr)
				if instance.Status == InstanceReady {
					_, err = idc.client.DeleteApplication(applicationString, true)
					if err != nil {
						for _, mContainerRes := range task.ContainerResChan {
							if mContainerRes.Instance == instance.InstanceIp { //可能有坑
								mContainerRes.ReleaseMsg = err.Error()
							}
						}
						errCount = errCount + 1
						if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
							task.ErrorMsg = "错误过多"
							task.FaultOutOfTolerant <- struct{}{}
							return
						} else {
							idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
							continue
						}
					}
					for _, mContainerRes := range task.ContainerResChan {
						if mContainerRes.Instance == instance.InstanceIp { //可能有坑
							mContainerRes.ReleaseMsg = "正在关闭容器"
						}
					}
					instance.Status = InstanceStopping
					instance.LastCheck = time.Now().Unix()
					index = index + 1
				}
				if instance.Status == InstanceStopping {
					var ok bool
					ok, err = utils.IsExistApplication(idc.client, applicationString)
					if err != nil {
						for _, mContainerRes := range task.ContainerResChan {
							if mContainerRes.Instance == instance.InstanceIp { //可能有坑
								mContainerRes.ReleaseMsg = err.Error()
							}
						}
						errCount = errCount + 1
						if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
							task.ErrorMsg = "错误过多"
							task.FaultOutOfTolerant <- struct{}{}
							return
						} else {
							idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
							continue
						}
					}
					if !ok {
						var mApplication *outMarathon.Application
						mApplication, err = utils.CreateMarathonAppFromJson(task.ReleaseTask.ReleaseConf.Service.MarathonConf)
						if err != nil {
							for _, mContainerRes := range task.ContainerResChan {
								if mContainerRes.Instance == instance.InstanceIp { //可能有坑
									mContainerRes.ReleaseMsg = err.Error()
								}
							}
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								task.FaultOutOfTolerant <- struct{}{}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						imageTag := fmt.Sprintf("%s:%s", task.ReleaseTask.ReleaseConf.Service.Code, task.ReleaseTask.ImageTag)
						mApplication.Container.Docker.SetParameter("image", imageTag)
						_, err = idc.client.CreateApplication(mApplication)
						if err != nil {
							for _, mContainerRes := range task.ContainerResChan {
								if mContainerRes.Instance == instance.InstanceIp { //可能有坑
									mContainerRes.ReleaseMsg = err.Error()
								}
							}
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								task.FaultOutOfTolerant <- struct{}{}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						for _, mContainerRes := range task.ContainerResChan {
							if mContainerRes.Instance == instance.InstanceIp { //可能有坑
								mContainerRes.ReleaseMsg = "正在启动容器"
							}
						}
						instance.Status = InstanceStarting
						instance.LastCheck = time.Now().Unix()
						index = index + 1
					} else {
						nowTime := time.Now().Unix()
						if nowTime-instance.LastCheck > task.ReleaseTask.ReleaseConf.TimeOut {
							for _, mContainerRes := range task.ContainerResChan {
								if mContainerRes.Instance == instance.InstanceIp { //可能有坑
									mContainerRes.ReleaseMsg = "等待关闭容器时超时"
								}
							}
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								task.FaultOutOfTolerant <- struct{}{}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						index = index + 1
					}
				}
				if instance.Status == InstanceStarting {
					var ok bool
					ok, err = idc.client.ApplicationOK(applicationString)
					if err != nil {
						for _, mContainerRes := range task.ContainerResChan {
							if mContainerRes.Instance == instance.InstanceIp { //可能有坑
								mContainerRes.ReleaseMsg = err.Error()
							}
						}
						errCount = errCount + 1
						if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
							task.ErrorMsg = "错误过多"
							task.FaultOutOfTolerant <- struct{}{}
							return
						} else {
							idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
							continue
						}
					}
					if ok {
						for _, mContainerRes := range task.ContainerResChan {
							if mContainerRes.Instance == instance.InstanceIp { //可能有坑
								mContainerRes.ReleaseMsg = "容器启动成功"
							}
						}
						idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
						continue
					} else {
						nowTime := time.Now().Unix()
						if nowTime-instance.LastCheck > task.ReleaseTask.ReleaseConf.TimeOut {
							for _, mContainerRes := range task.ContainerResChan {
								if mContainerRes.Instance == instance.InstanceIp { //可能有坑
									mContainerRes.ReleaseMsg = "等待启动容器时超时"
								}
							}
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								task.FaultOutOfTolerant <- struct{}{}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						index = index + 1
					}
				}
			}
			for count := len(idc.runningInstances); count <= task.ReleaseTask.ReleaseConf.IdcInnerParalle; {
				select {
				case newInstance := <-idc.instances:
					newRunningInstance := &ReleaseInstace{}
					newRunningInstance.InstanceIp = newInstance
					newRunningInstance.Status = InstanceReady
					idc.runningInstances = append(idc.runningInstances, newRunningInstance)
					count = len(idc.runningInstances)
				case <-time.After(time.Second):
					if len(idc.runningInstances) == 0 && idc.done == true {
						IdcReleaseInstances = append(IdcReleaseInstances[:outindex], IdcReleaseInstances[outindex+1:]...)
					} else {
						break
					}

				}
			}

		}
		<-time.After(time.Second)
	}

}
func releaseTaskFunc(task *ReleaseRoutine) {
	var err error
	clientChannel := make(chan *IdcReleaseInstance, task.ReleaseTask.ReleaseConf.IdcParalle)
	clients := make([]*IdcReleaseInstance, 3, 3)
	ips := make(map[string][]*models.Ip)
	for _, idc := range task.ReleaseTask.ReleaseConf.ReleaseIdc {
		var client outMarathon.Marathon
		client, err = utils.NewMarathonClient(idc.MarathonSerConf.Server, idc.MarathonSerConf.HttpBasicAuthUser, idc.MarathonSerConf.HttpBasicPassword)
		if err != nil {
			task.ErrorMsg = err.Error()
			return
		}
		mIdcReleaseInstance := &IdcReleaseInstance{}
		mIdcReleaseInstance.client = client
		mIdcReleaseInstance.idcCode = idc.IdcCode
		mIdcReleaseInstance.instances = make(chan *models.Ip, task.ReleaseTask.ReleaseConf.IdcInnerParalle)
		clients = append(clients, mIdcReleaseInstance)
		o := orm.NewOrm()
		var iplist []*models.Ip
		iplist, err = models.GetInstances(o, task.ReleaseTask.ReleaseConf.Service, idc)
		if err != nil {
			task.ErrorMsg = err.Error()
			return
		}
		ips[idc.IdcCode] = iplist
		for _, ip := range iplist {
			mContainerRes := &ContainerRes{}
			mContainerRes.Instance = ip
			task.ContainerResChan = append(task.ContainerResChan, mContainerRes)
		}
	}
	go idcTaskMonitor(clientChannel, task)
	count := 0
	for {
		alldone := true
		for _, client := range clients {
			if instances, ok := ips[client.idcCode]; ok {
				var endIndex int
				for index, instance := range instances {
					select {
					case client.instances <- instance:
						continue
					case <-time.After(time.Second):
						endIndex = index
						break
					case <-task.FaultOutOfTolerant:
						return
					}
				}
				instances = instances[endIndex:]
				if len(instances) == 0 {
					client.done = true
				} else {
					alldone = false
				}
			} else {
				client.done = true
			}

		}
		if count >= len(clients) && alldone == true {
			task.Done = true
			return
		}
		if count >= len(clients) {
			continue
		}
		select {
		case clientChannel <- clients[count]:
			count = count + 1
		case <-task.Abandon:
			for _, client := range clients {
				if client.done == false {
					client.done = true
				}
			}
			return
		case <-task.FaultOutOfTolerant:
			return
		case <-time.After(time.Second):
			continue

		}

	}

}

func AddTask(mReleaseTask *models.ReleaseTask) error {
	task := ReleaseRoutine{}
	task.ReleaseTask = mReleaseTask
	//...
	select {
	case TaskChannel <- &task:
		ReleaseRoutineTasks = append(ReleaseRoutineTasks, task)
		return nil
	default:
		return errors.New("Release TaskChannel is full")
	}
}

func CheckTaskStatus(mReleaseTask *models.ReleaseTask) ([]*ContainerRes, error) {
	for _, mReleaseRoutine := range ReleaseRoutineTasks {
		if mReleaseRoutine.ReleaseTask.Id == mReleaseTask.Id {
			return mReleaseRoutine.ContainerResChan, nil
		}
	}
	return nil, errors.New("task no found")
}

func AbandonTask(mReleaseTask *models.ReleaseTask) error {
	for _, mReleaseRoutine := range ReleaseRoutineTasks {
		if mReleaseRoutine.ReleaseTask.Id == mReleaseTask.Id {
			close(mReleaseRoutine.Abandon)
			return nil
		}
	}
	return errors.New("task no found")
}

func StartReleaseRoutinePool() {
	for i := 0; i < RoutinePoolSize; i++ {
		go func() {
			for {
				select {
				case task := <-TaskChannel:
					releaseTaskFunc(task)
				case <-Stop:
					return
				}
			}
		}()
	}
}

func StopReleaseRoutinePool() {
	close(Stop)
}
