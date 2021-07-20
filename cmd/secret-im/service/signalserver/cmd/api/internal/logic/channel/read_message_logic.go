package logic

import (
	"context"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"strconv"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ReadMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) ReadMessageLogic {
	return ReadMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadMessageLogic) ReadMessage(r *http.Request,req types.ReadMessageParams) error {
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

	// 标记消息确认
	err = storage.ChannelMessageAckDao{}.UpdateLastAckMessage(account.UUID, channelID, messageID)
	if err != nil {

		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	return nil
}
