// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"excel2config/internal/dao"
	"excel2config/internal/server/grpc"
	"excel2config/internal/server/http"
	"excel2config/internal/server/ws"
	"excel2config/internal/service"
)

// Injectors from wire.go:

func InitApp() (*App, func(), error) {
	redis, cleanup, err := dao.NewRedis()
	if err != nil {
		return nil, nil, err
	}
	memcache, cleanup2, err := dao.NewMC()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	db, cleanup3, err := dao.NewDB()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	client, cleanup4, err := dao.NewMongo()
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	daoDao, cleanup5, err := dao.New(redis, memcache, db, client)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	serviceService, cleanup6, err := service.New(daoDao)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	engine, err := http.New(serviceService)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	server, err := grpc.New(serviceService)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	wsServer, err := ws.New(daoDao)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	app, cleanup7, err := NewApp(serviceService, engine, server, wsServer)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	return app, func() {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
