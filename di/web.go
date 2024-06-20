package di

import (
	"github.com/Jayleonc/register/pkg/ginx"
	"github.com/Jayleonc/register/pkg/ginx/middleware"
	"github.com/Jayleonc/register/pkg/ginx/ratelimit"
	"github.com/Jayleonc/register/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc) *ginx.Server {
	engine := gin.Default()
	engine.Use(mdls...) // 使用中间件

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
