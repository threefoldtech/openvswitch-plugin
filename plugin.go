package main

import (
	"github.com/g8os/core0/base/plugin"
	"github.com/g8os/ovs-plugin/ovs"
)

var (
	Plugin = plugin.Commands{
		"bridge-add": ovs.BridgeAdd,
		"bridge-del": ovs.BridgeDelete,
	}
)

func main() {
	plugin.Plugin(Plugin)
}
