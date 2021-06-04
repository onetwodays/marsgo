package logic

import (
	"context"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/shared"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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
	if req.Timestamp==0{
		req.Timestamp = time.Now().UnixNano() / 1e6
	}

	for i,_:=range  req.Messages{

		//如果在线，直接通过websocket推送出去
		//todo:send to redis



		msg := &req.Messages[i]
		row:=&model.TMessages{}
		row.Type=int64(msg.Type)
		row.Source= sender
		row.SourceUuid=""
		row.SourceDevice=1
		row.Destination=req.Destination
		row.DestinationDevice=int64(msg.DestinationDeviceId)
		row.Timestamp=req.Timestamp
		row.Message=msg.Body
		row.Content=msg.Content
		row.Relay=msg.Relay
		row.Guid=""
		row.Ctime=time.Now()



		_,err:=l.svcCtx.MsgsModel.Insert(*row)
		if err!=nil{
			return nil, shared.NewCodeError(shared.ERRCODE_SQLINSERT,err.Error())
		}
	}

	return &types.PutMessagesRes{NeedsSync: true}, nil
}
