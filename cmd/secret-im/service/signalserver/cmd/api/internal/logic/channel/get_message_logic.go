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

type GetMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMessageLogic {
	return GetMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMessageLogic) GetMessage(r *http.Request,req types.GetMessageParams) (*types.OutgoingChannelMessage, error) {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}
	channelID:=req.ChannelId
	if !utils.IsValidUUID(channelID){
		return nil,shared.Status(http.StatusBadRequest,"invalid uuid")
	}
	messageID,err:=strconv.ParseInt(req.Id,10,64)
	if err!=nil{
		return nil,shared.Status(http.StatusBadRequest,err.Error())
	}

	// 校验用户权限
	participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil || participant.Left || participant.Kicked {

		return nil,shared.Status(http.StatusNotFound, ErrNotChannelParticipant(channelID, account.UUID).String())
	}

	// 获取指定消息
	message, err := storage.ChannelMessages{}.Get(channelID, messageID)
	if err != nil {

		return nil,shared.Status(http.StatusInternalServerError,err.Error())
	}
	result := newOutgoingChannelMessage(&message)

	return &result, nil
}
