package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/api/internal/model"
	"secret-im/service/signalserver/cmd/shared"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)
var (

)



type PutKeysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPutKeysLogic(ctx context.Context, svcCtx *svc.ServiceContext) PutKeysLogic {
	return PutKeysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutKeysLogic) PutKeys(adxName string,req types.PutKeysReqx) error {
	// todo: add your logic here and delete this line
	spk,err:=json.Marshal(req.SignedPreKey)
	if err!=nil{
		return shared.NewCodeError(shared.ERRCODE_JSONMARSHAL,err.Error())
	}
	signedPreKey:=string(spk)
	for i:=range req.PreKeys{
		key:=&model.TKeys{}
		key.Number=adxName
		key.DeviceId=1
		key.KeyId=req.PreKeys[i].KeyId
		key.PublicKey=req.PreKeys[i].PublicKey
		key.LastResort=0
		key.IdentityKey=req.IdentityKey
		key.SignedPrekey=signedPreKey

		_,err:=l.svcCtx.KeysModel.Insert(*key)
		if err!=nil{
			return shared.NewCodeError(shared.ERRCODE_SQLINSERT,err.Error())
		}

	}

	return nil
}
