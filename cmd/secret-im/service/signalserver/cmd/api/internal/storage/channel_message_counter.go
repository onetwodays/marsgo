package storage

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

// 频道消息计数管理器
type ChannelMessageCounter struct {
}

func (ChannelMessageCounter) countKey(channelID string) string {
	return fmt.Sprintf("channel_count::{%s}", channelID)
}

// 获取计数
func (m ChannelMessageCounter) Get(channelID string) (int64, error) {
	cmd := internal.client.Get(m.countKey(channelID))
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return 0, nil
		}
		return 0, cmd.Err()
	}
	return cmd.Int64()
}

func (m ChannelMessageCounter) GetChannels(channelIDs []string) (map[string]int64, error) {
	keys := make([]string, 0, len(channelIDs))
	for _, channelID := range channelIDs {
		keys = append(keys, m.countKey(channelID))
	}

	cmd := internal.client.MGet(keys...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	mapper := make(map[string]int64)
	for idx, item := range cmd.Val() {
		if item == nil {
			continue
		}
		id, err := strconv.ParseInt(item.(string), 10, 64)
		if err != nil {
			continue
		}
		mapper[channelIDs[idx]] = id
	}
	return mapper, nil
}

// 自增计数
func (m ChannelMessageCounter) Incr(channelID string) (int64, error) {
	cmd := internal.client.Incr(m.countKey(channelID))
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

func (m ChannelMessageCounter) IncrBy(channelID string, value int64) (int64, error) {
	cmd := internal.client.IncrBy(m.countKey(channelID), value)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Val(), nil
}

