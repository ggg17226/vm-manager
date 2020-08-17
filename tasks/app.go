package tasks

import (
	"sync"
	"vm-manager/Utils"
	"vm-manager/model"
)

var (
	vmMapMutex sync.Mutex
	vmMap      map[uint64]model.VmList

	startupStateMutex sync.Mutex
	startupState      = false

	VmRuntimeMapMutex sync.Mutex
	VmRuntimeMap      = make(map[uint64]Utils.VmRuntimePayload, 0)
)
