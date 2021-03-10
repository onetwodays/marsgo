package logic

import (
	"context"
	"math/rand"
	"secret-im/service/signalserver/cmd/api/util"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/shared"
	"strconv"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/db/redis"

	"github.com/tal-tech/go-zero/core/logx"

)

type GetSmsCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSmsCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSmsCodeLogic {
	return GetSmsCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSmsCodeLogic) containTestNumber(number string)(bool bool,code string){
	bool = false
	code = ""
	testDevices:=l.svcCtx.Config.TestDevices
	for i:=range testDevices{
		if testDevices[i].Number==number{
			bool =true
			code= testDevices[i].Code
			return
		}
	}
	return
}

func (l *GetSmsCodeLogic) verificationCode(number string) (code string){
	ok,code:=l.containTestNumber(number)
	if !ok{
		code = ""
		rand.Seed(time.Now().UnixNano())
		for i:=1;i<7;i++{
			code = code+strconv.Itoa(rand.Intn(9))
		}
	}
	return
}

func (l *GetSmsCodeLogic) GetSmsCode(req types.GetSmsCodeReq) (*types.GetSmsCodeResp, error) {
	// todo: add your logic here and delete this line
	if isVal:=util.IsValidPhoneNumber(req.Number);!isVal{
		return nil,shared.ErrorParam
	}
	code:= l.verificationCode(req.Number)

	redisKey:="pending_account::"+req.Number
	redis.RedisCacheManager().Set(redisKey,code,70*time.Second)
	_,err:=l.svcCtx.PendAccountsModel.Insert(model.TPendAccounts{
		Number: req.Number,
		VerificationCode: code,

	})
	if err!=nil{
		return nil, err
	}
	//  send sms ignore

	return &types.GetSmsCodeResp{
		SmsCode: code,
	}, nil
}
