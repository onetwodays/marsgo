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
	tProfilesFieldNames          = builderx.RawFieldNames(&TProfiles{})
	tProfilesRows                = strings.Join(tProfilesFieldNames, ",")
	tProfilesRowsExpectAutoSet   = strings.Join(stringx.Remove(tProfilesFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tProfilesRowsWithPlaceHolder = strings.Join(stringx.Remove(tProfilesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TProfilesModel interface {
		Insert(data TProfiles) (sql.Result, error)
		FindOne(id int64) (*TProfiles, error)
		FindOneByUUIDVersion(uuid string,version string) (*TProfiles, error)
		Update(data TProfiles) error
		Delete(id int64) error
		DeleteByUUID(uuid string) error
	}

	defaultTProfilesModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TProfiles struct {
		Id      int64  `db:"id"`
		Uuid    string `db:"uuid"`    // uuid
		Version string `db:"version"` // prekey version

		Name       string         `db:"name"`   // name
		Avatar     sql.NullString `db:"avatar"` // 头像文件相对路径
		Commitment sql.NullString `db:"commitment"`
		CreateTime time.Time      `db:"create_time"`
		UpdateTime time.Time      `db:"update_time"`
	}
)

func NewTProfilesModel(conn sqlx.SqlConn) TProfilesModel {
	return &defaultTProfilesModel{
		conn:  conn,
		table: "`t_profiles`",
	}
}

func (m *defaultTProfilesModel) Insert(data TProfiles) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, tProfilesRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Uuid, data.Version, data.Name, data.Avatar, data.Commitment)
	return ret, err
}

func (m *defaultTProfilesModel) FindOne(id int64) (*TProfiles, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tProfilesRows, m.table)
	var resp TProfiles
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

func (m *defaultTProfilesModel) FindOneByUUIDVersion(uuid string,version string) (*TProfiles, error){
	query := fmt.Sprintf("select %s from %s where `uuid` = ? and `version`=? limit 1", tProfilesRows, m.table)
	var resp TProfiles
	err := m.conn.QueryRow(&resp, query, uuid,version)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTProfilesModel) Update(data TProfiles) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tProfilesRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Uuid, data.Version, data.Name, data.Avatar, data.Commitment, data.Id)
	return err
}

func (m *defaultTProfilesModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *defaultTProfilesModel) DeleteByUUID(uuid string) error{
	query := fmt.Sprintf("delete from %s where `uuid` = ?", m.table)
	_, err := m.conn.Exec(query, uuid)
	return err
}
