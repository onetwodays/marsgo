package logic

import (
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"secret-im/service/signalserver/cmd/api/internal/logic"

	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/auth/helper"
	"secret-im/service/signalserver/cmd/api/internal/entities"

	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PutMsgsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutMsgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutMsgsLogic {
	return PutMsgsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutMsgsLogic) PutMsgs(r *http.Request, sender string, req types.PutMessagesReq) (*types.PutMessagesRes, error) {

	header := r.Header.Get(helper.UNIDENTIFIED)
	accessKey, _ := auth.NewAnonymous(header)

	source,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if accessKey == nil {
		if err!=nil {
			return nil, shared.Status(http.StatusUnauthorized, err.Error())
		}
	}

	destinationName := auth.NewAmbiguousIdentifier(req.Destination)

	if source != nil && source.IsFor(destinationName) {
		//todo 限制次数
	}

	var destination *entities.Account
	isSyncMessage := source != nil && source.IsFor(destinationName)
	if isSyncMessage {
		destination = source
	} else {
		_, destination, _ = storage.AccountManager{}.Get(destinationName)
	}
	if destination == nil {
		logx.Error("destination==nil")
		return nil, shared.Status(http.StatusNotFound, "destination==nil")
	}

	_, ok := helper.OptionalAccess{}.Verify(source, accessKey, destination)
	if !ok {
		reason:="helper.OptionalAccess{}.Verify(source,accessKey,destination) fail"
		logx.Error(reason)
		//return nil,shared.Status(code,reason)
	}

	//验证完整设备列表
	missingDevices, ok := l.validateCompleteDeviceList(destination, req.Messages, isSyncMessage)
	if !ok {
		jsb, _ := json.Marshal(missingDevices)
		logx.Error("验证完备设备列表失败:",string(jsb))
		//return nil,shared.Status(http.StatusConflict,string(jsb))
	}

	// 验证注册ID
	staleDevices, ok := l.validateRegistrationIds(destination, req.Messages)
	if !ok {
		jsb, _ := json.Marshal(staleDevices)
		logx.Error("验证注册id失败:",string(jsb))
		//return nil,shared.Status(http.StatusGone,string(jsb))
	}

	for i, _ := range req.Messages {
		msg := &req.Messages[i]

		destinationDevice, ok := destination.GetDevice(int64(msg.DestinationDeviceId))
		if ok {
			err := l.sendMessage(source, destination, &destinationDevice.Device, req.Timestamp, req.Online, msg)
			if err != nil {
				return nil, shared.Status(http.StatusNotFound, err.Error())
			}
		}
	}

	needsSync := !isSyncMessage && source != nil && source.GetEnabledDeviceCount() > 1

	return &types.PutMessagesRes{NeedsSync: needsSync}, nil
}

//
func (l *PutMsgsLogic) sendMessage(source,
	destinationAccount *entities.Account,
	destinationDevice *entities.Device,
	timestamp int64,
	online bool,
	incomingMessage *types.IncomingMessagex) error {

	now := time.Now().Unix()
	if timestamp == 0 {
		timestamp = now
	}
	messageBuilder := &textsecure.Envelope{
		Type:            textsecure.GetEnvelopeType(incomingMessage.Type),
		Timestamp:       uint64(timestamp),
		ServerTimestamp: uint64(now),
		ServerGuid: uuid.NewV4().String(), // 这里新生成一个,原来是推送redis不成功，才生成一个，这里提前生成
	}
	if source != nil {
		messageBuilder.Source = source.Number
		messageBuilder.SourceUuid = source.UUID
		messageBuilder.SourceDevice = uint32(source.AuthenticatedDevice.ID)
	}
	messageBody := textsecure.DecodeMessage(incomingMessage.Body)
	messageContent := textsecure.DecodeMessage(incomingMessage.Content)
	if messageBody != nil {
		messageBuilder.LegacyMessage = messageBody
	}
	if messageContent != nil {
		messageBuilder.Content = messageContent
	}
	delivered, err := l.svcCtx.PushSender.SendMessage(destinationAccount.Number, destinationDevice, messageBuilder, online)
	if err == nil {
		logx.Info("[put_msg_logic] send message success",
			" delivered:", delivered,
			" online:", online,
			" destination:", destinationAccount.Number,
			" FetchesMessages:", destinationDevice.FetchesMessages,
			" timestamp:", messageBuilder.GetTimestamp())
	} else {
		logx.Error("[Message] failed to send message reason:", err)
		if destinationDevice.ID == entities.DeviceMasterID {
			return err
		}
	}
	return nil
}

// 验证注册ID
func (l *PutMsgsLogic) validateRegistrationIds(
	account *entities.Account, messages []types.IncomingMessagex) ([]int64, bool) {
	var staleDevices []int64
	for _, message := range messages {
		device, ok := account.GetDevice(int64(message.DestinationDeviceId))
		if ok && message.DestinationRegistrationId > 0 && message.DestinationRegistrationId != device.RegistrationID {
			staleDevices = append(staleDevices, device.ID)
		}
	}

	if len(staleDevices) > 0 {
		return staleDevices, false
	}
	return nil, true
}

// 验证完整设备列表
func (l *PutMsgsLogic) validateCompleteDeviceList(account *entities.Account,
	messages []types.IncomingMessagex, isSyncMessage bool) (entities.MismatchedDevices, bool) {
	messageDeviceIDs := make(map[int64]struct{})
	accountDeviceIDs := make(map[int64]struct{})

	extraDeviceIDs := make(map[int]struct{})
	missingDeviceIDs := make(map[int64]struct{})

	for _, message := range messages {
		messageDeviceIDs[int64(message.DestinationDeviceId)] = struct{}{}
	}
	//logx.Info("消息中的设备id列表:",messageDeviceIDs)
	for _, device := range account.Devices {
		if device.IsEnabled() && !(isSyncMessage && device.ID == account.AuthenticatedDevice.ID) {
			accountDeviceIDs[device.ID] = struct{}{}

			if _, ok := messageDeviceIDs[device.ID]; !ok {
				missingDeviceIDs[device.ID] = struct{}{} //帐号中有，消息中没有
			}
		}
	}
	//logx.Info("帐号中有但不在消息中的可用设备id列表是:",missingDeviceIDs)
	//logx.Infof("帐号%s可用设备id列表是%v",account.Number,accountDeviceIDs)


	for _, message := range messages {
		if _, ok := accountDeviceIDs[int64(message.DestinationDeviceId)]; !ok {
			extraDeviceIDs[message.DestinationDeviceId] = struct{}{} //消息中有设备中没有
		}
	}
	//logx.Info("**消息中有，帐号中没有的可用设备id列表:",extraDeviceIDs)

	if len(missingDeviceIDs) > 0 || len(extraDeviceIDs) > 0 {
		devices := entities.MismatchedDevices{
			ExtraDevices:   make([]int64, 0, len(extraDeviceIDs)),
			MissingDevices: make([]int64, 0, len(missingDeviceIDs)),
		}
		for id := range extraDeviceIDs {
			devices.ExtraDevices = append(devices.ExtraDevices, int64(id))
		}
		for id := range missingDeviceIDs {
			devices.MissingDevices = append(devices.MissingDevices, id)
		}
		return devices, false
	}
	return entities.MismatchedDevices{}, true
}
