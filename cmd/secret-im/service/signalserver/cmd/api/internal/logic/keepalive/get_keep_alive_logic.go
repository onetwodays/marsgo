package logic

import (
	"context"
	"github.com/prometheus/common/log"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"secret-im/service/signalserver/cmd/api/internal/svc"
	"secret-im/service/signalserver/cmd/api/internal/svc/push"
	"secret-im/service/signalserver/cmd/api/shared"
	"secret-im/service/signalserver/cmd/api/websocket"
)

type GetKeepAliveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetKeepAliveLogic {
	return GetKeepAliveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetKeepAliveLogic) GetKeepAlive(r *http.Request) error {
	value := r.Context().Value("ws")
	if value == nil {
		return nil
	}
	context := value.(*websocket.SessionContext)
	if context.Device == nil {
		return nil
	}

	device := context.Device
	address := push.Address{Number: device.Number, DeviceID: device.Device.ID}

	if !context.PubSubManager.HasLocalSubscription(address.Serialize()) {
		reason:="[Keepalive] no local subscription found for: "+address.Serialize()
		log.Warnf(reason )
		context.Session.Close(1000, "OK")
		return shared.Status(http.StatusGone,reason)
	}
	return nil
}
