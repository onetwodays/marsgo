package logic

import (
	"context"
	"encoding/json"
	"secret-im/service/signalserver/cmd/api/db/redis"
	"secret-im/service/signalserver/cmd/api/util"

	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDirTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDirTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDirTokenLogic {
	return GetDirTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDirTokenLogic) GetDirToken(req types.GetDirTokenReq) (*types.GetDirTokenRes, error) {
	// todo: add your logic here and delete this line
	bt,err:=util.DecodeToken(req.Token)
	if err!=nil{
		return nil, err
	}
	value,err:= redis.RedisDirectoryManager().HGet("directory",string(bt[:])).Result()
	if err!=nil{
		return nil, err
	}
	ds:=types.GetDirTokenRes{}
	err=json.Unmarshal([]byte(value),&ds)
	if err!=nil{
		return nil, err
	}
	return &ds, nil
}
