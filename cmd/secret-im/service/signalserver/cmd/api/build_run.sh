go  build   -o ./deploy/signalserver.exe ./signalserver.go
echo "zh" |sudo -S chmod +x ./deploy/signalserver.exe
./deploy/signalserver.exe -f ./etc/signalserver-api.yaml