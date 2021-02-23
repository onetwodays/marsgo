package logic

import (
	"context"
	"strconv"

	"secret-im/model"
	"secret-im/shared"
	"secret-im/signalserver/internal/svc"
	"secret-im/signalserver/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

var errorUserNotFound = shared.NewCodeError(shared.UserNotFound,shared.CodeErrorMap[shared.UserNotFound])

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserInfoLogic {
	return UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(userId string) (*types.UserReply, error) {
	// todo: add your logic here and delete this line
	userInt,err:=strconv.ParseInt(userId,10,64)
	if err!=nil{
		return nil, err
	}
	userInfo,err:=l.svcCtx.UserModel.FindOne(userInt)
	switch err {
	case nil:
		return &types.UserReply{
			Id:       userInfo.Id,
			Username: userInfo.Name,
			Mobile:   userInfo.Mobile,
			Nickname: userInfo.Nickname,
			Gender:   userInfo.Gender,
		}, nil
	case model.ErrNotFound:
		return nil, errorUserNotFound
	default:
		return nil, err

	}


}
