package chat

import (
	"context"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mapping"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"strings"

	logic "secret-im/service/signalserver/cmd/api/internal/logic/textsecret_messages"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/textsecure"
)

func PutMsgHandler(req *textsecure.WebSocketMessage,svc *svc.ServiceContext,sender string) (*textsecure.WebSocketMessage,error){
	var putMesReq types.PutMessagesReq
	// Unmarshal to json object
	err:=mapping.UnmarshalJsonBytes(req.Request.Body,&putMesReq)
	if err!=nil {
		logx.Error("mapping.UnmarshalJsonBytes(reqPf.Request.Body) error:", err)
		return nil, err
	}
	logx.Info("发送方body=",string(req.Request.Body))

	logx.Infof("打印 发送方的消息=:%#v", putMesReq)
	if len(putMesReq.Destination)==0{
		strs:=strings.Split(req.Request.Path,`/`)
		if len(strs)==3{
			putMesReq.Destination=strs[len(strs)-1]
		}
	}


	// 交给logic处理
	l := logic.NewPutMsgsLogic(context.Background(), svc)
	recv, isOk := HasOne(putMesReq.Destination)
	putMsgRes, err := l.PutMsgs(sender, putMesReq,req.Request.Id,isOk)
	if err != nil {
		return nil, err
	}

	if isOk {
		for i := range putMsgRes.DestContent {
			recv.WriteOne(putMsgRes.DestContent[i])
			logx.Infof("msg from %s send to %s ok ", sender, putMesReq.Destination)
		}
	}
	body:=&entities.SendMessageResponse{
		NeedsSync: putMsgRes.NeedsSync,
	}
	return newWebSocketMessage(req,http.StatusOK,body),nil


}