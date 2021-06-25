package logic

import (
	"context"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMsgsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMsgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMsgsLogic {
	return GetMsgsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMsgsLogic) GetMsgs(who string,deviceId int64 ) (*types.GetPendingMsgsRes, error) {
	// todo: add your logic here and delete this line
	resp,err:=l.svcCtx.MsgsModel.FindManyByDst(who,deviceId)
	if err!=nil{
		return nil,shared.NewCodeError(shared.ERRCODE_SQLQUERY,err.Error())
	}

	list:=make([]types.OutcomingMessagex,len(resp))
	for i:=range resp{
		row:=&resp[i]
		item:=types.OutcomingMessagex{}
		item.Relay=row.Relay
		item.Content=row.Content
		item.Message=row.Message
		item.Type=int(row.Type)
		item.Relay=row.Relay
		item.SourceDevice=row.SourceDevice
		item.Source=row.Source
		item.SourceUuid=row.SourceUuid
		item.Cached=false
		item.Guid=row.Guid
		item.ServerTimestamp=row.Timestamp
		item.Timestamp=row.Timestamp
		list[i]=item
	}
	return &types.GetPendingMsgsRes{
		List: list,
		More: false,
	}, nil
}
