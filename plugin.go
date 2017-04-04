package main

import (
	"github.com/g8os/core0/base/plugin"
	"github.com/g8os/core0/base/pm/core"
)

const (
	Version = "1.0alpha"
)

func version(cmd *core.Command) (interface{}, error) {
	return Version, nil
}

var (
	//plugin manifest
	Manifest = plugin.Manifest{
		Domain:  "ovs",
		Version: plugin.Version_1,
	}

	Plugin = plugin.Commands{
		"version": version,
	}
)
