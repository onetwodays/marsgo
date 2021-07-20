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

type JoinChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) JoinChannelLogic {
	return JoinChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinChannelLogic) JoinChannel(r *http.Request,req types.JoinChannelParams) error {
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
	if !channel.Public{
		return shared.Status(http.StatusForbidden,ErrPrivateChannel(channelID).String())
	}

	// 是否频道成员
	participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil && !storage.IsNotFoundError(err) {
		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	if err == nil && !participant.Kicked && !participant.Left {
		return nil
	}

	// 添加频道成员
	nameMapper := map[string]string{
		account.UUID: "aaaaaa",
	}
	err = storage.Channels{}.AddParticipants(channelID, nameMapper)
	if err != nil {

		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 发送操作消息

	sendActionMessage(channelID, textsecure.MessageAction{
		Action: textsecure.MessageAction_ChannelAddParticipant, Participants: []string{account.UUID}})

	return nil
}
