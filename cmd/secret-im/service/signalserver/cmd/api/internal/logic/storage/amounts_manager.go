package storage

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/model"
	"secret-im/service/signalserver/cmd/shared"
)

// 帐号管理器
type AccountManager struct {

}

//创建帐号


//更新帐号
func (m AccountManager) Update(appAccount* entities.Account,db model.TAccountsModel)  error {
	jsb,err := json.Marshal(appAccount)
	if err!=nil{
		return shared.Status(http.StatusInternalServerError,err.Error())
	}
	return  db.UpdateDataByUuid(
		sql.NullString{
		String: string(jsb),
		Valid: true,
		},
		appAccount.UUID)
}

func (m AccountManager) Get (ambiguousIdentifier auth.AmbiguousIdentifier,db model.TAccountsModel) (*entities.Account,error){

	var dbAccount *model.TAccounts

	var err error

	if len(ambiguousIdentifier.UUID)!=0{
		dbAccount,err=db.FindOneByUuid(ambiguousIdentifier.UUID)
	}
	if dbAccount==nil{
		dbAccount,err = db.FindOneByNumber(ambiguousIdentifier.Number)
	}


	if dbAccount!=nil && err==nil{
		appAccount:=new(entities.Account)
		err = json.Unmarshal([]byte(dbAccount.Data.String),appAccount)
		if err == nil{
			return appAccount,nil
		}
	}


	return nil,err
}
