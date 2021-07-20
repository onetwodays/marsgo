package logic

import (
	"context"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/textsecure"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddParticipantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddParticipantLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddParticipantLogic {
	return AddParticipantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddParticipantLogic) AddParticipant(r *http.Request,req types.AddParticipantParams) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	channelID:=req.ChannelId
	if !utils.IsValidUUID(channelID){
		return shared.Status(http.StatusBadRequest,"invalid uuid")
	}

	// 获取用户信息
	_, err = storage.AccountManager{}.GetByUuid(req.UUID)
	if err != nil {
		if !storage.IsNotFoundError(err) {
			return shared.Status(http.StatusInternalServerError,err.Error())
		} else {
			return shared.Status(http.StatusNotFound, ErrAccountNotFound(req.UUID).String())
		}
		return nil
	}

	// 获取频道信息
	channel, err := storage.Channels{}.Get(channelID)
	if err != nil || channel.Deactivated {

		return shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())
	}

	// 是否频道成员
	participant, err := storage.ChannelParticipants{}.Get(channelID, req.UUID)
	if err != nil && !storage.IsNotFoundError(err) {

		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	if err == nil && !participant.Kicked && !participant.Left {

		return nil
	}

	// 校验用户权限
	if !channel.Public {
		operator, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
		if err != nil || operator.Left || operator.Kicked {

			return shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
		}
		if operator.AdminRights|storage.ChannelAdminRightInviteUsers == 0 {

			return shared.Status( http.StatusForbidden, ErrNoOperationPermission(channelID, account.UUID).String())
		}
	}

	// 添加频道成员
	nameMapper := map[string]string{
		req.UUID: req.Name,
	}
	err = storage.Channels{}.AddParticipants(channelID, nameMapper)
	if err != nil {
		logx.Error("[Channel::addUser] failed to add channel participant",
			" uuid:",    account.UUID,
			" channel:", channelID,
			" reason:",  err)

		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 发送操作消息

	sendActionMessage(channelID, textsecure.MessageAction{
		Action: textsecure.MessageAction_ChannelAddParticipant, Participants: []string{req.UUID}, Operator: account.UUID})


	return nil
}
