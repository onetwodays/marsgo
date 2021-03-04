package logic

import (
	"context"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PostMsgsPendingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPostMsgsPendingLogic(ctx context.Context, svcCtx *svc.ServiceContext) PostMsgsPendingLogic {
	return PostMsgsPendingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostMsgsPendingLogic) PostMsgsPending(req types.PostMsgsPendingReq,number string) (*types.PostMsgsPendingRes, error) {
	// todo: add your logic here and delete this line
	_,account,err:=l.svcCtx.GetOneAccountByNumber(number)
	if err!=nil{
		return nil, err
	}
	device:=&account.Devices[0]
	resp,err:=l.svcCtx.MsgsModel.FindMany(req.Destination,device.Id,req.PageSize,req.PageIndex)
	if err!=nil{
		return nil, err
	}


	return &types.PostMsgsPendingRes{
		PageIndex: req.PageIndex,
		PageSize: req.PageSize,
		Total: len(resp), //select count(1) 比较耗时
		List: resp,
	}, nil
}
