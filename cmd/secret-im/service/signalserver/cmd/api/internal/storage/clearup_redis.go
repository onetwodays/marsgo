package storage

import "secret-im/service/signalserver/cmd/api/config"

// 清理缓存
func ClearUpRedis() error {
	partition := config.AppConfig.PartitionID
	devices, err := DevicesManager{}.GetDevicesByPartition(partition)
	if err != nil {
		return err
	}

	for _, device := range devices {
		err = DevicesManager{}.SetOffline(device.UUID, device.DeviceID, partition)
		if err != nil {
			return err
		}
		err = ChannelParticipantsManager{}.SetOffline(device.UUID, device.DeviceID)
		if err != nil {
			return err
		}
	}
	return DevicesManager{}.ReleasePartition(partition)
}

