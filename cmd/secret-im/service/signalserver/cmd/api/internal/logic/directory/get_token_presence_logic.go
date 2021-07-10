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

type GetTokenPresenceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenPresenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTokenPresenceLogic {
	return GetTokenPresenceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenPresenceLogic) GetTokenPresence(r *http.Request,req types.GetTokenPresenceReq) (*types.ClientContact, error) {
	_,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}
	//todo limit
	contacts,err:=storage.DirectoryManager{}.Get([]string {req.Token})
	if err!=nil || len(contacts)==0{
		var reason string
		if err!=nil{
			reason=err.Error()
		}
		if len(contacts)==0{
			reason+=" no data"
		}
		return nil,shared.Status(http.StatusNotFound,reason)
	}
	return &contacts[0], nil
}
