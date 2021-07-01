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
	tPendingAccountsFieldNames          = builderx.RawFieldNames(&TPendingAccounts{})
	tPendingAccountsRows                = strings.Join(tPendingAccountsFieldNames, ",")
	tPendingAccountsRowsExpectAutoSet   = strings.Join(stringx.Remove(tPendingAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tPendingAccountsRowsWithPlaceHolder = strings.Join(stringx.Remove(tPendingAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TPendingAccountsModel interface {
		Insert(data TPendingAccounts) (sql.Result, error)
		FindOne(id int64) (*TPendingAccounts, error)
		FindOneByNumber(number string) (*TPendingAccounts, error)
		FindOneByVerificationCode(verificationCode string) (*TPendingAccounts, error)
		Update(data TPendingAccounts) error
		Delete(id int64) error
		DeleteByNumber(number string) error
	}

	defaultTPendingAccountsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TPendingAccounts struct {
		Id int64 `db:"id"` // pk

		Number           string       `db:"number"`            // 手机号
		VerificationCode string       `db:"verification_code"` // 验证码
		PushCode         string       `db:"push_code"`         // 推送码
		Timestamp        int64        `db:"timestamp"`
		CreateTime       time.Time    `db:"create_time"`
		UpdateTime       time.Time    `db:"update_time"`
		DeletedAt        sql.NullTime `db:"deleted_at"`
	}
)

func NewTPendingAccountsModel(conn sqlx.SqlConn) TPendingAccountsModel {
	return &defaultTPendingAccountsModel{
		conn:  conn,
		table: "`t_pending_accounts`",
	}
}

func (m *defaultTPendingAccountsModel) Insert(data TPendingAccounts) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, tPendingAccountsRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Number, data.VerificationCode, data.PushCode, data.Timestamp, data.DeletedAt)
	return ret, err
}

func (m *defaultTPendingAccountsModel) FindOne(id int64) (*TPendingAccounts, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tPendingAccountsRows, m.table)
	var resp TPendingAccounts
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

func (m *defaultTPendingAccountsModel) FindOneByNumber(number string) (*TPendingAccounts, error) {
	var resp TPendingAccounts
	query := fmt.Sprintf("select %s from %s where `number` = ? limit 1", tPendingAccountsRows, m.table)
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
func (m *defaultTPendingAccountsModel) DeleteByNumber(number string) error{
	query := fmt.Sprintf("delete from %s where `number` = ?", m.table)
	_, err := m.conn.Exec(query, number)
	return err
}

func (m *defaultTPendingAccountsModel) FindOneByVerificationCode(verificationCode string) (*TPendingAccounts, error) {
	var resp TPendingAccounts
	query := fmt.Sprintf("select %s from %s where `verification_code` = ? limit 1", tPendingAccountsRows, m.table)
	err := m.conn.QueryRow(&resp, query, verificationCode)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTPendingAccountsModel) Update(data TPendingAccounts) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tPendingAccountsRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Number, data.VerificationCode, data.PushCode, data.Timestamp, data.DeletedAt, data.Id)
	return err
}

func (m *defaultTPendingAccountsModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
