package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	tKeysFieldNames          = builderx.RawFieldNames(&TKeys{})
	tKeysRows                = strings.Join(tKeysFieldNames, ",")
	tKeysRowsExpectAutoSet   = strings.Join(stringx.Remove(tKeysFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tKeysRowsWithPlaceHolder = strings.Join(stringx.Remove(tKeysFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TKeysModel interface {
		Insert(data TKeys) (sql.Result, error)
		FindOne(id int64) (*TKeys, error)
		Update(data TKeys) error
		Delete(id int64) error
		DeleteMany(number string,deviceId int64) error
		FindMany (number string,deviceId int64) ([]TKeys,error)
		FindManyFirst (number string,deviceId int64) (*TKeys,error)
	}

	defaultTKeysModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TKeys struct {
		Id         int64  `db:"id"`
		Number     string `db:"number"`
		Keyid      int64  `db:"keyid"`
		Publickey  string `db:"publickey"`
		LastResort int64  `db:"last_resort"`
		Deviceid   int64  `db:"deviceid"`
	}
)

func NewTKeysModel(conn sqlx.SqlConn) TKeysModel {
	return &defaultTKeysModel{
		conn:  conn,
		table: "`t_keys`",
	}
}

func (m *defaultTKeysModel) Insert(data TKeys) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, tKeysRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Number, data.Keyid, data.Publickey, data.LastResort, data.Deviceid)
	return ret, err
}

func (m *defaultTKeysModel) FindOne(id int64) (*TKeys, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tKeysRows, m.table)
	var resp TKeys
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}


func (m *defaultTKeysModel) FindMany (number string,deviceId int64) ([]TKeys,error) {
	query := fmt.Sprintf("select %s from %s where 1=1 ", tKeysRows, m.table)
	if len(number)>0{
		query = fmt.Sprintf("%s and number='%s'  ",query,number)
	}
	if deviceId >=0{
		query = fmt.Sprintf("%s and deviceid=%d  ",query,deviceId)
	}
	query+=" order by id desc "
	var resp []TKeys
	err := m.conn.QueryRows(&resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}


func (m *defaultTKeysModel) FindManyFirst (number string,deviceId int64) (*TKeys,error){
	query := fmt.Sprintf("select %s from %s where 1=1 ", tKeysRows, m.table)
	if len(number)>0{
		query = fmt.Sprintf("%s and number='%s'  ",query,number)
	}
	if deviceId >=0{
		query = fmt.Sprintf("%s and deviceid=%d  ",query,deviceId)
	}
	query+=" order by id ASC limit 1 "
	var resp []TKeys
	err := m.conn.QueryRows(&resp, query)
	switch err {
	case nil:
		return &resp[0], nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}


func (m *defaultTKeysModel) Update(data TKeys) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tKeysRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Number, data.Keyid, data.Publickey, data.LastResort, data.Deviceid, data.Id)
	return err
}

func (m *defaultTKeysModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *defaultTKeysModel) DeleteMany(number string,deviceId int64) error{
	query := fmt.Sprintf("delete from %s where `number` = '%s' and deviceid=%d", m.table,number,deviceId)
	_, err := m.conn.Exec(query)
	return err

}


