package logic

import (
	"context"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/textsecure"
	"strconv"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteMessageLogic {
	return DeleteMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessageLogic) DeleteMessage(r *http.Request,req types.EditMessageParams) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	channelID:=req.ChannelId
	if !utils.IsValidUUID(channelID){
		return shared.Status(http.StatusBadRequest,"invalid uuid")
	}
	messageID,err:=strconv.ParseInt(req.Id,10,64)
	if err!=nil{
		return shared.Status(http.StatusBadRequest,err.Error())
	}

	// 校验用户权限
	participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil || participant.Left || participant.Kicked {

		return shared.Status(http.StatusNotFound, ErrNotChannelParticipant(channelID, account.UUID).String())
	}

	if participant.AdminRights|storage.ChannelAdminRightDeleteMessages == 0 {

		return shared.Status(http.StatusForbidden, ErrNoOperationPermission(channelID, account.UUID).String())
	}

	// 校验消息类型
	message, err := storage.ChannelMessages{}.Get(channelID, messageID)
	if err != nil {

		return shared.Status(http.StatusInternalServerError,err.Error())
	}
	if message.Type != model.ChannelMessageTypeNormal {

		return shared.Status(http.StatusBadRequest, ErrCannotModifyServiceMessage(channelID, messageID).String())
	}

	// 删除指定消息
	err = storage.ChannelMessages{}.Delete(channelID, messageID)
	if err != nil {
		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 发送操作消息

	sendActionMessage(channelID, textsecure.MessageAction{
		Action: textsecure.MessageAction_ChannelDeleteMessage, MessageId: messageID, Operator: account.UUID})


	return nil
}
