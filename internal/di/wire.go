// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"excel2config/internal/dao"
	"excel2config/internal/server/grpc"
	"excel2config/internal/server/http"
	"excel2config/internal/server/ws"
	"excel2config/internal/service"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, ws.New, NewApp))
}
