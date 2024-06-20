//go:build wireinject

package wire

import (
	"Jayleonc/gateway/di"
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

// logSet 用来注入日志服务
var logSet = wire.NewSet(
	di.InitLogClient,
	di.InitLogSender,
	di.InitLogger,
)

// httpSet 用来注入 http 客户端
var httpSet = wire.NewSet(
	di.InitHTTPClient,
)

// retrySet 用来注入重试任务
var retrySet = wire.NewSet(
	di.InitRetryScheduler,
	// 需要执行的任务

	di.InitJobs,
)

// kafkaSet 用来注入 kafka 服务
var kafkaSet = wire.NewSet(
	di.InitKafkaSaramaClient,
	di.NewSyncProducer,
	di.RegisterConsumers,
)

func InitWebServer() *App {
	wire.Build(
		thirdPartySet, webServerSet, logSet,
		httpSet, retrySet, kafkaSet,

		// 依赖注入
		wire.Struct(new(App), "*"),
	)

	return new(App)
}
