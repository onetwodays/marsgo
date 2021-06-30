package handlers

import "secret-im/service/signalserver/cmd/api/textsecure"

// 消息信息
type storedMessageInfo struct {
	id     int64
	cached bool
}

// 发送消息参数
type sendMessageRequest struct {
	message textsecure.Envelope
	info    *storedMessageInfo
	requery bool
}


// 鉴权session处理
type AuthenticatedHandler struct {
	connectionID  string
	//toBeSent      *queue.SyncQueue
	//context       *websocket.SessionContext
	//receiptSender *push.ReceiptSender
}

