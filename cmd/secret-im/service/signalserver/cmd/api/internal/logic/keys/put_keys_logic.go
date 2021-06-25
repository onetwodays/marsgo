package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/logic/storage"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

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


func (l *PutKeysLogic) PutKeys(req types.PutKeysReq,appAccount *entities.Account) error {
	// todo: add your logic here and delete this line
	updateAccount:=false
	device := appAccount.AuthenticatedDevice
	wasAccountEnabled :=appAccount.IsEnabled()

	if req.IdentityKey!=appAccount.IdentityKey{
		updateAccount = true
		appAccount.IdentityKey= req.IdentityKey
	}

	if device.SignedPreKey==nil || req.SignedPreKey!= *(device.SignedPreKey){
		updateAccount = true
		device.SignedPreKey = &req.SignedPreKey

	}

	if updateAccount {
		if err:= new(storage.AccountManager).Update(appAccount,l.svcCtx.AccountsModel);err!=nil{
			return shared.Status(http.StatusInternalServerError,err.Error())
		}

		if !wasAccountEnabled && appAccount.IsEnabled() {
		}
	}







	return nil
}
