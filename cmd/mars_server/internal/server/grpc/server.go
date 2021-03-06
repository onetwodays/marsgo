package grpc

import (
	pb "marsgo/cmd/mars_server/api"

	"marsgo/pkg/conf/paladin"
	"marsgo/pkg/net/rpc/warden"
)

// New new a grpc server.
func New(svc pb.DemoServer) (ws *warden.Server, err error) {
	var rc struct {
		Server *warden.ServerConfig
	}
	err = paladin.Get("grpc.toml").UnmarshalTOML(&rc)
	if err == paladin.ErrNotExist {
		err = nil
	}
	ws = warden.NewServer(rc.Server)
	pb.RegisterDemoServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}
