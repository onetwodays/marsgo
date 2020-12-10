package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	tPayTypeFieldNames          = builderx.FieldNames(&TPayType{})
	tPayTypeRows                = strings.Join(tPayTypeFieldNames, ",")
	tPayTypeRowsExpectAutoSet   = strings.Join(stringx.Remove(tPayTypeFieldNames, "create_time", "update_time"), ",")
	tPayTypeRowsWithPlaceHolder = strings.Join(stringx.Remove(tPayTypeFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
)

type (
	TPayTypeModel interface {
		Insert(data TPayType) (sql.Result, error)
		FindOne(id uint32) (*TPayType, error)
		FindAll() ([]TPayType,error)
		Update(data TPayType) error
		Delete(id uint32) error

	}

	defaultTPayTypeModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TPayType struct {
		Id         int64     `db:"id" json:"id"`   // pk
		Name       string    `db:"name" json:"name"` // desciption
		CreateTime time.Time `db:"create_time" json:"ctime"`
		UpdateTime time.Time `db:"update_time" json:"utime"`
	}
)

func NewTPayTypeModel(conn sqlx.SqlConn) TPayTypeModel {
	return &defaultTPayTypeModel{
		conn:  conn,
		table: "t_pay_type",
	}
}

func (m *defaultTPayTypeModel) Insert(data TPayType) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tPayTypeRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Id, data.Name)
	return ret, err
}

func (m *defaultTPayTypeModel) FindOne(id uint32) (*TPayType, error) {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", tPayTypeRows, m.table)
	var resp TPayType
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

func (m *defaultTPayTypeModel) FindAll() ([]TPayType,error){
	query := fmt.Sprintf("select %s from %s order by id desc", tPayTypeRows, m.table)
	var resp []TPayType
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

func (m *defaultTPayTypeModel) Update(data TPayType) error {
	query := fmt.Sprintf("update %s set %s where id = ?", m.table, tPayTypeRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Name, data.Id)
	return err
}

func (m *defaultTPayTypeModel) Delete(id uint32) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
