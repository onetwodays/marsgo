package logic

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/shared"
	"strings"
	"time"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)
var (
	errorUsernameUnRegister = shared.NewCodeError(shared.UsernameUnRegister, shared.CodeErrorMap[shared.UsernameUnRegister])
	errorIncorrectPassword  =  shared.NewCodeError(shared.IncorrectPassword, shared.CodeErrorMap[shared.IncorrectPassword])
)
type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.UserReply, error) {
	// todo: add your logic here and delete this line
	if len(strings.TrimSpace(req.Username)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		return nil, shared.ErrorParam
	}
	userInfo,err:=l.svcCtx.UserModel.FindOneByName(req.Username)
	switch err {
	case nil:
		if userInfo.Password!=req.Password{
			return nil, errorIncorrectPassword
		}else {
			now:=time.Now().Unix()
			accessExpire:=l.svcCtx.Config.Auth.AccessExpire
			jwtToken,err:= l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret,now,accessExpire,userInfo.Id)
			if err!=nil{
				return nil, err
			}
			return &types.UserReply{
				Id:       userInfo.Id,
				Username: userInfo.Name,
				Mobile:   userInfo.Mobile,
				Nickname: userInfo.Nickname,
				Gender:   userInfo.Gender,

				JwtToken: types.JwtToken{
					AccessToken: jwtToken,
					AccessExpire: now+accessExpire,
					RefreshAfter: now+accessExpire/2,

				},
			}, nil
		}
	case model.ErrNotFound:
		return nil, errorUsernameUnRegister
	default:
		return nil, err
	}


}


func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds //jwt 过期时间
	claims["iat"] = iat //jwt的签发时间
	claims["userId"] = userId //jwt携带用户id信息，用来跟header用户信息校验
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}