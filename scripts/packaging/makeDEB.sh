#!/bin/bash

### http://linuxconfig.org/easy-way-to-create-a-debian-package-and-local-package-repository

rm -rf ~/debbuild
mkdir -p ~/debbuild/DEBIAN
cp control ~/debbuild/DEBIAN

mkdir -p ~/debbuild/etc
cp ./fixtures/hydra-reverse-proxy.conf ~/debbuild/etc

mkdir -p ~/debbuild/etc/init.d
cp hydra-reverse-proxy-init.d.sh ~/debbuild/etc/init.d/hydra-reverse-proxy

mkdir -p ~/debbuild/usr/local
cp ../../bin/hydra-reverse-proxy  ~/debbuild/usr/local

chmod -R 644 ~/debbuild/usr/local
chmod 755 ~/debbuild/etc/init.d/hydra-reverse-proxy
chmod 755 ~/debbuild/usr/local/hydra-reverse-proxy

sudo chown -R root:root ~/debbuild/*

pushd ~
sudo dpkg-deb --build debbuild

popd
sudo mv ~/debbuild.deb hydra-reverse-proxy-0-1.x86_64.deb
