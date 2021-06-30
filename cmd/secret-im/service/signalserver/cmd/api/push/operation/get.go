package operation

import (
	"github.com/go-redis/redis"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/goredis"
	"strconv"
)

// 查询操作
type GetOperation struct {
	luaScript *goredis.LuaScript
}

// 创建实例
func NewGetOperation(client *redis.Client) (*GetOperation, error) {
	luaScript, err := goredis.NewLuaScript(client, "lua/apn/get.lua")
	if err != nil {
		return nil, err
	}
	op := GetOperation{luaScript: luaScript}
	return &op, nil
}

// 获取挂起设备
func (p GetOperation) GetPending(limit int) ([]string, error) {
	keys := []string{PendingNotificationsKey}
	result, err := p.luaScript.Exec(keys, utils.CurrentTimeMillis(), strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}

	pending := make([]string, 0)
	for _, val := range result.([]interface{}) {
		pending = append(pending, val.(string))
	}
	return pending, nil
}





