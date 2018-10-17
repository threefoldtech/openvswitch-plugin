package ovs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBindAddCmd(t *testing.T) {
	bond := BondAddArguments{
		Bridge: Bridge{
			Bridge: "backplane",
		},
		Port:  "bond0",
		Links: []string{"enp0s0", "enp0s1"},
		LACP:  false,
		Mode:  "balance-slb",
		Options: map[string]string{
			"updelay": "2000",
			"foo":     "bar",
		},
	}

	actual, err := bondAddCmd(bond)
	require.NoError(t, err)
	// expected := []string{"add-bond", "backplane", "bond0", "enp0s0", "enp0s1", "bond_mode=balance-slb", "other_config:updelay=2000"}
	expected := []string{"add-bond", "backplane", "bond0", "enp0s0", "enp0s1", "bond_mode=balance-slb", "other_config:updelay=2000,foo=bar"}
	assert.Equal(t, expected, actual)
}
