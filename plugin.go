package main

import (
	"encoding/json"
	"github.com/g8os/core0/base/plugin"
)

const (
	Version = "1.0alpha"
)

func version(input json.RawMessage) (interface{}, error) {
	return Version, nil
}

var (
	Plugin = plugin.Commands{
		"version": version,
	}
)

func main() {
	plugin.Plugin(Plugin)
}
