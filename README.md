> NOTE: This required core0 branch `container-plugins` checked out for this to build

## OVS container
Ovs container is a container with ovs capabilities. Once a container is created
from this image (with host_network=True), ovs capabilities will be plugged
into the container.

For more details see the [Introduction to Zero-OS](/docs/README.md) in the [`/docs`](/docs) documentation directory, which includes a comprehensive [table of contents](/docs/SUMMARY.md).

For more details see the [For more details see the [Introduction to the Open vSwitch Plugin](/docs/README.md) in the [`/docs`](/docs) documentation directory, which includes a [table of contents](/docs/SUMMARY.md).


### Building
to build simply run
```bash
make
```
This will create output directly `image` which is an ubuntu root with openvswitch installed.
It also will have the required `.startup.toml` and `.plugin.toml`

Once the `image` is created you cat upload it to `hub.gig.tech`

## Supported low level commands

#### ovs.bridge-add
Add a new VSwitch
```javascript
{
	"bridge": "name", //bridge name
}
```

### ovs.bridge-del
Delete a VSwitch
```javascript
{
	"bridge": "name", //bridge name
}
```

### ovs.port-add
Add a new port to a switch
```javascript
{
	"bridge": "name", //bridge name
	"port": "port name", //port name
	"vlan": 10, //vlan tag (0 for no vlan tag)
	"options": { //extra options to set on the port (optional)
		"type": "patch", //set type option
		"options:peer": "peer name", //set peer name option
	}
}
```

### ovs.port-del
Add a port to a bridge
```javascript
{
	"bridge": "name", //bridge name [optional]
	"port": "port name", //port name
}
```

### ovs.bond-add
Add a bond to a bridge
```javascript
{
	"bridge": "name", //bridge name
	"port": "bond-name",
	"links": ["eth0", "ethc1", "..."], //link names to bond
	"mode": "mode", //bond modes (active-backup, balance-slb, balance-tcp)
	"lacp": bool,
}
```

### ovs.set
Set attributes in a table
```javascript
{
	"table": "table-name", //table name (Interface, etc...)
	"record": "record-key", //record key (port name, etc...)
	"options": {
		"key", "value", // (ex: "type": "patch")
	}
}
```

## Supported high level commands
Those command abstract some networking operations

### ovs.vlan_ensure
Create a tagged bridge with given vlan tag
```javascript
{
	"master": "master-bridge-name", //ex: backplane
	"vlan": tag, //vlan tag (0==untagged, or 1 to 4094)
	"name": "create-bridge-name" //[optional], if not given bridge name will be "vlbr[tag]"
}
```
returns the created bridge name


### ovs.vxlan_ensure
Create a vxlan bridge with given vxlan id
```javascript
{
	"master": "master-bridge-name", //ex: vxbackend
	"vxlan": id, //vxlan id
}
```
returns the created bridge name
