package logic

import (
	"context"

	"encoding/json"

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
	// todo: add your logic here and delete this line

	header := r.Header.Get(helper.UNIDENTIFIED)
	accessKey, _ := auth.NewAnonymous(header)
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil && accessKey == nil {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return nil, shared.Status(http.StatusUnauthorized, reason)
	}
	source := appAccount.(*entities.Account)
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
		//reason:="helper.OptionalAccess{}.Verify(source,accessKey,destination) fail"
		//logx.Error(reason)
		//return nil,shared.Status(code,reason)
	}

	//验证完整设备列表
	missingDevices, ok := l.validateCompleteDeviceList(destination, req.Messages, isSyncMessage)
	if !ok {
		jsb, _ := json.Marshal(missingDevices)
		logx.Error(string(jsb))
		//return nil,shared.Status(http.StatusConflict,string(jsb))
	}

	// 验证注册ID
	staleDevices, ok := l.validateRegistrationIds(destination, req.Messages)
	if !ok {
		jsb, _ := json.Marshal(staleDevices)
		logx.Error(string(jsb))
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

		/*
			//如果在线，直接通过websocket推送出去
				//todo:send to redis

				row:=&model.TMessages{}
				row.Type=int64(msg.Type)
				row.Source= sender
				row.SourceUuid=account.Uuid
				row.SourceDevice=1
				row.Destination=req.Destination
				row.DestinationDevice=int64(msg.DestinationDeviceId)
				row.Timestamp=req.Timestamp
				row.Message="miss"//msg.Body
				row.Content=msg.Content  //
				row.Relay=msg.Relay
				row.Guid=utils.NewUuid() //消息的全局uuid
				row.Ctime=time.Now()

				if !online {
					_,err:=l.svcCtx.MsgsModel.Insert(*row)
					if err!=nil{
						return nil, shared.NewCodeError(shared.ERRCODE_SQLINSERT,err.Error())
					}
				} else {
					//if isOnline {
					if true {
						content,_ :=base64.StdEncoding.DecodeString(msg.Content)
						envelopePf:=&textsecure.Envelope{}
						envelopePf.Type=textsecure.GetEnvelopeType(msg.Type)
						envelopePf.SourceDevice=1
						envelopePf.Source=sender
						envelopePf.ServerGuid=row.Guid
						envelopePf.SourceUuid=account.Uuid
						envelopePf.ServerTimestamp=uint64(now)
						envelopePf.Timestamp=uint64(req.Timestamp)
						envelopePf.Relay=row.Relay
						//envelopePf.LegacyMessage=[]byte(row.Message)
						envelopePf.Content=content
						logx.Info("收件人的envelop:",envelopePf.String())
						contentPf,err:=proto.Marshal(envelopePf)
						if err!=nil{
							logx.Error("proto.Marshal(envelopePf):",err)
						}else{
							websocketReq:=&textsecure.WebSocketRequestMessage{}
							websocketReq.Id=msgId
							websocketReq.Headers=[]string{"X-Signal-Key: false","X-Signal-Timestamp:"+fmt.Sprintf("%d",now)}
							websocketReq.Path="/api/v1/message"
							websocketReq.Verb="PUT"
							websocketReq.Body=contentPf

							websocketMsg:=&textsecure.WebSocketMessage{}

							websocketMsg.Type=textsecure.WebSocketMessage_REQUEST
							websocketMsg.Request=websocketReq
							logx.Info("收件人最外层:",websocketMsg.String())
							msg,err:=proto.Marshal(websocketMsg)

							//pubsubMsg:=&textsecure.PubSubMessage{}
							//pubsubMsg.Type=textsecure.PubSubMessage_DELIVER
							//pubsubMsg.Content=contentPf
							//logx.Info("收件人最外层:",pubsubMsg.String())
							//msg,err:=proto.Marshal(pubsubMsg)


							if err!=nil{
								logx.Infof("proto.Marshal(websocketMsg) error:",err.Error() )
							}else{
								destContent=append(destContent,msg)

							}
						}
					}else{
						logx.Infof("destination is not online,not need websocket send,test brocast" )
					}
				}
		*/
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
		logx.Info("[Message] send message success",
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
	for _, device := range account.Devices {
		if device.IsEnabled() && !(isSyncMessage && device.ID == account.AuthenticatedDevice.ID) {
			accountDeviceIDs[device.ID] = struct{}{}

			if _, ok := messageDeviceIDs[device.ID]; !ok {
				missingDeviceIDs[device.ID] = struct{}{}
			}
		}
	}

	for _, message := range messages {
		if _, ok := accountDeviceIDs[int64(message.DestinationDeviceId)]; !ok {
			extraDeviceIDs[message.DestinationDeviceId] = struct{}{}
		}
	}

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
