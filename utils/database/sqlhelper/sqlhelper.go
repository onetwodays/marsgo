package sqlhelper

import (
    "context"
    "database/sql"
)

type DBConnectPool struct {
    pool *sql.DB
}

// 申请一个连接池
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
/*
  Except for the databasename, all values are optional. So the minimal DSN is:

/dbname
If you do not want to preselect a database, leave dbname empty:

/
This has the same effect as an empty DSN string:


 */
func NewDBConnectPool(driver,dsn string) (*DBConnectPool,error)  {
    db, err := sql.Open(driver, dsn)
    connectPool:=&DBConnectPool{
        pool:db ,
    }
    return connectPool,err
}

func (receiver *DBConnectPool) Close() error{
    return receiver.pool.Close()
}

func (receiver *DBConnectPool)  PingContext(ctx context.Context) error {
    return  receiver.pool.PingContext(ctx)
}

func (receiver *DBConnectPool)  Ping() error {
    return  receiver.pool.Ping()
}

func (receiver *DBConnectPool) Query(query string,args ...interface{}) (*sql.Rows,error){
    return receiver.Query(query,args)
}

func (receiver *DBConnectPool) QueryContext(ctx context.Context,query string,args ...interface{}) (*sql.Rows,error){
    return receiver.QueryContext(ctx,query,args)
}

func (receiver *DBConnectPool) QueryRow(query string,args ...interface{}) (*sql.Rows,error){
    return receiver.QueryRow(query,args)
}

func (receiver *DBConnectPool) QueryRowContext(ctx context.Context,query string,args ...interface{}) (*sql.Rows,error){
    return receiver.QueryRowContext(ctx,query,args)
}

func (receiver *DBConnectPool) Exec(query string,args ...interface{})(sql.Result,error){
    return receiver.Exec(query,args)
}

func (receiver *DBConnectPool) ExecContext(ctx context.Context,query string,args ...interface{})(sql.Result,error){
    return receiver.ExecContext(ctx,query,args)
}



/*
。如果需要批量插入一堆数据，就可以使用Prepared语句。
golang处理prepared语句有其独特的行为，了解其底层的实现，对于用好它十分重要
创建stmt的preprea方式是golang的一个设计，其目的是Prepare once, execute many times。
为了批量执行sql语句。但是通常会造成所谓的三次网络请求（ three network round-trips）。
即preparing executing和closing三次请求。
PrepareContext creates a prepared statement for later queries or executions.
Multiple queries or executions may be run concurrently from the returned statement.
The caller must call the statement's Close method when the statement is no longer needed
    func (s *Stmt) Close() error
    func (s *Stmt) Exec(args ...interface{}) (Result, error)
    func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (Result, error)
    func (s *Stmt) Query(args ...interface{}) (*Rows, error)
    func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error)
    func (s *Stmt) QueryRow(args ...interface{}) *Row
    func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d", squareNum)

	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d", squareNum)
}

 */
func (receiver *DBConnectPool) Prepare(query string) (*sql.Stmt, error){
    return receiver.Prepare(query)
}
func (receiver *DBConnectPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error){
    return receiver.PrepareContext(ctx,query)
}


/*
   func (tx *Tx) Commit() error
   func (tx *Tx) Exec(query string, args ...interface{}) (Result, error)
   func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
   func (tx *Tx) Prepare(query string) (*Stmt, error)
   func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error)
   func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error)
   func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
   func (tx *Tx) QueryRow(query string, args ...interface{}) *Row
   func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
   func (tx *Tx) Rollback() error
   func (tx *Tx) Stmt(stmt *Stmt) *Stmt
   func (tx *Tx) StmtContext(ctx context.Context, stmt *Stmt) *Stmt
 */

func (receiver *DBConnectPool) Begin() (*sql.Tx, error){
    return receiver.pool.Begin()
}
func (receiver *DBConnectPool) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error){
    return receiver.pool.BeginTx(ctx,opts)

}











