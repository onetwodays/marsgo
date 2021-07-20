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

type EditProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditProfileLogic {
	return EditProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditProfileLogic) EditProfile(r *http.Request,req types.EditChannelProfileParams) error {
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
	participant, err := storage.ChannelParticipants{}.Get(channelID, account.UUID)
	if err != nil || participant.Left || participant.Kicked {

		return shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
	}
	if participant.AdminRights|storage.ChannelAdminRightChangeInfo == 0 {

		return shared.Status(http.StatusForbidden, ErrNoOperationPermission(channelID, account.UUID).String())
	}

	// 更新频道资料
	err = storage.Channels{}.UpdateProfile(channelID, storage.ChannelProfile{
		Title: req.Title,
		Photo: req.Photo,
		About: req.About,
	})
	if err!=nil{
		logx.Error("[Channel::editProfile] failed to edit channel profile",
			" uuid:", account.UUID,
			" channel:", channelID,
			" reason:",  err)
		return  shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 发送操作信息
	if channel.Profile.Title!=req.Title{
		sendActionMessage(channelID,textsecure.MessageAction{
			Action: textsecure.MessageAction_ChannelEditTitle,
			Title: req.Title,
			Operator: account.UUID,
		})

	}
	if channel.Profile.Photo!=nil{
		if req.Photo==nil{
			sendActionMessage(channelID,textsecure.MessageAction{
				Action: textsecure.MessageAction_ChannelDeletePhoto,
				Operator: account.UUID,
			})
		}else if *channel.Profile.Photo!=*req.Photo{
			sendActionMessage(channelID,textsecure.MessageAction{
				Action: textsecure.MessageAction_ChannelEditPhoto,
				Photo: *req.Photo,
				Operator: account.UUID,

			})
		}
	}else if req.Photo!=nil{
		sendActionMessage(channelID,textsecure.MessageAction{
			Action: textsecure.MessageAction_ChannelEditPhoto,
			Photo: *req.Photo,
			Operator: account.UUID,
		})
	}

	return nil
}
