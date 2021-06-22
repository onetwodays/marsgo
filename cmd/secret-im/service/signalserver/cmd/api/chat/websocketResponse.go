package chat

import (
	"encoding/json"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

func newWebSocketMessage(req *textsecure.WebSocketMessage,code uint32,body interface{})  *textsecure.WebSocketMessage {
	webres:=&textsecure.WebSocketResponseMessage{}
	if body!=nil{
		bytes,err:= json.Marshal(body)
		if err==nil{
			webres.Body=bytes
		}else{
			logx.Errorf("newWebSocketResponseMessage->json.Marshal(body) error:%s,%#v",err.Error(),body)
			return nil
		}
	}
	webres.Id=req.Request.Id
	webres.Headers=[]string{"Content-Type:application/json"}
	webres.Status=code
	webres.Message=http.StatusText(int(webres.Status))

	wsMsg:=&textsecure.WebSocketMessage{}
	wsMsg.Type=textsecure.WebSocketMessage_RESPONSE
	wsMsg.Response=webres


	return wsMsg
}