package Utils

import (
	"github.com/reiver/go-telnet"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"vm-manager/model"
)

type VmRuntimePayload struct {
	Pid         int
	MonitorPort int
	KeyStr      string
	Cmd         string
	State       int //0-success 1-error
	VncType     int
	VncPort     int
	Name        string
}

func CheckVm(payload *VmRuntimePayload) {
	pid := FindPid(payload.KeyStr)
	if pid > 0 {
		payload.State = 0

	} else {
		payload.State = 1

	}
	payload.Pid = pid
}

func RunVm(payload *VmRuntimePayload) {

	pid := FindPid(payload.KeyStr)
	if pid > 0 {
		payload.State = 0
		payload.Pid = pid
		return
	}
	cmd := exec.Command("bash", "-c", payload.Cmd)
	errRunStart := cmd.Run()
	if errRunStart != nil {
		payload.State = 1
		payload.Pid = -1
	} else {
		payload.State = 0
		payload.Pid = FindPid(payload.KeyStr)
	}
}

func FindPid(keyStr string) int {
	cmd := exec.Command("bash", "-c", "ps -aux | grep qemu | grep "+keyStr+" | grep -v bash | awk '{print $2}'")
	out, err := cmd.CombinedOutput()

	if err != nil || out == nil {
		return -1
	}
	pidStr := strings.Replace(string(out), "\n", "", -1)
	if len(pidStr) < 1 {
		return -2
	}
	pid, errAtoi := strconv.Atoi(pidStr)
	if errAtoi != nil || pid < 1 {
		return -3
	}
	return pid
}

func ShutdownVm(port int) {
	conn, err := telnet.DialTo("127.0.0.1:" + strconv.Itoa(port))
	if err == nil {
		time.Sleep(1 * time.Second)
		conn.Write([]byte("{\"execute\":\"qmp_capabilities\"}\n"))

		conn.Write([]byte("{\"execute\":\"system_powerdown\"}\n"))
		b := make([]byte, 1024)
		for {
			n, errRead := conn.Read(b)
			if errRead != nil || n > 0 {
				break
			}
			time.Sleep(1)
		}
		conn.Close()
	}
}

func BuildVmRuntimePayload(vm *model.VmList) *VmRuntimePayload {
	keyStr := ""
	var monitorPort int
	monitorPort = 0

	vncPort := 0
	vncType := 0

	cmd := "/usr/bin/qemu-system-x86_64 -enable-kvm -cpu host,kvm=off "
	cmd += " -smp " + strconv.Itoa(int(vm.Cpu))
	cmd += " -m " + strconv.Itoa(int(vm.Mem)) + "M "

	for _, mac := range vm.VmMacList {
		cmd += " -net nic,model=virtio,macaddr=" + mac.Mac + " "
		if len(keyStr) < 1 {
			keyStr = mac.Mac
		}
	}
	for _, disk := range vm.VmDiskList {
		cmd += " -drive file=" + disk.DiskPath + ",if=virtio "
		if len(keyStr) < 1 {
			split := strings.Split(disk.DiskPath, "/")
			if len(split) > 0 {
				keyStr = split[len(split)-1]
			}
		}
	}
	cmd += " -net tap,script=/opt/image/qemu-ifup "
	for _, port := range vm.VmPortList {
		if port.Type == 0 {
			cmd += " -vga qxl -spice port=" + strconv.Itoa(int(port.Port)) + ",disable-ticketing "
			vncPort = int(port.Port)
			vncType = 0
		} else if port.Type == 1 {
			continue
		} else if port.Type == 2 {
			cmd += " -qmp tcp:0.0.0.0:" + strconv.Itoa(int(port.Port)) + ",server,nowait "
			if monitorPort == 0 {
				monitorPort = int(port.Port)
			}
		}
	}
	for _, pci := range vm.VmPciList {
		cmd += " -device vfio-pci,host=" + pci.PciId
	}
	cmd += " -daemonize "
	vmRuntimePayload := VmRuntimePayload{
		Cmd:         cmd,
		KeyStr:      keyStr,
		MonitorPort: monitorPort,
		VncPort:     vncPort,
		VncType:     vncType,
		Name:        vm.Name,
	}
	return &vmRuntimePayload
}
