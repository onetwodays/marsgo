package chat

import (
	"github.com/tal-tech/go-zero/core/logx"
	"secret-im/service/signalserver/cmd/api/internal/types"
)

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.



// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool
	ClientMap map[string]*Client

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	//
	SendSingle chan *types.SingleMessage
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		ClientMap:  make(map[string]*Client),
		SendSingle: make(chan *types.SingleMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case singleMsg:=<-h.SendSingle:
			if client,ok:=h.ClientMap[singleMsg.Id];ok{
				logx.Infof(">> send message %v %v %v ", client.Id, len(singleMsg.Message), singleMsg.Message)
				client.Send<- singleMsg.Message
			}

		case client := <-h.Register:
			h.Clients[client] = true
			h.ClientMap[client.Id]=client
			logx.Infof("websocket client (%s)  Registed",client.Id)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.ClientMap,client.Id)
				close(client.Send)
				logx.Infof("websocket client (%s) left",client.Id)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
					delete(h.ClientMap,client.Id)
				}
			}
		}
	}
}


