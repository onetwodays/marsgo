package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/api/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

type SetUserNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) SetUserNameLogic {
	return SetUserNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUserNameLogic) SetUserName(r *http.Request,req types.SetUserName) error {
	account,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return shared.Status(http.StatusUnauthorized,err.Error())
	}

	// 检查调用频率
	_,err=storage.UsernamesManager{}.Put(account.UUID,req.Username)
	if err!=nil{
		return shared.Status(http.StatusInternalServerError,err.Error())
	}

	return nil
}
