package operation

import (
	"fmt"
	"github.com/go-redis/redis"
	"secret-im/service/signalserver/cmd/api/goredis"
	"strconv"
)

// 删除操作
type RemoveOperation struct {
	removeByID     *goredis.LuaScript
	removeBySender *goredis.LuaScript
	removeByGUID   *goredis.LuaScript
	removeQueue    *goredis.LuaScript
}

// 创建实例
func NewRemoveOperation(client *redis.Client) (*RemoveOperation, error) {
	removeByID, err := goredis.NewLuaScript(client, "lua/remove_item_by_id.lua")
	if err != nil {
		return nil, err
	}
	removeBySender, err := goredis.NewLuaScript(client, "lua/remove_item_by_sender.lua")
	if err != nil {
		return nil, err
	}
	removeByGUID, err := goredis.NewLuaScript(client, "lua/remove_item_by_guid.lua")
	if err != nil {
		return nil, err
	}
	removeQueue, err := goredis.NewLuaScript(client, "lua/remove_queue.lua")
	if err != nil {
		return nil, err
	}

	op := RemoveOperation{
		removeByID:     removeByID,
		removeBySender: removeBySender,
		removeByGUID:   removeByGUID,
		removeQueue:    removeQueue,
	}
	return &op, nil
}

// 删除消息
func (p RemoveOperation) Remove(destination string, destinationDevice int64, id int64) error {
	key := NewKey(destination, destinationDevice)
	keys := []string{key.UserMessageQueue, key.UserMessageQueueMetadata, key.UserMessageQueueIndex()}
	_, err := p.removeByID.Exec(keys, strconv.FormatInt(id, 10))
	if err == redis.Nil {
		return nil
	}
	return err
}

// 根据发送者删除消息
func (p RemoveOperation) RemoveBySender(destination string, destinationDevice int64,
	sender string, timestamp int64) ([]byte, error) {

	key := NewKey(destination, destinationDevice)
	senderKey := fmt.Sprintf("{%s}::%d", sender, timestamp)

	keys := []string{key.UserMessageQueue, key.UserMessageQueueMetadata, key.UserMessageQueueIndex()}
	result, err := p.removeBySender.Exec(keys, senderKey)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return []byte(result.(string)), nil
}

// 根据GUID删除消息
func (p RemoveOperation) RemoveByGUID(destination string, destinationDevice int64, guid string) ([]byte, error) {
	key := NewKey(destination, destinationDevice)
	keys := []string{key.UserMessageQueue, key.UserMessageQueueMetadata, key.UserMessageQueueIndex()}
	result, err := p.removeByGUID.Exec(keys, guid)
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return []byte(result.(string)), nil
}

// 清理消息
func (p RemoveOperation) Clear(destination string, deviceID int64) error {
	key := NewKey(destination, deviceID)
	keys := []string{key.UserMessageQueue, key.UserMessageQueueMetadata, key.UserMessageQueueIndex()}
	_, err := p.removeQueue.Exec(keys)
	if err == redis.Nil {
		return redis.Nil
	}
	return err
}

