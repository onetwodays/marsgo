package chat

import (
	"github.com/tal-tech/go-zero/core/logx"
)


// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.

	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:  make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client.id]=client
			logx.Infof("websocket client (%s)  Registed",client.id)
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients,client.id)
				logx.Infof("websocket client (%s) left",client.id)
			}
		case message := <-h.broadcast:
			for _,c := range h.clients {
				err:=c.WriteOne(message)
				if err!=nil{
					logx.Errorf("send to %s broadcast failed by  ",c.id,err)
				}
			}
		}
	}
}


