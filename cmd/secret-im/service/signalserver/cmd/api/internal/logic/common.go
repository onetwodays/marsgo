package logic

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/tal-tech/go-zero/core/logx"
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

func Print(obj interface{}){

	jsb,_:=json.Marshal(obj)
	logx.Info("=====",string(jsb))
}

func GetUnidentifiedAccessChecksum(key []byte) string{
	mac := hmac.New(sha256.New,[]byte(key))
	mac.Write(make([]byte,32,32))
	sum:=mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}


