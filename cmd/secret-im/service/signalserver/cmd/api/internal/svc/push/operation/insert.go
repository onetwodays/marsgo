package operation

import (
	"github.com/go-redis/redis"
	"secret-im/service/signalserver/cmd/api/goredis"
)

// 插入操作
type InsertOperation struct {
	luaScript *goredis.LuaScript
}

// 创建实例
func NewInsertOperation(client *redis.Client) (*InsertOperation, error) {
	luaScript, err := goredis.NewLuaScript(client, "lua/apn/insert.lua")
	if err != nil {
		return nil, err
	}
	return &InsertOperation{luaScript: luaScript}, nil
}

// 插入设备
func (p InsertOperation) Insert(number string, deviceID int64, timestamp, interval int64) error {
	endpoint := GetEndpointKey(number, deviceID)

	keys := []string{PendingNotificationsKey, endpoint}
	_, err := p.luaScript.Exec(keys, timestamp, interval, number, deviceID)
	if err == redis.Nil {
		return nil
	}
	return err
}

