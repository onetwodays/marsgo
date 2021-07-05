package websocket

import (
	"bytes"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/pkg/queue"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"strconv"
	"time"
)


// 消息信息
type storedMessageInfo struct {
	id     int64
	cached bool
}

// 发送消息参数
type sendMessageRequest struct {
	message *textsecure.Envelope
	info    *storedMessageInfo
	requery bool
}

// 鉴权session处理
type AuthenticatedHandler struct {
	connectionID  string
	toBeSent      *queue.SyncQueue
	context       *SessionContext
	receiptSender *push.ReceiptSender
}

func constructEnvelopFromTypeMessage(message *types.OutcomingMessagex) *textsecure.Envelope {
	var jsb []byte
	var err error
	if len(message.Message)!=0{
		jsb,err=base64.StdEncoding.DecodeString(message.Message)
	}
	if err!=nil{

		return nil
	}

	jtt,err:=base64.StdEncoding.DecodeString(message.Content)
	if err!=nil{
		return nil
	}
	envelope := &textsecure.Envelope{
		Type:            textsecure.Envelope_Type(message.Type),
		Timestamp:       uint64(message.Timestamp),
		ServerTimestamp: uint64(message.ServerTimestamp),
		Source:          message.Source,
		SourceUuid: message.SourceUuid,
		SourceDevice:    uint32(message.SourceDevice),
		LegacyMessage: jsb,
		Content: jtt,
		Relay: message.Relay,
		ServerGuid: message.Guid,
	}
	return envelope
}

// 好像是群发消息事件
func (handler *AuthenticatedHandler) OnMessage(msg interface{}) {
	switch msg.(type) {
	case textsecure.ChannelEnvelope:
		handler.toBeSent.Push(msg)
	}
}

// 连接成功
func (handler *AuthenticatedHandler) OnWebSocketConnect(context *SessionContext) {
	// 初始状态
	if context.Device == nil {
		return
	}

	device := context.Device
	handler.context = context
	handler.toBeSent = queue.NewSyncQueue()
	handler.receiptSender = push.NewReceiptSender(context.PushSender)
	handler.connectionID = strconv.FormatInt(utils.SecureRandInt64(), 10)

	// 设备上线通知(断开其他连接)

	connectMessage := textsecure.PubSubMessage{}
	connectMessage.Type = textsecure.PubSubMessage_CONNECTED
	connectMessage.Content = []byte(handler.connectionID)
	address := push.Address{Number: device.Number, DeviceID: device.Device.ID}
	_, err := handler.context.PubSubManager.Publish(address.Serialize(), &connectMessage) //推到redis了
	if err != nil {
		handler.context.Session.Close(1000, "OK")
		return
	}

	// 标记设备上线


	partition := 1 //conf.GetServer().PartitionID
	storage.DevicesManager{}.SetOnline(device.UUID, device.Device.ID, partition)
	/*

	if true {
		storage.ChannelParticipantsManager{}.SetOnline(device.UUID, device.Device.ID, partition)
	}

	 */


	// 订阅设备消息
	if err = handler.context.PubSubManager.Subscribe(address.Serialize(), handler); err != nil {
		handler.context.Session.Close(1000, "OK")
		return
	}

	go handler.handleSendMessage()
	logx.Infof("number:%s,device_id:%d [Authenticated] device on-line",device.Number,device.Device.ID)

}

// 断开连接
func (handler *AuthenticatedHandler) OnWebSocketDisconnect() {
	if handler.context == nil {
		return
	}

	handler.toBeSent.Close()
	device := handler.context.Device
	address := push.Address{Number: device.Number, DeviceID: device.Device.ID}
	handler.context.PubSubManager.Unsubscribe(address.Serialize(), handler)


	partition := 1
	storage.DevicesManager{}.SetOffline(device.UUID, device.Device.ID, partition)
	/*
	if conf.GetServer().EnableChannel {
		storage.ChannelParticipantsManager{}.SetOffline(device.UUID, device.Device.ID)
	}

	 */
	logx.Infof("number:%s,device_id:%d [Authenticated] device off-line",device.Number,device.Device.ID)

}

// 从redis拉取消息后的处理流程
func (handler *AuthenticatedHandler) OnDispatchMessage(channel string, message *textsecure.PubSubMessage) {
	switch message.GetType() {
	// 推送消息
	case textsecure.PubSubMessage_DELIVER:
		var envelope textsecure.Envelope
		err := proto.Unmarshal(message.GetContent(), &envelope)
		if err == nil {
			handler.toBeSent.Push(sendMessageRequest{message: &envelope}) //放到自己的队列里面
		} else {
			logx.Errorf("[Authenticated] protobuf parse error,reason:%s",err.Error())

		}
	// 查询消息
	case textsecure.PubSubMessage_QUERY_DB:
		handler.processStoredMessages()
	// 连接成功
	case textsecure.PubSubMessage_CONNECTED:
		if bytes.Equal(message.Content, []byte(handler.connectionID)) {
			handler.context.Session.Close(1000, "OK")
		}
	default:
		logx.Errorf("[Authenticated] unknown pubsub message")

	}
}

// 订阅成功
func (handler *AuthenticatedHandler) OnDispatchSubscribed(channel string) {
	handler.processStoredMessages()
}

// 取消订阅
func (handler *AuthenticatedHandler) OnDispatchUnsubscribed(channel string) {
	handler.context.Session.Close(1000, "OK")
}

// 是否成功响应
func (handler *AuthenticatedHandler) isSuccessResponse(response *textsecure.WebSocketResponseMessage) bool {
	return response != nil && response.GetStatus() >= 200 && response.GetStatus() < 300
}

// 发送消息回执(是否已达)
func (handler *AuthenticatedHandler) sendDeliveryReceiptFor(message *textsecure.Envelope) {
	if len(message.Source) == 0 {
		return
	}

	device := handler.context.Device
	logx.Info(" [Authenticated] send delivery receipt","account:",device.Number," timestamp:",message.GetTimestamp())

	err := handler.receiptSender.SendReceipt(
		device.Number, device.UUID, device.Device.ID, message.GetSource(), int64(message.GetTimestamp()))
	if err != nil {
		logx.Error("[Authenticated] failed to send receipt ","source:",device.Number,"destination:", message.GetSource(),"timestamp:",message.GetTimestamp(),"reason:",err)

	}
}

// 消息持久化
func (handler *AuthenticatedHandler) requeueMessage(message *textsecure.Envelope) {
	device := handler.context.Device
	logx.Info("account:", device.Number,"[Authenticated] requeue message")


	err := handler.context.PushSender.GetWebsocketSender().QueueMessage(device.Number, device.Device.ID, message)
	if err != nil {
		logx.Error("[Authenticated] failed to requeue message ","account:",device.Number,"device:",device.Device.ID,"timestamp:",message.Timestamp,"reason:",err,)

	}

	err = handler.context.PushSender.SendQueuedNotification(device.Number, &device.Device)

}

// 拉取未读消息
func (handler *AuthenticatedHandler) processStoredMessages() {
	device := handler.context.Device

	logx.Info("account:", device.Number,"[Authenticated] process stored messages")


	messages, err := storage.MessagesManager{}.GetMessagesForDevice(device.Number, device.Device.ID)
	if err != nil {
		logx.Error("[Authenticated] failed to get messages for device,reason:",err)
		return
	}

	for idx, message := range messages.List {

		envelope := constructEnvelopFromTypeMessage(&message)
		if envelope!=nil{
			continue
		}

		request := sendMessageRequest{
			message: envelope,
			info:    &storedMessageInfo{id: message.Id, cached: message.Cached},
			requery: idx == len(messages.List)-1 && messages.More,
		}
		handler.toBeSent.Push(request)
	}

	if !messages.More {
		handler.context.Clientx.SendRequest("PUT", "/api/v1/queue/empty", nil, nil)
	}


}

// 发送消息
func (handler *AuthenticatedHandler) sendMessage(message *textsecure.Envelope, info *storedMessageInfo, requery bool) {
	var body []byte
	var header string
	device := handler.context.Device
	if len(device.Device.SignalingKey) ==0 {
		header = "X-Signal-Key: false"
		body, _ = proto.Marshal(message)
	} else {
		body, _ = proto.Marshal(message) // 暂时这样
		/*
		header = "X-Signal-Key: true"
		encrypted, err := entities.NewEncryptedOutgoingMessage(&message, *device.Device.SignalingKey)
		if err != nil {
			log.WithFields(log.Fields{
				"reason": err,
			}).Error("[Authenticated] failed to encrypt message with use signaling key")
		}
		body = encrypted.Serialized

		 */
	}

	messageType := "message"
	isReceipt := message.GetType() == textsecure.Envelope_RECEIPT
	if isReceipt {
		messageType = "message receipt"
	}

	logx.Info("account:", device.Number,"timestamp:",message.GetTimestamp(),"[Authenticated] deliver ",messageType)



	future := handler.context.Clientx.SendRequest("PUT", "/api/v1/message", []string{header}, body)
	response, err := future.Get(time.Second * 10)
	if err != nil {
		// 消息持久化
		if info == nil {
			handler.requeueMessage(message)
		}
		logx.Info("[Authenticated] failed to send request ","reason:",err,"timestamp:", message.GetTimestamp() )

		handler.context.Session.Close(1001, "Request error")
		return
	}
	logx.Info("[Authenticated] deliver " + messageType + " ack","account:",   device.Number,"timestamp:", message.GetTimestamp(),)



	if handler.isSuccessResponse(response) {
		if info != nil {
			// 删除消息记录

			err = storage.MessagesManager{}.Delete(device.Number, device.Device.ID, info.id, info.cached)
			if err != nil {
				logx.Error("[Authenticated] failed to delete message:",
					" number:",device.Number,
					" device:", device.Device.ID,
					" id:",     info.id,
					" cached:", info.cached,
					" reason:", err)

			}


		}

		if !isReceipt {
			handler.sendDeliveryReceiptFor(message)
		}

		// 推送未读消息
		if requery {
			handler.processStoredMessages()
		}
	} else if info == nil {
		handler.requeueMessage(message)
	}
}

// 发送频道消息
func (handler *AuthenticatedHandler) sendChannelMessage(message *textsecure.ChannelEnvelope) {
	var body []byte
	var header string
	device := handler.context.Device

	header = "X-Signal-Key: false"
	body, _ = proto.Marshal(message)

	logx.Info("[Authenticated] deliver channel message ",
		" account:",    device.Number,
		" device:",     device.Device.ID,
		" channel_id:", message.GetChannelId(),
		" message_id:", message.GetId())



	future := handler.context.Clientx.SendRequest("PUT", "/api/v1/channel/message", []string{header}, body)
	_, err := future.Get(time.Second * 10)
	if err != nil {
		logx.Error("[Authenticated] failed to send request: ",
			" account:",    device.Number,
			" device:",   device.Device.ID,
			" channel_id:", message.GetChannelId(),
			" message_id:", message.GetId(),
			" reason:",     err)
		handler.context.Session.Close(1001, "Request error")
		return
	}

	logx.Info("[Authenticated] deliver channel message ack ",
		" account:",   device.Number,
		" device:",  device.Device.ID,
		" channel_id:", message.GetChannelId(),
		" message_id:", message.GetId())


}

// 处理发送消息
func (handler *AuthenticatedHandler) handleSendMessage() {
	for {
		request := handler.toBeSent.Pop()
		if request == nil {
			break
		}
		switch request.(type) {
		case textsecure.ChannelEnvelope:
			data := request.(*textsecure.ChannelEnvelope)
			handler.sendChannelMessage(data)
		case sendMessageRequest:
			data := request.(sendMessageRequest)
			handler.sendMessage(data.message, data.info, data.requery)
		}
	}
}

