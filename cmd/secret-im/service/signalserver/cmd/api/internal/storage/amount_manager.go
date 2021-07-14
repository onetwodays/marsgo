package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	pkgUtils "secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/api/internal/types"
	shared "secret-im/service/signalserver/cmd/api/shared"
)

func DbAccount2AppAccount(dbAccount *model.TAccounts) (*entities.Account, error) {
	if dbAccount == nil {
		return nil, errors.New("param is nil pointer")
	}
	appAccount := new(entities.Account)
	err := json.Unmarshal([]byte(dbAccount.Data.String), appAccount)
	if err != nil {

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

func (m AccountManager) Get(ambiguousIdentifier *auth.AmbiguousIdentifier) (*model.TAccounts,*entities.Account, error) {

	var dbAccount *model.TAccounts

	var err error

	if len(ambiguousIdentifier.UUID) != 0 {
		dbAccount, err = internal.accountDB.FindOneByUuid(ambiguousIdentifier.UUID)
	}
	if dbAccount == nil {
		dbAccount, err = internal.accountDB.FindOneByNumber(ambiguousIdentifier.Number)
	}
	if err != nil {
		return nil,nil, err
	}
	appAccount,err:= DbAccount2AppAccount(dbAccount)
	if err!=nil{
		return nil, nil, err
	}
	return dbAccount,appAccount,nil
}

func (m AccountManager) GetByNumber(number string) (*entities.Account, error) {

	dbAccount, err := internal.accountDB.FindOneByNumber(number)
	if err != nil {
		return nil, err
	}
	return DbAccount2AppAccount(dbAccount)
}

func (m AccountManager) GetByNumbers(numbers []string) (map[string]entities.Account, error) {

	mapper := make(map[string]entities.Account)
	/*
	var nocache []string
	accounts, err := m.redisGetByNumbers(numbers)
	if err != nil {
		nocache = numbers
	} else {
		for idx, account := range accounts {
			if account == nil {
				nocache = append(nocache, numbers[idx])
				continue
			}
			mapper[account.Number] = *account
		}
	}

	if len(nocache) == 0 {
		return mapper, nil
	}

	 */

	list, err := internal.accountDB.FindManyByNumbers(numbers)
	if err != nil {
		return nil, err
	}
	for i:= range list {
		//m.redisSet(&account)
		item,err:= DbAccount2AppAccount(&list[i])
		if err!=nil{
			return nil, err
		}
		mapper[list[i].Uuid]=*item
	}
	return mapper, nil
}



func (m AccountManager) GetByUuid(uuid string) (*entities.Account, error) {

	dbAccount, err := internal.accountDB.FindOneByUuid(uuid)
	if err != nil {
		return nil, err
	}
	return DbAccount2AppAccount(dbAccount)
}


func (m AccountManager) GetByUUIDs(ids []string) (map[string]entities.Account, error) {
	mapper := make(map[string]entities.Account)
	/*
	var nocache []string

	accounts, err := m.redisGetByUUIDs(ids)
	if err != nil {
		nocache = ids
	} else {
		for idx, account := range accounts {
			if account == nil {
				nocache = append(nocache, ids[idx])
				continue
			}
			mapper[account.UUID] = *account
		}
	}

	if len(nocache) == 0 {
		return mapper, nil
	}

	 */

	list, err := internal.accountDB.FindManyByUuids(ids)
	if err != nil {
		return nil, err
	}
	for i:= range list {
		//m.redisSet(&account)
		item,err:= DbAccount2AppAccount(&list[i])
		if err!=nil{
			return nil, err
		}
		mapper[list[i].Uuid]=*item
	}
	return mapper, nil
}

func (m AccountManager)  CreateDBAccount(number, password,userAgent string, accountAttributes *types.AccountAttributes) (*model.TAccounts, error) {

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
