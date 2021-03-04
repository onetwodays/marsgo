package logic

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/util"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginChainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginChainLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginChainLogic {
	return LoginChainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginChainLogic) LoginChain(req types.ChainLoginReq) (*types.ChainLoginRes, error) {
	// todo: add your logic here and delete this line
	activePublicKey:=""
	ownwerPublicKey:=""
	accountResp,_err:=l.svcCtx.EosApi.GetAccount(eos.AccountName(req.Name))
	if _err!=nil{
		return nil, _err
	}
	for index :=range accountResp.Permissions{
		if accountResp.Permissions[index].PermName=="active"{
			activePublicKey =accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
		}
		if accountResp.Permissions[index].PermName=="owner"{
			ownwerPublicKey =accountResp.Permissions[index].RequiredAuth.Keys[0].PublicKey.String() //默认取第一个作为public key
		}
	}
	logx.Infof("activePublicKey=%s,ownwerPublicKey=%s",activePublicKey,ownwerPublicKey)
	digest:=util.Sha256Str(req.Name)
	s,_err:=ecc.NewSignature(req.Sign)
	if _err!=nil{
		return nil, nil
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

	return &types.ChainLoginRes{
		JwtToken: types.JwtToken{
			AccessToken: jwtToken,
			AccessExpire: now+accessExpire,
			RefreshAfter: now+accessExpire/2,

		},
	}, nil
}
func (l *LoginChainLogic) getJwtToken(secretKey string, now, seconds int64, name string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds //jwt 过期时间
	claims["iat"] = now //jwt的签发时间
	claims["name"] = name //jwt携带用户id信息，用来跟header用户信息校验
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}