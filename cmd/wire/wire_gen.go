// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/Jayleonc/register/di"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := di.InitRedis()
	v := di.InitGinMiddlewares(cmdable)
	server := di.InitWebServer(v)
	app := &App{
		Web: server,
	}
	return app
}

// wire.go:

// thirdPartySet 用来注入第三方依赖
var thirdPartySet = wire.NewSet(di.InitDB, di.InitRedis)

// webServerSet 用来注入 web 服务
var webServerSet = wire.NewSet(di.InitWebServer, di.InitGinMiddlewares)
