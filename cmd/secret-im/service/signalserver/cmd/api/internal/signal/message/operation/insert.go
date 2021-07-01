package operation

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"secret-im/service/signalserver/cmd/api/goredis"
	"secret-im/service/signalserver/cmd/api/textsecure"

)

// 插入操作
type InsertOperation struct {
	insert *goredis.LuaScript
}

// 创建实例
func NewInsertOperation(client *redis.Client) (*InsertOperation, error) {
	insert, err := goredis.NewLuaScript(client, "lua/insert_item.lua")
	if err != nil {
		return nil, err
	}
	return &InsertOperation{insert: insert}, nil
}

// 插入消息
func (p InsertOperation) Insert(guid, destination string, destinationDevice,
	timestamp int64, message *textsecure.Envelope) error {

	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	sender := "nil"
	if len(message.Source) != 0 {
		sender = fmt.Sprintf("{%s}::%d", message.GetSource(), message.GetTimestamp())
	}

	key := NewKey(destination, destinationDevice)
	keys := []string{key.UserMessageQueue, key.UserMessageQueueMetadata, key.UserMessageQueueIndex()}
	args := []interface{}{data, timestamp, sender, guid}
	_, err = p.insert.Exec(keys, args...)
	return err
}

