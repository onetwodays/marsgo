package logic

import (
	"context"
	"encoding/json"
	"errors"
	"secret-im/service/signalserver/cmd/api/db/redis"
	"secret-im/service/signalserver/cmd/api/util"
	"secret-im/service/signalserver/cmd/model"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ConfirmCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) ConfirmCodeLogic {
	return ConfirmCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmCodeLogic) ConfirmCode(req types.ConfirmVerificationCodeReq) (*types.ConfirmVerificationCodeRes, error) {
	// todo: add your logic here and delete this line

	redisKey:="pending_account::"+req.PhoneNumber
	code, err := redis.RedisDirectoryManager().Get(redisKey).Result()
	if err!=nil{
		return nil, errors.New("redis data not found")
	}
	if req.VerificationCode!=code{
		return nil, errors.New("verification code check fail")
	}
    salt:=util.GenSalt()
	device:= types.Device{
		Id:1,
		Salt:salt,
		AccountAttributes: types.AccountAttributes{
			SignalingKey: req.AccountAttributes.SignalingKey,
			FetchesMessages: req.AccountAttributes.FetchesMessages,
			RegistrationId: req.AccountAttributes.RegistrationId,
			Name: req.AccountAttributes.Name,
			Voice: req.AccountAttributes.Voice,
			Video: req.AccountAttributes.Video,
		},
		AuthToken: util.GenAuthKey(salt,req.Password),

	}

	account:=types.Account{}
	account.Number = req.PhoneNumber
	account.Devices=[]types.Device{device}
	accountString,_:=json.Marshal(account)
	//把旧的删除
	l.svcCtx.AccountsModel.DeleteByNumber(req.PhoneNumber)

	//新添加的
	_,err=l.svcCtx.AccountsModel.Insert(model.TAccounts{
		Number: req.PhoneNumber,
		Data: string(accountString),
	})
	if err!=nil{
		return nil, err
	}
	l.svcCtx.PendAccountsModel.DeleteByNumber(req.PhoneNumber)

	return &types.ConfirmVerificationCodeRes{
		Account: string(accountString),
	}, nil
}
