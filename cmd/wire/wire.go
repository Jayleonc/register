//go:build wireinject

package wire

import (
	"github.com/Jayleonc/register/di"
	"github.com/Jayleonc/register/internal/web"
	"github.com/google/wire"
)

// thirdPartySet 用来注入第三方依赖
var thirdPartySet = wire.NewSet(
	di.InitDB,
	di.InitRedis,
	//di.InitKafkaSaramaClient,
	//di.NewSyncProducer,
)

// webServerSet 用来注入 web 服务
var webServerSet = wire.NewSet(
	di.InitWebServer,
	di.InitGinMiddlewares,
)

// configCenterSet 用来注入配置中心
var configCenterSet = wire.NewSet(
	web.NewConfigHandler,
	di.InitEtcdClient,
	di.InitConfigClient,
)

func InitWebServer() *App {
	wire.Build(
		thirdPartySet, webServerSet,
		configCenterSet,
		// 依赖注入
		wire.Struct(new(App), "*"),
	)

	return new(App)
}
