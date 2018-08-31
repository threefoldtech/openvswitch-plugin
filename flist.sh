
#!/bin/bash

set -ex

# setting up locales
if ! grep -q ^en_US /etc/locale.gen; then
    echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
    locale-gen
fi

apt-get update
apt-get install -y openvswitch-switch wget


wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz

mkdir -p $HOME/go/bin
mkdir -p $HOME/go/pkg
mkdir -p $HOME/go/src
mkdir -p $HOME/src/github.com/threefoldtech
mkdir -p /var/lib/corex/plugins
mkdir -p /run/openvswitch
mkdir -p /var/run/openvswitch

cp -r /openvswitch-plugin $HOME/src/github.com/threefoldtech

export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go

BINARY=ovs-plugin

cd $HOME/src/github.com/threefoldtech
go build -o $BINARY
cp $BINARY image/var/lib/corex/plugins
cp startup.toml image/.startup.toml
cp plugin.toml image/.plugin.toml

tar -cpzf "/tmp/archives/ovs.tar.gz" --exclude tmp --exclude dev --exclude sys --exclude proc /
