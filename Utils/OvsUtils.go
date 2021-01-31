package Utils

import (
	"os/exec"
)

const (
	OvsBrName = "br-dpdk"
)

func CheckAndCreatePort(name string, path string) error {
	cmdStr := "ovs-vsctl list-br | grep " + OvsBrName + " && ovs-vsctl list-ports " + OvsBrName + " | grep " + name + " ||" +
		" ovs-vsctl add-port " + OvsBrName + " " + name + " -- " +
		"set Interface " + name + " type=dpdkvhostuserclient " +
		"options:vhost-server-path=" + path

	cmd := exec.Command("bash", "-c", cmdStr)
	return cmd.Run()
}
