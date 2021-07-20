package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) SendMessageLogic {
	return SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(r *http.Request,req types.SendMessageParams) (*types.ChannelMessageID, error) {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}
	channelID:=req.ChannelId
	// 获取频道信息
	channel, err := storage.Channels{}.Get(channelID)
	if err != nil || channel.Deactivated {
		return nil,shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())
	}
	/*
	if !channel.Public{
		return nil,shared.Status(http.StatusForbidden,ErrPrivateChannel(channelID).String())
	}*/
	// 校验用户权限
	participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil || participant.Left || participant.Kicked {

		return nil,shared.Status(http.StatusNotFound, ErrNotChannelParticipant(channelID, account.UUID).String())
	}
	/*
	if !c.rateLimiters.MessagesLimiter.Validate(account.Number, 1) {
		respond.Status(w, http.StatusRequestEntityTooLarge)
		return
	}

	 */

	if participant.Banned {

		return nil ,shared.Status(http.StatusForbidden, ErrUserIsBanned(account.UUID).String())
	}

	// 保存消息记录
	message := storage.ChannelMessage{
		ChannelID:       channelID,
		Type:            model.ChannelMessageTypeNormal,
		Source:          &account.UUID,
		SourceDevice:    &account.AuthenticatedDevice.ID,
		Content:         []byte(req.Content),
		Relay:           req.Relay,
		Timestamp:       req.Timestamp,
		ServerTimestamp: time.Now().Unix(),
	}
	err = storage.ChannelMessages{}.Insert(&message)
	if err != nil {

		return nil,shared.Status(http.StatusInternalServerError,err.Error())
	}
	storage.ChannelMessagesManager{}.SetLatestMessage(channelID, &message) //保存在redis

	// 发送消息到设备
	sendChannelMessageToDevices(&message)


	return &types.ChannelMessageID{
		ChannelID: channelID,
		MessageID: message.MessageID,

	}, nil
}


