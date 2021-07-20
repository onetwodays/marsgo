package model

import (
	"database/sql"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	tAccountsFieldNames             = builderx.RawFieldNames(&TAccounts{})
	tAccountsRows                   = strings.Join(tAccountsFieldNames, ",")
	tAccountsRowsExpectAutoSet      = strings.Join(stringx.Remove(tAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tAccountsRowsWithPlaceHolder    = strings.Join(stringx.Remove(tAccountsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	//tAccountsRowsWithPlaceHolderNew = strings.Join(stringx.Remove(tAccountsFieldNames, "`id`", "`create_time`", "`delete_time`"), "=?,") + "=?"
)

type (
	TAccountsModel interface {
		Insert(data TAccounts) (sql.Result, error)
		FindOne(id int64) (*TAccounts, error)
		FindOneByNumber(number string) (*TAccounts, error)
		FindManyByNumbers (numbers []string) ([]TAccounts,error)
		FindOneByUuid(uuid string) (*TAccounts, error)
		FindManyByUuids (uuids []string) ([]TAccounts,error)
		Update(data TAccounts) error
		UpdateDataByUuid(data sql.NullString,uuid string) error
		Delete(id int64) error
	}

	defaultTAccountsModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TAccounts struct {
		Id         int64          `db:"id"` // pk
		Number     string         `db:"number"`
		Uuid       string         `db:"uuid"`
		Data       sql.NullString `db:"data"`
		CreateTime time.Time      `db:"create_time"`
		UpdateTime time.Time      `db:"update_time"`
		DeleteTime sql.NullTime   `db:"delete_time"`
	}
)

func NewTAccountsModel(conn sqlx.SqlConn) TAccountsModel {
	return &defaultTAccountsModel{
		conn:  conn,
		table: "`t_accounts`",
	}
}

func (m *defaultTAccountsModel) Insert(data TAccounts) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, tAccountsRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Number, data.Uuid, data.Data, data.DeleteTime)
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

func (m *defaultTAccountsModel) FindOneByUuid(uuid string) (*TAccounts, error) {
	var resp TAccounts
	query := fmt.Sprintf("select %s from %s where `uuid` = ? limit 1", tAccountsRows, m.table)
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

func (m *defaultTAccountsModel) FindManyByUuids (uuids []string) ([]TAccounts,error){
	var resp []TAccounts
	var ps []string
	for i:=range uuids{
		ps=append(ps,fmt.Sprintf(`'%s'`,uuids[i]))
	}
	query := fmt.Sprintf("select %s from %s where `uuid` in   ", tAccountsRows, m.table)
	query+= ` (`+strings.Join(ps,",")+`) `
	logx.Info(query)
	err:=m.conn.QueryRows(&resp,query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTAccountsModel) FindManyByNumbers (numbers []string) ([]TAccounts,error){
	var resp []TAccounts
	query := fmt.Sprintf("select %s from %s where `number` in ?  ", tAccountsRows, m.table)
	err:=m.conn.QueryRows(&resp,query,` (`+strings.Join(numbers,",")+`) `)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTAccountsModel) Update(data TAccounts) error {
	query := fmt.Sprintf("update %s set %s,update_time=%s where `id` = ?", m.table, tAccountsRowsWithPlaceHolder,"CURRENT_TIMESTAMP")
	_, err := m.conn.Exec(query, data.Number, data.Uuid, data.Data, data.DeleteTime, data.Id)
	return err
}

func (m *defaultTAccountsModel) UpdateDataByUuid(data sql.NullString,uuid string) error  {
	query := fmt.Sprintf("update %s set data=?,update_time=CURRENT_TIMESTAMP where `uuid` = ?", m.table)
	_, err := m.conn.Exec(query, data,uuid)
	return err

}



func (m *defaultTAccountsModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
