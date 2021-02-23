 goctl api go -api cmd/signalserver.api -dir . --style go_zero #signalserver目录下执行
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="user" -dir ./model
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="t_msg" -dir ./model
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="t_pay_type" -dir ./model
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="t_payment_account" -dir ./model

go run signalserver.go  -f ./etc/signalserver-api.yaml

