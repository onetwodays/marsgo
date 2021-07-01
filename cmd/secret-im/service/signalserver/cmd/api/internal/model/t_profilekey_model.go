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
	tProfilekeyFieldNames          = builderx.RawFieldNames(&TProfilekey{})
	tProfilekeyRows                = strings.Join(tProfilekeyFieldNames, ",")
	tProfilekeyRowsExpectAutoSet   = strings.Join(stringx.Remove(tProfilekeyFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tProfilekeyRowsWithPlaceHolder = strings.Join(stringx.Remove(tProfilekeyFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TProfilekeyModel interface {
		Insert(data TProfilekey) (sql.Result, error)
		FindOne(id int64) (*TProfilekey, error)
		FindOneByAccountName(accountName string) (*TProfilekey, error)
		Update(data TProfilekey) error
		Delete(id int64) error
	}

	defaultTProfilekeyModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TProfilekey struct {
		Id          int64  `db:"id"`
		AccountName string `db:"account_name"`
		ProfileKey  string `db:"profile_key"`
	}
)

func NewTProfilekeyModel(conn sqlx.SqlConn) TProfilekeyModel {
	return &defaultTProfilekeyModel{
		conn:  conn,
		table: "`t_profilekey`",
	}
}

func (m *defaultTProfilekeyModel) Insert(data TProfilekey) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tProfilekeyRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.AccountName, data.ProfileKey)
	return ret, err
}

func (m *defaultTProfilekeyModel) FindOne(id int64) (*TProfilekey, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tProfilekeyRows, m.table)
	var resp TProfilekey
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

func (m *defaultTProfilekeyModel) FindOneByAccountName(accountName string) (*TProfilekey, error) {
	var resp TProfilekey
	query := fmt.Sprintf("select %s from %s where `account_name` = ? limit 1", tProfilekeyRows, m.table)
	err := m.conn.QueryRow(&resp, query, accountName)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTProfilekeyModel) Update(data TProfilekey) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tProfilekeyRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.AccountName, data.ProfileKey, data.Id)
	return err
}

func (m *defaultTProfilekeyModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
