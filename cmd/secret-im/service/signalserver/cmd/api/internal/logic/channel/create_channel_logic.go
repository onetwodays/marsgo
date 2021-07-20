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
	"secret-im/service/signalserver/cmd/api/textsecure"
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

	//群主
	currAccount,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}

	//群里所有的成员的帐号信息入参
	addInputParticipants,accountMapper,uuids,err:=l.getGroupUseAccounts(currAccount,req)
	if err!=nil{
		return nil,err
	}

	//要保村到数据库的群成员信息
	participants,channelID:=l.createChannelParticipants(currAccount,addInputParticipants,accountMapper)


	// 要插入频道信息
	channel:=storage.Channel{
		ChannelID: channelID,
		Creator: currAccount.UUID,
		Profile: storage.ChannelProfile{
			Title: req.Title,
		},
		Public: req.Public,
		Date: time.Now().Unix(),
	}

	// 群和群成员入库
	if err:=new(storage.Channels).Insert(&channel,participants);err!=nil{
		logx.Error("[Channel::createChannel] failed to create channel",
			" uuid:",  currAccount.UUID,
			" reason:", err,)
		return nil, shared.Status(http.StatusInternalServerError,err.Error())

	}

	//发送操作消息

	uuids = utils.StringSlice{}.Distinct(uuids)

	sendActionMessage(channelID, textsecure.MessageAction{
		Action: textsecure.MessageAction_ChannelCreate, Title: req.Title, Participants: uuids})

	return newChannelEntity(channel,true,false,false),nil
}



func (l *CreateChannelLogic) getGroupUseAccounts(currAccount *entities.Account,req types.ChannelCreationInfo) ([]types.ChannelInputParticipant,map[string]entities.Account,[]string,error)  {

	// 获取用户列表
	var ids []string
	var numbers []string
	// 群主加入群组
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
	// 得到这些人的帐号
	accountMapper := make(map[string]entities.Account)

	if len(ids) > 0 {
		accounts, err := storage.AccountManager{}.GetByUUIDs(ids)
		if err != nil {
			logx.Error("[Channel::createChannel] failed to get accounts:",
				" uuid:",   currAccount.UUID,
				" numbers:", numbers,
				" reason:", err)
			return nil,nil,nil,shared.Status(http.StatusInternalServerError,err.Error())
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
			return nil,nil,nil,shared.Status(http.StatusInternalServerError,err.Error())
		}
		for key, val := range accounts {
			ids=append(ids,key)
			accountMapper[key] = val
		}
	}
	return addParticipants,accountMapper,ids,nil


}

func (l *CreateChannelLogic) createChannelParticipants(currAccount *entities.Account,addParticipants []types.ChannelInputParticipant,accountMapper map[string]entities.Account) ([]storage.ChannelParticipant,string)  {

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
	return participants,channelID

}


