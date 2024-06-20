package di

import (
	"Jayleonc/gateway/internal/web"
	"Jayleonc/gateway/pkg/ginx"
	"Jayleonc/gateway/pkg/ginx/middleware"
	"Jayleonc/gateway/pkg/ginx/ratelimit"
	"Jayleonc/gateway/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, u *web.UserHandler) *ginx.Server {
	engine := gin.Default()
	engine.Use(mdls...) // 使用中间件

	// 注册路由
	u.RegisterRoutes(engine)

	server := &ginx.Server{
		Engine: engine,
		Addr:   viper.GetString("http.addr"),
	}

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
