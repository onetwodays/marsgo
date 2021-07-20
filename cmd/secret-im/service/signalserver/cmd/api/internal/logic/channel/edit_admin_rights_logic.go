package logic

import (
	"context"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type EditAdminRightsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditAdminRightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditAdminRightsLogic {
	return EditAdminRightsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditAdminRightsLogic) EditAdminRights(r *http.Request,req types.EditAdminRightParams) error {
	currAccount,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}
	if !utils.IsValidUUID(req.ChannelId){
		return shared.Status(http.StatusBadRequest,"invaild uuid format")
	}
	if req.UserID==currAccount.UUID{
		return shared.Status(http.StatusForbidden,ErrNoOperationPermission(req.ChannelId,currAccount.UUID).String())
	}
	// 获取频道信息
	channel,err:=storage.Channels{}.Get(req.ChannelId)
	if err!=nil || channel.Deactivated{
		return shared.Status(http.StatusNotFound,ErrChannelNotFound(req.ChannelId).String())
	}
	if req.UserID==channel.Creator{
		return shared.Status(http.StatusForbidden,ErrNoOperationPermission(req.ChannelId,req.UserID).String())
	}
	// 是否频道成员

	participant, err := storage.ChannelParticipants{}.Get(req.ChannelId, req.UserID)
	if err != nil || participant.Kicked || participant.Left {

		return shared.Status(http.StatusNotFound, ErrNotChannelParticipant(req.ChannelId, req.UserID).String())
	}

	// 校验管理权限
	operator, err := storage.ChannelParticipants{}.Get(req.ChannelId, currAccount.UUID)
	if err != nil || operator.Left || operator.Kicked {

		return shared.Status(http.StatusForbidden, ErrNotChannelParticipant(req.ChannelId, currAccount.UUID).String())
	}
	if participant.AdminRights|storage.ChannelAdminRightAddAdmins == 0 {

		return shared.Status(http.StatusForbidden, ErrNoOperationPermission(req.ChannelId, currAccount.UUID).String())
	}

	// 更新管理权限
	adminRights := storage.NewChannelAdminRightsFromEntity(req.AdminRights)
	err = storage.ChannelParticipants{}.UpdateAdminRights(req.ChannelId, req.UserID, adminRights)

	return nil
}
