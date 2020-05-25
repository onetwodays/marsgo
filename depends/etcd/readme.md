goreman -f ./Procfile.learner start

go get github.com/mattn/goreman

export ETCDCTL_API=3

etcdctl get foo --print-value-only # 仅仅输出值