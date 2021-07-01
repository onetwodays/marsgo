package message

import (
	uuid "github.com/satori/go.uuid"
	"secret-im/service/signalserver/cmd/api/internal/signal"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

// 消息管理器
type MessagesManager struct {
}

// 插入消息
func (MessagesManager) Insert(destination string, destinationDevice int64, message *textsecure.Envelope) error {
	guid := uuid.NewV4().String()
	return signal.SC.MessageCacheOper.Insert(guid, destination, destinationDevice, message)
}

// 获取设备消息
func (MessagesManager) GetMessagesForDevice(
	destination string, destinationDevice int64) (types.GetPendingMsgsRes, error) {
	messages, err := Messages{}.Load(destination, destinationDevice, ResultSetChunkSize)
	if err != nil {
		return types.GetPendingMsgsRes{}, err
	}

	if len(messages) < ResultSetChunkSize {
		cachedMessages, err := signal.SC.MessageCacheOper.Get(destination, destinationDevice, ResultSetChunkSize-len(messages))
		if err != nil {
			return types.GetPendingMsgsRes{}, err
		}
		messages = append(messages, cachedMessages...)
	}

	result := types.GetPendingMsgsRes{
		List: messages,
		More:     len(messages) >= ResultSetChunkSize,
	}
	return result, nil
}

// 清理用户消息
func (MessagesManager) Clear(destination string) error {
	account, err := AccountsManager{}.GetByNumber(destination)
	if err != nil {
		return err
	}
	for _, device := range account.Devices {
		err = MessagesCache{}.ClearDevice(destination, device.ID)
		if err != nil {
			return err
		}
	}
	return Messages{}.Clear(destination)
}

// 清理设备消息
func (MessagesManager) ClearDevice(destination string, deviceID int64) error {
	err := MessagesCache{}.ClearDevice(destination, deviceID)
	if err != nil && err != redis.Nil {
		return err
	}
	return Messages{}.ClearDevice(destination, deviceID)
}

// 删除消息
func (MessagesManager) Delete(destination string, deviceID, id int64, cached bool) error {
	if !cached {
		return Messages{}.Remove(destination, id)
	}
	return MessagesCache{}.Remove(destination, deviceID, id)
}

// 根据GUID删除消息
func (MessagesManager) DeleteByGUID(destination string, deviceID int64,
	guid string) (*entities.OutgoingMessageEntity, error) {

	removed, err := MessagesCache{}.RemoveByGUID(destination, deviceID, guid)
	if err != nil {
		return nil, err
	}
	if removed != nil {
		return removed, nil
	}

	message, err := Messages{}.RemoveByGUID(destination, guid)
	if err != nil {
		return nil, err
	}
	entity := constructEntityFromMessage(message)
	return &entity, nil
}

// 根据发送者删除消息
func (MessagesManager) DeleteBySender(destination string, destinationDevice int64,
	source string, timestamp int64) (*entities.OutgoingMessageEntity, error) {

	removed, err := MessagesCache{}.RemoveBySender(destination, destinationDevice, source, timestamp)
	if err != nil {
		return nil, err
	}
	if removed != nil {
		return removed, nil
	}

	message, err := Messages{}.RemoveBySender(destination, destinationDevice, source, timestamp)
	if err != nil {
		return nil, err
	}
	entity := constructEntityFromMessage(message)
	return &entity, nil
}
