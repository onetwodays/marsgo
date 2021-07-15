package logic

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type CreateChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateChannelLogic {
	return CreateChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateChannelLogic) CreateChannel(r *http.Request,req types.ChannelCreationInfo) (*types.Channel, error) {
	currAccount,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}

	// 获取用户列表
	var ids []string
	var numbers []string
	addParticipants := []types.ChannelInputParticipant{{
		UUID: currAccount.UUID,
		Name: "超级管理员",
	}}
	addParticipants = append(addParticipants, req.Participants...)

	for _, participant := range addParticipants {
		identifier := auth.NewAmbiguousIdentifier(participant.UUID) //兼容number和uuid
		if  len(identifier.UUID)!= 0 {
			ids = append(ids, identifier.UUID)
		} else {
			numbers = append(numbers, identifier.Number)
		}
	}
	accountMapper := make(map[string]entities.Account)

	if len(ids) > 0 {
		accounts, err := storage.AccountManager{}.GetByUUIDs(ids)
		if err != nil {
			logx.Error("[Channel::createChannel] failed to get accounts:",
				" uuid:",   currAccount.UUID,
				" numbers:", numbers,
				" reason:", err)
			return nil,shared.Status(http.StatusInternalServerError,err.Error())
		}
		for key, val := range accounts {
			accountMapper[key] = val
		}
	}

	if len(numbers) > 0 {
		accounts, err := storage.AccountManager{}.GetByNumbers(numbers)
		if err != nil {
			logx.Error("[Channel::createChannel] failed to get accounts:",
				" uuid:",   currAccount.UUID,
				" ids:",    ids,
				" reason:", err)
			return nil,shared.Status(http.StatusInternalServerError,err.Error())
		}
		for key, val := range accounts {
			ids=append(ids,key)
			accountMapper[key] = val
		}
	}


	// 生成用户列表
	set := make(map[string]struct{})
	channelID := uuid.NewV4().String()
	participants := make([]storage.ChannelParticipant, 0, len(accountMapper))
	for _, participant := range addParticipants {
		var ok bool
		var account entities.Account
		identifier := auth.NewAmbiguousIdentifier(participant.UUID)
		if len(identifier.UUID) != 0 {
			account, ok = accountMapper[identifier.UUID]
		} else {
			account, ok = accountMapper[identifier.Number]
		}
		if !ok {
			continue
		}

		if _, ok = set[currAccount.UUID]; ok {
			continue
		}
		set[currAccount.UUID] = struct{}{}

		participant := storage.ChannelParticipant{
			ParticipantID: storage.ParticipantID{
				UserID:    account.UUID,
				ChannelID: channelID,
			},
			Name: participant.Name,
			Date: time.Now().Unix(),
		}
		if account.UUID == currAccount.UUID {
			participant.AdminRights = storage.DefaultChannelAdminRights()
		}
		participants = append(participants, participant)
	}

	// 插入频道信息

	channel:=storage.Channel{
		ChannelID: channelID,
		Creator: currAccount.UUID,
		Profile: storage.ChannelProfile{
			Title: req.Title,
		},
		Public: req.Public,
		Date: time.Now().Unix(),
	}

	if err:=new(storage.Channels).Insert(&channel,participants);err!=nil{
		logx.Error("[Channel::createChannel] failed to create channel",
			" uuid:",  currAccount.UUID,
			" reason:", err,)
		return nil, shared.Status(http.StatusInternalServerError,err.Error())

	}

	//发送操作消息

	ids = utils.StringSlice{}.Distinct(ids)
	/*
	sendActionMessage(channelID, textsecure.MessageAction{
		Action: textsecure.MessageAction_ChannelCreate, Title: req.Title, Participants: ids})

	 */

	return &types.Channel{

	}, nil
}
