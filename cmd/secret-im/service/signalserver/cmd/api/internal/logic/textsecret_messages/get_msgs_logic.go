package logic

import (
	"context"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"secret-im/service/signalserver/cmd/api/textsecure"
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

func (l *GetMsgsLogic) GetMsgs(who string, deviceId int64) (*types.GetPendingMsgsRes, error) {
	// todo: add your logic here and delete this line
	resp, err := l.svcCtx.MsgsModel.FindManyByDst(who, deviceId)
	if err != nil {
		return nil, shared.NewCodeError(shared.ERRCODE_SQLQUERY, err.Error())
	}

	list := make([]types.OutcomingMessagex, len(resp))
	for i := range resp {


		row := &resp[i]

		envelop :=&textsecure.Envelope{
			Type  :textsecure.GetEnvelopeType(int(row.Type)),
			Source :row.Source,
			SourceUuid:row.SourceUuid,
			SourceDevice:uint32(row.SourceDevice),
			Relay:row.Relay,
			Timestamp:uint64(row.Timestamp),
			//LegacyMessage:   []byte        `protobuf:"bytes,6,opt,name=legacyMessage,proto3" json:"legacyMessage,omitempty"` // Contains an encrypted DataMessage XXX -- Remove after 10/01/15
			//Content         []byte        `protobuf:"bytes,8,opt,name=content,proto3" json:"content,omitempty"`             // Contains an encrypted Content
			ServerGuid:row.Guid,
			ServerTimestamp:uint64(row.Timestamp),
		}
		b,_:=base64.StdEncoding.DecodeString(row.Content)  //解码放在proto里面
		envelop.Content=b
		jsb,_:=proto.Marshal(envelop)
		item := types.OutcomingMessagex{}
		item.Relay = row.Relay
		item.Content = string(jsb)//row.Content //string(b)
		item.Message = row.Message
		item.Type = int(row.Type)
		item.Relay = row.Relay
		item.SourceDevice = row.SourceDevice
		item.Source = row.Source
		item.SourceUuid = row.SourceUuid
		item.Cached = false
		item.Guid = row.Guid
		item.ServerTimestamp = row.Timestamp
		item.Timestamp = row.Timestamp
		list[i] = item
	}
	return &types.GetPendingMsgsRes{
		List: list,
		More: false,
	}, nil
}
