package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"net/http"
	pkgUtils "secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

type VerifyAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) VerifyAccountLogic {
	return VerifyAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}


func (l *VerifyAccountLogic) VerifyAccount(reqheader http.Header,req types.VerifyAccountReq) (*types.VerifyAccountRes, error) {

	// 1.校验head
	authorizationHeader,err:=l.checkAuthHead(reqheader)
	if err!=nil {
		return nil, err
	}
	number := authorizationHeader.Identifier.Number

	// 2.验证手机收到的验证码
	/*
	err = l.checkCode(number,req.VerificationCode)
	if err!=nil{
		return nil, err
	}

	 */

	accountAttr :=&types.AccountAttributes{
		SignalingKey:req.SignalingKey,
		FetchesMessages:req.FetchesMessages,
		RegistrationID:req.RegistrationID,
		Name :req.Name,
		Pin :req.Pin,
		RegistrationLock:req.RegistrationLock,
		UnidentifiedAccessKey:req.UnidentifiedAccessKey,
		UnrestrictedUnidentifiedAccess:req.UnrestrictedUnidentifiedAccess,
		Capabilities:req.Capabilities,
	}



	// 3.查询一下帐号表，
	dbAccount,err:=l.svcCtx.AccountsModel.FindOneByNumber(number)
	if err!=nil && err!=sqlx.ErrNotFound{
		reason:=fmt.Sprintf("error:%s,%s is  exist",err.Error(),number)
		logx.Error(reason)
		return nil, shared.Status(http.StatusInternalServerError,reason)
	}

	//4.帐号存在时,验证pin码，pin码验证通过之后，要更新一下帐号系统
	if err==nil && dbAccount!=nil {
		isOk,err := l.checkPin(number,dbAccount,accountAttr)
		if !isOk {
			return nil,err
		}
	}

	userAgent := reqheader.Get("User-Agent")
	dbNewaccount,err:=storage.AccountManager{}.CreateDBAccount(number,authorizationHeader.Password,userAgent,accountAttr)

	err = l.createOrUpdate(dbAccount!=nil,dbNewaccount)
	if err != nil {
		return nil, err
	}

	//删除用过的验证码
	/*
		err  = l.deleteCode(number)
		if err!=nil{
			return nil, err
		}
	*/


	return &types.VerifyAccountRes{UUID:dbNewaccount.Uuid }, nil
}





func (l *VerifyAccountLogic) checkAuthHead(reqheader http.Header) (*auth.AuthorizationHeader,error){
	authorizationHeaderValue := reqheader.Get("Authorization") //拿这个头解析
	header, err := new(auth.AuthorizationHeader).FromFullHeader(authorizationHeaderValue)
	if err != nil {
		reason:=fmt.Sprintf("Authorization Basic check fail:%s",err.Error())
		return nil, shared.Status(http.StatusBadRequest,reason)
	}

	if len(header.Identifier.Number) == 0 {
		reason:=fmt.Sprintf("Authorization Basic number is empty:%+v",header)
		return nil, shared.Status(http.StatusBadRequest,reason)
	}
	return header,nil

}

func (l *VerifyAccountLogic) checkCode(number string,code string) error {
	record,err:=l.svcCtx.PendAccountsModel.FindOneByNumber(number)
	if err!=nil{
		return shared.Status(http.StatusForbidden,err.Error())
	}

	storedVerficationCode:=auth.StoredVerificationCode{
		Code: record.VerificationCode,
		Timestamp: record.Timestamp,
		PushCode: record.PushCode,
	}
	if !storedVerficationCode.IsValid(code){
		reason:=fmt.Sprintf("verification code(%s) invalid,right(%s)",code,storedVerficationCode.Code)
		return shared.Status(http.StatusForbidden,reason)
	}
	return nil
}





func (l *VerifyAccountLogic) deleteCode(number string) error {
	err  := l.svcCtx.PendAccountsModel.DeleteByNumber(number)
	if err!=nil {
		reason:=fmt.Sprintf("l.svcCtx.PendAccountsModel.DeleteByNumber(number) error:%s",err.Error())
		logx.Error(reason)
		return shared.Status(http.StatusInternalServerError,reason)
	}
	return nil
}



func (l *VerifyAccountLogic) createOrUpdate(isUpdate bool,dbAccount *model.TAccounts) error {
	var err error
	if isUpdate{
		err = l.svcCtx.AccountsModel.Update(*dbAccount)
	}else{
		_,err=l.svcCtx.AccountsModel.Insert(*dbAccount)
	}
	if err!=nil {
		return shared.Status(http.StatusInternalServerError,err.Error())
	}
	return nil
}

func (l *VerifyAccountLogic) checkPin(number string,
	                                  dbAccount *model.TAccounts,
	                                  accountAttributes *types.AccountAttributes) (bool,error) {


	appAccount:=&entities.Account{}
	err:=json.Unmarshal([]byte(dbAccount.Data.String),appAccount)
	if err!=nil{
		return false, shared.Status(http.StatusInternalServerError,err.Error())
	}

	//如果帐号有pin码或者有注册锁，且距第一次登录<7 天，
	if (len(appAccount.Pin)!=0 || len(appAccount.RegistrationLock)!=0) &&
		pkgUtils.CurrentTimeMillis()-appAccount.GetLastSeen() <pkgUtils.DaysToMillis(7) {

		//如果帐号有pin码或者有注册锁，且距第一次登录<7 天，
		// c.rateLimiters.VerifyLimiter.Clear(number)

		var credentials *auth.ExternalServiceCredentials
		timeRemaining := pkgUtils.DaysToMillis(7) - (pkgUtils.CurrentTimeMillis() - appAccount.GetLastSeen())

		if len(appAccount.RegistrationLock) != 0 && len(appAccount.RegistrationLockSalt) != 0 {
			credentials = l.svcCtx.BackupCredentialsGenerator.GenerateFor(number)
		}

		registrationLockFailure:= entities.RegistrationLockFailure {
			TimeRemaining: timeRemaining,
			BackupCredentials: credentials, //返回帐号已有的东西
		}
		jsb,_:=json.Marshal(registrationLockFailure)
		reason:=string(jsb)
		//客户端没有上传pin和RegistrationLock
		if len(accountAttributes.Pin) == 0 && len(accountAttributes.RegistrationLock)==0 {
			// 注册锁定失败
			logx.Errorf("告诉客户端注册锁定:%s",reason)
			return false, shared.Status(http.StatusLocked,reason)
		}

		/*
			if !c.rateLimiters.PinLimiter.Validate(number, 1) {
				respond.Status(w, http.StatusRequestEntityTooLarge)
				return
			}
		*/
		var pinMatches bool
		if len(appAccount.RegistrationLock)!=0 && len (appAccount.RegistrationLockSalt) !=0 {
			authenticationCredentials := auth.AuthenticationCredentials{
				Salt:                      appAccount.RegistrationLockSalt, //帐号里保存的盐
				HashedAuthenticationToken: appAccount.RegistrationLock,     //帐号里保存的hash值
			}
			pinMatches = authenticationCredentials.Verify(accountAttributes.RegistrationLock) //请求里的的
		} else {
			pinMatches = appAccount.Pin == accountAttributes.Pin
		}

		if !pinMatches{
			logx.Infof("%s","pin not match")
			return false, shared.Status(http.StatusLocked,reason)
		}
		//c.rateLimiters.PinLimiter.Clear(number)
	}

	return true,nil

}


// 创建帐号
func (l *VerifyAccountLogic) createDBAccount(number, password,userAgent string, accountAttributes *types.AccountAttributes) (*model.TAccounts, error) {

	device := new(entities.DeviceFull)
	device.ID = entities.DeviceMasterID
	device.SetAuthenticationCredentials(auth.NewAuthenticationCredentials(password))
	device.SignalingKey = accountAttributes.SignalingKey
	device.FetchesMessages = accountAttributes.FetchesMessages
	device.RegistrationID = accountAttributes.RegistrationID
	device.Name = accountAttributes.Name
	device.Capabilities = accountAttributes.Capabilities
	device.Created = pkgUtils.CurrentTimeMillis()
	device.LastSeen = pkgUtils.TodayInMillis()
	device.UserAgent = userAgent

	account := new(entities.Account)
	account.Number= number
	account.UUID = uuid.NewV4().String()
	account.AddDevice(device)
	account.Pin = accountAttributes.Pin
	account.UnidentifiedAccessKey = accountAttributes.UnidentifiedAccessKey
	account.UnrestrictedUnidentifiedAccess = accountAttributes.UnrestrictedUnidentifiedAccess

	jsb,err:=json.Marshal(account)
	if err!=nil{
		return nil, err
	}

	return &model.TAccounts{
		Number: number,
		Uuid: account.UUID,
		Data: sql.NullString{String: string(jsb),Valid: true},

	},nil
}
