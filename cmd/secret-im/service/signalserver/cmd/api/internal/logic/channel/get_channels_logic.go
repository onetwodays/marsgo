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

type GetChannelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetChannelsLogic {
	return GetChannelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChannelsLogic) GetChannels(r *http.Request,req types.ChannelsParams) (*types.Channels, error) {
	currAccount,err:= logic.GetSourceAccount(r,l.svcCtx.AccountsModel)
	if err!=nil{
		return nil,shared.Status(http.StatusUnauthorized,err.Error())
	}
	if req.Limit<=0{
		req.Limit=10
	}
	//
	channels:=make([]types.Channel,0)
	deactivatedChannels := make([]string, 0)
	for {
		if len(channels)>=req.Limit{
			break
		}
		maxID:=req.MaxID
		if len(channels)>0{
			maxID=&channels[len(channels)-1].ID
		}
		result,err:=storage.ChannelJoinedDao{}.Filter(currAccount.UUID,maxID,req.Limit)
		if err!=nil{
			break
		}
		if len(result)==0{
			break
		}

		// 获取频道信息
		channelIdsSet := make([]string, 0, len(result))
		for _, item := range result {
			channelIdsSet = append(channelIdsSet, item.ChannelID)
		}
		channelSlice, err := storage.Channels{}.GetList(channelIdsSet)
		if err != nil {

			return nil,shared.Status(http.StatusInternalServerError,err.Error())
		}
		channelMapper := make(map[string]storage.Channel)
		for _, item := range channelSlice {
			channelMapper[item.ChannelID] = item
		}

		// 生成频道列表
		for _, item := range result {
			channel, ok := channelMapper[item.ChannelID]
			if !ok {
				continue
			}
			if channel.Deactivated {
				deactivatedChannels = append(deactivatedChannels, channel.ChannelID)
				continue
			}
			channels = append(channels, *newChannelEntity(channel, true, item.Left, item.Kicked))
		}
	}
	if len(deactivatedChannels)>0{
		storage.ChannelJoinedDao{}.UpdateDeactivated(currAccount.UUID,deactivatedChannels)
	}

	// 填充频道消息属性
	err = fillChannelMessageAttributes(channels, currAccount.UUID)
	if err != nil {
		return nil,shared.Status(http.StatusInternalServerError,err.Error())
	}
	return &types.Channels{Channels: channels}, nil
}
