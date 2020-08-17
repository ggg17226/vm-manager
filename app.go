package main

import (
	log "github.com/sirupsen/logrus"
	"runtime"
	"vm-manager/Config"
	"vm-manager/LogUtils"
	"vm-manager/controller"
	"vm-manager/model"
	"vm-manager/tasks"
)

func main() {
	runtime.GOMAXPROCS(2)

	LogUtils.InitLog()
	Config.InitConfig()

	model.InitDb(Config.AppConfig.Db.Username, Config.AppConfig.Db.Password, Config.AppConfig.Db.Host, Config.AppConfig.Db.DbName)
	defer model.DBClient.Close()

	log.Info("init success")

	go tasks.BuildVmListTask()

	go tasks.VmManagerTask()

	controller.InitGin(Config.AppConfig.Server.Listen)
}
