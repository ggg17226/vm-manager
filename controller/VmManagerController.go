package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vm-manager/Utils"
	"vm-manager/tasks"
)

func StartupVm(context *gin.Context) {
	tasks.VmRuntimeMapMutex.Lock()
	defer tasks.VmRuntimeMapMutex.Unlock()

	id := context.Param("id")

	idNum, parseErr := strconv.ParseUint(id, 10, 64)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "param error",
		})
	}

	if val, ok := tasks.VmRuntimeMap[idNum]; ok {
		if val.State == 0 {
			context.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"msg":    "already running",
			})
		} else {
			tmp := tasks.VmRuntimeMap[idNum]
			go Utils.RunVm(&tmp)
			context.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"msg":    "starting vm",
			})
		}
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "vm not exist",
		})
	}
}

func ShutdownVm(context *gin.Context) {
	tasks.VmRuntimeMapMutex.Lock()
	defer tasks.VmRuntimeMapMutex.Unlock()

	id := context.Param("id")

	idNum, parseErr := strconv.ParseUint(id, 10, 64)
	if parseErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "param error",
		})
	}

	if val, ok := tasks.VmRuntimeMap[idNum]; ok {
		if val.State == 1 {
			context.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"msg":    "already shutdown",
			})
		} else {
			tmp := tasks.VmRuntimeMap[idNum]
			go Utils.ShutdownVm(tmp.MonitorPort)
			context.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"msg":    "stopping vm",
			})
		}
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "vm not exist",
		})
	}
}
