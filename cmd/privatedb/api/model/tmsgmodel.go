package model

import (
    "database/sql"
    "fmt"
    "github.com/tal-tech/go-zero/core/logx"
    "privatedb/api/internal/types"
    "strings"
    "time"

    "github.com/tal-tech/go-zero/core/stores/sqlc"
    "github.com/tal-tech/go-zero/core/stores/sqlx"
    "github.com/tal-tech/go-zero/core/stringx"
    "github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
    tMsgFieldNames          = builderx.FieldNames(&TMsg{})
    tMsgRows                = strings.Join(tMsgFieldNames, ",")
    tMsgRowsExpectAutoSet   = strings.Join(stringx.Remove(tMsgFieldNames, "id", "create_time", "update_time"), ",")
    tMsgRowsWithPlaceHolder = strings.Join(stringx.Remove(tMsgFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
)

type (
    TMsgModel interface {
        Insert(data TMsg) (sql.Result, error)
        FindOne(id int64) (*TMsg, error)
        Update(data TMsg) error
        Delete(id int64) error
        FindMany(req types.GetmsgsReq) ([]TMsg, error)
    }

    defaultTMsgModel struct {
        conn  sqlx.SqlConn
        table string
    }

    TMsg struct {
        Id         int64     `db:"id" json:"id"`
        Pair       string    `db:"pair" json:"pair"`
        DealId     int64     `db:"deal_id" json:"deal_id"`
        Sender     string    `db:"sender" json:"sender"`
        Receiver   string    `db:"receiver" json:"receiver"`
        Context    string    `db:"context" json:"context"`
        CreateTime time.Time `db:"create_time" json:"ctime"`
        UpdateTime time.Time `db:"update_time" json:"utime"`
        Status     int64     `db:"status" json:"status"` // 1.not send 2 sender

    }
)

func NewTMsgModel(conn sqlx.SqlConn) TMsgModel {
    return &defaultTMsgModel{
        conn:  conn,
        table: "t_msg",
    }
}

func (m *defaultTMsgModel) Insert(data TMsg) (sql.Result, error) {
    query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, tMsgRowsExpectAutoSet)
    logx.Info(query)
    logx.Info("111111")
    ret, err := m.conn.Exec(query, data.Pair, data.DealId, data.Sender, data.Receiver, data.Context, data.Status)
    return ret, err
}

func (m *defaultTMsgModel) FindOne(id int64) (*TMsg, error) {
    query := fmt.Sprintf("select %s from %s where id = ? limit 1", tMsgRows, m.table)
    var resp TMsg
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

// select * from t_msg where pair='?' and deal_id=? and sender='?' and receiver='?' and status=?;
func (m *defaultTMsgModel) FindMany(req types.GetmsgsReq) ([]TMsg, error) {

    query := fmt.Sprintf("select %s from %s where pair=? and deal_id=?  ", tMsgRows, m.table)
    if req.Status!=0{
        query = fmt.Sprintf("%s and status=%d  ", query,req.Status)
    }
    query+=" order by id desc limit ? offset ?"


    var resp []TMsg
    err := m.conn.QueryRows(&resp, query, req.Pair, req.Dealid,req.PageSize,req.PageIndex*req.PageSize)
    switch err {
    case nil:
        return resp, nil
    case sqlc.ErrNotFound:
        return nil, ErrNotFound
    default:
        return nil, err
    }
}

func (m *defaultTMsgModel) Update(data TMsg) error {
    query := fmt.Sprintf("update %s set %s where id = ?", m.table, tMsgRowsWithPlaceHolder)
    _, err := m.conn.Exec(query, data.Pair, data.DealId, data.Sender, data.Receiver, data.Context, data.Status, data.Id)
    return err
}

func (m *defaultTMsgModel) Delete(id int64) error {
    query := fmt.Sprintf("delete from %s where id = ?", m.table)
    _, err := m.conn.Exec(query, id)
    return err
}
