#!/usr/bin/env bash
GOVERSION="go1.14.linux-amd64"
cd ~
wget https://dl.google.com/go/${GOVERSION}.tar.gz
tar -C /usr/local -xzf ${GOVERSION}.tar.gz

mkdir go-path
echo "ulimit -c unlimited" >> ~/.bashrc
echo "export GOROOT=/usr/lib/${GOVERSION}/go" >> ~/.bashrc
echo "export GOPATH=~/go-path" >> ~/.bashrc
echo "export GOBIN=$GOROOT/bin" >> ~/.bashrc
echo "export PATH=$PATH:$GOBIN:$GOPATH/bin"  >> ~/.bashrc
echo "export ETCDCTL_API=3"  >> ~/.bashrc
echo "export GOPROXY=https://goproxy.cn"  >> ~/.bashrc
echo "export export GO111MODULE=on"  >> ~/.bashrc
source .bashrc