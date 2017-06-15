package ovs

import (
	"encoding/json"
	"fmt"
	"strings"
)

type VLanEnsureArguments struct {
	Master string `json:"master"`
	VLan   uint16 `json:"vlan"`
	Name   string `json:"name"`
}

func (v *VLanEnsureArguments) Validate() error {
	if v.Master == "" {
		return fmt.Errorf("master bridge not specified")
	}
	if v.VLan < 0 || v.VLan >= 4095 { //0 for untagged
		return fmt.Errorf("invalid vlan tag")
	}
	return nil
}

func portsList(br string) ([]string, error) {
	output, err := vsctl("list-ports", br)
	if err != nil {
		return nil, err
	}

	return strings.Fields(output), nil
}

func portToBridge(port string) (string, bool) {
	out, err := vsctl("port-to-br", port)
	if err != nil {
		return "", false
	}

	return strings.TrimSpace(out), true
}

func in(l []string, a string) bool {
	for _, i := range l {
		if i == a {
			return true
		}
	}
	return false
}

func VLanEnsure(args json.RawMessage) (interface{}, error) {
	//abstract method to ensure a bridge exists that has this vlan tag.
	var vlan VLanEnsureArguments
	if err := json.Unmarshal(args, &vlan); err != nil {
		return nil, err
	}

	if err := vlan.Validate(); err != nil {
		return nil, err
	}

	if !bridgeExists(vlan.Master) {
		return nil, fmt.Errorf("master bridge does not exist")
	}

	name := vlan.Name

	portName := fmt.Sprintf("vlbr%dp", vlan.VLan)
	portPeerName := fmt.Sprintf("vlbr%din", vlan.VLan)

	if br, ok := portToBridge(portName); ok {
		if br != vlan.Master {
			return nil, fmt.Errorf("reassigning vlang tag to another master bridge is not allowed")
		}
	}

	if br, ok := portToBridge(portPeerName); ok {
		//peer already exists.
		if name == "" {
			return br, nil
		} else if br != name {
			return nil, fmt.Errorf("reassigning vlan tag to another bridge not allowed")
		} else {
			//we already validated this setup.
			return name, nil
		}
	}

	if name == "" {
		name = fmt.Sprintf("vlbr%d", vlan.VLan)
	}

	if err := bridgeAdd(name, nil); err != nil {
		return nil, err
	}

	//add port in master
	if err := portAdd(&PortAddArguments{
		Bridge: Bridge{
			Bridge: vlan.Master,
		},
		Port: portName,
		VLan: vlan.VLan,
		Options: map[string]string{
			"type":         "patch",
			"options:peer": portPeerName,
		},
	}); err != nil {
		return nil, err
	}

	//connect port to vlan bridge
	if err := portAdd(&PortAddArguments{
		Bridge: Bridge{
			Bridge: name,
		},
		Port: portPeerName,
		Options: map[string]string{
			"type":         "patch",
			"options:peer": portName,
		},
	}); err != nil {
		return nil, err
	}

	return name, nil
}
