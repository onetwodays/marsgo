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

type GetParticipantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetParticipantLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetParticipantLogic {
	return GetParticipantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetParticipantLogic) GetParticipant(r *http.Request,req types.GetParicipantParams) (*types.ChannelParticipant, error) {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil, shared.Status(http.StatusUnauthorized,err.Error())
	}
	// 获取频道信息
	channelID:=req.ChannelId
	userID:=req.Id
	channel, err := storage.Channels{}.Get(channelID)
	if err != nil  {
		return nil,shared.Status(http.StatusNotFound, ErrChannelNotFound(channelID).String())
	}

	// 查询成员信息
	users := utils.StringSlice{}.Distinct([]string{account.UUID, userID})
	participants, err := storage.ChannelParticipants{}.GetList(channelID, users)
	if err != nil {
		return nil,shared.Status(http.StatusInternalServerError,err.Error())
	}

	// 校验用户权限
	participant, ok := getParticipant(participants, userID)
	if !ok || participant.Left || participant.Kicked {

		return nil,shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, userID).String())
	}
	if !channel.Public {
		participant, ok = getParticipant(participants, account.UUID)
		if !ok || participant.Left || participant.Kicked {

			return nil,shared.Status(http.StatusForbidden, ErrNotChannelParticipant(channelID, account.UUID).String())
		}
	}

	// 获取用户头像
	targetAccount := account
	if userID != account.UUID {
		targetAccount, err = storage.AccountManager{}.GetByUuid(userID)
		if err != nil {

			return nil,shared.Status(http.StatusInternalServerError,err.Error())
		}
	}

	// 返回成员信息
	result := types.ChannelParticipant{
		UUID:         userID,
		Name:         participant.Name,
		Avatar:       targetAccount.Avatar,
		AvatarDigest: targetAccount.AvatarDigest,
		Banned:       participant.Banned,
	}
	if participant.AdminRights > 0 {
		adminRights := participant.AdminRights.ToEntity()
		result.AdminRights = &adminRights
	}

	return &result, nil
}

func getParticipant(participants []storage.ChannelParticipant, userID string) (storage.ChannelParticipant, bool) {
	for _, participant := range participants {
		if participant.UserID == userID {
			return participant, true
		}
	}
	return storage.ChannelParticipant{}, false
}

