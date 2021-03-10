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
	tMessagesFieldNames          = builderx.RawFieldNames(&TMessages{})
	tMessagesRows                = strings.Join(tMessagesFieldNames, ",")
	tMessagesRowsExpectAutoSet   = strings.Join(stringx.Remove(tMessagesFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	tMessagesRowsWithPlaceHolder = strings.Join(stringx.Remove(tMessagesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	TMessagesModel interface {
		Insert(data TMessages) (sql.Result, error)
		FindOne(id int64) (*TMessages, error)
		FindMany(destination string,deviceId int64,pageSize,pageIndex int) ([]TMessages,error)
		Update(data TMessages) error
		Delete(id int64) error
	}

	defaultTMessagesModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TMessages struct {
		Relay             string    `db:"relay"`   // 类似cdn
		Source            string    `db:"source"`  // 源手机号
		Message           string    `db:"message"` // 消息
		Content           string    `db:"content"`
		Destination       string    `db:"destination"` // 目的手机号
		CreateTime        time.Time `db:"create_time"`
		Id                int64     `db:"id"`                 // pk
		Type              int64     `db:"type"`               // 消息类型
		Tm                int64     `db:"tm"`                 // unix  时间戳
		SourceDevice      int64     `db:"source_device"`      // 源手机号绑定的设备id
		DestinationDevice int64     `db:"destination_device"` // 源手机号绑定的设备id
	}
)

func NewTMessagesModel(conn sqlx.SqlConn) TMessagesModel {
	return &defaultTMessagesModel{
		conn:  conn,
		table: "`t_messages`",
	}
}

func (m *defaultTMessagesModel) Insert(data TMessages) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, tMessagesRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Relay, data.Source, data.Message, data.Content, data.Destination, data.Type, data.Tm, data.SourceDevice, data.DestinationDevice)
	return ret, err
}

func (m *defaultTMessagesModel) FindOne(id int64) (*TMessages, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", tMessagesRows, m.table)
	var resp TMessages
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

func (m *defaultTMessagesModel) FindMany(destination string,deviceId int64,pageSize,pageIndex int) ([]TMessages,error){
	query := fmt.Sprintf("select %s from %s where destination=? and destination_device=?  ", tMessagesRows, m.table)

	query+=" order by create_time ASC limit ? offset ?"


	var resp []TMessages
	err := m.conn.QueryRows(&resp, query, destination, deviceId,pageSize,pageIndex*pageSize)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTMessagesModel) Update(data TMessages) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, tMessagesRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Relay, data.Source, data.Message, data.Content, data.Destination, data.Type, data.Tm, data.SourceDevice, data.DestinationDevice, data.Id)
	return err
}

func (m *defaultTMessagesModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
