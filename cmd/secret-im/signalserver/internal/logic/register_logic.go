package logic

import (
	"context"
	"secret-im/model"

	"secret-im/signalserver/internal/svc"
	"secret-im/signalserver/internal/types"
	"secret-im/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

var (
	errorDuplicateUserName = shared.NewCodeError(shared.DuplicateUsername, shared.CodeErrorMap[shared.DuplicateUsername])
	errorDuplicateUserMobile = shared.NewCodeError(shared.DuplicateMobile, shared.CodeErrorMap[shared.DuplicateMobile])
)


type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterLogic {
	return RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req types.RegisterReq) error {
	// todo: add your logic here and delete this line
	_,err:=l.svcCtx.UserModel.FindOneByName(req.Username)
	if err!=nil{
		return errorDuplicateUserName
	}
	_,err =l.svcCtx.UserModel.FindOneByMobile(req.Mobile)
	if err!=nil{
		return errorDuplicateUserMobile
	}
	_,err= l.svcCtx.UserModel.Insert(model.User{
		Name:req.Username,
		Password: req.Password,
		Mobile: req.Mobile,
		Gender: "boy",
		Nickname: "admin",
	})

	return nil
}
