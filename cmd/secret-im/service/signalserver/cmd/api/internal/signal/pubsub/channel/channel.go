package channel

import "secret-im/service/signalserver/cmd/api/textsecure"

// 分发通道
type DispatchChannel interface {
	// 接收消息
	OnDispatchMessage(channel string, message textsecure.PubSubMessage)
	// 订阅成功
	OnDispatchSubscribed(channel string)
	// 取消订阅
	OnDispatchUnsubscribed(channel string)
}
