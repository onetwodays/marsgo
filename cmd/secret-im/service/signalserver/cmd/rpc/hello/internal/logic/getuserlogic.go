package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/rpc/hello/hello_rpc"
	"secret-im/service/signalserver/cmd/rpc/hello/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *hello_rpc.IdReq) (*hello_rpc.UserInfoReply, error) {
	// todo: add your logic here and delete this line
	one,err := l.svcCtx.UserModel.FindOne(in.Id)
	if err !=nil{
		return nil, err
	}

	return &hello_rpc.UserInfoReply{
		Id:one.Id,
		Name:one.Name,
		Number: one.Mobile,
		Gender: one.Gender,

	}, nil
}
