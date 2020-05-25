git config --global --unset http.proxy
git config --global --unset https.proxy

git config --global http.proxy  http://jiantuo:jt666%40fg@fg.hopex.com:13128
git config --global https.proxy http://jiantuo:jt666%40fg@fg.hopex.com:13128
ulimit -c unlimited 
export GOROOT=/usr/lib/go
export GOPATH=/home/iliu/project/go-repos/PATH
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN:$GOROOT/bin
export ETCDCTL_API=3
export GO111MODULE=on 
#export GOPROXY=https://mirrors.aliyun.com/goproxy/
#export GOPROXY=https://goproxy.io
export GOPROXY=https://goproxy.cn
# source .bashrc