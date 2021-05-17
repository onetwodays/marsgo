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
	tKeysFieldNames          = builderx.RawFieldNames(&TKeys{})
	tKeysRows                = strings.Join(tKeysFieldNames, ",")
	tKeysRowsExpectAutoSet   = strings.Join(stringx.Remove(tKeysFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tKeysRowsWithPlaceHolder = strings.Join(stringx.Remove(tKeysFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TKeysModel interface {
		Insert(data TKeys) (sql.Result, error)
		FindOne(id int64) (*TKeys, error)
		FindMany(number string, deviceId int64) ([]TKeys, error)
		Update(data TKeys) error
		Delete(id int64) error
	}

	defaultTKeysModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TKeys struct {
		Id           int64     `db:"id"`
		Number       string    `db:"number"`
		KeyId        int64     `db:"key_id"`
		PublicKey    string    `db:"public_key"`
		LastResort   int64     `db:"last_resort"`
		DeviceId     int64     `db:"device_id"`
		CreateTime   time.Time `db:"create_time"`
		SignedPrekey string    `db:"signed_prekey"`
		IdentityKey  string    `db:"identity_key"`
	}
)

func NewTKeysModel(conn sqlx.SqlConn) TKeysModel {
	return &defaultTKeysModel{
		conn:  conn,
		table: "`t_keys`",
	}
}

func (m *defaultTKeysModel) Insert(data TKeys) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, tKeysRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Number, data.KeyId, data.PublicKey, data.LastResort, data.DeviceId, data.SignedPrekey, data.IdentityKey)
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

func (m * defaultTKeysModel) FindMany(number string, deviceId int64) ([]TKeys, error){
	query := fmt.Sprintf("select %s from %s where `number` = ? and `device_id`=? ", tKeysRows, m.table)
	var resp []TKeys
	err := m.conn.QueryRows(&resp, query, number,deviceId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTKeysModel) Update(data TKeys) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tKeysRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Number, data.KeyId, data.PublicKey, data.LastResort, data.DeviceId, data.SignedPrekey, data.IdentityKey, data.Id)
	return err
}

func (m *defaultTKeysModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
