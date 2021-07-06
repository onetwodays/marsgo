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
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
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

		r, err := m.basicAuth(r, true)
		if err != nil {
			e := shared.Status(http.StatusUnauthorized, err.Error())
			logx.Error("CheckBasicAuthMiddleware fail:", e)
			httpx.Error(w, e)
			return
		}
		next(w, r)
	}
}

func (m *CheckBasicAuthMiddleware) basicAuth(r *http.Request, enabledRequired bool) (*http.Request, error) {
	appAccount, err := m.BasicAuthByHeader(r, enabledRequired)
	if err != nil {
		return r, err
	}
	// 保存到context里面，可以往下传递
	ctx := r.Context()
	ctx = context.WithValue(ctx, shared.HttpReqContextAccountKey, appAccount)
	r = r.WithContext(ctx)
	return r, nil
}

func (m *CheckBasicAuthMiddleware) BasicAuthByUserPasswd(user, passwd string, enabledRequired bool, ignorePassword ...bool) (*entities.Account, error) {

	header, err := new(auth.AuthorizationHeader).FromUserAndPassword(user, passwd)
	if err != nil {
		return nil, err
	}
	return m.basicAuthForHeader(header, enabledRequired, ignorePassword...)
}

func (m *CheckBasicAuthMiddleware) BasicAuthByHeader(r *http.Request, enabledRequired bool, ignorePassword ...bool) (*entities.Account, error) {

	header := new(auth.AuthorizationHeader)

	ctx := r.Context().Value("ws")
	if ctx != nil {
		logx.Info("1.通过websocket进来的请求")
		uuid := r.Context().Value(shared.CONTENTKEYUUID)
		deviceId := r.Context().Value(shared.CONTENTKEYDEVICEID)
		if uuid != nil && deviceId != nil {
			logx.Infof("2.从context拿uuid=%s和deviceid=%d来校验头", uuid, deviceId)
			header.Identifier = auth.AmbiguousIdentifier{
				UUID: uuid.(string),
			}
			header.DeviceID = deviceId.(int64)
			logx.Info("3.websocket请求context,basicAuthForHeader.....")
			return m.basicAuthForHeader(header, enabledRequired, ignorePassword...)
		}
	}
	_, err := header.FromFullHeader(r.Header.Get(shared.AuthorizationHeader))
	if err != nil {
		logx.Error("header.FromFullHeader fail:", err, " 再次尝试自定义头")
		//return nil, err
	}

	// 这里容错一下，防止有些头没有传进来
	if len(header.Identifier.UUID) == 0 && len(header.Identifier.Number) == 0 {

		logx.Error("header.Identifier.UUID 与header.Identifier.Number 都是空，已从x-user-name 取值，作为用户标识 ")
		header.Identifier.Number = r.Header.Get(shared.HEADADXUSERNAME)
		logx.Info(shared.HEADADXUSERNAME, "=", header.Identifier.Number)
	}
	if len(header.Identifier.Number) == 0 {
		err = errors.New("x-user-name 值为空，报错")
		logx.Error(err)
		return nil, err
	}
	if header.DeviceID == 0 {
		logx.Error("设备id=0,，已默认是1")
		header.DeviceID = entities.DeviceMasterID
	}
	return m.basicAuthForHeader(header, enabledRequired, ignorePassword...)
}

// 更新活动日期
func (m *CheckBasicAuthMiddleware) updateLastSeen(dbAccount *model.TAccounts,
	appAccount *entities.Account,
	device *entities.DeviceFull) error {
	if device.LastSeen == utils.TodayInMillis() {
		return nil
	}
	device.LastSeen = utils.TodayInMillis()
	jsb, _ := json.Marshal(appAccount)
	dbAccount.Data = sql.NullString{String: string(jsb), Valid: true}
	return m.model.Update(*dbAccount)

}

// 标头基本身份验证
func (m *CheckBasicAuthMiddleware) basicAuthForHeader(header *auth.AuthorizationHeader, enabledRequired bool, ignorePassword ...bool) (*entities.Account, error) {

	// 校验虚拟账号
	/*
		account, ok := authVirtualAccount(header, ignorePassword...)
		if ok {
			return account, true
		}

	*/

	// 获取帐号信息
	dbAccount, appAccount, err := storage.AccountManager{}.Get(&header.Identifier)
	if err != nil {
		return nil, err
	}

	// 验证设备信息
	number := header.Identifier.Number
	device, ok := appAccount.GetDevice(header.DeviceID)
	if !ok {
		reason := fmt.Sprintf("account(%s) has no header setting device (%d)", number, header.DeviceID)
		return nil, errors.New(reason)
	}
	// todo
	enabledRequired = false
	if enabledRequired {
		if !device.IsEnabled() {
			reason := fmt.Sprintf("account(%s) has  header setting device (id=%d) is not enable", number, header.DeviceID)
			return nil, errors.New(reason)
		}

		if !appAccount.IsEnabled() {
			reason := fmt.Sprintf("account(%s) has  master device (id=%d) is not enable ", number, header.DeviceID)
			return nil, errors.New(reason)
		}
	}

	// 验证用户密码

	if len(ignorePassword) == 0 || !ignorePassword[0] {

		if !device.GetAuthenticationCredentials().Verify(header.Password) {
			reason := fmt.Sprintf("account(%s) check header passwd (%s) fail，but close check ", number, header.Password)
			logx.Error(reason)
			//todo
			//return nil, errors.New(reason)
		}
	}

	appAccount.AuthenticatedDevice = device
	err = m.updateLastSeen(dbAccount, appAccount, device)
	if err != nil {
		return nil, err
	}
	return appAccount, nil
}
