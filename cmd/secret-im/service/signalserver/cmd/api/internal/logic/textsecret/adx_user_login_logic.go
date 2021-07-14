package logic

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/marsofsnow/eos-go"
	"github.com/marsofsnow/eos-go/ecc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/util"


	"secret-im/service/signalserver/cmd/api/internal/model"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AdxUserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdxUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) AdxUserLoginLogic {
	return AdxUserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdxUserLoginLogic) AdxUserLogin(req types.AdxUserLoginReq,userAgent string) (*types.AdxUserLoginRes, error) {
	// todo: add your logic here and delete this line
	if err := l.checkEosUserValid(req.Account, req.Sign); err != nil {
		return nil, shared.Status(http.StatusInternalServerError,err.Error())
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, req.Account)
	if err != nil {
		return nil, shared.Status(http.StatusInternalServerError,err.Error())
	}

	accountAttr := &types.AccountAttributes{
		SignalingKey:                   req.SignalingKey,
		FetchesMessages:                req.FetchesMessages,
		RegistrationID:                 req.RegistrationID,
		Name:                           req.Name,
		Pin:                            req.Pin,
		RegistrationLock:               req.RegistrationLock,
		UnidentifiedAccessKey:          req.UnidentifiedAccessKey,
		UnrestrictedUnidentifiedAccess: req.UnrestrictedUnidentifiedAccess,
		Capabilities: types.DeviceCapabilities{
			UUID: req.Capabilities.UUID,
		},
	}

	//新用户的话，创建一个，不是新用户的话，返回已经存在的uuid
	dbAccount, err := l.svcCtx.AccountsModel.FindOneByNumber(req.Account)
	if err != nil && err != sqlx.ErrNotFound {
		return nil, shared.Status(http.StatusInternalServerError,err.Error())
	}



	dbNewaccount, err := storage.AccountManager{}.CreateDBAccount(req.Account, "1234",userAgent , accountAttr)

	//err = l.createOrUpdate(dbAccount != nil, dbNewaccount)
	err = l.createOrUpdate(false, dbNewaccount)
	if err != nil {
		return nil, err
	}
	var uuid string
	var isNew bool
	if dbAccount == nil {
		uuid = dbNewaccount.Uuid
		isNew = true
	} else {
		uuid = dbAccount.Uuid
		isNew = false
	}

	return &types.AdxUserLoginRes{
		JwtTokenAdx: types.JwtTokenAdx{
			AccessToken:  jwtToken,
			AccessExpire: now + accessExpire,
			RefreshAfter: now + accessExpire/2,
		},
		Uuid:  uuid,
		IsNew: isNew,
	}, nil
}

func (l *AdxUserLoginLogic) createOrUpdate(isUpdate bool, dbAccount *model.TAccounts) error {
	var err error
	if isUpdate {
		err = l.svcCtx.AccountsModel.Update(*dbAccount)
	} else {
		_, err = l.svcCtx.AccountsModel.Insert(*dbAccount)
	}
	if err != nil {
		return shared.Status(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (l *AdxUserLoginLogic) checkEosUserValid(userName, userNameSign string) error {
	activePublicKey := ""
	ownwerPublicKey := ""
	accountResp, _err := l.svcCtx.EosApi.GetAccount(eos.AccountName(userName))
	if _err != nil {
		return _err
	}
	//Unmarshal: public key should start with ["PUB_K1_" | "PUB_R1_" | "PUB_WA_"] (or the old "EOS")

	for index := range accountResp.Permissions {
		if accountResp.Permissions[index].PermName == "active" {
			activePublicKey = accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
		}
		if accountResp.Permissions[index].PermName == "owner" {
			ownwerPublicKey = accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
		}
	}

	logx.Info("activePublicKey=%s,ownwerPublicKey=%s", activePublicKey, ownwerPublicKey)
	digest := util.Sha256Str(userName)
	s, _err := ecc.NewSignature(userNameSign)
	if _err != nil {
		return _err
	}

	pubkey, _err := ecc.NewPublicKey(activePublicKey)
	if _err != nil {
		return _err
	}

	if !s.Verify(digest, pubkey) {
		return errors.New("验证EOS签名失败")
	}
	return nil

	/*
	   private_key:="5K8bmn8AMNewSzgB3VNnz7pahVVLTF7LaksnF8tjoSPVcvS2xDw"
	   public_key:= "EOS8GnxsyzChJF4pobmcvYan3qv1ovw66ypYNiJZatSWKkejfEqg4"

	   privkey,_:=ecc.NewPrivateKey(private_key)
	   pubkey,_:=ecc.NewPublicKey(public_key)
	   digest:=tool.Sha256Str(userName)
	   signature,_:=privkey.Sign(digest) //对hash值签名
	   fmt.Println("签名的内容是：",signature.String(),"|")

	   s,_:=ecc.NewSignature(signature.String())
	   fmt.Println("验签结果:",s.Verify(digest,pubkey))

	*/
}

func (l *AdxUserLoginLogic) getJwtToken(secretKey string, now, seconds int64, name string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds         //jwt 过期时间
	claims["iat"] = now                   //jwt的签发时间
	claims[shared.JWTADXUSERNAME] = name //jwt携带用户id信息，用来跟header用户信息校验
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
