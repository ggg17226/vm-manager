package tasks

import (
	"runtime"
	"time"
	"vm-manager/Utils"
)

func VmManagerTask() {
	time.Sleep(10 * time.Second)
	startupStateMutex.Lock()
	defer startupStateMutex.Unlock()
	for {
		vmMapMutex.Lock()
		if vmMap != nil && len(vmMap) > 0 {
			if !startupState {
				startupState = true
				VmRuntimeMapMutex.Lock()
				for _, vm := range vmMap {
					if vm.AutoStartup == 1 {
						VmRuntimeMap[vm.Id] = *Utils.BuildVmRuntimePayload(&vm)
					}
				}
				for id, payload := range VmRuntimeMap {
					Utils.RunVm(&payload)
					VmRuntimeMap[id] = payload
				}
				VmRuntimeMapMutex.Unlock()
			} else {
				VmRuntimeMapMutex.Lock()
				for _, vm := range vmMap {
					VmRuntimeMap[vm.Id] = *Utils.BuildVmRuntimePayload(&vm)
				}
				for id, payload := range VmRuntimeMap {
					Utils.CheckVm(&payload)
					VmRuntimeMap[id] = payload
				}
				VmRuntimeMapMutex.Unlock()
			}
		}
		vmMapMutex.Unlock()
		runtime.Gosched()
		time.Sleep(5 * time.Second)
	}
}
