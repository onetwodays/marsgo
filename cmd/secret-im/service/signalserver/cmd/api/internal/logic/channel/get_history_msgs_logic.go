package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetHistoryMsgsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHistoryMsgsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetHistoryMsgsLogic {
	return GetHistoryMsgsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHistoryMsgsLogic) GetHistoryMsgs(r *http.Request,req types.ChannelMessagesParams) (*types.ChannelMessages, error) {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}
	channelID:=req.ChannelId
	// 获取频道信息
	channel, err := storage.Channels{}.Get(channelID)
	if err != nil {
		return nil,shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())
	}

	// 校验用户身份
	if !channel.Public {
		participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
		if err != nil || participant.Left || participant.Kicked {

			return nil,shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
		}
	}

	// 查询频道消息
	messages, err := storage.ChannelMessages{}.GetMessages(
		channelID, req.Start, req.Start+req.Limit-1)
	if err != nil {
		return nil,shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 返回消息列表
	result := make([]types.OutgoingChannelMessage, 0, len(messages))
	for idx := range messages {
		result = append(result, newOutgoingChannelMessage(&messages[idx]))
	}

	return &types.ChannelMessages{}, nil
}
