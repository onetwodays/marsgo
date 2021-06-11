package chat

import (
	"github.com/tal-tech/go-zero/core/logx"
	"sync"
)

var (
	// Registered clients.
	rwlock sync.RWMutex
	clients map[string]*Client = make(map[string]*Client)
)


// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {


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

	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			comeOne(c)
		case c := <-h.unregister:
			leaveOne(c)
		case message := <-h.broadcast:
			broadcast(message)
		}
	}
}

func HasOne(cname string) (*Client,bool) {
	rwlock.RLock()
	defer func() {
		rwlock.RUnlock()
	}()
	c,ok:=clients[cname]
	return c,ok
}


func comeOne(c *Client){
	rwlock.Lock()
	defer func() {
		rwlock.Unlock()
	}()
	clients[c.id]=c
	logx.Infof("websocket client (%s)  Registed",c.id)
}

func leaveOne(c *Client){
	rwlock.Lock()
	defer func() {
		rwlock.Unlock()
	}()
	if _, ok := clients[c.id]; ok {
		delete(clients,c.id)
		logx.Infof("websocket client (%s) left",c.id)
	}
}

func broadcast(msg []byte){
	rwlock.Lock()
	defer func() {
		rwlock.Unlock()
	}()
	for _,c := range clients {
		err:=c.WriteOne(msg)
		if err!=nil{
			logx.Errorf("send to %s broadcast failed by %s  ",c.id,err.Error())
		}
	}
}

func broadcastExclude(msg []byte,excludeClient string){
	rwlock.Lock()
	defer func() {
		rwlock.Unlock()
	}()
	for _,c := range clients {
		if c.id!=excludeClient{
			err:=c.WriteOne(msg)
			if err!=nil{
				logx.Errorf("send to %s broadcast failed by %s ",c.id,err.Error())
			}
		}

	}
}






