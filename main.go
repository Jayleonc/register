package main

import (
	"Jayleonc/register/sdk"
	"context"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"time"
)

func main() {
	// 创建 etcd 客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}

	// 构建全局的 InterfaceBuilder
	globalInterfaceBuilder := sdk.NewInterfaceBuilder()
	// 使用 Gin 框架创建 HTTP 服务器
	router := gin.Default()

	// 配置路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 演示如何使用配置中心
	ctx := context.Background()
	configCenter := sdk.NewConfigCenter(client)
	err = configCenter.PutConfig(ctx, "example_key", "example_value")
	if err != nil {
		log.Fatalf("Failed to put config: %v", err)
	}
	value, err := configCenter.GetConfig(ctx, "example_key")
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}
	log.Printf("Config value: %s", value)

	// 监听配置变更
	go func() {
		watchCh, err := configCenter.WatchConfig(ctx, "example_key")
		if err != nil {
			log.Fatalf("Failed to watch config: %v", err)
		}
		for newValue := range watchCh {
			log.Printf("Config changed: %s = %s", "example_key", newValue)
		}
	}()

	// 创建服务器实例
	server := sdk.NewServer(
		"registry-service",
		sdk.WithRegistry(client),
		sdk.WithRegistryTimeout(10*time.Second),
		sdk.WithHTTPServer(&http.Server{
			Handler: router,
			Addr:    ":8080",
		}),
	)

	// 启动一个 goroutine 监听所有服务的变更事件
	go func() {
		subscribeCh, err := server.SubscribeAllServices()
		if err != nil {
			log.Fatalf("Failed to subscribe to service changes: %v", err)
		}
		for event := range subscribeCh {
			switch event.Type {
			case "REGISTER":
				log.Printf("Service registered: %s at %s", event.Instance.Name, event.Instance.Address)
			case "PUT":
				log.Printf("Service renewed: %s at %s", event.Instance.Name, event.Instance.Address)
			case "DELETE":
				log.Printf("Service down: %s at %s", event.Instance.Name, event.Instance.Address)
			}
		}
	}()

	// 启动 HTTP 服务器
	if err := server.Start(":8080", globalInterfaceBuilder); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
