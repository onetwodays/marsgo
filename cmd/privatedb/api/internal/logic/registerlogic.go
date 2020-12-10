package logic

import (
	"context"
	"privatedb/api/model"

	"privatedb/api/internal/svc"
	"privatedb/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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
	_,err:=l.svcCtx.UserModer.FindOneByName(req.Username)
	if err==nil{
		return errorDuplicateUsername
	}
	_,err = l.svcCtx.UserModer.FindOneByMobile(req.Mobile)
	if err==nil{
		return errorDuplicateMobile
	}
	_,err = l.svcCtx.UserModer.Insert(model.User{
		Name: req.Username,
		Password: req.Password,
		Mobile: req.Mobile,
		Gender: "男",
		Nickname: "admin",

	})

	return nil
}
