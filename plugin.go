package main

import (
	"github.com/g8os/core0/base/plugin"
	"github.com/g8os/ovs-plugin/ovs"
)

var (
	commands = plugin.Commands{
		"bridge-add":  ovs.BridgeAdd,
		"bridge-del":  ovs.BridgeDelete,
		"port-add":    ovs.PortAdd,
		"port-del":    ovs.PortDel,
		"bond-add":    ovs.BondAdd,
		"vtep-ensure": ovs.VtepEnsure,
		"vtep-del":    ovs.VtepDelete,
		"set":         ovs.Set,
	}
)

func main() {
	plugin.Plugin(commands)
}
