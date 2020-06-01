// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	pb "marsgo/cmd/mars_server/api"
	"marsgo/cmd/mars_server/internal/dao"
	"marsgo/cmd/mars_server/internal/server/grpc"
	"marsgo/cmd/mars_server/internal/server/http"
	"marsgo/cmd/mars_server/internal/service"

	"github.com/google/wire"
)

var daoProvider = wire.NewSet(dao.New, dao.NewDB, dao.NewRedis, dao.NewMC)
var serviceProvider = wire.NewSet(service.New, wire.Bind(new(pb.DemoServer), new(*service.Service)))

func InitApp() (*App, func(), error) {
	panic(wire.Build(daoProvider, serviceProvider, http.New, grpc.New, NewApp))
}
