package controllers

import (
	"errors"
	"fmt"
	"time"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	outMarathon "github.com/gambol99/go-marathon"
	"github.com/xianyouQ/go-dockermgr/models"
	"github.com/xianyouQ/go-dockermgr/utils"
)

const (
	InstanceReady = iota
	InstanceStopping
	InstanceStarting
	InstanceFinish
	InstanceFail
)

type ContainerRes struct {
	Instance   *models.Ip
	ReleaseMsg string
	Status     int
	LastCheck  int64
	IdcCode    string
}

type ReleaseRoutine struct {
	ReleaseTask      *models.ReleaseTask
	ContainerResChan []*ContainerRes
	abandon          chan struct{}
	Status           int
	ErrorMsg         string
}
type IdcReleaseInstance struct {
	idcCode          string
	client           outMarathon.Marathon
	runningInstances []*ContainerRes
}

var (
	TaskChannel         chan *ReleaseRoutine
	RoutinePoolSize     int
	Stop                chan struct{}
	ReleaseRoutineTasks []*ReleaseRoutine
)

func init() {
	TaskChannelSize, err := beego.AppConfig.Int("task_channel_maxsize")
	if err != nil {
		panic("task_channel_maxsize must be integer")

	}
	TaskChannel = make(chan *ReleaseRoutine, TaskChannelSize)
	ReleaseRoutineTasks = make([]*ReleaseRoutine, 0, TaskChannelSize)
	Stop = make(chan struct{})
	RoutinePoolSize, err = beego.AppConfig.Int("routinepool_size")
	if err != nil {
		panic("task_channel_maxsize must be  integer")
	}
	StartReleaseRoutinePool()
}

func releaseTaskFunc(task *ReleaseRoutine) {
	//task.Status = TaskRunning
	var err error
	clients := make([]*IdcReleaseInstance, 0, 3)
	ips := make(map[string][]*ContainerRes)
	o := orm.NewOrm()
	for _, idc := range task.ReleaseTask.ReleaseConf.ReleaseIdc {
		var client outMarathon.Marathon
		client, err = utils.NewMarathonClient(idc.MarathonSerConf.Server, idc.MarathonSerConf.HttpBasicAuthUser, idc.MarathonSerConf.HttpBasicPassword)
		if err != nil {
			task.ErrorMsg = err.Error()
			//task.Status = TaskFail
			task.ReleaseTask.TaskStatus = models.Failed
			task.ReleaseTask.ReleaseMsg = err.Error()
			_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
			if err != nil {
				logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
			}
			return
		}
		mIdcReleaseInstance := &IdcReleaseInstance{}
		mIdcReleaseInstance.client = client
		mIdcReleaseInstance.idcCode = idc.IdcCode
		var iplist []*models.Ip
		iplist, err = models.GetInstances(o, task.ReleaseTask.Service, idc)
		if err != nil {
			task.ErrorMsg = err.Error()
			//task.Status = TaskFail
			task.ReleaseTask.TaskStatus = models.Failed
			task.ReleaseTask.ReleaseMsg = err.Error()
			_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
			if err != nil {
				logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
			}
			return
		}

		for _, ip := range iplist {
			mContainerRes := &ContainerRes{}
			mContainerRes.Instance = ip
			mContainerRes.IdcCode = idc.IdcCode
			mContainerRes.Status = InstanceReady
			task.ContainerResChan = append(task.ContainerResChan, mContainerRes)
			ips[idc.IdcCode] = append(ips[idc.IdcCode], mContainerRes)
		}
		clients = append(clients, mIdcReleaseInstance)
	}
	var errCount int
	IdcReleaseInstances := make([]*IdcReleaseInstance, 0, 3)
	for {
		length := len(IdcReleaseInstances)
		diff := task.ReleaseTask.ReleaseConf.IdcParalle - length
		if diff > 0 && len(clients) > diff {
			IdcReleaseInstances = append(IdcReleaseInstances, clients[:diff]...)
			clients = clients[diff:]
		} else if diff > 0 && len(clients) > 0 {
			IdcReleaseInstances = append(IdcReleaseInstances, clients...)
			clients = clients[:0]
		} else if diff == task.ReleaseTask.ReleaseConf.IdcParalle && len(clients) == 0 {
			//task.Status = TaskSuccess
			task.ReleaseTask.TaskStatus = models.Success
			_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus")
			if err != nil {
				logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
			}
			return
		}

		for outindex := 0; outindex < len(IdcReleaseInstances); {
			idc := IdcReleaseInstances[outindex]
			for index := 0; index < len(idc.runningInstances); {
				instance := idc.runningInstances[index]
				applicationString := fmt.Sprintf("/%s/%s", task.ReleaseTask.Service.Code, instance.Instance.IpAddr)
				if instance.Status == InstanceReady {
					_, err = idc.client.DeleteApplication(applicationString, true)
					if err != nil {
						logs.GetLogger("ReleaseRoutinePool").Println(err.Error())
						instance.ReleaseMsg = err.Error()
						instance.Status = InstanceFail
						errCount = errCount + 1
						if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
							task.ErrorMsg = "错误过多"
							//task.Status = TaskFail
							task.ReleaseTask.TaskStatus = models.Failed
							task.ReleaseTask.ReleaseMsg = "错误过多"
							_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
							if err != nil {
								logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
							}
							return
						} else {
							idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
							continue
						}
					}
					instance.ReleaseMsg = "正在关闭容器"
					instance.Status = InstanceStopping
					instance.LastCheck = time.Now().Unix()
					index = index + 1
					continue
				}
				if instance.Status == InstanceStopping {
					var ok bool
					ok, err = utils.CheckIfDeployment(idc.client, applicationString)
					if err != nil {
						logs.GetLogger("ReleaseRoutinePool").Println(err.Error())
						instance.ReleaseMsg = err.Error()
						instance.Status = InstanceFail
						errCount = errCount + 1
						if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
							task.ErrorMsg = "错误过多"
							//task.Status = TaskFail
							task.ReleaseTask.TaskStatus = models.Failed
							task.ReleaseTask.ReleaseMsg = "错误过多"
							_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
							if err != nil {
								logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
							}
							return
						} else {
							idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
							continue
						}
					}
					if !ok {
						var mApplication *outMarathon.Application
						mApplication, err = utils.CreateMarathonAppFromJson(task.ReleaseTask.Service.MarathonConf)
						if err != nil {
							logs.GetLogger("ReleaseRoutinePool").Println(err.Error())
							instance.ReleaseMsg = err.Error()
							instance.Status = InstanceFail
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								//task.Status = TaskFail
								task.ReleaseTask.TaskStatus = models.Failed
								task.ReleaseTask.ReleaseMsg = "错误过多"
								_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
								if err != nil {
									logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
								}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						imageTag := fmt.Sprintf("%s:%s", task.ReleaseTask.Service.Code, task.ReleaseTask.ImageTag)
						mApplication.Container.Docker.Image = imageTag
						mApplication.ID = applicationString
						mApplication.Container.Docker.EmptyParameters()
						mApplication.Container.Docker.AddParameter("ip", instance.Instance.IpAddr)
						mApplication.Container.Docker.AddParameter("mac-address", instance.Instance.MacAddr)
						mApplication.Container.Docker.AddParameter("net", "dockerbr0")
						_, err = idc.client.CreateApplication(mApplication)
						if err != nil {
							logs.GetLogger("ReleaseRoutinePool").Println(err.Error())
							instance.ReleaseMsg = err.Error()
							instance.Status = InstanceFail
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								//task.Status = TaskFail
								task.ReleaseTask.TaskStatus = models.Failed
								task.ReleaseTask.ReleaseMsg = "错误过多"
								_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
								if err != nil {
									logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
								}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						instance.ReleaseMsg = "正在启动容器"
						instance.Status = InstanceStarting
						instance.LastCheck = time.Now().Unix()
						index = index + 1
						continue
					} else {
						nowTime := time.Now().Unix()
						if nowTime-instance.LastCheck > task.ReleaseTask.ReleaseConf.TimeOut {
							instance.ReleaseMsg = "等待关闭容器时超时"
							instance.Status = InstanceFail
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								//task.Status = TaskFail
								task.ReleaseTask.TaskStatus = models.Failed
								task.ReleaseTask.ReleaseMsg = "错误过多"
								_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
								if err != nil {
									logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
								}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						index = index + 1
						continue
					}
				}
				if instance.Status == InstanceStarting {
					var ok bool
					ok, err = idc.client.ApplicationOK(applicationString)
					if err != nil {
						logs.GetLogger("ReleaseRoutinePool").Println(err.Error())
						instance.ReleaseMsg = err.Error()
						instance.Status = InstanceFail
						errCount = errCount + 1
						if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
							task.ErrorMsg = "错误过多"
							//task.Status = TaskFail
							task.ReleaseTask.TaskStatus = models.Failed
							task.ReleaseTask.ReleaseMsg = "错误过多"
							_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
							if err != nil {
								logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
							}
							return
						} else {
							idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
							continue
						}
					}
					if ok {
						instance.ReleaseMsg = "容器启动成功"
						instance.Status = InstanceFinish
						idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
						continue
					} else {
						nowTime := time.Now().Unix()
						if nowTime-instance.LastCheck > task.ReleaseTask.ReleaseConf.TimeOut {
							instance.ReleaseMsg = "等待启动容器时超时"
							instance.Status = InstanceFail
							errCount = errCount + 1
							if errCount >= task.ReleaseTask.ReleaseConf.FaultTolerant {
								task.ErrorMsg = "错误过多"
								//task.Status = TaskFail
								task.ReleaseTask.TaskStatus = models.Failed
								task.ReleaseTask.ReleaseMsg = "错误过多"
								_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus", "ReleaseMsg")
								if err != nil {
									logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
								}
								return
							} else {
								idc.runningInstances = append(idc.runningInstances[:index], idc.runningInstances[index+1:]...)
								continue
							}
						}
						index = index + 1
						continue
					}
				}
			}
			innerLength := len(idc.runningInstances)
			innerDiff := task.ReleaseTask.ReleaseConf.IdcInnerParalle - innerLength
			if innerDiff > 0 && len(ips[idc.idcCode]) > innerDiff {
				idc.runningInstances = append(idc.runningInstances, ips[idc.idcCode][:innerDiff]...)
				ips[idc.idcCode] = ips[idc.idcCode][innerDiff:]
			} else if innerDiff > 0 && len(ips[idc.idcCode]) > 0 {
				idc.runningInstances = append(idc.runningInstances, ips[idc.idcCode]...)
				ips[idc.idcCode] = ips[idc.idcCode][:0]
			} else if innerDiff == task.ReleaseTask.ReleaseConf.IdcInnerParalle && len(ips[idc.idcCode]) == 0 {
				IdcReleaseInstances = append(IdcReleaseInstances[:outindex], IdcReleaseInstances[outindex+1:]...)
			}
			select {
			case <-task.abandon:
				//task.Status = TaskAbandon
				task.ReleaseTask.TaskStatus = models.Abandon
				_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "TaskStatus")
				if err != nil {
					logs.GetLogger("ReleaseRoutinePool").Printf("update Task status fail,detail:%s", err.Error())
				}
				return
			case <-time.After(time.Millisecond * 1000):
			}
		}

	}

}

func AddTask(mReleaseTask *models.ReleaseTask) error {
	task := ReleaseRoutine{}
	task.ReleaseTask = mReleaseTask
	//task.Status = TaskReady
	select {
	case TaskChannel <- &task:
		ReleaseRoutineTasks = append(ReleaseRoutineTasks, &task)
		return nil
	default:
		return errors.New("Release TaskChannel is full")
	}
}

func CheckTaskStatus(mReleaseTask *models.ReleaseTask) (*ReleaseRoutine, error) {
	for _, mReleaseRoutine := range ReleaseRoutineTasks {
		if mReleaseRoutine.ReleaseTask.Id == mReleaseTask.Id {
			mReleaseTask.Service = mReleaseRoutine.ReleaseTask.Service
			return mReleaseRoutine, nil
		}
	}
	return nil, errors.New("task no found")
}

func AbandonTask(mReleaseTask *models.ReleaseTask) error {
	for _, mReleaseRoutine := range ReleaseRoutineTasks {
		if mReleaseRoutine.ReleaseTask.Id == mReleaseTask.Id {
			close(mReleaseRoutine.abandon)
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
					var err error
					var testjson []byte
					testjson, err = json.Marshal(task.ContainerResChan)
					if err != nil {
						logs.GetLogger("ReleaseRoutinePool").Printf("save release result fail,detail:%s", err.Error())
					}
					task.ReleaseTask.ReleaseResult = string(testjson)
					o := orm.NewOrm()
					_, err = models.CreateOrUpdateRelease(o, task.ReleaseTask, "ReleaseResult")
					if err != nil {
						logs.GetLogger("ReleaseRoutinePool").Printf("save release result fail,detail:%s", err.Error())
					}
					for index, mReleaseRoutine := range ReleaseRoutineTasks {
						if mReleaseRoutine == task {
							ReleaseRoutineTasks = append(ReleaseRoutineTasks[:index], ReleaseRoutineTasks[index+1:]...)
						}
					}
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
