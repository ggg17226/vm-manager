package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vm-manager/tasks"
)

func GetStatus(context *gin.Context) {
	tasks.VmRuntimeMapMutex.Lock()
	defer tasks.VmRuntimeMapMutex.Unlock()
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"vmInfo": tasks.VmRuntimeMap,
	})
}
