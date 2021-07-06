package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

type PutKeysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutKeysLogic {
	return PutKeysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}


func (l *PutKeysLogic) PutKeys(r *http.Request,req types.PutKeysReq) error {

	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)


	updateAccount:=false
	device := account.AuthenticatedDevice
	wasAccountEnabled :=account.IsEnabled()

	if req.IdentityKey!=account.IdentityKey{
		updateAccount = true
		account.IdentityKey= req.IdentityKey
	}

	if device.SignedPreKey==nil || req.SignedPreKey!= *(device.SignedPreKey){
		updateAccount = true
		device.SignedPreKey = &req.SignedPreKey
	}

	if updateAccount {
		if err:= new(storage.AccountManager).Update(account);err!=nil{
			return shared.Status(http.StatusInternalServerError,err.Error())
		}

		if !wasAccountEnabled && account.IsEnabled() {
		}
	}

	// 把prekeys批量插入数据库
	for i:=range req.PreKeys{
		err:=l.svcCtx.PreKeysInsertor.Insert(account.Number,device.ID,req.PreKeys[i].KeyId,req.PreKeys[i].PublicKey)
		if err!=nil{
			return shared.Status(http.StatusInternalServerError,err.Error())
		}
	}
	return nil
}
