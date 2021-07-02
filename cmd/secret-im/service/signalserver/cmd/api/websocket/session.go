package websocket

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/api/websocket/factory"
	"sync"
	"sync/atomic"
	"time"
)


const (
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
	maxMessageLimit = 4096
)

// 接收消息处理
type MessageHandler func(int64, []byte)

// 连接关闭处理
type CloseHandler func(int64, int, string) error

//请求客户端
type Clientx struct {
	session    *Session
	requestMap sync.Map
}

// 处理响应消息
func (client *Clientx) handleResponse(message *textsecure.WebSocketResponseMessage) {
	id := int64(message.GetId())
	value, ok := client.requestMap.Load(id)
	if !ok {
		return
	}
	future := value.(*Future)
	future.setResult(message, nil)
	client.requestMap.Delete(id)
	if client.session.context.Device != nil {
		logx.Info("[Client] recv response from ", client.session.context.Device.Number,
			" WebSocketResponseMessage id:", id)
	}
}

// 发送请求
func (client *Clientx) SendRequest(verb, path string, headers []string, body []byte) *Future {
	future := newFuture()
	if client.session.IsClosed(){
		err:=errors.New("closed")
		future.setResult(nil,err)
		return future
	}
	requestID := utils.SecureRandInt64()
	client.requestMap.Store(requestID, future)

	requestMessage := factory.CreateRequest(uint64(requestID), verb, path, headers, body)
	data, err := proto.Marshal(requestMessage)
	if err != nil {
		client.requestMap.Delete(requestID)
		future.setResult(nil, err)
		return future
	}

	if err = client.session.Send(data); err != nil {
		client.requestMap.Delete(requestID)
		future.setResult(nil, err)
	}

	if client.session.context.Device != nil {

		logx.Info("[Client] send request to ",client.session.context.Device.Number,
			" id:",   requestID,
			" verb:", verb,
			" path:", path)
	}
	return future
}

func (client *Clientx) CancelAll() {
	err := errors.New("canceled")
	client.requestMap.Range(func(key, value interface{}) bool {
		future := value.(*Future)
		future.setResult(nil, err)
		client.requestMap.Delete(key)
		return true
	})
}

// 会话信息
type Session struct {
	id           int64
	closed       int32
	toBeSent     chan []byte
	conn         *websocket.Conn
	client       *Clientx
	router       http.Handler
	context      *SessionContext
	handler      SessionHandler
	closeHandler CloseHandler
}

type Options struct {
	id            int64
	router        http.Handler
	conn          *websocket.Conn
	handler       SessionHandler
	closeHandler  CloseHandler
	pushSender    *push.Sender
	pubSubManager *pubsub.Manager
	device        *ConnectedDevice
}

func newSession(option Options) *Session {
	session := &Session{
		id:           option.id,
		conn:         option.conn,
		router:       option.router,
		handler:      option.handler,      // 全局函数
		closeHandler: option.closeHandler, //session_manager.onClosed
		toBeSent:     make(chan []byte, 256),
	}
	session.client = &Clientx{session: session}
	session.context = &SessionContext{ //这个是为了在http.Request传送
		Device:        option.device,
		Session:       session,
		Clientx:       session.client,
		PushSender:    option.pushSender,
		PubSubManager: option.pubSubManager,
	}
	return session
}

// 发送数据
func (session *Session) Send(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if session.IsClosed() {
		return errors.New("closed")
	}
	session.toBeSent <- data
	return nil
}

// 是否已经关闭
func (session *Session) IsClosed() bool {
	return atomic.LoadInt32(&session.closed) == 1
}

// 获取对端地址
func (session *Session) RemoteAddr() string {
	return session.conn.RemoteAddr().String()
}

// 传递消息
func (session *Session) DeliverMessage(msg interface{}) {
	session.handler.OnMessage(msg)
}

// 关闭Session
func (session *Session) Close(code int, text string) error {
	data := websocket.FormatCloseMessage(code, text)
	if len(data) == 0 {
		return errors.New("invalid code")
	}

	var err error
	if atomic.CompareAndSwapInt32(&session.closed, 0, 1) {
		close(session.toBeSent)
		deadline := time.Now().Add(writeWait)
		err = session.conn.WriteControl(websocket.CloseMessage, data, deadline)
		if err != nil {
			logx.Errorf("id=%d,reason=%s,write control message error", session.id, err.Error())
		}
		if err = session.conn.Close(); err != nil {
			logx.Errorf("id=%d,reason=%s,close session connection error", session.id, err.Error())
		}
	}
	return err
}

// 处理读操作
func (session *Session) handleRead() {
	session.conn.SetReadLimit(maxMessageLimit)
	session.conn.SetReadDeadline(time.Now().Add(pongWait))
	session.conn.SetPongHandler(func(appData string) error {
		session.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	var code int
	var text string
	defer func() {
		session.onWebsocketClose(code, text)
	}()
	for {
		messageType, data, err := session.conn.ReadMessage()
		if err != nil {
			closeErr, ok := err.(*websocket.CloseError)
			if ok {
				code, text = closeErr.Code, closeErr.Text
			} else {
				code, text = websocket.CloseTryAgainLater, err.Error()
			}
		}
		switch messageType {
		case websocket.PingMessage:
			session.conn.WriteMessage(websocket.PongMessage, nil)
		case websocket.TextMessage:
			logx.Info("recv text message:", string(data))
		case websocket.BinaryMessage:
			session.onWebsocketBinary(data)
		}
	}
}

//处理写操作
func (session *Session) handleWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case data, ok := <-session.toBeSent:
			if !ok {
				return
			}
			session.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := session.conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				continue
			}

		case <-ticker.C:
			session.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := session.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				continue
			}
		}
	}
}

// 处理请求消息
func (session *Session) handleRequest(request *textsecure.WebSocketRequestMessage) {
	number := "null"
	if session.context.Device != nil {
		number = session.context.Device.Number
	}
	logx.Infof("{handle request from account:%s,id:%d,verb:%s,path:%s", number, request.GetId(), request.GetVerb(), request.GetPath())
	message := HandleHTTPRequest(session.context, session.router, request)
	if data, err := proto.Marshal(message); err != nil {
		session.Send(data)
	}
}

// 接收消息处理
func (session *Session) onWebsocketBinary(data []byte) {
	var message textsecure.WebSocketMessage
	if err := proto.Unmarshal(data, &message); err != nil {
		session.Close(1018, "Badly formatted")
		return
	}

	switch message.GetType() {
	case textsecure.WebSocketMessage_REQUEST:
		if message.GetRequest() == nil {
			break
		}
		session.handleRequest(message.GetRequest())
		return

	case textsecure.WebSocketMessage_RESPONSE:
		if message.GetResponse() == nil {
			break
		}
		session.client.handleResponse(message.GetResponse())
		return
	}
	session.Close(1018, "Badly formatted")

}

func (session *Session) onWebsocketClose(code int, text string) error {
	if atomic.CompareAndSwapInt32(&session.closed, 0, 1) {
		close(session.toBeSent)
	}
	session.client.CancelAll()
	session.handler.OnWebSocketDisconnect() //解除订阅
	return session.closeHandler(session.id, code, text) //从会话容器删掉此会话
}
