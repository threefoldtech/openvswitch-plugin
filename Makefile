path = "/sbin:/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin"

all: ovs ovs.so
	sudo mkdir -p image/var/lib/corex/plugins
	sudo cp ovs.so image/var/lib/corex/plugins
	sudo cp startup.toml image/.startup.toml

ovs.so:
	go build -buildmode=plugin -o ovs.so

image:
	sudo debootstrap xenial image http://ftp.belnet.be/ubuntu.com/ubuntu/

ovs: image
	sudo chroot image /bin/bash -c 'PATH=${path} apt-get update'
	sudo chroot image /bin/bash -c 'PATH=${path} apt-get install -y openvswitch-switch'

.PHONY: ovs ovs.so