#goctl model mysql datasource -url="hopexdev:devhopex@tcp(127.0.0.1:3306)/privatedb" -table="user"  -dir .
#goctl model mysql ddl -c -src book.sql -dir . -style go_zero
goctl model mysql ddl  -src t_pend_accounts.sql -dir . -style go_zero