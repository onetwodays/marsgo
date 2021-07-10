package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/shared"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProfileByUserNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileByUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProfileByUserNameLogic {
	return GetProfileByUserNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileByUserNameLogic) GetProfileByUserName(r *http.Request,req types.GetProfileByUserName) (*types.Profile, error) {
	username := strings.ToLower(req.UserName)
	uuid,err:=storage.UsernamesManager{}.GetUUIDForUsername(username)
	if err!=nil{
		return nil, shared.Status(http.StatusNotFound,err.Error())
	}
	dstAccount,err:=storage.AccountManager{}.GetByUuid(uuid)
	if err!=nil {
		return nil, shared.Status(http.StatusNotFound,err.Error())
	}

	mac := hmac.New(sha256.New, []byte(dstAccount.UnidentifiedAccessKey))
	//mac.Write([]byte(prefix))
	sum := mac.Sum(nil)
	if len(sum) > 32 {
		sum = sum[:32]
	}


	return &types.Profile{
		IdentityKey:dstAccount.IdentityKey,
		Name:dstAccount.Name,
		Avatar :dstAccount.Avatar,
		UnidentifiedAccess:base64.StdEncoding.EncodeToString(sum), //todo 这里是计算一个hash值的
		UnrestrictedUnidentifiedAccess:dstAccount.UnrestrictedUnidentifiedAccess,
		Capabilities:types.UserCapabilities{
			Uuid: true,
			Gv2: true,
		},
		UserName:username,
		Uuid:dstAccount.UUID,
		Credential: "",
	}, nil
}

