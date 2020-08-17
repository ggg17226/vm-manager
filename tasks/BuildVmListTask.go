package tasks

import (
	"runtime"
	"time"
	"vm-manager/model"
)

func BuildVmListTask() {
	for {
		vmMapMutex.Lock()
		vmMap = model.GetAllVm()
		vmMapMutex.Unlock()
		runtime.Gosched()
		time.Sleep(30 * time.Second)
	}
}
