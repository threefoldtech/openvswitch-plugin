#!/bin/bash
set -e

copypkg() {
   pkgname=$1
   target=$2
   for file in $(dpkg -L $pkgname); do
      if [ -d "$file" ]; then
         continue
      fi
      dirname=$(dirname $file)
      if [[ "$dirname" == /usr/share/man* ]]; then
         continue
      fi
      if [[ "$dirname" == /usr/share/doc* ]]; then
         continue
      fi
      targetdir="${target}${dirname}"
      mkdir -p "$targetdir"
      cp $file "$targetdir"
      if file "$file" | grep dynamic; then
         lddcopy "$file" "$target"
      fi
   done
}

apt-get update
apt-get install -y openvswitch-switch curl git

if ! which lddcopy; then
    pushd /tmp
    git clone --depth=1 https://github.com/maxux/lddcopy.git
    cp /tmp/lddcopy/lddcopy.sh /usr/local/bin/lddcopy
    chmod +x /usr/local/bin/lddcopy
    popd
fi

if ! which go; then
    curl https://dl.google.com/go/go1.11.linux-amd64.tar.gz > /tmp/go1.11.linux-amd64.tar.gz
    tar -C /usr/local -xzf /tmp/go1.11.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
fi
mkdir -p /gopath
export GOPATH=/gopath

OVSPLUGIN=$GOPATH/src/github.com/threefoldtech/openvswitch-plugin

go get -u -v -d github.com/threefoldtech/openvswitch-plugin


TARGET=/tmp/target-ovs

mkdir -p /tmp/archives
rm -rf "$TARGET"
mkdir "$TARGET"
mkdir -p "$TARGET/var/lib/corex/plugins"
mkdir -p "$TARGET/bin"
mkdir -p "$TARGET/etc/openvswitch"
mkdir -p "$TARGET/run/openvswitch"
mkdir -p "$TARGET/var/run/openvswitch"
mkdir -p "$TARGET/tmp"

pushd $OVSPLUGIN
go build -o "$TARGET/var/lib/corex/plugins/ovs-plugin"
cp startup.toml "$TARGET/.startup.toml"
cp plugin.toml "$TARGET/.plugin.toml"
popd


copypkg openvswitch-switch "$TARGET"
copypkg openvswitch-common "$TARGET"
pushd "$TARGET/usr/sbin"
ln -s /usr/lib/openvswitch-switch/ovs-vswitchd
popd

cd "$TARGET"
tar -czf /tmp/archives/ovs.tar.gz .
