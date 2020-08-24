package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitGin(listenAddress string) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(func(context *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		context.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		//请求方式
		reqMethod := context.Request.Method
		//请求路由
		reqUrl := context.Request.RequestURI
		//状态码
		statusCode := context.Writer.Status()
		//请求ip
		clientIP := context.ClientIP()
		var statusPrefix int32
		statusPrefix = int32(statusCode / 100)
		if statusPrefix < 2 || statusPrefix > 4 {
			log.WithFields(log.Fields{
				"status_code":  statusCode,
				"latency_time": latencyTime,
				"client_ip":    clientIP,
				"req_method":   reqMethod,
				"req_uri":      reqUrl,
				"type":         "req error",
			}).Error()
		}
	})
	router.Use(cors.Default())
	router.GET("/status", GetStatus)

	router.GET("/vm/:id/start", StartupVm)

	router.GET("/vm/:id/shutdown", ShutdownVm)

	router.Run(listenAddress)
}
