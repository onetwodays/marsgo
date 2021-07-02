package storage

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/pkg/utils-tools"
	operation "secret-im/service/signalserver/cmd/api/internal/storage/message/operation"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"time"
)

// 单次返回消息数量
const ResultSetChunkSize = 100
// 消息缓存
type MessagesCache struct {
}


// 序列化消息
type SerializedMessage struct {
	ID   int64
	Data []byte
}

// 构造消息实体
func constructEntityFromEnvelope(id int64, envelope *textsecure.Envelope) types.OutcomingMessagex {

	return types.OutcomingMessagex {
		Id:              id,
		Cached:          true,
		Guid:            envelope.ServerGuid,
		Type:            int(envelope.GetType()),
		Relay:           envelope.Relay,
		Timestamp:       int64(envelope.GetTimestamp()),
		Source:          envelope.Source,
		SourceUuid:      envelope.SourceUuid,
		SourceDevice:    int64(envelope.SourceDevice),
		Message:         base64.StdEncoding.EncodeToString(envelope.LegacyMessage),
		Content:         base64.StdEncoding.EncodeToString(envelope.Content),
		ServerTimestamp: int64(envelope.GetServerTimestamp()),
	}
}


// 插入消息
func (MessagesCache) Insert(guid string, destination string, destinationDevice int64, message *textsecure.Envelope) error {
	message.ServerGuid = guid
	return internal.insertOperation.Insert(guid, destination, destinationDevice, utils.CurrentTimeMillis(), message)
}

// 获取消息
func (MessagesCache) Get(destination string, destinationDevice int64,
	limit int) ([]types.OutcomingMessagex, error) {

	key := operation.NewKey(destination, destinationDevice)
	items, err := internal.getOperation.GetItems(key.UserMessageQueue, key.UserMessageQueuePersistInProgress, limit)
	if err != nil {
		return nil, err
	}

	results := make([]types.OutcomingMessagex, 0, len(items))
	for _, item := range items {
		var message textsecure.Envelope
		if err = proto.Unmarshal(item.Data, &message); err != nil {
			logx.Error("[MessageCache] failed to parse envelope"," reason:",err)
		}else{
			results = append(results, constructEntityFromEnvelope(item.ID, &message))
		}
	}
	return results, nil
}

// 获取持久化队列
func (MessagesCache) GetQueuesToPersist(delayTimeMillis int64) ([]string, error) {
	maxTime := utils.CurrentTimeMillis() - delayTimeMillis
	return internal.getOperation.GetQueues(operation.Key{}.UserMessageQueueIndex(), maxTime, 100)
}

// 从队列获取消息
func (MessagesCache) GetFromQueue(number string, deviceID int64) ([]SerializedMessage, error) {
	key := operation.NewKey(number, deviceID)
	cmd := internal.client.ZRangeWithScores(key.UserMessageQueue, 0, ResultSetChunkSize)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	messages := make([]SerializedMessage, 0, len(cmd.Val()))
	for _, item := range cmd.Val() {
		messages = append(messages, SerializedMessage{
			ID:   int64(item.Score),
			Data: []byte(item.Member.(string)),
		})
	}
	return messages, nil
}

// 队列锁定
func (MessagesCache) QueueLock(number string, deviceID int64) error {
	key := operation.NewKey(number, deviceID)
	return internal.client.Set(key.UserMessageQueuePersistInProgress, 1, time.Second*30).Err()
}

// 队列解锁
func (MessagesCache) QueueUnlock(number string, deviceID int64) error {
	key := operation.NewKey(number, deviceID)
	return internal.client.Del(key.UserMessageQueuePersistInProgress).Err()
}

// 删除消息
func (MessagesCache) Remove(destination string, destinationDevice, id int64) error {
	return internal.removeOperation.Remove(destination, destinationDevice, id)
}

// 根据GUID删除消息
func (MessagesCache) RemoveByGUID(destination string, destinationDevice int64,
	guid string) (*types.OutcomingMessagex, error) {

	serialized, err := internal.removeOperation.RemoveByGUID(destination, destinationDevice, guid)
	if err != nil {
		return nil, err
	}
	if serialized == nil {
		return nil, nil
	}

	var envelope textsecure.Envelope
	if err = proto.Unmarshal(serialized, &envelope); err != nil {
		return nil, err
	}
	entity := constructEntityFromEnvelope(0, &envelope)
	return &entity, nil
}

// 根据发送者删除消息
func (MessagesCache) RemoveBySender(destination string, destinationDevice int64,
	sender string, timestamp int64) (*types.OutcomingMessagex, error) {

	serialized, err := internal.removeOperation.RemoveBySender(destination, destinationDevice, sender, timestamp)
	if err != nil {
		return nil, err
	}
	if serialized == nil {
		return nil, nil
	}

	var envelope textsecure.Envelope
	if err = proto.Unmarshal(serialized, &envelope); err != nil {
		return nil, err
	}
	entity := constructEntityFromEnvelope(0, &envelope)
	return &entity, nil
}

// 清理设备消息
func (MessagesCache) ClearDevice(destination string, deviceID int64) error {
	return internal.removeOperation.Clear(destination, deviceID)
}



