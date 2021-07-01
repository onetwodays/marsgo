package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/shared"
)

type CheckBasicAuthMiddleware struct {
	model model.TAccountsModel
}

func NewCheckBasicAuthMiddleware(model model.TAccountsModel) *CheckBasicAuthMiddleware {
	return &CheckBasicAuthMiddleware{
		model: model,
	}
}

func (m *CheckBasicAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r,err:=m.BasicAuth(r,false)
		if err!=nil{
			httpx.Error(w,shared.Status(http.StatusUnauthorized,err.Error()))
			return
		}
		next(w, r)
	}
}

func (m *CheckBasicAuthMiddleware) BasicAuth(r *http.Request, enabledRequired bool)  (*http.Request,error){

	appAccount,err,isOk:=m.BasicAuthForHeader(r,enabledRequired)

	if !isOk {
		return r,err
	}
	// 保存到context里面，可以往下传递
	ctx :=r.Context()
	ctx = context.WithValue(ctx,shared.HttpReqContextAccountKey,appAccount)
	r = r.WithContext(ctx)


	return r,nil
}

// 更新活动日期
func (m *CheckBasicAuthMiddleware) updateLastSeen(dbAccount *model.TAccounts,
	                                              appAccount * entities.Account,
	                                              device * entities.DeviceFull) error {
	if device.LastSeen == utils.TodayInMillis() {
		return nil
	}
	device.LastSeen = utils.TodayInMillis()
	jsb,_ := json.Marshal(appAccount)
	dbAccount.Data=sql.NullString{String: string(jsb),Valid: true}
	return m.model.Update(*dbAccount)

}
// 标头基本身份验证
func (m *CheckBasicAuthMiddleware) BasicAuthForHeader(r *http.Request,
	enabledRequired bool, ignorePassword ...bool) ( *entities.Account,error, bool) {
	authorizationHeader := r.Header.Get(shared.AuthorizationHeader)
	header, err := new(auth.AuthorizationHeader).FromFullHeader(authorizationHeader)
	//如果不是basic auth，现在放行，去http找头，密码保持为空。
	if err != nil {
		return nil,err,false
	}
	// 校验虚拟账号
	/*
	account, ok := authVirtualAccount(header, ignorePassword...)
	if ok {
		return account, true
	}

	 */

	// 获取帐号信息
	var dbaccount *model.TAccounts

	if len(header.Identifier.UUID) != 0 {
		dbaccount,err= m.model.FindOneByUuid(header.Identifier.UUID)
		if err!=nil{
			logx.Errorf("m.model.FindOneByUuid(%s):%s",header.Identifier.UUID,err)
		}
	}

	if (err!=nil || dbaccount==nil) && len(header.Identifier.Number) != 0  {
		dbaccount,err= m.model.FindOneByNumber(header.Identifier.Number)
	}
	if err!=nil {
		reason:=fmt.Sprintf("m.model.FindOneByNumber(%s):%s",header.Identifier.Number,err)
		logx.Error(reason)
		return nil,errors.New(reason),false
	}

	///////////////////////
	appAccount:=&entities.Account{}
	err= json.Unmarshal([]byte(dbaccount.Data.String),appAccount)
	if err!=nil {
		return nil,err,false
	}


	// 验证设备信息
	device, ok := appAccount.GetDevice(header.DeviceID)
	if !ok {
		reason := fmt.Sprintf("account(%s) has no header setting device (%d)",dbaccount.Number,header.DeviceID)
		return nil,errors.New(reason), false
	}

	if enabledRequired {
		if !device.IsEnabled() {
			reason := fmt.Sprintf("account(%s) has  header setting device (id=%d) is not enable",dbaccount.Number,header.DeviceID)
			return nil,errors.New(reason), false
		}

		if !appAccount.IsEnabled() {
			reason := fmt.Sprintf("account(%s) has  master device (id=%d) is not enable ",dbaccount.Number,header.DeviceID)
			return nil,errors.New(reason), false
		}
	}

	// 验证用户密码
	/*
	if len(ignorePassword) == 0 || !ignorePassword[0] {
		if !device.GetAuthenticationCredentials().Verify(header.Password) {
			reason := fmt.Sprintf("account(%s) check header passwd (%s) fail ",dbaccount.Number,header.Password)
			return nil,errors.New(reason), false
		}
	}

	 */
	appAccount.AuthenticatedDevice = device
	err=m.updateLastSeen(dbaccount,appAccount, device)
	if err!=nil{
		return nil,err,false
	}
	return appAccount,nil, true
}






