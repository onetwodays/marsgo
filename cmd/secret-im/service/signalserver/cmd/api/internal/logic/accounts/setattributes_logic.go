package logic

import (
	"context"
	"net/http"
	"secret-im/pkg/utils-tools"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SetattributesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetattributesLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetattributesLogic {
	return SetattributesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetattributesLogic) Setattributes(r *http.Request,req types.SetAttributesReq) error {
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil  {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)
	userAgent := r.Header.Get("User-Agent")

	// 更新帐号信息
	device := account.AuthenticatedDevice
	device.FetchesMessages = req.FetchesMessages
	device.Name = req.Name
	device.LastSeen = utils.TodayInMillis()
	device.Capabilities = req.Capabilities
	device.RegistrationID = req.RegistrationID
	device.SignalingKey = req.SignalingKey
	device.UserAgent = userAgent

	if len(req.Pin) > 0 {
		account.Pin = req.Pin
	} else if len(req.RegistrationLock) > 0 {
		credentials := auth.NewAuthenticationCredentials(req.RegistrationLock)
		account.RegistrationLock = credentials.HashedAuthenticationToken
		account.RegistrationLockSalt = credentials.Salt
	} else {
		account.Pin = ""
		account.RegistrationLock = ""
		account.RegistrationLockSalt = ""
	}

	account.UnidentifiedAccessKey = req.UnidentifiedAccessKey
	account.UnrestrictedUnidentifiedAccess = req.UnrestrictedUnidentifiedAccess

	if err := new(storage.AccountManager).Update(account); err != nil {
		return shared.Status(http.StatusInternalServerError,err.Error())

	}

	return nil
}
