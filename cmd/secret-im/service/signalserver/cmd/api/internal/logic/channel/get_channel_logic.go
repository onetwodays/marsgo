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

type GetChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetChannelLogic {
	return GetChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChannelLogic) GetChannel(r *http.Request,req types.GetChannelParams) (*types.Channel, error) {
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
	if account == nil && !channel.Public {

		return nil,shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, "").String())
	}

	// 获取成员信息
	var participant storage.ChannelParticipant
	if !channel.Public || account != nil {
		participant, err = storage.ChannelParticipants{}.Get(channelID, account.UUID)
		if err != nil {
			if !storage.IsNotFoundError(err) {

				return nil,shared.Status(http.StatusInternalServerError,err.Error())
			}
			return nil,shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
		}

		if !channel.Public && (participant.Left || participant.Kicked) {

			return nil,shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
		}
	}

	// 返回频道信息
	isParticipant := len(participant.UserID) > 0
	channelSlice := []types.Channel{*newChannelEntity(channel, isParticipant, participant.Left, participant.Kicked)}
	if account == nil {
		fillChannelMessageAttributes(channelSlice)
	} else {
		fillChannelMessageAttributes(channelSlice, account.UUID)
	}


	return &channelSlice[0], nil
}
