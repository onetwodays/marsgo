
#go get github.com/mattn/goreman  安装goreman

goreman start
# etcdctl put mykey "this is awesome"
# etcdctl get mykey
# goreman -f ./Procfile.learner start #Follow the steps in Procfile.learner to add a learner node to the cluster. Start the learner node with