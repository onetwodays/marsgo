package storage

import (
	"fmt"
	"github.com/go-redis/redis"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"strconv"
	"strings"
)

// 设备管理器
type DevicesManager struct {
}

func (DevicesManager) deviceKey(uuid string) string {
	return fmt.Sprintf("online::{%s}", uuid)
}

func (DevicesManager) partitionKey(partition int) string {
	return fmt.Sprintf("partition::{%d}", partition)
}

// 设置离线
func (m DevicesManager) SetOffline(uuid string, deviceId int64, partition int) error {
	pipeline := internal.client.Pipeline()
	pipeline.Del(m.deviceKey(uuid), strconv.FormatInt(deviceId, 10))
	pipeline.SRem(m.partitionKey(partition), fmt.Sprintf("%s:%d", uuid, deviceId))
	_, err := pipeline.Exec()
	return err
}

// 设置在线
func (m DevicesManager) SetOnline(uuid string, deviceId int64, partition int) error {
	pipeline := internal.client.Pipeline()
	pipeline.HSet(m.deviceKey(uuid), strconv.FormatInt(deviceId, 10), partition)
	pipeline.SAdd(m.partitionKey(partition), fmt.Sprintf("%s:%d", uuid, deviceId))
	_, err := pipeline.Exec()
	return err
}

// 获取设备区域
func (m DevicesManager) GetPartitions(uuid string) ([]entities.DevicePartition, error) {
	cmd := internal.client.HGetAll(m.deviceKey(uuid))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	result := make([]entities.DevicePartition, 0, len(cmd.Val()))
	for key, value := range cmd.Val() {
		partition, err := strconv.Atoi(value)
		if err != nil {
			continue
		}

		deviceID, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			continue
		}

		result = append(result, entities.DevicePartition{
			UUID:      uuid,
			DeviceID:  deviceID,
			Partition: partition,
		})
	}
	return result, nil
}

// 获取在线设备
func (m DevicesManager) GetOnlineDevices(users []string) ([]entities.DevicePartition, error) {
	pipeline := internal.client.Pipeline()
	cmds := make([]*redis.StringStringMapCmd, 0, len(users))
	for _, uuid := range users {
		cmds = append(cmds, pipeline.HGetAll(m.deviceKey(uuid)))
	}
	_, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	result := make([]entities.DevicePartition, 0, len(cmds))
	for idx, cmd := range cmds {
		for key, value := range cmd.Val() {
			partition, err := strconv.Atoi(value)
			if err != nil {
				continue
			}

			deviceID, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				continue
			}

			result = append(result, entities.DevicePartition{
				UUID:      users[idx],
				DeviceID:  deviceID,
				Partition: partition,
			})
		}
	}
	return result, nil
}

// 获取设备列表
func (m DevicesManager) GetDevicesByPartition(partition int) ([]entities.DevicePartition, error) {
	cmd := internal.client.SMembers(m.partitionKey(partition))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	result := make([]entities.DevicePartition, 0, len(cmd.Val()))
	for _, s := range cmd.Val() {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			continue
		}

		uuid := parts[0]
		deviceID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}
		result = append(result, entities.DevicePartition{
			UUID:      uuid,
			DeviceID:  deviceID,
			Partition: partition,
		})
	}
	return result, nil
}

// 释放区域缓存
func (m DevicesManager) ReleasePartition(partition int) error {
	return internal.client.Del(m.partitionKey(partition)).Err()
}
