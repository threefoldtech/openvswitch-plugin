package main

import (
	"github.com/zero-os/0-core/base/plugin"
	"github.com/zero-os/openvswitch-plugin/ovs"
)

var (
	commands = plugin.Commands{
		"bridge-add":   ovs.BridgeAdd,
		"bridge-del":   ovs.BridgeDelete,
		"port-add":     ovs.PortAdd,
		"port-del":     ovs.PortDel,
		"bond-add":     ovs.BondAdd,
		"vtep-ensure":  ovs.VtepEnsure,
		"vtep-del":     ovs.VtepDelete,
		"vlan-ensure":  ovs.VLanEnsure,
		"vxlan-ensure": ovs.VXLanEnsure,
		"set":          ovs.Set,
	}
)

func main() {
	plugin.Plugin(commands)
}
