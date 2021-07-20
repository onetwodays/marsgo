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

type DeleteChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteChannelLogic {
	return DeleteChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteChannelLogic) DeleteChannel(r *http.Request,req types.DeleteChannelParams) error {
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

	// 校验用户权限
	if channel.Creator != account.UUID {
		return shared.Status(http.StatusForbidden, ErrNoOperationPermission(channelID, account.UUID).String())
	}

	// 停用频道
	err = storage.Channels{}.Deactivate(channelID)
	if err != nil {

		return shared.Status(http.StatusInternalServerError,err.Error())
	}


	return nil
}
