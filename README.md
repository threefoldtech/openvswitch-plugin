> NOTE: This required core0 branch `container-plugins` checked out for this to build

## OVS container
Ovs container is a container with ovs capabilities. Once a container is created
from this image (with host_network=True), ovs capabilities will be plugged
into the container.

### Building
to build simply run
```bash
make
```
This will create output directly `image` which is an ubuntu root with openvswitch installed.
It also will have the required `.startup.toml` and `.plugin.toml`

Once the `image` is created you cat upload it to `hub.gig.tech`

### Supported commands
> Work in progress

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