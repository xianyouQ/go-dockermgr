package controllers
import (
	"github.com/xianyouQ/go-dockermgr/models"
    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)



type  ContainerRes struct {

}

type ReleaseRoutine struct {
    ReleaseTask models.ReleaseTask
    ContainerResChan *chan ContainerRes
}
var (
    TaskChannel chan ReleaseRoutine
    RoutinePoolSize int
)

func init () {
    TaskChannelSize,err := beego.AppConfig.Int("task_channel_maxsize")
    if err != nil {
        logs.Critical("task_channel_maxsize must be integer")
    }
    TaskChannel = make(chan ReleaseRoutine, TaskChannelSize)
    RoutinePoolSize,err = beego.AppConfig.Int("routinepool_size")
    if err != nil {
        logs.Critical("task_channel_maxsize must be  integer")
    }
}


func releaseTaskFunc(task ReleaseRoutine) {

}

func StartReleaseRoutinePool() {
    for i := 0 ; i < RoutinePoolSize ; i++ {
        go func() {
            for {
                select {
                    case task := <- TaskChannel:
                        releaseTaskFunc(task)
                }
            }
        }()
    }
}
