package cassa

import (

	"secret-im/pkg/utils-tools"

	"github.com/gocassa/gocassa"
	"github.com/gocql/gocql"
)

// 模型接口
type Table interface {
	TableName() string
	Keys() gocassa.Keys
}

// 连接选项
type Options struct {
	NodeIPs  []string
	KeySpace string
	Username string
	Password string
	NumConns int
}

// 连接实例
type Conn struct {
	gocassa.KeySpace
	sess   *gocql.Session
	tables map[string]gocassa.Table
}

// 关闭连接
func (conn Conn) Close() {
	conn.sess.Close()
}

// 确保表存在
func (conn Conn) Ensure(tables []Table) error {
	for i := 0; i < len(tables); i++ {
		table := tables[i]
		t := conn.Table(table.TableName(), table, table.Keys())
		if err := t.CreateIfNotExist(); err != nil {
			return err
		}
		conn.tables[utils.GetObjectType(table)] = t
	}
	return nil
}

// 指定模型
func (conn Conn) Model(table Table) gocassa.Table {
	return conn.tables[utils.GetObjectType(table)]
}

// 获取表名
func (conn Conn) TableName(table Table) string {
	return conn.KeySpace.Name() + "." + conn.Model(table).Name()
}

// 生成查询对象
func (conn Conn) Query(stmt string, values ...interface{}) *gocql.Query {
	return conn.sess.Query(stmt, values...)
}

// 打开连接
func Open(options Options) (Conn, error) {
	if options.NumConns == 0 {
		options.NumConns = 2
	}

	qe, sess, err := newGoCQLBackend(options)
	if err != nil {
		return Conn{}, err
	}
	conn := gocassa.NewConnection(qe)

	err = conn.CreateKeySpace(options.KeySpace)
	if err != nil {
		if _, ok := err.(*gocql.RequestErrAlreadyExists); !ok {
			return Conn{}, err
		}
	}

	keySpace := conn.KeySpace(options.KeySpace)
	keySpace.DebugMode(false)

	ret := Conn{
		KeySpace: keySpace,
		sess:     sess,
		tables:   make(map[string]gocassa.Table),
	}
	return ret, nil
}

func newGoCQLBackend(options Options) (gocassa.QueryExecutor, *gocql.Session, error) {
	cluster := gocql.NewCluster(options.NodeIPs...)
	cluster.NumConns = options.NumConns
	cluster.Consistency = gocql.One
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: options.Username,
		Password: options.Password,
	}
	sess, err := cluster.CreateSession()
	if err != nil {
		return nil, nil, err
	}
	return gocassa.GoCQLSessionToQueryExecutor(sess), sess, nil
}
