package ovs

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

const (
	Binary = "ovs-vsctl"
)

func vsctl(args ...string) (string, error) {
	cmd := exec.Command(Binary, args...)
	data, err := cmd.Output()
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %s", err, string(err.Stderr))
		}
		return "", err
	}

	return string(data), nil
}

func bridgeExists(name string) bool {
	if _, err := vsctl("br-exists", name); err != nil {
		return false
	}

	return true
}

type Bridge struct {
	Bridge  string            `json:"bridge"`
	Options map[string]string `json:"options"`
}

func (b *Bridge) Validate() error {
	if b.Bridge == "" {
		return fmt.Errorf("bridge name is not set")
	}
	return nil
}

func bridgeAdd(name string, options map[string]string) error {
	commands := []string{"add-br", name}
	for key, value := range options {
		commands = append(commands, "--", "set", "Bridge", name, fmt.Sprintf("%s=%s", key, value))
	}
	_, err := vsctl(commands...)
	return err
}

func BridgeAdd(args json.RawMessage) (interface{}, error) {
	var bridge Bridge
	if err := json.Unmarshal(args, &bridge); err != nil {
		return nil, err
	}

	if err := bridge.Validate(); err != nil {
		return nil, err
	}

	return nil, bridgeAdd(bridge.Bridge, bridge.Options)
}

func BridgeDelete(args json.RawMessage) (interface{}, error) {
	var bridge Bridge
	if err := json.Unmarshal(args, &bridge); err != nil {
		return nil, err
	}

	if err := bridge.Validate(); err != nil {
		return nil, err
	}

	return vsctl("del-br", bridge.Bridge)
}

type PortAddArguments struct {
	Bridge
	Port    string            `json:"port"`
	VLan    uint16            `json:"vlan"`
	Options map[string]string `json:"options"`
}

func (p *PortAddArguments) Validate() error {
	if err := p.Bridge.Validate(); err != nil {
		return err
	}
	if p.Port == "" {
		return fmt.Errorf("missing port name")
	}
	return nil
}

func portAdd(port *PortAddArguments) error {
	var err error
	if port.VLan == 0 {
		_, err = vsctl("add-port", port.Bridge.Bridge, port.Port)
	} else {
		_, err = vsctl("add-port", port.Bridge.Bridge, port.Port, fmt.Sprintf("tag=%d", port.VLan))
	}

	if err != nil {
		return err
	}

	//setting options
	if len(port.Options) != 0 {
		return set(&SetArguments{
			Table:  "Interface",
			Record: port.Port,
			Values: port.Options,
		})
	}

	return nil
}

func PortAdd(args json.RawMessage) (interface{}, error) {
	var port PortAddArguments
	if err := json.Unmarshal(args, &port); err != nil {
		return nil, err
	}

	if err := port.Validate(); err != nil {
		return nil, err
	}

	return nil, portAdd(&port)
}

type PortDelArguments struct {
	Bridge
	Port string `json:"port"`
}

func (p *PortDelArguments) Validate() error {
	if p.Port == "" {
		return fmt.Errorf("missing port name")
	}
	return nil
}

func PortDel(args json.RawMessage) (interface{}, error) {
	var port PortDelArguments
	if err := json.Unmarshal(args, &port); err != nil {
		return nil, err
	}

	if err := port.Validate(); err != nil {
		return nil, err
	}
	var err error
	if port.Bridge.Bridge == "" {
		_, err = vsctl("del-port", port.Port)
	} else {
		_, err = vsctl("del-port", port.Bridge.Bridge, port.Port)
	}

	return nil, err
}

type SetArguments struct {
	Table  string            `json:"table"`
	Record string            `json:"record"`
	Values map[string]string `json:"values"`
}

func (s *SetArguments) Validate() error {
	if s.Table == "" {
		return fmt.Errorf("missing table name")
	}
	if s.Record == "" {
		return fmt.Errorf("missing record")
	}
	if len(s.Values) == 0 {
		return fmt.Errorf("no values to set")
	}

	return nil
}

func set(s *SetArguments) error {
	args := []string{"set", s.Table, s.Record}
	for key, value := range s.Values {
		args = append(args, fmt.Sprintf("%s=%s", key, value))
	}

	_, err := vsctl(args...)
	return err
}

func Set(args json.RawMessage) (interface{}, error) {
	var s SetArguments
	if err := json.Unmarshal(args, &s); err != nil {
		return nil, err
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}

	return nil, set(&s)
}

type BondMode string

const (
	BondModeActiveBackup = BondMode("active-backup")
	BondModeBalanceSLB   = BondMode("balance-slb")
	BondModeBalanceTCP   = BondMode("balance-tcp")
)

type BondAddArguments struct {
	Bridge
	Port  string   `json:"port"`
	Links []string `json:"links"`
	Mode  BondMode `json:"mode"`
	LACP  bool     `json:"lacp"`
}

func (b *BondAddArguments) Validate() error {
	if err := b.Bridge.Validate(); err != nil {
		return err
	}

	if b.Port == "" {
		return fmt.Errorf("missing port name")
	}

	if len(b.Links) <= 1 {
		return fmt.Errorf("need more than one link to bond")
	}

	return nil
}

func BondAdd(args json.RawMessage) (interface{}, error) {
	var bond BondAddArguments
	if err := json.Unmarshal(args, &bond); err != nil {
		return nil, err
	}

	if err := bond.Validate(); err != nil {
		return nil, err
	}
	mode := bond.Mode
	if mode == BondMode("") {
		mode = BondModeBalanceSLB
	}
	a := []string{"add-bond", bond.Bridge.Bridge, bond.Port}
	a = append(a, bond.Links...)
	if bond.LACP {
		a = append(a, "lacp=active")
	}

	a = append(a, fmt.Sprintf("bond_mode=%v", mode))

	for k, v := range bond.Options {
		a = append(a, fmt.Sprintf("%v=%v", k, v))
	}

	_, err := vsctl(a...)
	return nil, err
}
