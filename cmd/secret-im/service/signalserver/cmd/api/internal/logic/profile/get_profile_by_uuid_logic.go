package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/auth/helper"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProfileByUuidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileByUuidLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProfileByUuidLogic {
	return GetProfileByUuidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileByUuidLogic) GetProfileByUuid(r *http.Request,req types.GetProfileByUUID) (*types.Profile, error) {
	header := r.Header.Get(helper.UNIDENTIFIED)
	accessKey, _ := auth.NewAnonymous(header)
	source,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if accessKey == nil && source==nil {
		return nil, shared.Status(http.StatusUnauthorized, "")
	}
	dstAccount,err:=storage.AccountManager{}.GetByUuid(req.Uuid)
	if err!=nil{
		return nil,shared.Status(http.StatusNotFound,err.Error())
	}

	code, ok := helper.OptionalAccess{}.Verify(source, accessKey, dstAccount)
	if !ok {
		reason:="helper.OptionalAccess{}.Verify(source,accessKey,destination) fail"
		logx.Error(reason)
		return nil,shared.Status(code,reason)
	}


	versionedProfile,err:=storage.ProfileManager{}.Get(req.Uuid,req.Version)
	if err!=nil {
		return nil, shared.Status(http.StatusNotFound,err.Error())
	}
	username,_:=storage.UsernamesManager{}.GetUsernameForUUID(dstAccount.UUID)
	name:=versionedProfile.Name
	avatar:=versionedProfile.Avatar
	credential:=""

	return &types.Profile{
		IdentityKey:dstAccount.IdentityKey,
		Name:name,
		Avatar :avatar,
		UnidentifiedAccess:"", //todo 这里是计算一个hash值的
		UnrestrictedUnidentifiedAccess:dstAccount.UnrestrictedUnidentifiedAccess,
		Capabilities:types.UserCapabilities{
			Uuid: true,
			Gv2: true,
		},
		UserName:username,
		Uuid:dstAccount.UUID,
		Credential: credential,

	}, nil
}
