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
	tPendAccountsFieldNames          = builderx.RawFieldNames(&TPendAccounts{})
	tPendAccountsRows                = strings.Join(tPendAccountsFieldNames, ",")
	tPendAccountsRowsExpectAutoSet   = strings.Join(stringx.Remove(tPendAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tPendAccountsRowsWithPlaceHolder = strings.Join(stringx.Remove(tPendAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TPendAccountsModel interface {
		Insert(data TPendAccounts) (sql.Result, error)
		FindOne(id int64) (*TPendAccounts, error)
		FindOneByNumber(number string) (*TPendAccounts, error)
		Update(data TPendAccounts) error
		Delete(id int64) error
		DeleteByNumber(number string) error
	}

	defaultTPendAccountsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TPendAccounts struct {
		Id               int64     `db:"id"`     // 自增id
		Number           string    `db:"number"` // phone number
		VerificationCode string    `db:"verification_code"`
		CreateTime       time.Time `db:"create_time"`
	}
)

func NewTPendAccountsModel(conn sqlx.SqlConn) TPendAccountsModel {
	return &defaultTPendAccountsModel{
		conn:  conn,
		table: "`t_pend_accounts`",
	}
}

func (m *defaultTPendAccountsModel) Insert(data TPendAccounts) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, tPendAccountsRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Number, data.VerificationCode)
	return ret, err
}

func (m *defaultTPendAccountsModel) FindOne(id int64) (*TPendAccounts, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tPendAccountsRows, m.table)
	var resp TPendAccounts
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

func (m *defaultTPendAccountsModel) FindOneByNumber(number string) (*TPendAccounts, error) {
	var resp TPendAccounts
	query := fmt.Sprintf("select %s from %s where `number` = ? limit 1", tPendAccountsRows, m.table)
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

func (m *defaultTPendAccountsModel) Update(data TPendAccounts) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tPendAccountsRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Number, data.VerificationCode, data.Id)
	return err
}

func (m *defaultTPendAccountsModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
func (m *defaultTPendAccountsModel) DeleteByNumber(number string) error {
	query := fmt.Sprintf("delete from %s where `number` = '?'", m.table)
	_, err := m.conn.Exec(query, number)
	return err
}
