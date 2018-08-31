#!/bin/bash
set -ex

apt-get update
apt-get install -y openvswitch-switch wget locales

# setting up locales
if ! grep -q ^en_US /etc/locale.gen; then
    echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
    locale-gen
fi

wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz

export GOPATH=$HOME/go
mkdir -p $GOPATH/bin
mkdir -p $GOPATH/pkg
mkdir -p $GOPATH/src
mkdir -p $GOPATH/src/github.com/threefoldtech
mkdir -p /var/lib/corex/plugins
mkdir -p /run/openvswitch
mkdir -p /var/run/openvswitch

cp -r /openvswitch-plugin $GOPATH/src/github.com/threefoldtech

export PATH=$PATH:/usr/local/go/bin

BINARY=ovs-plugin

cd $GOPATH/src/github.com/threefoldtech/openvswitch-plugin
go build -o $BINARY
cp $BINARY /var/lib/corex/plugins
cp startup.toml /.startup.toml
cp plugin.toml /.plugin.toml

mkdir -p /tmp/archives
tar -cpzf "/tmp/archives/ovs.tar.gz" --exclude tmp --exclude dev --exclude sys --exclude proc /
