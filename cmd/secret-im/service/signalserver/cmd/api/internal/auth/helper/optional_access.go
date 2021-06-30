package helper

import (
	"bytes"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"strconv"
)

var UNIDENTIFIED = "Unidentified-Access-Key"

// 可选访问方式
type OptionalAccess struct {
}

// 验证访问权限
func (OptionalAccess) Verify(
	requestAccount *entities.Account,
	accessKey *auth.Anonymous, targetAccount *entities.Account) (int, bool) {

	if requestAccount != nil && targetAccount != nil && targetAccount.IsEnabled() {
		return http.StatusOK, true
	}

	if requestAccount != nil && (targetAccount == nil || (targetAccount != nil && !targetAccount.IsEnabled())) {
		return http.StatusNotFound, false
	}

	if accessKey != nil && targetAccount != nil && targetAccount.IsEnabled() &&
		targetAccount.UnrestrictedUnidentifiedAccess {
		return http.StatusOK, true
	}

	if accessKey != nil &&
		targetAccount != nil &&
		len(targetAccount.UnidentifiedAccessKey) != 0 &&
		targetAccount.IsEnabled() &&
		bytes.Equal(accessKey.UnidentifiedSenderAccessKey, []byte(targetAccount.UnidentifiedAccessKey)) {
		return http.StatusOK, true
	}
	return http.StatusUnauthorized, false
}

// 验证设备访问权限
func (oa OptionalAccess) VerifyDevices(
	requestAccount *entities.Account,
	accessKey *auth.Anonymous, targetAccount *entities.Account, deviceSelector string) (int, bool) {
	code, ok := oa.Verify(requestAccount, accessKey, targetAccount)
	if !ok {
		return code, false
	}

	if deviceSelector == "*" {
		return http.StatusOK, true
	}

	deviceID, err := strconv.ParseInt(deviceSelector, 10, 64)
	if err != nil {
		return http.StatusUnprocessableEntity, false
	}

	_, ok = targetAccount.GetDevice(deviceID)
	if ok {
		return http.StatusOK, true
	}

	if requestAccount != nil {
		return http.StatusNotFound, true
	} else {
		return http.StatusUnauthorized, true
	}
}

