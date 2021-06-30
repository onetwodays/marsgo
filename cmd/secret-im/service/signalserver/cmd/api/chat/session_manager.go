package chat

import (
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/shared"

	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"sync"
	"sync/atomic"
)

type Key struct {
	UUID     string
	DeviceID int64
}

// Session管理器
type SessionManager struct {
	cleaning     bool
	sequence     int64
	router       http.Handler
	shutdownChan chan struct{}

	mutex         sync.Mutex
	deviceMapper  map[Key]int64
	sessionMapper map[int64]*Session

	//pushSender         *push.Sender
	//pubSubManager      *pubsub.Manager
	makeSessionHandler MakeSessionHandler //是个函数
	ctx                *svc.ServiceContext
}

// 创建Session管理器
func NewSessionManager(ctx *svc.ServiceContext, router http.Handler,
/*pushSender *push.Sender, pubSubManager *pubsub.Manager,*/
	makeSessionHandler MakeSessionHandler) *SessionManager {
	return &SessionManager{
		router: router,
		//pushSender:         pushSender,
		//pubSubManager:      pubSubManager,
		makeSessionHandler: makeSessionHandler,
		deviceMapper:       make(map[Key]int64),
		sessionMapper:      make(map[int64]*Session),
		shutdownChan:       make(chan struct{}),
		ctx:                ctx,
	}
}

// 接受连接
func (manager *SessionManager) HandleAccept(w http.ResponseWriter, r *http.Request) {
	if manager.cleaning {
		httpx.Error(w, shared.Status(http.StatusServiceUnavailable, ""))
		return
	}

	account, pass := authenticate(r, manager.ctx)
	if !pass {
		httpx.Error(w, shared.Status(http.StatusUnauthorized, ""))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Errorf("[Websocket] upgrade failed,reason:%s", err)
		httpx.Error(w, shared.Status(http.StatusInternalServerError, err.Error()))
		return
	}

	session := manager.onCreated(conn, account)
	logx.Infof("[Websocket] connection accepted,id:%d", session.id)

	go session.handleRead()
	go session.handleWrite()
	//这是一个回调函数，
	session.handler.OnWebSocketConnect(session.context)
}

// 关闭服务
func (manager *SessionManager) Shutdown() {
	manager.cleaning = true
	sessions := make([]*Session, 0)

	manager.mutex.Lock()
	for _, session := range manager.sessionMapper {
		sessions = append(sessions, session)
	}
	manager.mutex.Unlock()

	for _, session := range sessions {
		session.Close(1000, "OK")
	}
	if len(sessions) == 0 {
		close(manager.shutdownChan)
	}
	<-manager.shutdownChan
}

// 获取Session
func (manager *SessionManager) GetByDevice(uuid string, deviceID int64) (*Session, bool) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	key := Key{
		UUID:     uuid,
		DeviceID: deviceID,
	}
	id, ok := manager.deviceMapper[key]
	if !ok {
		return nil, false
	}

	session, ok := manager.sessionMapper[id]
	if !ok {
		return nil, false
	}
	return session, true
}

// 获取下一ID
func (manager *SessionManager) nextID() int64 {
	return atomic.AddInt64(&manager.sequence, 1)
}

// 关闭Session时的回调函数
func (manager *SessionManager) onClosed(id int64, code int, text string) error {
	manager.mutex.Lock()
	session, ok := manager.sessionMapper[id]
	if !ok {
		manager.mutex.Unlock()
		return nil
	}
	delete(manager.sessionMapper, id)
	count := len(manager.sessionMapper)
	finished := manager.cleaning && count == 0

	device := session.context.Device
	if device != nil {
		key := Key{
			UUID:     device.UUID,
			DeviceID: device.Device.ID,
		}
		delete(manager.deviceMapper, key)
	}
	manager.mutex.Unlock()

	if finished {
		close(manager.shutdownChan)
	}

	if ok && session != nil {
		logx.Infof("[Websocket] connection closed[id:%d,count:%d,code:%d,text:%s]", id, count, code, text)
	}
	return nil
}

// 创建Session
func (manager *SessionManager) onCreated(conn *websocket.Conn, account *entities.Account) *Session {
	id := manager.nextID()
	var device *ConnectedDevice
	if account != nil {
		device = &ConnectedDevice{
			Number: account.Number,
			UUID:   account.UUID,
			Device: account.AuthenticatedDevice.Device,
		}
	}

	session := newSession(Options{
		id:     id,
		router: manager.router,
		conn:   conn,
		device: device,
		//pushSender:    manager.pushSender,
		//pubSubManager: manager.pubSubManager,
		handler:      manager.makeSessionHandler(),
		closeHandler: manager.onClosed,
	})
	manager.mutex.Lock()
	manager.sessionMapper[session.id] = session
	if device != nil {
		key := Key{
			UUID:     device.UUID,
			DeviceID: device.Device.ID,
		}
		manager.deviceMapper[key] = session.id
	}
	manager.mutex.Unlock()
	return session
}
