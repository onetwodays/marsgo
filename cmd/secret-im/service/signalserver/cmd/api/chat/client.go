// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/timex"
	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/api/util"
	"strings"
	"sync"

	//"github.com/golang/protobuf/proto"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	//"secret-im/service/signalserver/cmd/api/textsecure"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {return true},
}

var xhub *Hub

func SetHub(hub *Hub)  {
	xhub = hub
}

func Ts(){

}

func init()  {
	xhub=NewHub()
	go xhub.Run()

}

func Broadcast(msg []byte){
	if xhub!=nil{
		xhub.broadcast<-msg
	}
}


// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id string
	conn *websocket.Conn
	outputQueue chan []byte
	rwlock sync.RWMutex
	isClosed bool
	hub *Hub  //管理所有的连接
	src *svc.ServiceContext
}

func (c *Client) handleMsg(msg []byte) {
	logx.Infof("==========开始处理 [%s] ws消息===============",c.id)
	startTime:=timex.Now()
	defer func(){
		logx.WithDuration(timex.Since(startTime)).Infof("==========完成处理 [%s] ws消息===============",c.id)
	}()

	reqPf:=&textsecure.WebSocketMessage{}
	err:=proto.Unmarshal(msg,reqPf)
	if err!=nil{
		logx.Error("proto.Unmarshal(msg,reqPf) error:",err)


		return
	}
	if reqPf.Type!=textsecure.WebSocketMessage_REQUEST{
		logx.Errorf("the websocket msg (%s) is not  WebSocketMessage_REQUEST:",reqPf.String())
		return
	}

	logx.Info("收到的websocket请求报文是:",reqPf.String())
	path:=reqPf.Request.Path
	var reply *textsecure.WebSocketMessage =nil

	if path=="/v1/keepalive"{
		reply,err= GetKeepAliveResponse(reqPf)
	}else if strings.Contains(path,`/v1/messages/`){
		reply,err= PutMsgHandler(reqPf,c.src,c.id)
	}

	if reply!=nil {
		logx.Infof("回复发件人(%s)消息: %s",c.id,reply.String())
		pf,err:=proto.Marshal(reply)
		if err!=nil{
			logx.Error("reply to from ->proto.Marshal(reply) error:",err)
		}else{
			c.WriteOne(pf)
		}
	}else{
		logx.Info("server error,get websocket relpy is nil",err)
	}
}

type wsConnReq struct {
	Login    string `form:"login,optional"`
	Password string `form:"password,optional"`
}

// serveWs handles websocket requests from the peer.
func WsConnectHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r){
			err := shared.NewCodeError(shared.ERRCODE_WSCONNECTERR, "is not websocket upfgrade")
			httpx.Error(w, err)
		}
		adxName:= r.Header.Get(shared.HEADADXUSERNAME)
		if len(adxName)==0 {
			adxName= util.GenSalt()
		}


		var req wsConnReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error("httpx.Parse /v1/websocket/?login= error",err.Error(),"->匿名建立websocket连接")

		}else{
			adxName = req.Login
		}
		if _,exist:=HasOne(adxName);exist { //已经建立一个连接了，不运行在建立一个
			err := shared.NewCodeError(shared.ERRCODE_WSCONNECTDUP, "one account only has one ws connection")
			httpx.Error(w, err)

		}


		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Info(err)
			err:= shared.NewCodeError(shared.ERRCODE_WSCONNECTERR,err.Error())
			httpx.Error(w, err)

		}else{
			client := &Client{
				id:adxName,
				conn: conn,
				outputQueue: make(chan []byte, maxMessageSize),
				isClosed: false,
				hub: xhub,
				src:ctx,
			}
			client.hub.register <- client // mutex
			go client.writePump()
			go client.readPump()

		}
	}

}




func (c *Client) close(){

	c.rwlock.Lock()
	defer c.rwlock.Unlock()
	if !c.isClosed{
		close(c.outputQueue)
		c.isClosed=true
		c.conn.Close()
		c.hub.unregister <- c
	}

}




// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		// 在读的时候处理客户端离开情况
		c.close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, message, err := c.conn.ReadMessage()
		logx.Info("websocket  收到的消息格式是(1:text,2:binary,8:close,9:ping,10:pong):",messageType)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message //读到的消息放到广播里面, to write redis,测试用
		c.handleMsg(message)

		/*
		msg:=new(textsecure.OutMessage)  //通过websocket来发送请求
		err=proto.Unmarshal(message,msg)
		if err != nil {
			break
		}
		logx.Infof(">>>id  %v >>>>>> request messageContent >>>> %v", c.Id, msg)

		if msg.GetType()==textsecure.OutMessage_CIPHERTEXT{

		}
		*/


	}
}



// 写入一条消息
func (c *Client) WriteOne(msg []byte) (err error) {
	//判断一下是否关闭
	c.rwlock.RLock()
	defer c.rwlock.RUnlock()
	if c.isClosed{
		err = errors.New("write to ws failed by connection is closed")
	}else{
		select {
		case c.outputQueue <- msg: //最多缓存1024条消息 ，如果超出1024条消息还没有发出去，关闭
			//logx.Infof("msg send to %s ok",c.id)
		default:   //发生拥堵时关闭client
			c.close()
			err = errors.New("write to ws failed by chan buffer full ")
		}
	}

	return
}



// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
// To improve efficiency under high load, the writePump function coalesces pending chat messages
// in the send channel to a single WebSocket message.
//This reduces the number of system calls and the amount of data sent over the network.



func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod) //每隔一定时间发ping
	defer func() {
		ticker.Stop()
		c.conn.Close() //这里只是关闭连接，并不是c.close()的调用
	}()

	for {
		select {
		case message, ok := <-c.outputQueue:

			if !ok { //已经关闭了
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			//w, err := c.conn.NextWriter(websocket.TextMessage)
			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				logx.Errorf("ws sender(%s) get NextWriter error:%s ",c.id,err.Error())
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.outputQueue)
			for i := 0; i < n; i++ {
				//w.Write(newline) //如果是二进制，这行请注释
				w.Write(<-c.outputQueue)
			}
			if err := w.Close(); err != nil { //写完关闭w，不是关闭连接
				logx.Errorf("ws sender(%s) close got NextWriter error:%s",c.id,err.Error())
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}





