// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"bytes"
	"errors"
	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/util"
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



// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id string
	conn *websocket.Conn
	outputQueue chan []byte
	rwlock sync.RWMutex
	isClosed bool
	hub *Hub  //管理所有的连接

}

// serveWs handles websocket requests from the peer.
func WsConnectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Info(err)
		return
	}
	client := &Client{
		hub: xhub,
		conn: conn,
		outputQueue: make(chan []byte, maxMessageSize),
		id:util.GenSalt(),
		isClosed: false,

	}
	client.hub.register <- client // mutex
	go client.writePump()
	go client.readPump()
}

func AdxWsConnectHandler(adxName string,w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Info(err)
			httpx.Error(w,err) // tell client  an error of type HandshakeError
		}else{
			logx.Info("receive handshake from ",conn.RemoteAddr().String())
			client := &Client{
				hub: xhub,
				conn: conn,
				outputQueue: make(chan []byte, maxMessageSize),
				id:adxName,
				isClosed: false,
			}
			client.hub.register <- client
			go client.writePump()
			go client.readPump()
		}

}

func (c *Client) close(){
	close(c.outputQueue)
	c.rwlock.Lock()
	defer c.rwlock.Unlock()
	c.isClosed=true
	c.conn.Close()
	c.hub.unregister <- c
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
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message //读到的消息放到广播里面, to write redis

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
			logx.Infof("msg(%s) send to %s ok",string(msg),c.id)
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
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.outputQueue:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.outputQueue)
			for i := 0; i < n; i++ {
				w.Write(newline) //如果是二进制，这行请注释
				w.Write(<-c.outputQueue)
			}
			if err := w.Close(); err != nil { //写完关闭w，不是关闭连接
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





