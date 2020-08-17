package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	DBClient *gorm.DB
)

func InitDb(username string, password string, dbAddress string, dbName string) {
	var err error
	DBClient, err = gorm.Open("mysql", username+":"+password+"@("+dbAddress+")/"+dbName+
		"?charset=utf8mb4,utf8&collation=utf8mb4_general_ci&parseTime=True&loc=Asia%2FShanghai&readTimeout=120s&"+
		"tls=false&writeTimeout=120s")
	if err != nil {
		log.WithField("err", err).WithField("op", "init db client error").Fatal()
	}
	DBClient.DB().SetMaxOpenConns(50)
	DBClient.DB().SetMaxIdleConns(10)
	DBClient.DB().SetConnMaxLifetime(time.Hour)
}

func GetAllVm() map[uint64]VmList {
	vmMap := make(map[uint64]VmList)

	vmList := getVmList()
	vmDisk := getAllVmDisk()
	vmMac := getAllVmMac()
	vmPci := getAllVmPci()
	vmPort := getAllVmPort()

	for _, vm := range vmList {
		vm.VmDiskList = make([]VmDisk, 0)
		vm.VmMacList = make([]VmMac, 0)
		vm.VmPciList = make([]VmPci, 0)
		vm.VmPortList = make([]VmPort, 0)
		vmMap[vm.Id] = vm
	}
	for _, disk := range vmDisk {
		vm := vmMap[disk.VmId]
		vmDiskList := vm.VmDiskList
		vmDiskList = append(vmDiskList, disk)
		vm.VmDiskList = vmDiskList
		vmMap[vm.Id] = vm
	}

	for _, mac := range vmMac {
		vm := vmMap[mac.VmId]
		vmMacList := vm.VmMacList
		vmMacList = append(vmMacList, mac)
		vm.VmMacList = vmMacList
		vmMap[vm.Id] = vm
	}

	for _, pci := range vmPci {
		vm := vmMap[pci.VmId]
		vmPciList := vm.VmPciList
		vmPciList = append(vmPciList, pci)
		vm.VmPciList = vmPciList
		vmMap[vm.Id] = vm
	}

	for _, port := range vmPort {
		vm := vmMap[port.VmId]
		vmPortList := vm.VmPortList
		vmPortList = append(vmPortList, port)
		vm.VmPortList = vmPortList
		vmMap[vm.Id] = vm
	}

	return vmMap
}

func getAllVmDisk() []VmDisk {
	vmDiskList := make([]VmDisk, 0)
	var lastId uint64
	lastId = 0
	for {
		tmpVmDisk := make([]VmDisk, 0)
		DBClient.Where("id > ?", lastId).Limit(100).Find(&tmpVmDisk)
		if tmpVmDisk == nil || len(tmpVmDisk) < 1 {
			break
		}
		for _, t := range tmpVmDisk {
			vmDiskList = append(vmDiskList, t)
			if t.Id > lastId {
				lastId = t.Id
			}
		}
	}
	return vmDiskList
}

func getAllVmMac() []VmMac {
	vmMapList := make([]VmMac, 0)
	var lastId uint64
	lastId = 0
	for {
		tmpVmMac := make([]VmMac, 0)
		DBClient.Where("id > ?", lastId).Limit(100).Find(&tmpVmMac)
		if tmpVmMac == nil || len(tmpVmMac) < 1 {
			break
		}
		for _, t := range tmpVmMac {
			vmMapList = append(vmMapList, t)
			if t.Id > lastId {
				lastId = t.Id
			}
		}
	}
	return vmMapList
}

func getAllVmPci() []VmPci {
	vmPciList := make([]VmPci, 0)
	var lastId uint64
	lastId = 0
	for {
		tmpVmPci := make([]VmPci, 0)
		DBClient.Where("id > ?", lastId).Limit(100).Find(&tmpVmPci)
		if tmpVmPci == nil || len(tmpVmPci) < 1 {
			break
		}
		for _, t := range tmpVmPci {
			vmPciList = append(vmPciList, t)
			if t.Id > lastId {
				lastId = t.Id
			}
		}
	}
	return vmPciList
}

func getAllVmPort() []VmPort {
	vmPortList := make([]VmPort, 0)
	var lastId uint64
	lastId = 0
	for {
		tmpVmPort := make([]VmPort, 0)
		DBClient.Where("id > ?", lastId).Limit(100).Find(&tmpVmPort)
		if tmpVmPort == nil || len(tmpVmPort) < 1 {
			break
		}
		for _, t := range tmpVmPort {
			vmPortList = append(vmPortList, t)
			if t.Id > lastId {
				lastId = t.Id
			}
		}
	}
	return vmPortList
}

func getVmList() []VmList {
	vmList := make([]VmList, 0)
	var lastId uint64
	lastId = 0
	for {
		tmpVmList := make([]VmList, 0)
		DBClient.Where("id > ?", lastId).Limit(100).Find(&tmpVmList)
		if tmpVmList == nil || len(tmpVmList) < 1 {
			break
		}
		for _, t := range tmpVmList {
			vmList = append(vmList, t)
			if t.Id > lastId {
				lastId = t.Id
			}
		}
	}
	return vmList
}
