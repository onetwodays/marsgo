package dao

import (
	"context"

	"marsgo/cmd/mars_server/internal/model"
	"marsgo/pkg/conf/paladin"
	"marsgo/pkg/database/sql"
)

func NewDB() (db *sql.DB, err error) {
	var cfg struct {
		Client *sql.Config //指针分配空间怎么办？
	}
	if err = paladin.Get("db.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = sql.NewMySQL(cfg.Client)
	return
}

func (d *dao) RawArticle(ctx context.Context, id int64) (art *model.Article, err error) {
	// get data from db
	return
}
