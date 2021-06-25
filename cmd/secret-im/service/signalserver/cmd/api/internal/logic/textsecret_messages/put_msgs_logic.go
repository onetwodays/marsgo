package logic

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/shared"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/utils"
)

type PutMsgsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutMsgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutMsgsLogic {
	return PutMsgsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutMsgsLogic) PutMsgs(sender string,req types.PutMessagesReq,msgId uint64) (*types.PutMessagesRes, error) {
	// todo: add your logic here and delete this line


	account,err:=l.svcCtx.AccountsModel.FindOneByNumber(sender)
	if err!=nil{
		return nil, err
	}

	now:= time.Now().UnixNano() / 1e6
	if req.Timestamp==0{
		req.Timestamp = now
	}
	var destContent [][]byte


	for i,_:=range  req.Messages{

		//如果在线，直接通过websocket推送出去
		//todo:send to redis
		msg := &req.Messages[i]
		row:=&model.TMessages{}
		row.Type=int64(msg.Type)
		row.Source= sender
		row.SourceUuid=account.Uuid
		row.SourceDevice=1
		row.Destination=msg.Destination
		row.DestinationDevice=int64(msg.DestinationDeviceId)
		row.Timestamp=req.Timestamp
		row.Message=msg.Body
		row.Content=msg.Content
		row.Relay=msg.Relay
		row.Guid=utils.NewUuid() //消息的全局uuid
		row.Ctime=time.Now()

		_,err:=l.svcCtx.MsgsModel.Insert(*row)
		if err!=nil{
			return nil, shared.NewCodeError(shared.ERRCODE_SQLINSERT,err.Error())
		}else{
			//if isOnline {
			if true {

				content,_ :=base64.StdEncoding.DecodeString(msg.Content)


				envelopePf:=&textsecure.Envelope{}
				envelopePf.Type=textsecure.GetEnvelopeType(msg.Type)
				envelopePf.SourceDevice=1
				envelopePf.Source=sender
				envelopePf.ServerGuid=row.Guid
				envelopePf.SourceUuid=account.Uuid
				envelopePf.ServerTimestamp=uint64(now)
				envelopePf.Timestamp=uint64(req.Timestamp)
				envelopePf.Relay=row.Relay
				//envelopePf.LegacyMessage=[]byte(row.Message)
				envelopePf.Content=content
				logx.Info("收件人的envelop:",envelopePf.String())
				contentPf,err:=proto.Marshal(envelopePf)
				if err!=nil{
					logx.Error("proto.Marshal(envelopePf):",err)
				}else{
					websocketReq:=&textsecure.WebSocketRequestMessage{}
					websocketReq.Id=msgId
					websocketReq.Headers=[]string{"X-Signal-Key: false","X-Signal-Timestamp:"+fmt.Sprintf("%d",now)}
					websocketReq.Path="/api/v1/message"
					websocketReq.Verb="PUT"
					websocketReq.Body=contentPf

					websocketMsg:=&textsecure.WebSocketMessage{}

					websocketMsg.Type=textsecure.WebSocketMessage_REQUEST
					websocketMsg.Request=websocketReq
					logx.Info("收件人最外层:",websocketMsg.String())
					msg,err:=proto.Marshal(websocketMsg)

					//pubsubMsg:=&textsecure.PubSubMessage{}
					//pubsubMsg.Type=textsecure.PubSubMessage_DELIVER
					//pubsubMsg.Content=contentPf
					//logx.Info("收件人最外层:",pubsubMsg.String())
					//msg,err:=proto.Marshal(pubsubMsg)


					if err!=nil{
						logx.Infof("proto.Marshal(websocketMsg) error:",err.Error() )
					}else{
						destContent=append(destContent,msg)

					}
				}
			}else{
				logx.Infof("destination is not online,not need websocket send,test brocast" )
			}
		}
	}

	return &types.PutMessagesRes{NeedsSync: false,DestContent: destContent}, nil
}
