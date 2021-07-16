package localcache

import (
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"sync"
	"sync/atomic"
)

// 全局对象
var OnlineParticipants *OnlineParticipantsManager

func init() {
	OnlineParticipants = NewOnlineParticipantsManager(60 * 10)
}

// 在线成员
type onlineParticipants struct {
	channelID    string
	version      int64
	participants atomic.Value
	mutex        sync.Mutex
}

// 获取缓存
func (cache *onlineParticipants) Get() (map[int][]entities.DevicePartition, error) {
	version, err := storage.ChannelParticipantsManager{}.GetVersion(cache.channelID)
	if err != nil {
		return nil, err
	}

	if version == atomic.LoadInt64(&cache.version) {
		value := cache.participants.Load()
		if value == nil {
			return nil, nil
		}
		return value.(map[int][]entities.DevicePartition), nil
	}

	participants, version, err := storage.ChannelParticipantsManager{}.GetOnlineParticipants(cache.channelID)
	if err != nil {
		return nil, err
	}

	cache.mutex.Lock()
	cache.version = version
	cache.participants.Store(participants)
	cache.mutex.Unlock()
	return participants, nil
}

// 在线成员管理
type OnlineParticipantsManager struct {
	mapper sync.Map
	timer  *CountdownTimer
}

// 创建在线成员管理器
func NewOnlineParticipantsManager(expiration int64) *OnlineParticipantsManager {
	var manager OnlineParticipantsManager
	onExpired := func(channelID string) {
		manager.mapper.Delete(channelID)
	}
	manager.timer = NewCountdownTimer(expiration, onExpired)
	return &manager
}

// 获取缓存
func (manager *OnlineParticipantsManager) Get(channelID string) (map[int][]entities.DevicePartition, error) {
	defer manager.timer.Update(channelID)

	value, ok := manager.mapper.Load(channelID)
	if ok {
		participants := value.(*onlineParticipants)
		return participants.Get()
	}

	participants := &onlineParticipants{
		channelID: channelID,
	}
	manager.mapper.Store(channelID, participants)
	return participants.Get()
}

// 设置过期时间
func (manager *OnlineParticipantsManager) SetExpiration(seconds int64) {
	manager.timer.SetExpiration(seconds)
}

