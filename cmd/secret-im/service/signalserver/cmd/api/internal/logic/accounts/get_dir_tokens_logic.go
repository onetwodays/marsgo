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

type GetDirTokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDirTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDirTokensLogic {
	return GetDirTokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDirTokensLogic) GetDirTokens(req types.GetDirTokensReq) (*types.GetDirTokensRes, error) {
	// todo: add your logic here and delete this line
	list :=make([]types.GetDirTokenRes,len(req.Tokens),len(req.Tokens))
	for i:=range req.Tokens{
		bt,err:=util.DecodeToken(req.Tokens[i])
		if err!=nil{
			continue
		}
		value,err:= redis.RedisDirectoryManager().HGet("directory",string(bt[:])).Result()
		if err!=nil{
			continue
		}
		ds:=types.GetDirTokenRes{}
		err=json.Unmarshal([]byte(value),&ds)
		if err!=nil{
			continue
		}
		list=append(list,ds)
	}



	return &types.GetDirTokensRes{List: list}, nil
}
