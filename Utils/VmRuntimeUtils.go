package Utils

import (
	"github.com/prometheus/procfs"
	"github.com/reiver/go-telnet"
	log "github.com/sirupsen/logrus"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"vm-manager/Config"
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

type ProcessInfoPayload struct {
	Pid     int
	CmdLine []string
}

var (
	qemuProcesses           []ProcessInfoPayload
	qemuProcessesLastUpdate int64
	qemuProcessesUpdateLock sync.Mutex
)

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

	defer setFlushNow()

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

	flushProcesses()

	if qemuProcesses != nil && len(qemuProcesses) > 0 {
		for _, process := range qemuProcesses {
			if FindInArray(&(process.CmdLine), keyStr) {
				return process.Pid
			}
		}
	}
	return -2

	//	cmd := exec.Command("bash", "-c", "ps -aux | grep qemu | grep "+keyStr+" | grep -v bash | awk '{print $2}'")
	//out, err := cmd.CombinedOutput()
	//
	//if err != nil || out == nil {
	//	return -1
	//}
	//pidStr := strings.Replace(string(out), "\n", "", -1)
	//if len(pidStr) < 1 {
	//	return -2
	//}
	//pid, errAtoi := strconv.Atoi(pidStr)
	//if errAtoi != nil || pid < 1 {
	//	return -3
	//}
	//return pid
}

func ShutdownVm(port int) {
	defer setFlushNow()
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

	switch Config.AppConfig.Host.NetType {
	case "dpdk":
		memId := "mem-" + vm.Name
		cmd += " -object memory-backend-file,id=" + memId + ",size=" + strconv.Itoa(int(vm.Mem)) + "M,mem-path=/dev/hugepages,share=on "
		cmd += " -mem-prealloc "
		cmd += " -numa node,memdev=" + memId + " "
		break
	default:
		break
	}

	for _, mac := range vm.VmMacList {

		switch Config.AppConfig.Host.NetType {
		case "dpdk":
			tmpMacStr := strings.ReplaceAll(mac.Mac, ":", "")
			socketPath := "/run/openvswitch/vhost-" + vm.Name + "-" + tmpMacStr
			CheckAndCreatePort("vhost-"+vm.Name+"-"+tmpMacStr, socketPath)

			charId := "char-" + vm.Name + "-" + tmpMacStr
			cmd += " -chardev socket,id=" + charId + ",path=" + socketPath + ",server "
			netDevId := "ndev-" + vm.Name + "-" + tmpMacStr
			cmd += " -netdev type=vhost-user,id=" + netDevId + ",chardev=" + charId + ",vhostforce "
			cmd += " -device virtio-net-pci,mac=" + mac.Mac + ",netdev=" + netDevId + " "
			break
		case "tap":
			cmd += " -net nic,model=virtio,macaddr=" + mac.Mac + " "
			if len(keyStr) < 1 {
				keyStr = mac.Mac
			}
			break
		default:
			log.WithField("err", "unknown net type "+Config.AppConfig.Host.NetType).WithField("op", "build vm runtime payload").Fatal()
			break
		}

	}
	if "tap" == Config.AppConfig.Host.NetType {
		cmd += " -net tap,script=/opt/image/qemu-ifup "
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

func flushProcesses() {
	qemuProcessesUpdateLock.Lock()
	defer qemuProcessesUpdateLock.Unlock()

	if qemuProcesses == nil || len(qemuProcesses) == 0 || math.Abs(float64(time.Now().Unix()-qemuProcessesLastUpdate)) > 10 {
		list, err := getQemuProcList()
		if err == nil {
			qemuProcesses = *list
		} else {
			qemuProcesses = make([]ProcessInfoPayload, 0)
		}
		qemuProcessesLastUpdate = time.Now().Unix()
	}
}

func setFlushNow() {
	qemuProcessesUpdateLock.Lock()
	defer qemuProcessesUpdateLock.Unlock()
	qemuProcessesLastUpdate = 0
}

func getQemuProcList() (*[]ProcessInfoPayload, error) {
	procArr, err := procfs.AllProcs()
	if err != nil {
		return nil, err
	}

	result := make([]ProcessInfoPayload, 0)

	for _, proc := range procArr {
		stat, err := proc.Stat()
		if err != nil {
			continue
		}
		cmdLine, err := proc.CmdLine()
		if err != nil {
			continue
		}
		if FindInArray(&cmdLine, "qemu-system") {
			payload := ProcessInfoPayload{
				Pid:     stat.PID,
				CmdLine: cmdLine,
			}
			result = append(result, payload)
		}
	}
	return &result, nil

}

func FindInArray(arr *[]string, str string) bool {
	for _, s := range *arr {
		if strings.Contains(s, str) {
			return true
		}
	}
	return false
}
