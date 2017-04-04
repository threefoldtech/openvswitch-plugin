path = "/sbin:/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin"
BINARY = ovs-plugin

all: ovs $(BINARY)
	sudo mkdir -p image/var/lib/corex/plugins
	sudo cp $(BINARY) image/var/lib/corex/plugins
	sudo cp startup.toml image/.startup.toml
	sudo cp plugin.toml image/.plugin.toml

$(BINARY):
	go build -o $(BINARY)

image:
	sudo debootstrap xenial image http://ftp.belnet.be/ubuntu.com/ubuntu/

ovs: image
	sudo chroot image /bin/bash -c 'PATH=${path} apt-get update'
	sudo chroot image /bin/bash -c 'PATH=${path} apt-get install -y openvswitch-switch'

.PHONY: ovs $(BINARY)