#!/bin/bash

### http://tecadmin.net/create-rpm-of-your-own-script-in-centosredhat/#

sudo yum install rpm-build rpmdevtools
rm -rf ~/rpmbuild
rpmdev-setuptree

mkdir ~/rpmbuild/SOURCES/hydra-reverse-proxy-0
cp ./fixtures/hydra-reverse-proxy.conf  ~/rpmbuild/SOURCES/hydra-reverse-proxy-0
cp hydra-reverse-proxy-init.d.sh ~/rpmbuild/SOURCES/hydra-reverse-proxy-0
cp ../../bin/hydra-reverse-proxy ~/rpmbuild/SOURCES/hydra-reverse-proxy-0

cp hydra-reverse-proxy.spec ~/rpmbuild/SPECS

pushd ~/rpmbuild/SOURCES/
tar czf hydra-reverse-proxy-0.1.tar.gz hydra-reverse-proxy-0/
cd ~/rpmbuild
rpmbuild -ba SPECS/hydra-reverse-proxy.spec

popd
cp ~/rpmbuild/RPMS/x86_64/hydra-reverse-proxy-0-1.x86_64.rpm .