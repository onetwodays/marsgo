package logic

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/auth"
	"secret-im/service/signalserver/cmd/api/internal/auth/helper"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"strings"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProfileByUuidCredentiaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}


func NewGetProfileByUuidCredentiaLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProfileByUuidCredentiaLogic {
	return GetProfileByUuidCredentiaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileByUuidCredentiaLogic) GetProfileByUuidCredentia(r *http.Request,req types.GetProfileByUUIDCredentia) (*types.Profile, error) {


	source,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)

	header := r.Header.Get(helper.UNIDENTIFIED)
	accessKey, _ := auth.NewAnonymous(header)
	if accessKey == nil && source==nil {
		return nil, shared.Status(http.StatusUnauthorized, err.Error())
	}
    // 得到目的帐号
	dstAccount,err:=storage.AccountManager{}.GetByUuid(req.Uuid)
	if err!=nil {
		return nil,shared.Status(http.StatusNotFound,err.Error())
	}

	code, ok := helper.OptionalAccess{}.Verify(source, accessKey, dstAccount)
	if !ok {
		return nil,shared.Status(code,"helper.OptionalAccess{}.Verify(source,accessKey,destination) fail")
	}


	versionedProfile,err:=storage.ProfileManager{}.Get(req.Uuid,req.Version)
	if err!=nil {
		return nil, shared.Status(http.StatusNotFound,err.Error())
	}
	username,_:=storage.UsernamesManager{}.GetUsernameForUUID(dstAccount.UUID)
	name:=dstAccount.Name
	avatar:=dstAccount.Avatar

	if versionedProfile!=nil{
		name=versionedProfile.Name
		avatar=versionedProfile.Avatar
	}

	credential,err:=l.getProfileCredential(req.CredentialRequest,versionedProfile,req.Uuid)
	if err!=nil{
		return nil,shared.Status(http.StatusInternalServerError,err.Error())
	}
	unidentifiedAccessKey,_:=base64.StdEncoding.DecodeString(dstAccount.UnidentifiedAccessKey)
	resp:= &types.Profile{
		IdentityKey:dstAccount.IdentityKey,
		Name:name,
		Avatar :avatar,
		UnidentifiedAccess:logic.GetUnidentifiedAccessChecksum(unidentifiedAccessKey), //todo 这里是计算一个hash值的
		UnrestrictedUnidentifiedAccess:dstAccount.UnrestrictedUnidentifiedAccess,
		Capabilities:types.UserCapabilities{
			Uuid: true,
			Gv2: true,
		},
		UserName:username,
		Uuid:dstAccount.UUID,
		Credential: credential,

	}
	logic.Print(resp)

	return resp,  nil
}

func (l *GetProfileByUuidCredentiaLogic)getProfileCredential(credentialRequest string , profile *entities.VersionedProfile,uuid string ) (string,error) {
	// 1.credentialRequest 如果为空
	if len(credentialRequest)==0 || profile==nil{
		return "",nil
	}
	commitment,_:=base64.StdEncoding.DecodeString(profile.Commitment)
	credentialRequest=strings.ReplaceAll(credentialRequest,"(byte)0x","")
	request,err:=hex.DecodeString(credentialRequest)

	/*
	sp,err:= zkgroup.GenerateServerSecretParams()
	if err!=nil{
		return "", err
	}
	logx.Infof("====(%s)========",base64.StdEncoding.EncodeToString(sp))
	po:=zkgroup.NewServerZkProfileOperations(sp)
	 */

	b,err:=l.svcCtx.ServerZkProfileOperations.IssueProfileKeyCredential(request,[]byte(uuid),commitment)
	if err!=nil{
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b),nil



}


