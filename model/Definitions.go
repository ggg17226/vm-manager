package model

import "time"

type VmList struct {
	Id          uint64     `gorm:"PRIMARY_KEY;AUTO_INCREMENT;Column:id"`
	Name        string     `gorm:"type:varchar(128);Column:name;NOT NULL"`
	Cpu         int32      `gorm:"Column:cpu;NOT NULL"`
	Mem         int32      `gorm:"Column:mem;NOT NULL"`
	AutoStartup int32      `gorm:"Column:auto_startup;NOT NULL"`
	Status      int32      `gorm:"Column:status;NOT NULL"`
	CreatedAt   *time.Time `gorm:"Column:create_time"`
	UpdatedAt   *time.Time `gorm:"Column:update_time"`
	VmDiskList  []VmDisk
	VmMacList   []VmMac
	VmPciList   []VmPci
	VmPortList  []VmPort
}

type VmDisk struct {
	Id       uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT;Column:id"`
	VmId     uint64 `gorm:"Column:vm_id;NOT NULL"`
	DiskPath string `gorm:"type:varchar(1024);Column:disk_path;NOT NULL"`
}

type VmMac struct {
	Id   uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT;Column:id"`
	VmId uint64 `gorm:"Column:vm_id;NOT NULL"`
	Mac  string `gorm:"type:char(17);Column:mac;NOT NULL"`
}

type VmPci struct {
	Id    uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT;Column:id"`
	VmId  uint64 `gorm:"Column:vm_id;NOT NULL"`
	PciId string `gorm:"type:varchar(32);Column:pci_id;NOT NULL"`
}

type VmPort struct {
	Id   uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT;Column:id"`
	VmId uint64 `gorm:"Column:vm_id;NOT NULL"`
	Type int32  `gorm:"Column:type;NOT NULL"`
	Port int32  `gorm:"Column:port;NOT NULL"`
}
