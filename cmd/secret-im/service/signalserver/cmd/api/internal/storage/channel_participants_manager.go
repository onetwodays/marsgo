package storage


import (
	"fmt"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

// 频道在线成员管理器
type ChannelParticipantsManager struct {
}

func (ChannelParticipantsManager) channelKey(channelID string) string {
	return fmt.Sprintf("channel::{%s}", channelID)
}

// 加入频道
func (m ChannelParticipantsManager) join(channelID string, devices []entities.DevicePartition) error {
	pipeline := internal.client.Pipeline()
	key := m.channelKey(channelID)
	for _, device := range devices {
		field := fmt.Sprintf("%s:%d", device.UUID, device.DeviceID)
		pipeline.HSet(key, field, device.Partition)
	}
	pipeline.HIncrBy(key, "version", 1)

	_, err := pipeline.Exec()
	return err
}

// 离开频道
func (m ChannelParticipantsManager) leave(channelID, userID string) error {
	var cursor uint64
	fields := make([]string, 0)
	key := m.channelKey(channelID)
	for {
		var err error
		var keys []string
		keys, cursor, err = internal.client.HScan(key, cursor, userID, 100).Result()
		if err != nil {
			return err
		}
		fields = append(fields, keys...)
		if cursor == 0 {
			break
		}
	}
	if len(fields) == 0 {
		return nil
	}
	return internal.client.HDel(key, fields...).Err()
}

// 停用频道
func (m ChannelParticipantsManager) deactivate(channelID string) error {
	return internal.client.Del(m.channelKey(channelID)).Err()
}

// 设置用户离线
func (m ChannelParticipantsManager) SetOffline(userID string, deviceID int64) error {
	channels, err := ChannelJoinedDao{}.GetChannels(userID)
	if err != nil {
		return err
	}

	pipeline := internal.client.Pipeline()
	for _, channel := range channels {
		key := m.channelKey(channel)
		field := fmt.Sprintf("%s:%d", userID, deviceID)
		pipeline.HDel(key, field)
		pipeline.HIncrBy(key, "version", 1)
	}

	_, err = pipeline.Exec()
	return err
}

// 设置用户上线
func (m ChannelParticipantsManager) SetOnline(uuid string, deviceID int64, partition int) error {
	channels, err := ChannelJoinedDao{}.GetChannels(uuid)
	if err != nil {
		return err
	}

	pipeline := internal.client.Pipeline()
	for _, channel := range channels {
		key := m.channelKey(channel)
		field := fmt.Sprintf("%s:%d", uuid, deviceID)
		pipeline.HSet(key, field, partition)
		pipeline.HIncrBy(key, "version", 1)
	}

	_, err = pipeline.Exec()
	return err
}

// 获取缓存版本
func (m ChannelParticipantsManager) GetVersion(channelID string) (int64, error) {
	cmd := internal.client.HGet(m.channelKey(channelID), "version")
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return 0, nil
		}
		return 0, cmd.Err()
	}
	return cmd.Int64()
}

// 获取频道在线成员
func (m ChannelParticipantsManager) GetOnlineParticipants(channelID string) (map[int][]entities.DevicePartition, int64, error) {
	cmd := internal.client.HGetAll(m.channelKey(channelID))
	if cmd.Err() != nil {
		return nil, 0, cmd.Err()
	}

	var version int64
	mapper := make(map[int][]entities.DevicePartition)
	for key, value := range cmd.Val() {
		if key == "version" {
			version, _ = strconv.ParseInt(value, 10, 64)
			continue
		}

		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			continue
		}

		userID := parts[0]
		deviceID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}

		partition, err := strconv.Atoi(value)
		if err != nil {
			continue
		}

		device := entities.DevicePartition{
			UUID:      userID,
			DeviceID:  deviceID,
			Partition: partition,
		}
		users, ok := mapper[partition]
		if ok {
			users = append(users, device)
		} else {
			users = []entities.DevicePartition{device}
		}
		mapper[partition] = users
	}
	return mapper, version, nil
}
