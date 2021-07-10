package logic

import (
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/middleware"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/shared"
)

func GetSourceAccount(r *http.Request,db model.TAccountsModel) (*entities.Account,error){
	var account *entities.Account
	var err error
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil  {
		checkBasicAuth := middleware.NewCheckBasicAuthMiddleware(db)
		account, err = checkBasicAuth.BasicAuthByHeader(r, true)
	}else {
		account= appAccount.(*entities.Account)
	}

	return account,err
}

