package dao

import (
	"context"

	"marsgo/pkg/cache/redis"
	"marsgo/pkg/conf/paladin"
	"marsgo/pkg/log"
)

func NewRedis() (r *redis.Redis, err error) {
	var cfg struct {
		Client *redis.Config
	}
	if err = paladin.Get("redis.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	r = redis.NewRedis(cfg.Client)
	return
}

func (d *dao) PingRedis(ctx context.Context) (err error) {
	if _, err = d.redis.Do(ctx, "SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}