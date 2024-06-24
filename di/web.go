package di

import (
	"github.com/Jayleonc/register/internal/pkg/ginx"
	"github.com/Jayleonc/register/internal/pkg/ginx/middleware"
	"github.com/Jayleonc/register/internal/pkg/ginx/ratelimit"
	"github.com/Jayleonc/register/internal/pkg/limiter"
	"github.com/Jayleonc/register/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, etcdClient *clientv3.Client, configHandler *web.ConfigHandler) *ginx.Server {
	engine := gin.Default()
	engine.Use(mdls...) // 使用中间件

	configHandler.RegisterRoutes(engine)

	server := &ginx.Server{
		Engine: engine,
		Addr:   viper.GetString("http.addr"),
		Client: etcdClient,
	}

	server.Engine.Static("/static", "./static")

	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
		}),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		middleware.RecoveryWithLogger(),
	}
}
