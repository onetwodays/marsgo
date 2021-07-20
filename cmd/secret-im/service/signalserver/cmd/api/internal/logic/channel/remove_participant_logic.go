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

type RemoveParticipantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveParticipantLogic(ctx context.Context, svcCtx *svc.ServiceContext) RemoveParticipantLogic {
	return RemoveParticipantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveParticipantLogic) RemoveParticipant(r *http.Request,req types.RemoveParticipantParams) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	channelID:=req.ChannelId
	userID:=req.Id
	// 获取频道信息
	channel, err := storage.Channels{}.Get(channelID)
	if err != nil || channel.Deactivated {
		return shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())
	}

	// 是否加入频道
	participant, err := storage.ChannelParticipants{}.Get(channelID, userID)
	if err != nil {
		if !storage.IsNotFoundError(err) {
			return shared.Status(http.StatusInternalServerError,err.Error())
		}
		return shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())

	}
	if participant.Left || participant.Kicked {

		return nil
	}
	if userID == channel.Creator {

		return shared.Status(http.StatusForbidden, ErrNoOperationPermission(channelID, account.UUID).String())
	}

	// 校验用户权限
	operator, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil || operator.Left || operator.Kicked {

		return shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
	}
	if operator.AdminRights|storage.ChannelAdminRightBanUsers == 0 {

		return shared.Status(http.StatusForbidden, ErrNoOperationPermission(channelID, account.UUID).String())
	}

	// 移除频道成员
	err = storage.Channels{}.RemoveParticipant(channelID, userID, true)
	if err != nil {

		logx.Error("[Channel::deleteUser] failed to remove channel participant",
			" uuid:",    account.UUID,
			" channel:", channelID,
			" user:", userID,
			" reason:",  err)

		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 发送操作消息

	sendActionMessage(channelID, textsecure.MessageAction{
		Action: textsecure.MessageAction_ChannelDeleteParticipant, Participants: []string{userID}, Operator: account.UUID})


	return nil
}
