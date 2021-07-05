package logic

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/model"
	shared "secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/util"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
	pkgUtils "secret-im/pkg/utils-tools"
)

type GetCodeReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCodeReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCodeReqLogic {
	return GetCodeReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCodeReqLogic) GetCodeReq(req types.GetCodeReq,locale string) error {
	// todo: add your logic here and delete this line
	logx.Infof("input:[Accept-Language:%s,%+v]",locale,req)

	//校验手机号的有效性
	if !util.IsValidNumber(req.Number){
		reason:=fmt.Sprintf("%s is not valid number",req.Number)
		logx.Error(reason)
		return shared.Status(http.StatusBadRequest,reason)
	}
	//检查调用频率
	switch req.Transport {
	case "sms":
		//return shared.Status(http.StatusRequestEntityTooLarge,"Frequently called")
	case "voice":
		return shared.Status(http.StatusNotImplemented,"voice not implemented,sorry")
	default:
		reason:=fmt.Sprintf("%s unprocess",req.Transport)
		return shared.Status(http.StatusUnprocessableEntity,reason)
	}

	//保存帐号信息
	verificationCode := generateVerificationCode(req.Number)
	storedVerificationCode := auth.StoredVerificationCode{
		Code: verificationCode.VerificationCode,
		Timestamp: pkgUtils.CurrentTimeMillis(),
	}
	item:=&model.TPendingAccounts{
		Number: req.Number,
		VerificationCode: storedVerificationCode.Code,
		Timestamp: storedVerificationCode.Timestamp,
		PushCode: "",
		DeletedAt: sql.NullTime{Time: time.Time{},Valid: false},

	}
	_,err:=l.svcCtx.PendAccountsModel.Insert(*item)
	if err!=nil{
		reason:=fmt.Sprintf(" insert db (%+v)发生错误：%s",err.Error())
		logx.Error(reason)
		return shared.Status(http.StatusInternalServerError,reason)
	}
	//发送验证码
	if req.Transport== "sms"{

	}else{
		return shared.Status(http.StatusNotImplemented,"!sms send verificationCode not implemenntnt")
	}
	logx.Infof("number:%s,code:%s",req.Number,verificationCode.VerificationCode)



	//返回响应
	return nil









	return nil
}


// 生成验证码
func generateVerificationCode(number string) util.VerificationCode {
	/*
	code, ok := conf.IsTestAccount(number)
	if ok {
		return utils.VerificationCode{VerificationCode: code}
	}

	 */

	n := big.NewInt(int64(pkgUtils.SecureRandIntN(900000)))
	randomInt := int(n.Add(n, big.NewInt(100000)).Int64())
	return util.NewVerificationCode(randomInt)
}

