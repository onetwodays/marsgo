package websocket

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

var Authenticated *SessionManager

func InitWebsocketEnv(ctx *svc.ServiceContext, router http.Handler) {
	Authenticated = NewSessionManager(ctx,
		                              router,
		                              func() SessionHandler {
		                                    return new(AuthenticatedHandler)
		                              })


}

func WsAcceptHandler(w http.ResponseWriter, r *http.Request)  {
	Authenticated.HandleAccept(w,r)
}