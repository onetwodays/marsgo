goctl api go -api privatedb.api -dir .
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="user" -dir ./model
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="t_msg" -dir ./model
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="t_pay_type" -dir ./model
goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/pb" -table="t_payment_account" -dir ./model

go run privatedb.go  -f ./etc/privatedb-api.yaml

