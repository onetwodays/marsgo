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
	tUsernamesFieldNames          = builderx.RawFieldNames(&TUsernames{})
	tUsernamesRows                = strings.Join(tUsernamesFieldNames, ",")
	tUsernamesRowsExpectAutoSet   = strings.Join(stringx.Remove(tUsernamesFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tUsernamesRowsWithPlaceHolder = strings.Join(stringx.Remove(tUsernamesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TUsernamesModel interface {
		Insert(data TUsernames) (sql.Result, error)
		FindOne(id int64) (*TUsernames, error)
		FindOneByUsername(username string) (*TUsernames, error)
		FindOneByUuid(uuid string) (*TUsernames, error)
		Update(data TUsernames) error
		Delete(id int64) error
		DeleteByUuid(uuid string) error
	}

	defaultTUsernamesModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TUsernames struct {
		Id         int64     `db:"id"`
		Uuid       string    `db:"uuid"`
		Username   string    `db:"username"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewTUsernamesModel(conn sqlx.SqlConn) TUsernamesModel {
	return &defaultTUsernamesModel{
		conn:  conn,
		table: "`t_usernames`",
	}
}

func (m *defaultTUsernamesModel) Insert(data TUsernames) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tUsernamesRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Uuid, data.Username)
	return ret, err
}

func (m *defaultTUsernamesModel) FindOne(id int64) (*TUsernames, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tUsernamesRows, m.table)
	var resp TUsernames
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

func (m *defaultTUsernamesModel) FindOneByUsername(username string) (*TUsernames, error) {
	var resp TUsernames
	query := fmt.Sprintf("select %s from %s where `username` = ? limit 1", tUsernamesRows, m.table)
	err := m.conn.QueryRow(&resp, query, username)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTUsernamesModel) FindOneByUuid(uuid string) (*TUsernames, error) {
	var resp TUsernames
	query := fmt.Sprintf("select %s from %s where `uuid` = ? limit 1", tUsernamesRows, m.table)
	err := m.conn.QueryRow(&resp, query, uuid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTUsernamesModel) Update(data TUsernames) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tUsernamesRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Uuid, data.Username, data.Id)
	return err
}

func (m *defaultTUsernamesModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *defaultTUsernamesModel)  DeleteByUuid(uuid string) error{
	query := fmt.Sprintf("delete from %s where `uuid` = ?", m.table)
	_, err := m.conn.Exec(query, uuid)
	return err
}
