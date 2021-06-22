package chat

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/textsecure"
)



func GetKeepAliveResponse(req *textsecure.WebSocketMessage) (*textsecure.WebSocketMessage,error) {
	return  newWebSocketMessage(req,http.StatusOK,nil),nil
}