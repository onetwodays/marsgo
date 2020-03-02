#!/usr/bin/env bash
GOVERSION="go1.14.linux-amd64"
cd ~
wget https://dl.google.com/go/${GOVERSION}.tar.gz
sudo tar -C /usr/local -xzf ${GOVERSION}.tar.gz

mkdir -p ~/go-path/src
mkdir -p ~/go-path/bin
mkdir -p ~/go-path/pkg

echo "#git config  --global  http.proxy     http://jiantuo:jt666%40fg@fg.hopex.com:13128" >> ~/.bashrc
echo "#git config  --global  https.proxy   https://jiantuo:jt666%40fg@fg.hopex.com:13128" >> ~/.bashrc
echo "ulimit -c unlimited" >> ~/.bashrc
echo "export GOROOT=/usr/local/go" >> ~/.bashrc
echo "export GOPATH=~/go-path" >> ~/.bashrc
echo "export GOBIN=$GOROOT/bin" >> ~/.bashrc
echo "export PATH=$PATH:$GOBIN:$GOPATH/bin"  >> ~/.bashrc
echo "export ETCDCTL_API=3"  >> ~/.bashrc
echo "export export GO111MODULE=on"  >> ~/.bashrc
echo "export GOPROXY=https://goproxy.cn"  >> ~/.bashrc
echo "#export GOPROXY=https://mirrors.aliyun.com/goproxy/"  >> ~/.bashrc
echo "#export GOPROXY=https://goproxy.io"  >> ~/.bashrc
source ~/.bashrc