package controllers
import (
	"github.com/xianyouQ/go-dockermgr/models"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"errors"
    "github.com/xianyouQ/go-dockermgr/utils"
    outMarathon "github.com/xianyouQ/go-marathon"
	"github.com/astaxie/beego/orm"
)



type  ContainerRes struct {

}

type ReleaseRoutine struct {
    ReleaseTask *models.ReleaseTask
    ContainerResChan  []*ContainerRes
    Abandon chan struct{}
    ErrorMsg  string
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
}



func releaseTaskFunc(task *ReleaseRoutine) {
    var err error
    clients := make([]outMarathon.Marathon,3,3)
    ips := make ([][]*models.Ip,3,3)
    for _,idc := range task.ReleaseTask.ReleaseConf.ReleaseIdc {
        var client outMarathon.Marathon
        client,err = utils.NewMarathonClient(idc.MarathonSerConf.Server,idc.MarathonSerConf.HttpBasicAuthUser,idc.MarathonSerConf.HttpBasicPassword)
        if err != nil {
            task.ErrorMsg = err.Error()
            return
        }
        clients = append(clients,client)
        o := orm.NewOrm()
        var ip []*models.Ip
        ip , err = models.GetInstances(o,task.ReleaseTask.ReleaseConf.Service,idc)
        if err != nil {
            task.ErrorMsg = err.Error()
            return
        }
        ips = append(ips,ip)
    }

    for {


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

/*
func AbandonTask(mReleaseTask *models.ReleaseTask) error {

}
*/

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
