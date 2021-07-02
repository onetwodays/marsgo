package svc

import (
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

// 死信处理
type DeadLetterHandler struct{
	//ctx *svc.ServiceContext
}

// 接收消息
func (DeadLetterHandler) OnDispatchMessage(channel string, message *textsecure.PubSubMessage) {
	logx.Infof("[DeadLetterHandler] handling dead letter to: %s", channel)

	address, err := push.NewAddress(channel)

	if err != nil {
		logx.Error("[DeadLetterHandler] invalid websocket address")
		return
	}

	if message.GetType() != textsecure.PubSubMessage_DELIVER {
		return
	}

	var envelope textsecure.Envelope
	err = proto.Unmarshal(message.GetContent(), &envelope)
	if err != nil {
		logx.Info("[DeadLetterHandler] bad pubsub message")
		return
	}


	err = storage.MessagesManager{}.Insert(address.Number, address.DeviceID, &envelope)
	if err != nil {
		logx.Info("[DeadLetterHandler] failed to storage message"," channel:",channel)
	}


}

// 订阅成功
func (DeadLetterHandler) OnDispatchSubscribed(channel string) {
	logx.Infof("channel:%s [DeadLetterHandler] subscription notice!",channel)

}

// 取消订阅
func (DeadLetterHandler) OnDispatchUnsubscribed(channel string) {
	logx.Infof("channel:%s [DeadLetterHandler] unsubscribe notice!",channel)

}
