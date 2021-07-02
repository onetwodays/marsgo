package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/shared"
)

func DbAccount2AppAccount(dbAccount *model.TAccounts) (*entities.Account, error) {
	if dbAccount == nil {
		return nil, errors.New("param is nil pointer")
	}
	appAccount := new(entities.Account)
	err := json.Unmarshal([]byte(dbAccount.Data.String), appAccount)
	if err == nil {
		logx.Error("json.Unmarshal([]byte(dbAccount.Data.String)) error:", err)
		return nil, err
	}
	return appAccount, nil
}

// AccountManager 帐号管理器
type AccountManager struct {
}

//创建帐号

// Update 更新帐号
func (m AccountManager) Update(appAccount *entities.Account) error {
	jsb, err := json.Marshal(appAccount)
	if err != nil {
		return shared.Status(http.StatusInternalServerError, err.Error())
	}
	return internal.accountDB.UpdateDataByUuid(
		sql.NullString{
			String: string(jsb),
			Valid:  true,
		},
		appAccount.UUID)
}

func (m AccountManager) Get(ambiguousIdentifier auth.AmbiguousIdentifier) (*entities.Account, error) {

	var dbAccount *model.TAccounts

	var err error

	if len(ambiguousIdentifier.UUID) != 0 {
		dbAccount, err = internal.accountDB.FindOneByUuid(ambiguousIdentifier.UUID)
	}
	if dbAccount == nil {
		dbAccount, err = internal.accountDB.FindOneByNumber(ambiguousIdentifier.Number)
	}
	if err != nil {
		return nil, err
	}
	return DbAccount2AppAccount(dbAccount)
}

func (m AccountManager) GetByNumber(number string) (*entities.Account, error) {

	dbAccount, err := internal.accountDB.FindOneByNumber(number)
	if err != nil {
		return nil, err
	}
	return DbAccount2AppAccount(dbAccount)
}
