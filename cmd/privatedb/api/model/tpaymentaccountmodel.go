package model

import (
	"database/sql"
	"fmt"
	"privatedb/api/internal/types"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	tPaymentAccountFieldNames          = builderx.FieldNames(&TPaymentAccount{})
	tPaymentAccountRows                = strings.Join(tPaymentAccountFieldNames, ",")
	tPaymentAccountRowsExpectAutoSet   = strings.Join(stringx.Remove(tPaymentAccountFieldNames, "id", "create_time", "update_time"), ",")
	tPaymentAccountRowsWithPlaceHolder = strings.Join(stringx.Remove(tPaymentAccountFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
)

type (
	TPaymentAccountModel interface {
		Insert(data TPaymentAccount) (sql.Result, error)
		FindOne(id int64) (*TPaymentAccount, error)
		Update(data TPaymentAccount) error
		Delete(id int64) error
		FindMany(req types.GetPayAccountReq) ([]TPaymentAccount, error)
	}

	defaultTPaymentAccountModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TPaymentAccount struct {
		Id              uint32     `db:"id" json:"pmt_id"` // primary key
		Username        string    `db:"username" json:"username"`
		PaymentTypeId   uint32     `db:"payment_type_id" json:"pmt_type" ` // pay type
		PaymentTypeName string    `db:"payment_type_name" json:"pmet_type_name"`
		IsActived       uint8     `db:"is_actived" json:"is_actived"` // 1:actived 0:not actived
		Detail          string    `db:"detail" json:"ciphertext"`
		CreateTime      time.Time `db:"create_time" json:"ctime"`
		UpdateTime      time.Time `db:"update_time" json:"utime"`

	}
)

func NewTPaymentAccountModel(conn sqlx.SqlConn) TPaymentAccountModel {
	return &defaultTPaymentAccountModel{
		conn:  conn,
		table: "t_payment_account",
	}
}

func (m *defaultTPaymentAccountModel) Insert(data TPaymentAccount) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, tPaymentAccountRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Username,data.PaymentTypeId, data.PaymentTypeName, data.IsActived, data.Detail)
	return ret, err
}

func (m *defaultTPaymentAccountModel) FindOne(id int64) (*TPaymentAccount, error) {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", tPaymentAccountRows, m.table)
	var resp TPaymentAccount
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


func (m *defaultTPaymentAccountModel) FindMany(req types.GetPayAccountReq) ([]TPaymentAccount, error){
	query := fmt.Sprintf("select %s from %s where 1=1  ", tPaymentAccountRows, m.table)

	if len(req.Username)>0{
		query = fmt.Sprintf("%s and username='%s' ",query,req.Username)
	}

	if req.PaymentType!=0{
		query = fmt.Sprintf("%s and payment_type_id=%d ",query,req.PaymentType)
	}

	if req.IsActived!=0{
		query = fmt.Sprintf("%s and is_actived=%d ",query,req.IsActived)
	}
	query+=" order by id desc "

	var resp []TPaymentAccount
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

func (m *defaultTPaymentAccountModel) Update(data TPaymentAccount) error {
	query := fmt.Sprintf("update %s set %s,update_time=now() where id = ?", m.table, tPaymentAccountRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Username,data.PaymentTypeId, data.PaymentTypeName, data.IsActived, data.Detail,data.Id)
	return err
}

func (m *defaultTPaymentAccountModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
