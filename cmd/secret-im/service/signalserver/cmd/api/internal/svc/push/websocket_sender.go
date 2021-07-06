package push

import (
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

// websocket_sender 就是简单的把envolpe推给redis，由redis的一个携程把数据拉下来，根据地址推送出去
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
	message *textsecure.Envelope, channel ChannelType, online bool) (bool, error) {
	address := Address{
		Number:   number,
		DeviceID: deviceID,
	}


	content, err := proto.Marshal(message)
	if err != nil {
		return false, err
	}
	pubSubMessage := &textsecure.PubSubMessage{
		Type:    textsecure.PubSubMessage_DELIVER,
		Content: content,
	}
	// 只是写到redis里面.
	n, err := sender.pubSubManager.Publish(address.Serialize(), pubSubMessage)
	// 推redis成功
	if err == nil && n > 0 {
		logx.Info("[ws_sender]消息已经推送至redis中，key is ",address.Serialize())
		return true, nil
	}

	if err!=nil{
		logx.Errorf("[ws_sender] push key(%s) fail:%s",address.Serialize(),err.Error())
	}


	// 如果不在线的话
	if !online {
		logx.Errorf("[Websocket_sender]redis推送不成功，n=%d,online=%v,sender.QueueMessage",n,online)
		err = sender.QueueMessage(number, deviceID, message)
	}
	return false, err
}

// 发送调配消息
func (sender *WebsocketSender) SendProvisioningMessage(address ProvisioningAddress, body []byte) (bool, error) {

	pubSubMessage := &textsecure.PubSubMessage{
		Type:    textsecure.PubSubMessage_DELIVER,
		Content: body,
	}
	n, err := sender.pubSubManager.Publish(address.Serialize(), pubSubMessage)
	return n > 0, err
}

// 保存消息到redis
func (sender *WebsocketSender) QueueMessage(number string, deviceID int64, message *textsecure.Envelope) error {
	address := Address{
		Number:   number,
		DeviceID: deviceID,
	}


	err := storage.MessagesManager{}.Insert(number, deviceID, message)
	if err != nil {
		return err
	}



	pubSubMessage := &textsecure.PubSubMessage{Type: textsecure.PubSubMessage_QUERY_DB}
	_, err = sender.pubSubManager.Publish(address.Serialize(), pubSubMessage)
	return err
}

