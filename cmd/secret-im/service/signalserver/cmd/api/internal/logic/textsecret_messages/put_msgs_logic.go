package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/api/chat"
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

func (l *PutMsgsLogic) PutMsgs(sender string,req types.PutMessagesReq) (*types.PutMessagesRes, error) {
	// todo: add your logic here and delete this line
	now:= time.Now().UnixNano() / 1e6
	if req.Timestamp==0{
		req.Timestamp = now
	}
	c,_:=chat.HasOne(req.Destination)

	for i,_:=range  req.Messages{

		//如果在线，直接通过websocket推送出去
		//todo:send to redis
		msg := &req.Messages[i]
		row:=&model.TMessages{}
		row.Type=int64(msg.Type)
		row.Source= sender
		row.SourceUuid=""
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
				envelope:=&types.Envelope{}
				envelope.Xtype=types.EnvelopeTypeCiphertext
				envelope.Source=row.Source
				envelope.SourceDevice=1
				envelope.SourceUuid=row.SourceUuid
				envelope.Relay=row.Relay
				envelope.Timestamp=(uint64)(row.Timestamp)
				envelope.LegacyMessage=row.Message
				envelope.Content=row.Content
				envelope.ServerGuid=row.Guid
				envelope.ServerTimestamp=uint64(now)

				pubsubMsg:=&types.PubsubMessage{}
				pubsubMsg.Xtype=types.PubSubTypeDELIVER
				pubsubMsg.Content=*envelope
				msg,err:=json.Marshal(pubsubMsg)
				if err!=nil{
					logx.Infof("json.Marshal(pubsubMsg) error:",err.Error() )
				}else {
					if c!=nil{
						c.WriteOne(msg)
					}else{

						chat.Broadcast(msg)
					}
				}









			}else{
				logx.Infof("destination is not online,not need websocket send,test brocast" )

			}
		}
	}

	return &types.PutMessagesRes{NeedsSync: true}, nil
}
