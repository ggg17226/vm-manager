/**
日志配置
*/
package LogUtils

import (
	rotateLogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
	"vm-manager/Utils"
)

var (
	logPath  = "logs"
	useLog   = "both"
	logLevel = log.InfoLevel
)

func InitLog() {
	if !Utils.FileExist(logPath) {
		os.Mkdir(logPath, 0777)
	}
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})
	// 设置日志格式为文本格式
	//log.SetFormatter(&log.TextFormatter{})

	/**
	设置日志格式
	*/
	logFilePath := filepath.Join(logPath, "vm-manager.log")
	rotateWriter, _ := rotateLogs.New(
		logFilePath+".%Y-%m-%d-%H",
		rotateLogs.WithLinkName("vm-manager.log"),
		rotateLogs.WithRotationTime(12*time.Hour),
	)
	/**
	  配置日志输出目标
	*/
	switch useLog {
	default:
		log.SetOutput(os.Stdout)
		break
	case "both":
		mw := io.MultiWriter(os.Stdout, rotateWriter)
		log.SetOutput(mw)
		break
	case "file":
		log.SetOutput(rotateWriter)
		break
	case "std":
		log.SetOutput(os.Stdout)
		break
	}

	log.SetLevel(logLevel)
}
