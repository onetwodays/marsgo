package push

import (
	"errors"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub"

	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

// 消息发送器
type Sender struct {
	redisOperation  *RedisOperation
	websocketSender *WebsocketSender
}

// 创建消息发送器
func NewPushSender(pubSubManager *pubsub.Manager, redisOperation *RedisOperation) *Sender {
	websocketSender := WebsocketSender{
		pubSubManager: pubSubManager,
	}
	return &Sender{
		redisOperation:  redisOperation,
		websocketSender: &websocketSender,
	}
}

// 获取ws发送器
func (sender *Sender) GetWebsocketSender() *WebsocketSender {
	return sender.websocketSender
}

// 发送消息
func (sender *Sender) SendMessage(number string, device *entities.Device,
	message *textsecure.Envelope, online bool) (bool, error) {

	if  len(device.GcmID)!=0 {
		return sender.sendGcmMessage(number, device, message, online)
	} else if len(device.ApnID)!=0 {
		return sender.sendApnMessage(number, device, message, online)
	} else if device.FetchesMessages {
		return sender.websocketSender.SendMessage(number, device.ID, message, ChannelTypeWEB, online)
	} else {
		return false, errors.New("not implemented")
	}
}

// 发送通知
func (sender *Sender) SendQueuedNotification(number string, device *entities.Device) (err error) {
	if len(device.GcmID)!=0 {
		err = sender.sendGcmNotification(number, device)
	} else if len(device.ApnID)!=0 {
		err = sender.sendApnNotification(number, device, true)
	} else if !device.FetchesMessages {
		return errors.New("no notification possible")
	}
	return err
}

// 发送GCM消息
func (sender *Sender) sendGcmMessage(number string, device *entities.Device,
	message *textsecure.Envelope, online bool) (bool, error) {

	delivered, err := sender.websocketSender.SendMessage(number, device.ID, message, ChannelTypeAPN, online)
	if err != nil {
		return false, err
	}

	if !delivered && message.GetType() != textsecure.Envelope_RECEIPT && !online {
		sender.sendGcmNotification(number, device)
		return false, nil
	}
	return delivered, nil
}

// 发送GCM通知
func (sender *Sender) sendGcmNotification(number string, device *entities.Device) error {
	return nil
}

// 发送APN消息
func (sender *Sender) sendApnMessage(number string, device *entities.Device,
	message *textsecure.Envelope, online bool) (bool, error) {

	delivered, err := sender.websocketSender.SendMessage(number, device.ID, message, ChannelTypeAPN, online)
	if err != nil {
		return false, err
	}

	if !delivered && message.GetType() != textsecure.Envelope_RECEIPT && !online {
		sender.sendApnNotification(number, device, false)
		return false, nil
	}
	return delivered, nil
}

// 发送APN通知
func (sender *Sender) sendApnNotification(number string, device *entities.Device, newOnly bool) error {
	if newOnly {
		ok, _ := sender.redisOperation.IsScheduled(number, device.ID)
		if ok {
			return nil
		}
	}

	apnMessage := ApnMessage{
		Number:   number,
		DeviceID: device.ID,
	}
	if len(device.VoipApnID)==0 {
		apnMessage.ApnID = device.ApnID
	} else {
		apnMessage.IsVoip = true
		apnMessage.ApnID = device.VoipApnID
		sender.redisOperation.Schedule(number, device.ID)
	}

	return AddToApnMessageQueue(apnMessage)
}

