package signal

import (
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/signal/message"
	"secret-im/service/signalserver/cmd/api/internal/signal/pubsub"
	"secret-im/service/signalserver/cmd/api/internal/signal/pubsub/dispatch"
	"secret-im/service/signalserver/cmd/api/internal/signal/push"
	"secret-im/service/signalserver/cmd/api/internal/signal/websocket"
	"secret-im/service/signalserver/cmd/api/internal/signal/websocket/handlers"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

var SC *SignalContext

type SignalContext struct {
	MessageCacheOper *message.MessagesCache
	RedisOperation   *push.RedisOperation
	Dispatcher       *dispatch.RedisDispatchManager
	PubSubManager    *pubsub.Manager
	PushSender       *push.Sender
	SM               *websocket.SessionManager
	Svc              *svc.ServiceContext
}

func NewSignalContext(ctx *svc.ServiceContext, router http.Handler) *SignalContext {

	if ctx == nil || ctx.RedisClient == nil {
		logx.Error("serviceContent or redis client   not nil")
		return nil
	}

	messageCacheOper, err := message.NewMessagesCache(ctx.RedisClient)
	if err != nil {
		logx.Errorf("NewMessagesCache fail,reason:%s", err.Error())
		return nil
	}

	redisOperation, err := push.NewRedisOperation(ctx.RedisClient)
	if err != nil {
		logx.Errorf("NewRedisOperation fail,reason:%s", err.Error())
		return nil
	}

	dispatcher := dispatch.NewRedisDispatchManager(ctx.RedisClient, 128, new(handlers.DeadLetterHandler))
	pubSubManager := pubsub.NewManager(dispatcher)
	pushSender := push.NewPushSender(pubSubManager, redisOperation)

	authenticated := websocket.NewSessionManager(ctx, router, pushSender, pubSubManager, func() websocket.SessionHandler {
		return new(handlers.AuthenticatedHandler)
	})

	return &SignalContext{
		// redis推送相关
		MessageCacheOper: messageCacheOper,
		RedisOperation:   redisOperation,
		Dispatcher:       dispatcher,
		PubSubManager:    pubSubManager,
		PushSender:       pushSender,
		Svc:              ctx,
		SM:               authenticated,
	}
}
