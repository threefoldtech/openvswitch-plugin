package ovs

import (
	"encoding/json"
	"fmt"
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
	if name == "" {
		name = fmt.Sprintf("vlbr%d", vlan.VLan)
	}

	if bridgeExists(name) {
		//TODO: validate that the port is also added and of correct vlan tag.
		return name, nil
	}

	if err := bridgeAdd(name); err != nil {
		return nil, err
	}

	//add port in master
	if err := portAdd(&PortAddArguments{
		Bridge: Bridge{
			Bridge: vlan.Master,
		},
		Port: fmt.Sprintf("vlbr%dp", vlan.VLan),
		VLan: vlan.VLan,
		Options: map[string]string{
			"type":         "patch",
			"options:peer": fmt.Sprintf("vlbr%din", vlan.VLan),
		},
	}); err != nil {
		return nil, err
	}

	//connect port to vlan bridge
	if err := portAdd(&PortAddArguments{
		Bridge: Bridge{
			Bridge: name,
		},
		Port: fmt.Sprintf("vlbr%din", vlan.VLan),
		Options: map[string]string{
			"type":         "patch",
			"options:peer": fmt.Sprintf("vlbr%dp", vlan.VLan),
		},
	}); err != nil {
		return nil, err
	}

	return name, nil
}
