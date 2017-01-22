package controllers
import (
	"github.com/xianyouQ/go-dockermgr/models"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"errors"
    "github.com/xianyouQ/go-dockermgr/utils"
    outMarathon "github.com/xianyouQ/go-marathon"
	"github.com/astaxie/beego/orm"
	"time"
)



type  ContainerRes struct {
    Instance *models.Ip
    ReleaseMsg string
}

type ReleaseRoutine struct {
    ReleaseTask *models.ReleaseTask
    ContainerResChan  []*ContainerRes
    Abandon chan struct{}
    Done bool
    ErrorMsg  string
}
type IdcReleaseInstance struct {
    idcCode string
    client outMarathon.Marathon
    runningInstances []*models.Ip
    done bool 
    instances chan models.Ip
}
var (
    TaskChannel chan *ReleaseRoutine
    RoutinePoolSize int
    Stop chan struct{}
    ReleaseRoutineTasks []ReleaseRoutine
)

func init() {
    TaskChannelSize,err := beego.AppConfig.Int("task_channel_maxsize")
    if err != nil {
        logs.Critical("task_channel_maxsize must be integer")
    }
    TaskChannel = make(chan *ReleaseRoutine, TaskChannelSize)
    ReleaseRoutineTasks = make([]ReleaseRoutine,TaskChannelSize)
    Stop = make(chan struct{})
    RoutinePoolSize,err = beego.AppConfig.Int("routinepool_size")
    if err != nil {
        logs.Critical("task_channel_maxsize must be  integer")
    }
    StartReleaseRoutinePool()
}




func idcTaskMonitor(clients <-chan *IdcReleaseInstance,task *ReleaseRoutine) {
    IdcReleaseInstances := make([]*IdcReleaseInstance,3,3)
    for {
        if len(IdcReleaseInstances) < task.ReleaseTask.ReleaseConf.IdcParalle {
            select {
                case newIdcRelease := <- clients:
                    IdcReleaseInstances = append(IdcReleaseInstances,newIdcRelease)
                case <- time.After(time.Second):
                    if task.Done == true {
                        return
                    }
                    continue
            }
        }
        

    }

}
func releaseTaskFunc(task *ReleaseRoutine) {
    var err error
    clientChannel := make(chan *IdcReleaseInstance,task.ReleaseTask.ReleaseConf.IdcParalle)
    clients := make([]*IdcReleaseInstance,3,3)
    ips := make (map[string][]*models.Ip)
    for _,idc := range task.ReleaseTask.ReleaseConf.ReleaseIdc {
        var client outMarathon.Marathon
        client,err = utils.NewMarathonClient(idc.MarathonSerConf.Server,idc.MarathonSerConf.HttpBasicAuthUser,idc.MarathonSerConf.HttpBasicPassword)
        if err != nil {
            task.ErrorMsg = err.Error()
            return
        }
        mIdcReleaseInstance := &IdcReleaseInstance{}
        mIdcReleaseInstance.client = client
        mIdcReleaseInstance.idcCode = idc.IdcCode
        mIdcReleaseInstance.instances = make(chan models.Ip,task.ReleaseTask.ReleaseConf.IdcInnerParalle)
        clients = append(clients,mIdcReleaseInstance)
        o := orm.NewOrm()
        var iplist []*models.Ip
        iplist , err = models.GetInstances(o,task.ReleaseTask.ReleaseConf.Service,idc)
        if err != nil {
            task.ErrorMsg = err.Error()
            return
        }
        ips[idc.IdcCode] = iplist
        for _,ip :=range iplist {
            mContainerRes := &ContainerRes{}
            mContainerRes.Instance = ip
            task.ContainerResChan = append(task.ContainerResChan,mContainerRes)
        }
    }
    go idcTaskMonitor(clientChannel,task)
    count := 0
    for  {
        alldone := true
        for _,client := range clients{
            if instances,ok := ips[client.idcCode]; ok {
                var endIndex int
                for index,instance := range instances {
                    select {
                        case client.instances <- *instance:
                            continue
                        case <- time.After(time.Second):
                            endIndex = index
                            break
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
            break
        }
        if count >= len(clients) {
            continue
        }
        select {
            case clientChannel <- clients[count]:
                count = count + 1
            case <- task.Abandon:
                for _,client := range clients{
                    if client.done == false {
                        client.done = true
                    }
                }
                return
            case <- time.After(time.Second):
                continue
            
        }

    }

}

func AddTask(mReleaseTask *models.ReleaseTask)  error {
    task := ReleaseRoutine{}
    task.ReleaseTask = mReleaseTask
    //...
    select {
        case TaskChannel <- &task:
            ReleaseRoutineTasks = append(ReleaseRoutineTasks,task)
            return nil
        default:
            return errors.New("Release TaskChannel is full")
    }
}


func CheckTaskStatus(mReleaseTask *models.ReleaseTask)  ([]*ContainerRes,error){
     for _,mReleaseRoutine := range ReleaseRoutineTasks {
         if mReleaseRoutine.ReleaseTask.Id == mReleaseTask.Id {
             return mReleaseRoutine.ContainerResChan,nil
         }
     }
     return nil,errors.New("task no found")
}


func AbandonTask(mReleaseTask *models.ReleaseTask) error {
     for _,mReleaseRoutine := range ReleaseRoutineTasks {
         if mReleaseRoutine.ReleaseTask.Id == mReleaseTask.Id {
             close(mReleaseRoutine.Abandon)
             return nil
         }
     }
     return errors.New("task no found")
}


func StartReleaseRoutinePool() {
    for i := 0 ; i < RoutinePoolSize ; i++ {
        go func() {
            for {
                select {
                    case task := <- TaskChannel:
                        releaseTaskFunc(task)
                    case <- Stop:
                        return
                }
            }
        }()
    }
}

func StopReleaseRoutinePool() {
    close(Stop)
}
