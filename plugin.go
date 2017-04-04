package main

import (
	"github.com/g8os/core0/base/plugin"
	"github.com/g8os/ovs-plugin/ovs"
)

var (
	commands = plugin.Commands{
		"bridge-add": ovs.BridgeAdd,
		"bridge-del": ovs.BridgeDelete,
		"port-add":   ovs.PortAdd,
		"set":        ovs.Set,
	}
)

func main() {
	plugin.Plugin(commands)
}
