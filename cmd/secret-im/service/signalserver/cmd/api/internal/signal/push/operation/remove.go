package operation

import (
	"github.com/go-redis/redis"
	"secret-im/service/signalserver/cmd/api/goredis"
)

// 删除操作
type RemoveOperation struct {
	luaScript *goredis.LuaScript
}

// 创建实例
func NewRemoveOperation(client *redis.Client) (*RemoveOperation, error) {
	luaScript, err := goredis.NewLuaScript(client, "lua/apn/remove.lua")
	if err != nil {
		return nil, err
	}
	op := RemoveOperation{luaScript: luaScript}
	return &op, nil
}

// 根据端点移除设备
func (p RemoveOperation) RemoveByEndpoint(endpoint string) (bool, error) {
	if endpoint == PendingNotificationsKey {
		return false, nil
	}

	keys := []string{PendingNotificationsKey, endpoint}
	val, err := p.luaScript.Exec(keys)
	if err != nil {
		return false, err
	}
	return val.(int64) > 0, nil
}

// 移除设备
func (p RemoveOperation) Remove(number string, deviceID int64) (bool, error) {
	endpoint := GetEndpointKey(number, deviceID)
	return p.RemoveByEndpoint(endpoint)
}

