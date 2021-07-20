package storage

import (
	"encoding/json"
	"fmt"
	"secret-im/pkg/utils-tools"
)

// 频道消息管理器
type ChannelMessagesManager struct {
}

// 获取最新消息
func (m ChannelMessagesManager) GetLatestMessages(channelIDs []string) (map[string]ChannelMessage, error) {
	channelIDs = utils.StringSlice{}.Distinct(channelIDs)
	mapper := make(map[string]ChannelMessage)
	messages, err := m.redisMGet(channelIDs)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		if message == nil {
			continue
		}
		mapper[message.ChannelID] = *message
	}
	return mapper, nil
}

// 设置最新消息
func (m ChannelMessagesManager) SetLatestMessage(channelID string, message *ChannelMessage) error {
	return m.redisSet(channelID, message)
}

func (ChannelMessagesManager) redisKey(channelID string) string {
	return fmt.Sprintf("channel::{%s}::latest", channelID)
}

func (m ChannelMessagesManager) redisSet(channelID string, message *ChannelMessage) error {
	jsb, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return internal.client.Set(m.redisKey(channelID), jsb, 0).Err()
}

func (m ChannelMessagesManager) redisMGet(channelIDs []string) ([]*ChannelMessage, error) {
	keys := make([]string, 0, len(channelIDs))
	for _, channelID := range channelIDs {
		keys = append(keys, m.redisKey(channelID))
	}

	results, err := internal.client.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]*ChannelMessage, 0, len(results))
	for _, item := range results {
		if item == nil {
			messages = append(messages, nil)
			continue
		}

		var message ChannelMessage
		err = json.Unmarshal([]byte(item.(string)), &message)
		if err != nil {
			messages = append(messages, nil)
			continue
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

