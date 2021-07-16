package logic

import (
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/storage/localcache"
	"secret-im/service/signalserver/cmd/api/queue"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"time"
)

// 发送操作消息
func SendActionMessage(channelID string, action textsecure.MessageAction) error {
	data, err := proto.Marshal(&action)
	if err != nil {
		return err
	}

	message := storage.ChannelMessage{
		ChannelID:       channelID,
		Type:            model.ChannelMessageTypeService,
		Action:          data,
		Timestamp:       time.Now().Unix(),
		ServerTimestamp: time.Now().Unix(),
	}
	err = storage.ChannelMessages{}.Insert(&message)
	if err != nil {
		return err
	}
	return sendChannelMessageToDevices(&message)
}

// 发送消息到设备
func sendChannelMessageToDevices(message *storage.ChannelMessage) error {
	// 获取在线成员
	t := time.Now()
	mapper, err := localcache.OnlineParticipants.Get(message.ChannelID)
	if err != nil {
		return err
	}

	logx.Error("[Channel::sendMessage] get channel online participants"," consume:", time.Since(t))


	// 构造消息内容
	timestamp := uint64(message.Timestamp)
	serverTimestamp := uint64(message.ServerTimestamp)
	typ := textsecure.ChannelEnvelope_MESSAGE
	if message.Type == model.ChannelMessageTypeService {
		typ = textsecure.ChannelEnvelope_MESSAGE_SERVICE
	}

	envelope := textsecure.ChannelEnvelope{
		Type:            typ,
		Id:              message.MessageID,
		ChannelId:       message.ChannelID,
		SourceUuid:      *message.Source,
		Relay:           *message.Relay,
		Deleted:         message.Deleted,
		Timestamp:       timestamp,
		ServerTimestamp: serverTimestamp,
	}
	if typ == textsecure.ChannelEnvelope_MESSAGE {
		envelope.Content = message.Content
	} else if typ == textsecure.ChannelEnvelope_MESSAGE_SERVICE {
		var action textsecure.MessageAction
		if err = proto.Unmarshal(message.Action, &action); err == nil {
			envelope.Action = &action
		}
	}

	if message.SourceDevice != nil {
		sourceDevice := uint32(*message.SourceDevice)
		envelope.SourceDevice = sourceDevice
	}

	content, err := proto.Marshal(&envelope)
	if err != nil {
		return err
	}

	// 路由消息到区域服务器
	channelMessage := queue.SendMessageToDevice_CHANNEL_MESSAGE
	for participant, devices := range mapper {
		request := queue.SendMessageToDevice{
			Type:    channelMessage,
			Content: content,
		}

		for _, device := range devices {
			request.Devices = append(request.Devices, &queue.Device{
				Id:   device.DeviceID,
				Uuid: device.UUID,
			})
		}
		queue.Publish(queue.SendToDeviceTopic(participant), &request)
	}
	return nil
}

