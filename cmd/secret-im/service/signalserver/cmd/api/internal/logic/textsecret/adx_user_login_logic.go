package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	eos "github.com/marsofsnow/eos-go"
	ecc "github.com/marsofsnow/eos-go/ecc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/utils"
	"secret-im/service/signalserver/cmd/api/util"

	//"secret-im/service/signalserver/cmd/api/util"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/shared"
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

func (l *AdxUserLoginLogic) AdxUserLogin(req types.AdxUserLoginReq) (*types.AdxUserLoginRes, error) {
	// todo: add your logic here and delete this line


	activePublicKey:=""
	ownwerPublicKey:=""
	accountResp,_err:=l.svcCtx.EosApi.GetAccount(eos.AccountName(req.Name))
	if _err!=nil {
		return nil, _err
	}
//Unmarshal: public key should start with ["PUB_K1_" | "PUB_R1_" | "PUB_WA_"] (or the old "EOS")

	for index :=range accountResp.Permissions{
		if accountResp.Permissions[index].PermName=="active"{
			activePublicKey =accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
		}
		if accountResp.Permissions[index].PermName=="owner"{
			ownwerPublicKey =accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
		}
	}

	logx.Info("activePublicKey=%s,ownwerPublicKey=%s",activePublicKey,ownwerPublicKey)
	digest:= util.Sha256Str(req.Name)
	s,_err:=ecc.NewSignature(req.Sign)
	if _err!=nil{
		return nil, _err
	}

	pubkey,_err:=ecc.NewPublicKey(activePublicKey)
	if _err!=nil{
		return nil, _err
	}

	if !s.Verify(digest,pubkey){
		return nil, errors.New("验证签名失败")
	}
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
	now:=time.Now().Unix()
	accessExpire:=l.svcCtx.Config.Auth.AccessExpire
	jwtToken,err:= l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret,now,accessExpire,req.Name)
	if err!=nil{
		return nil, err
	}


	//新用户的话，创建一个，不是新用户的话，返回已经存在的uuid
	account,err:=l.svcCtx.AccountsModel.FindOneByNumber(req.Name)
	if err!=nil && err!=sqlx.ErrNotFound{
		return nil,err
	}

	uuid:= ""
	isNew:=false
	if err!=nil && err==sqlx.ErrNotFound{
		uuid = utils.NewUuid()
		account:=&model.TAccounts{
			Number: req.Name,
			Uuid:uuid,
			Data: sql.NullString{String: "",Valid: false},
			DeleteTime: sql.NullTime{Time: time.Time{},Valid: false},
		}
		_,err:=l.svcCtx.AccountsModel.Insert(*account)
		if err!=nil{
			return nil, err
		}
		isNew=true

	}else{
		uuid = account.Uuid
	}


	return &types.AdxUserLoginRes{
		JwtTokenAdx: types.JwtTokenAdx{
			AccessToken: jwtToken,
			AccessExpire: now+accessExpire,
			RefreshAfter: now+accessExpire/2,
		},
		Uuid: uuid,
		IsNew: isNew,
	}, nil
}
func (l *AdxUserLoginLogic) getJwtToken(secretKey string, now, seconds int64, name string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds //jwt 过期时间
	claims["iat"] = now //jwt的签发时间
	claims[shared.JWTADXUSERNAME] = name //jwt携带用户id信息，用来跟header用户信息校验
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}