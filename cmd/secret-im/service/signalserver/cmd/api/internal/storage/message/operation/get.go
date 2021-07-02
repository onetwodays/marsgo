package operation

import (
	"github.com/go-redis/redis"
	"secret-im/service/signalserver/cmd/api/goredis"
	"strconv"
)

// 对组
type Pair struct {
	ID   int64
	Data []byte
}

// 查询操作
type GetOperation struct {
	getItems  *goredis.LuaScript
	getQueues *goredis.LuaScript
}

// 创建实例
func NewGetOperation(client *redis.Client) (*GetOperation, error) {
	getItems, err := goredis.NewLuaScript(client, "lua/get_items.lua")
	if err != nil {
		return nil, err
	}
	getQueues, err := goredis.NewLuaScript(client, "lua/get_queues_to_persist.lua")
	if err != nil {
		return nil, err
	}
	op := GetOperation{getItems: getItems, getQueues: getQueues}
	return &op, nil
}

// 获取消息
func (p GetOperation) GetItems(queue, lock string, limit int) ([]Pair, error) {
	keys := []string{queue, lock}
	result, err := p.getItems.Exec(keys, strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}

	var pair Pair
	var items []Pair
	for idx, val := range result.([]interface{}) {
		if idx%2 == 0 {
			pair.Data = []byte(val.(string))
		} else {
			pair.ID, _ = strconv.ParseInt(val.(string), 10, 64)
			items = append(items, pair)
		}
	}
	return items, nil
}

// 获取队列
func (p GetOperation) GetQueues(queue string, maxTimeMillis int64, limit int) ([]string, error) {
	keys := []string{queue}
	args := []interface{}{maxTimeMillis, limit}
	result, err := p.getQueues.Exec(keys, args...)
	if err != nil {
		return nil, err
	}

	l := result.([]interface{})
	queues := make([]string, 0, len(l))
	for _, key := range l {
		queues = append(queues, key.(string))
	}
	return queues, nil
}

