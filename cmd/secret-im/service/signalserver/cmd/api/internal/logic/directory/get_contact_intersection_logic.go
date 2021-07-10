package logic

import (
	"context"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/logic"
	"secret-im/service/signalserver/cmd/api/internal/storage"
	"secret-im/service/signalserver/cmd/api/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetContactIntersectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContactIntersectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetContactIntersectionLogic {
	return GetContactIntersectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContactIntersectionLogic) GetContactIntersection(r *http.Request,req types.ClientContactTokens) (*types.ClientContacts, error) {
	_,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}
	// todo limit

	// 返回联系人交集
	intersection, err := storage.DirectoryManager{}.Get(req.Contacts)
	if err!=nil{
		return nil,shared.Status(http.StatusNotFound,err.Error())
	}

	return &types.ClientContacts{
		Contacts: intersection,
	}, nil
}
