go  build   -o ./build/bookstore.exe ./bookstore.go
echo "zh" |sudo -S chmod +x ./build/bookstore.exe
./build/bookstore.exe -f ./etc/bookstore.yaml
# ETCDCTL_API=3 etcdctl get bookstore.rpc   --prefix
