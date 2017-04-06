package ovs

import (
	"encoding/json"
	"fmt"
)

type VXLanEnsureArguments struct {
	Master string `json:"master"`
	VXLan  uint   `json:"vxlan"`
}

func (v *VXLanEnsureArguments) Validate() error {
	if v.Master == "" {
		return fmt.Errorf("master bridge not specified")
	}

	return nil
}

//Creates a bridge that connectes to master bridge over a vtep with the given vxlan id
func VXLanEnsure(args json.RawMessage) (interface{}, error) {
	//abstract method to ensure a bridge exists that has this vlan tag.
	var vxlan VXLanEnsureArguments
	if err := json.Unmarshal(args, &vxlan); err != nil {
		return nil, err
	}

	if err := vxlan.Validate(); err != nil {
		return nil, err
	}

	if !bridgeExists(vxlan.Master) {
		return nil, fmt.Errorf("master bridge does not exist")
	}

	vtep, err := vtepEnsure(&VTepEnsureArguments{
		Bridge: Bridge{vxlan.Master},
		VNID:   vxlan.VXLan,
	})

	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("vxlbr%d", vxlan.VXLan)
	if bridgeExists(name) {
		//TODO: validate that the port is also added and of correct vlan tag.
		return name, nil
	}

	if err := bridgeAdd(name); err != nil {
		return nil, err
	}

	//add port in vxlan bridge
	if err := portAdd(&PortAddArguments{
		Bridge: Bridge{name},
		Port:   vtep,
	}); err != nil {
		return nil, err
	}

	return name, nil
}
