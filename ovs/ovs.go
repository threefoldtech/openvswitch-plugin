package ovs

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"os/exec"
)

const (
	Binary = "ovs-vsctl"
)

var (
	log = logging.MustGetLogger("ovs")
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

type Bridge struct {
	Bridge string `json:"bridge"`
}

func (b *Bridge) Validate() error {
	if b.Bridge == "" {
		return fmt.Errorf("bridge name is not set")
	}
	return nil
}

func BridgeAdd(args json.RawMessage) (interface{}, error) {
	var bridge Bridge
	if err := json.Unmarshal(args, &bridge); err != nil {
		return nil, err
	}

	if err := bridge.Validate(); err != nil {
		return nil, err
	}

	return vsctl("add-br", bridge.Bridge)
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
