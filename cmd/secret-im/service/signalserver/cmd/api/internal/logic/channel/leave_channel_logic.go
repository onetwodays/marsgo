package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/textsecure"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type LeaveChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLeaveChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) LeaveChannelLogic {
	return LeaveChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LeaveChannelLogic) LeaveChannel(r *http.Request,req types.LeaveChannelParams) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	channelID:=req.ChannelId
	// 获取频道信息
	channel, err := storage.Channels{}.Get(channelID)
	if err != nil || channel.Deactivated {
		return shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())
	}
	if account.UUID==channel.Creator{
		return shared.Status(http.StatusForbidden,ErrNoOperationPermission(channelID,account.UUID).String())
	}
	// 校验用户权限
	participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil || participant.Kicked || participant.Left {

		return shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
	}

	// 移除频道成员
	err = storage.Channels{}.RemoveParticipant(channelID, account.UUID, false)
	if err != nil {

		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 发送操作消息

	sendActionMessage(channelID, textsecure.MessageAction{
		Action:textsecure.MessageAction_ChannelDeleteParticipant, Participants: []string{account.UUID}})

	return nil
}
