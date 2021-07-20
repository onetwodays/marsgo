package websocket

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/timex"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/svc/pubsub"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/api/textsecure"
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
	sessionName   string
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
	account       *entities.Account
}

func newSession(option Options) *Session {
	session := &Session{
		id:           option.id,
		conn:         option.conn,
		router:       option.router,
		handler:      option.handler,      // 全局函数
		closeHandler: option.closeHandler, //session_manager.onClosed
		toBeSent:     make(chan []byte, 1024),

	}
	session.client = &Clientx{session: session}
	session.context = &SessionContext{ //这个是为了在http.Request传送
		Device:        option.device,
		Session:       session,
		Clientx:       session.client,
		PushSender:    option.pushSender,
		PubSubManager: option.pubSubManager,
		account:       option.account,
	}
	if session.context.account==nil{
		session.sessionName=session.conn.RemoteAddr().String()+":anonymous"
	}else {
		session.sessionName=session.conn.RemoteAddr().String()+":"+session.context.account.Number
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

//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
//c.hub.broadcast <- message //读到的消息放到广播里面, to write redis,测试用

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
			logx.Error("websocket读消息发生错误:",err.Error())
			closeErr, ok := err.(*websocket.CloseError)
			if ok {
				code, text = closeErr.Code, closeErr.Text
			} else {
				code, text = websocket.CloseTryAgainLater, err.Error()
			}
			break
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
			//w, err := c.conn.NextWriter(websocket.TextMessage)
			w, err := session.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				logx.Errorf("ws sender(%d) get NextWriter error:%s ",session.id,err.Error())
				return
			}
			w.Write(data)

			// Add queued chat messages to the current websocket message.
			n := len(session.toBeSent)
			for i := 0; i < n; i++ {
				//w.Write(newline) //如果是二进制，这行请注释
				w.Write(<-session.toBeSent)
			}
			if err := w.Close(); err != nil { //写完关闭w，不是关闭连接
				logx.Errorf("ws sender(%d) close got NextWriter error:%s",session.id,err.Error())
				return
			}

			/*
			err := session.conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				continue
			}

			 */

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
	startTime:=timex.Now()

	/*
	number := "null"
	if session.context.Device != nil {
		number = session.context.Device.Number
	}
	 */
	//logx.Infof("handle ws request from account:%s,id:%d,verb:%s,path:%s", )
	message := HandleHTTPRequest(session.context, session.router, request)


	if data, err := proto.Marshal(message); err == nil {
		session.Send(data)

	}else{
		logx.Error("(session *Session) handleRequest:",err)
	}
	if request.Path!="/v1/keepalive" {
		logx.WithDuration(timex.Since(startTime)).Infof("完成 [%s] wsreq处理 [verb:%s,path:%s,code:%d,id:%d,(%s)]",
			session.sessionName,
			request.GetVerb(),
			request.GetPath(),
			message.GetResponse().Status,
			request.GetId(),
			message.String())
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
