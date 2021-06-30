package chat

import "secret-im/service/signalserver/cmd/api/internal/entities"

// 已连接的设备信息

type ConnectedDevice struct {
	Number string
	UUID  string
	Device entities.Device
}


//  上下文信息
type SessionContext struct {
	Clientx *Clientx
	Session *Session
	Device *ConnectedDevice

}

// 消息处理接口
type SessionHandler interface {
	// 传递消息，比如消息的转发，在sessio.Develidy调用
	OnMessage(interface{})
	// websocket链接建立时调用
	OnWebSocketConnect(ctx *SessionContext)
	//  在 session.onWebsocketClose调用
	OnWebSocketDisconnect()
}

// 创建处理程序
type MakeSessionHandler func() SessionHandler