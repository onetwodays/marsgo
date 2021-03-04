package chat

/*
import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// Manager 所有 websocket 信息
type WsManager struct {
	Group map[string]map[string]*WsClient
	groupCount, clientCount uint
	Lock sync.Mutex
	Register, UnRegister chan *WsClient
	Message	chan *MessageData
	GroupMessage chan *GroupMessageData
	BroadCastMessage chan *BroadCastMessageData
}

// Client 单个 websocket 信息
type WsClient struct {
	Id, Group string
	Socket *websocket.Conn
	Message   chan []byte
}

// messageData 单个发送数据信息
type MessageData struct {
	Id, Group string
	Message    []byte
}

// groupMessageData 组广播数据信息
type GroupMessageData struct {
	Group string
	Message    []byte
}

// 广播发送数据信息
type BroadCastMessageData struct {
	Message    []byte
}

// 读信息，从 websocket 连接直接读取数据
func (c *WsClient) Read() {
	defer func() {
		WebsocketManager.UnRegister <- c
		log.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		log.Printf("client [%s] receive message: %s", c.Id, string(message))
		c.Message <- message
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		log.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			log.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("client [%s] write message: %s", c.Id, string(message))
			err := c.Socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Printf("client [%s] writemessage err: %s", c.Id, err)
			}
		}
	}
}

// 启动 websocket 管理器
func (manager *Manager) Start() {
	log.Printf("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			log.Printf("client [%s] connect", client.Id)
			log.Printf("register client [%s] to group [%s]", client.Id, client.Group)

			manager.Lock.Lock()
			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*Client)
				manager.groupCount += 1
			}
			manager.Group[client.Group][client.Id] = client
			manager.clientCount += 1
			manager.Lock.Unlock()

		// 注销
		case client := <-manager.UnRegister:
			log.Printf("unregister client [%s] from group [%s]", client.Id, client.Group)
			manager.Lock.Lock()
			if _, ok := manager.Group[client.Group]; ok {
				if _, ok := manager.Group[client.Group][client.Id]; ok {
					close(client.Message)
					delete(manager.Group[client.Group], client.Id)
					manager.clientCount -= 1
					if len(manager.Group[client.Group]) == 0 {
						//log.Printf("delete empty group [%s]", client.Group)
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
				}
			}
			manager.Lock.Unlock()

			// 发送广播数据到某个组的 channel 变量 Send 中
			//case data := <-manager.boardCast:
			//	if groupMap, ok := manager.wsGroup[data.GroupId]; ok {
			//		for _, conn := range groupMap {
			//			conn.Send <- data.Data
			//		}
			//	}
		}
	}
}

// 处理单个 client 发送数据
func (manager *Manager) SendService() {
	for {
		select {
		case data := <-manager.Message:
			if groupMap, ok := manager.Group[data.Group]; ok {
				if conn, ok := groupMap[data.Id]; ok {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for {
		select {
		// 发送广播数据到某个组的 channel 变量 Send 中
		case data := <-manager.GroupMessage:
			if groupMap, ok := manager.Group[data.Group]; ok {
				for _, conn := range groupMap {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 处理广播数据
func (manager *Manager) SendAllService() {
	for {
		select {
		case data := <-manager.BroadCastMessage:
			for _, v := range manager.Group {
				for _, conn := range v {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, message []byte) {
	data := &MessageData{
		Id: id,
		Group: group,
		Message:    message,
	}
	manager.Message <- data
}

// 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group: group,
		Message:    message,
	}
	manager.GroupMessage <- data
}

// 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message:    message,
	}
	manager.BroadCastMessage <- data
}

// 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}
*/
