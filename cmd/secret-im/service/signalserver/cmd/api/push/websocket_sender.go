package push

import (
	"github.com/golang/protobuf/proto"
	"secret-im/service/signalserver/cmd/api/pubsub"
	"secret-im/service/signalserver/cmd/api/textsecure"

)

// 通道类型
type ChannelType string

var (
	ChannelTypeAPN ChannelType = "APN"
	ChannelTypeGCM ChannelType = "GCM"
	ChannelTypeWEB ChannelType = "WEB"
)

// Websocket发送器
type WebsocketSender struct {
	pubSubManager *pubsub.Manager
}

// 发送消息
func (sender *WebsocketSender) SendMessage(number string, deviceID int64,
	message textsecure.Envelope, channel ChannelType, online bool) (bool, error) {
	address := Address{
		Number:   number,
		DeviceID: deviceID,
	}

	typ := textsecure.PubSubMessage_DELIVER
	content, err := proto.Marshal(&message)
	if err != nil {
		return false, err
	}
	pubSubMessage := textsecure.PubSubMessage{
		Type:    typ,
		Content: content,
	}

	n, err := sender.pubSubManager.Publish(address.Serialize(), pubSubMessage)
	if err == nil && n > 0 {
		return true, nil
	}

	if !online {
		err = sender.QueueMessage(number, deviceID, message)
	}
	return false, err
}

// 发送调配消息
func (sender *WebsocketSender) SendProvisioningMessage(address ProvisioningAddress, body []byte) (bool, error) {
	typ := textsecure.PubSubMessage_DELIVER
	pubSubMessage := textsecure.PubSubMessage{
		Type:    typ,
		Content: body,
	}
	n, err := sender.pubSubManager.Publish(address.Serialize(), pubSubMessage)
	return n > 0, err
}

// 保存到消息队列
func (sender *WebsocketSender) QueueMessage(number string, deviceID int64, message textsecure.Envelope) error {
	address := Address{
		Number:   number,
		DeviceID: deviceID,
	}

    /*
	err := storage.MessagesManager{}.Insert(number, deviceID, message)
	if err != nil {
		return err
	}
     */

	typ := textsecure.PubSubMessage_QUERY_DB
	pubSubMessage := textsecure.PubSubMessage{Type: typ}
	_, err := sender.pubSubManager.Publish(address.Serialize(), pubSubMessage)
	return err
}

