package ovs

import (
	"encoding/json"
	"fmt"
)

type VXLanEnsureArguments struct {
	Master string `json:"master"`
	VXLan  uint   `json:"vxlan"`
	Name   string `json:"name"`
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
		Bridge: Bridge{vxlan.Master, nil},
		VNID:   vxlan.VXLan,
	})

	if err != nil {
		return nil, err
	}

	name := vxlan.Name

	if br, ok := portToBridge(vtep); ok {
		if name == "" {
			return br, nil
		} else if br != name {
			return nil, fmt.Errorf("reassigning vxlan tag to another bridge not allowed")
		} else {
			return name, nil
		}
	}

	if name == "" {
		name = fmt.Sprintf("vxlbr%d", vxlan.VXLan)
	}

	if err := bridgeAdd(name, nil); err != nil {
		return nil, err
	}

	//add port in vxlan bridge
	if err := portAdd(&PortAddArguments{
		Bridge: Bridge{name, nil},
		Port:   vtep,
	}); err != nil {
		return nil, err
	}

	return name, nil
}
