package pubsub

import (
	"secret-im/service/signalserver/cmd/api/internal/signal/pubsub/channel"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

// 调度管理器
type DispatchManager interface {
	Shutdown()
	HasSubscription(name string) bool
	Publish(name string, message *textsecure.PubSubMessage) (int64, error)
	Subscribe(name string, dispatchChannel channel.DispatchChannel) error
	Unsubscribe(name string, dispatchChannel channel.DispatchChannel) error
}

// 发布订阅管理器
type Manager struct {
	dispatchManager DispatchManager
}

// 创建发布订阅管理器
func NewManager(dispatchManager DispatchManager) *Manager {
	return &Manager{dispatchManager: dispatchManager}
}

// 是否订阅
func (manager *Manager) HasLocalSubscription(address string) bool {
	return manager.dispatchManager.HasSubscription(address)
}

// 发布消息
func (manager *Manager) Publish(address string, message *textsecure.PubSubMessage) (int64, error) {
	return manager.dispatchManager.Publish(address, message)
}

// 订阅消息
func (manager *Manager) Subscribe(address string, channel channel.DispatchChannel) error {
	return manager.dispatchManager.Subscribe(address, channel)
}

// 取消订阅
func (manager *Manager) Unsubscribe(address string, channel channel.DispatchChannel) error {
	return manager.dispatchManager.Unsubscribe(address, channel)
}
