package logic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/entities"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"strings"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetProfileLogic {
	return SetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetProfileLogic) SetProfile(r *http.Request,req types.CreateProfileRequest) (*types.ProfileAvatarUploadAttributes, error) {
	appAccount := r.Context().Value(shared.HttpReqContextAccountKey)
	if appAccount == nil  {
		reason := "check basic auth fail ,may by the handler not use middle"
		logx.Error(reason)
		return nil,shared.Status(http.StatusUnauthorized, reason)
	}
	account := appAccount.(*entities.Account)
	currentProfile,err:=storage.ProfileManager{}.Get(account.UUID,req.Version)
	if err!=nil && err!=sqlx.ErrNotFound {
		return nil, shared.Status(http.StatusInternalServerError,err.Error())
	}

	var avatar string =""
	if req.Avatar{
		avatar=generateAvatarObjectName()
	}

	// 已经处理的更新和插入
	err =  new(storage.ProfileManager).Set(account.UUID,&entities.VersionedProfile{
		Version: req.Version,
		Name: req.Name,
		Avatar: avatar,
		Commitment: req.Commitment,
	})
	if err!=nil{
		return nil,shared.Status(http.StatusInternalServerError,err.Error())
	}
	var res *types.ProfileAvatarUploadAttributes=nil
	if req.Avatar{
		var currentAvatar string
		if currentProfile!=nil &&
			len(currentProfile.Avatar)!=0 &&
			strings.HasPrefix(currentProfile.Avatar,"profiles/"){
		                  	currentAvatar=currentProfile.Avatar
		}

		if len(currentAvatar)==0 &&
			len(account.Avatar)!=0 &&
			strings.HasPrefix(account.Avatar,"profiles/"){
			currentAvatar=account.Avatar
		}
		if len(currentAvatar)!=0{
			//删除currentAvatar

		}
		res=&types.ProfileAvatarUploadAttributes{
			Key: avatar,
			Credential:"",
			Policy: "",
			Acl: "private",
			Date:time.Now().Format("20060102T150405Z07:00"),
			Signature: "",
		}
	}
	account.Name=req.Name
	if len(avatar)!=0{
		account.Avatar=avatar
	}
	err=new(storage.AccountManager).Update(account)
	if err!=nil{
		return nil, shared.Status(http.StatusInternalServerError,err.Error())
	}
	return res, nil
}


// 生成头像对象名
func generateAvatarObjectName() string {
	var object [16]byte
	rand.Read(object[:])
	return "profiles/" + base64.RawURLEncoding.EncodeToString(object[:])
}
