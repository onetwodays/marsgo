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

type SetUserNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetUserNameLogic {
	return SetUserNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUserNameLogic) SetUserName(r *http.Request,req types.SetUserNameParams) error {
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
	// 判断用户权限
	participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil || participant.Kicked || participant.Left {

		return shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())
	}

	// 更新显示名称
	err = storage.ChannelParticipants{}.UpdateName(channelID, account.UUID, req.Name)
	if err != nil {
		return shared.Status(http.StatusForbidden, ErrNoOperationPermission(channelID, account.UUID).String())
	}

	return nil
}
