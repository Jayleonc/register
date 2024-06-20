package main

import (
	"Jayleonc/register/internal/core/resolver"
	"Jayleonc/register/sdk"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"time"
)

func main() {
	port := ":8081"
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}

	router := gin.Default()

	// 构建全局的 InterfaceBuilder
	globalInterfaceBuilder := sdk.NewInterfaceBuilder()

	// 配置路由，使用 HandlerBuilder 添加参数和返回值信息
	router.GET("/health", sdk.NewHandlerBuilder(globalInterfaceBuilder, "GET", "/health").
		AddReturn("status", "string").
		Build(func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		}),
	)

	router.GET("/logs/:id", sdk.NewHandlerBuilder(globalInterfaceBuilder, "GET", "/logs/:id").
		AddParam("id", "string").
		AddReturn("logs", "array").
		Build(func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(200, gin.H{
				"logs": []string{"log1 for " + id, "log2 for " + id},
			})
		}),
	)

	router.POST("/logs", sdk.NewHandlerBuilder(globalInterfaceBuilder, "POST", "/logs").
		AddReturn("status", "string").
		Build(func(c *gin.Context) {
			var req struct {
				Log string `json:"log"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "success"})
		}),
	)

	s := sdk.NewServer(
		"user_service",
		sdk.WithRegistry(client),
		sdk.WithRegistryTimeout(10*time.Second),
		sdk.WithHTTPServer(&http.Server{
			Handler: router,
			Addr:    port,
		}),
		sdk.WithResolver(resolver.NewEtcdResolver(client)),
	)

	if err := s.Start(port, globalInterfaceBuilder); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
