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
	tAccountsFieldNames          = builderx.RawFieldNames(&TAccounts{})
	tAccountsRows                = strings.Join(tAccountsFieldNames, ",")
	tAccountsRowsExpectAutoSet   = strings.Join(stringx.Remove(tAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tAccountsRowsWithPlaceHolder = strings.Join(stringx.Remove(tAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TAccountsModel interface {
		Insert(data TAccounts) (sql.Result, error)
		FindOne(id int64) (*TAccounts, error)
		FindOneByNumber(number string) (*TAccounts, error)
		Update(data TAccounts) error
		Delete(id int64) error
		DeleteByNumber(number string) error
	}

	defaultTAccountsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TAccounts struct {
		Data   string `db:"data"` // account json byte
		Id     int64  `db:"id"`
		Number string `db:"number"` // phonenumber
	}
)

func NewTAccountsModel(conn sqlx.SqlConn) TAccountsModel {
	return &defaultTAccountsModel{
		conn:  conn,
		table: "`t_accounts`",
	}
}

func (m *defaultTAccountsModel) Insert(data TAccounts) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tAccountsRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Data, data.Number)
	return ret, err
}

func (m *defaultTAccountsModel) FindOne(id int64) (*TAccounts, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tAccountsRows, m.table)
	var resp TAccounts
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

func (m *defaultTAccountsModel) FindOneByNumber(number string) (*TAccounts, error) {
	var resp TAccounts
	query := fmt.Sprintf("select %s from %s where `number` = ? limit 1", tAccountsRows, m.table)
	err := m.conn.QueryRow(&resp, query, number)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTAccountsModel) Update(data TAccounts) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tAccountsRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Data, data.Number, data.Id)
	return err
}

func (m *defaultTAccountsModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *defaultTAccountsModel) DeleteByNumber(number string) error {
	query := fmt.Sprintf("delete from %s where `number` = '?'", m.table)
	_, err := m.conn.Exec(query, number)
	return err
}
