package storage

import (
	"encoding/base64"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"time"
)

// MessagesManager 消息管理器
type MessagesManager struct {
}


func constructEntityFromMessage(message *model.TMessages) *types.OutcomingMessagex {
	return &types.OutcomingMessagex{
		Id:              message.Id,
		Cached:          false,
		Guid:            message.Guid,
		Type:            int(message.Type),
		Relay:           message.Relay,
		Timestamp:       message.Timestamp,
		Source:          message.Source,
		SourceUuid:      message.SourceUuid,
		SourceDevice:    message.SourceDevice,
		Message:         message.Message,
		Content:         message.Content,
		ServerTimestamp: message.ServerTimestamp,
	}
}

// 插入消息
func (MessagesManager) Insert(destination string, destinationDevice int64, message *textsecure.Envelope) error {
	guid := uuid.NewV4().String()
	return new(MessagesCache).Insert(guid, destination, destinationDevice, message)
}

// 获取设备消息
func (MessagesManager) GetMessagesForDevice(
	destination string, destinationDevice int64) (types.GetPendingMsgsRes, error) {

	messages, err := internal.msgDB.FindManyByDst(destination, destinationDevice, ResultSetChunkSize)
	if err != nil {
		return types.GetPendingMsgsRes{}, err
	}
	list := make([]types.OutcomingMessagex, 0, len(messages))
	for i := range messages {
		list = append(list, *constructEntityFromMessage(&messages[i]))
	}

	if len(messages) < ResultSetChunkSize {
		cachedMessages, err := new(MessagesCache).Get(destination, destinationDevice, ResultSetChunkSize-len(messages))
		if err != nil {
			return types.GetPendingMsgsRes{}, err
		}
		list = append(list, cachedMessages...)
	}

	result := types.GetPendingMsgsRes{
		List: list,
		More: len(messages) >= ResultSetChunkSize,
	}
	return result, nil
}

// 清理用户消息
func (MessagesManager) Clear(destination string) error {

	appAccount, err := AccountManager{}.GetByNumber(destination)
	if err != nil {
		return err
	}
	// 先清空redis
	for _, device := range appAccount.Devices {
		err = MessagesCache{}.ClearDevice(destination, device.ID)
		if err != nil {
			return err
		}
	}
	// 在清空database
	return internal.msgDB.DeleteManyByDestination(destination)
}

// 清理设备消息
func (MessagesManager) ClearDevice(destination string, deviceID int64) error {

	err := MessagesCache{}.ClearDevice(destination, deviceID)
	if err != nil && err != redis.Nil {
		return err
	}
	return internal.msgDB.DeleteManyByDestinationDeviceId(destination, deviceID)
}

// 删除消息
func (MessagesManager) Delete(destination string, deviceID, id int64, cached bool) error {
	if !cached {
		return internal.msgDB.Remove(destination, id)
	}
	return MessagesCache{}.Remove(destination, deviceID, id)
}

// 根据GUID删除消息
func (MessagesManager) DeleteByGUID(destination string, deviceID int64,
	guid string) (*types.OutcomingMessagex, error) {

	removed, err := MessagesCache{}.RemoveByGUID(destination, deviceID, guid)
	if err != nil {
		return nil, err
	}
	if removed != nil {
		return removed, nil
	}

	msg, err := internal.msgDB.RemoveByGUID(destination, guid)
	if err != nil {
		return nil, err
	}
	entity := constructEntityFromMessage(msg)
	return entity, nil
}

// 根据发送者删除消息
func (MessagesManager) DeleteBySender(destination string, destinationDevice int64,
	source string, timestamp int64) (*types.OutcomingMessagex, error) {

	removed, err := MessagesCache{}.RemoveBySender(destination, destinationDevice, source, timestamp)
	if err != nil {
		return nil, err
	}
	if removed != nil {
		return removed, nil
	}

	msg, err := internal.msgDB.RemoveBySender(destination, destinationDevice, source, timestamp)
	if err != nil {
		return nil, err
	}
	entity := constructEntityFromMessage(msg)
	return entity, nil
}
func (MessagesManager) Store(guid string,message *textsecure.Envelope,destination string,destinationDevice int64)(*model.TMessages,error){
	recode:=&model.TMessages{
		Guid: guid,
		Destination: destination,
		DestinationDevice: destinationDevice,
		Type: int64(message.GetType()),
		Relay: message.GetRelay(),
		Timestamp: int64(message.GetTimestamp()),
		ServerTimestamp: int64(message.GetServerTimestamp()),
		Source: message.GetSource(),
		SourceDevice: int64(message.GetSourceDevice()),
		Message: base64.StdEncoding.EncodeToString(message.GetLegacyMessage()),
		Content: base64.StdEncoding.EncodeToString(message.GetContent()),
		Ctime: time.Now(),
	}

	_,err:=internal.msgDB.Insert(*recode)
	if err!=nil{
		return nil,err
	}

	return recode,nil





}