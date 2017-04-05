package ovs

import (
	"encoding/json"
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
)

const (
	MultiCastGroup = 239
)

type VTepEnsureArguments struct {
	Bridge
	VNID uint `json:"vnid"`
}

func (t *VTepEnsureArguments) Validate() error {
	if err := t.Bridge.Validate(); err != nil {
		return err
	}
	if t.VNID == 0 {
		return fmt.Errorf("invalid nid")
	}
	return nil
}

func getGroupForVNID(vnid uint) net.IP {
	//VNID is 24 bit, that fits the last 3 octet of the MC group IP
	id := (vnid / 256) + 1

	ip := fmt.Sprintf("%d.%d.%d.%d",
		MultiCastGroup,
		id&0x00ff0000>>16,
		id&0x0000ff00>>8,
		id&0x000000ff,
	)

	return net.ParseIP(ip)
}

func VtepEnsure(args json.RawMessage) (interface{}, error) {
	var vtep VTepEnsureArguments
	if err := json.Unmarshal(args, &vtep); err != nil {
		return nil, err
	}

	if err := vtep.Validate(); err != nil {
		return nil, err
	}

	dev, err := netlink.LinkByName(vtep.Bridge.Bridge)

	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("vtep-%d", vtep.VNID)
	link, err := netlink.LinkByName(name)

	if err == nil {
		if link.Type() != "vxlan" {
			return nil, fmt.Errorf("invalid device type got '%s'", link.Type())
		}

		if link.(*netlink.Vxlan).VtepDevIndex != dev.Attrs().Index {
			return nil, fmt.Errorf("recrating same VNID with different vxbridge")
		}

		return nil, nil
	}

	vxlan := &netlink.Vxlan{
		LinkAttrs: netlink.LinkAttrs{
			Name:   name,
			Flags:  net.FlagBroadcast | net.FlagMulticast,
			TxQLen: -1,
		},
		VxlanId:      int(vtep.VNID),
		Group:        getGroupForVNID(vtep.VNID),
		VtepDevIndex: dev.Attrs().Index,
		Learning:     true,
	}

	if err := netlink.LinkAdd(vxlan); err != nil {
		return nil, err
	}

	return nil, netlink.LinkSetUp(vxlan)
}

type VTepDeleteArguments struct {
	VNID uint `json:"vnid"`
}

func (t *VTepDeleteArguments) Validate() error {
	if t.VNID == 0 {
		return fmt.Errorf("invalid nid")
	}
	return nil
}

func VtepDelete(args json.RawMessage) (interface{}, error) {
	var vtep VTepDeleteArguments
	if err := json.Unmarshal(args, &vtep); err != nil {
		return nil, err
	}

	if err := vtep.Validate(); err != nil {
		return nil, err
	}

	name := fmt.Sprintf("vtep-%d", vtep.VNID)
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	return nil, netlink.LinkDel(link)
}
