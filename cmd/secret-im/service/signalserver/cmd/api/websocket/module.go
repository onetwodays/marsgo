package websocket

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/svc"
)

var authenticated *SessionManager

func InitWebsocketEnv(ctx *svc.ServiceContext, router http.Handler) {
	authenticated = NewSessionManager(ctx,
		                              router,
		                              func() SessionHandler {
		                                    return new(AuthenticatedHandler)
		                              })


}

func WsAcceptHandler(w http.ResponseWriter, r *http.Request)  {
	authenticated.HandleAccept(w,r)
}