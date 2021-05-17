package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"
	"secret-im/service/signalserver/cmd/shared"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetKeysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetKeysLogic {
	return GetKeysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKeysLogic) GetKeys(req types.GetKeysReq) (*types.GetKeysResx, error) {
	// todo: add your logic here and delete this line
	keys,err:=l.svcCtx.KeysModel.FindMany(req.Identifier,req.DeviceId)
	if err!=nil{
		return nil, shared.NewCodeError(shared.ERRCODE_SQLQUERY,err.Error())
	}
	if len(keys)>0{
		devices:=make([]types.PreKeyResponseItem,len(keys))
		for i:=range keys{
			devices[i].DeviceId=1
			devices[i].RegistrationId=1
			devices[i].PreKey=types.PreKeyx{
				KeyId: keys[i].KeyId,
				PublicKey: keys[i].PublicKey,
			}


			err=json.Unmarshal([]byte(keys[i].SignedPrekey),&devices[i].SignedPrekey) // 为什么这里要传指针呢？
			if err!=nil{
				return nil, shared.NewCodeError(shared.ERRCODE_JSONUNMARSHAL,err.Error())
			}

		}
		return &types.GetKeysResx{
			IdentityKey: keys[0].IdentityKey,
			Devices: devices,

		}, nil
	}


	return &types.GetKeysResx{

	}, nil
}
