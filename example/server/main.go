package main

import (
	"git.daochat.cn/service/registry/registry"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"time"
)

func main() {
	port := ":8080"
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
	}

	router := gin.Default()

	// 配置路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/logs/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{"logs": []string{"log1 for " + id, "log2 for " + id}})
	})

	router.POST("/logs", func(c *gin.Context) {
		var req struct {
			Log string `json:"log"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "success"})
	})

	// 创建服务器实例
	s := registry.NewServer(
		"user_service",
		registry.WithRegistry(client),
		registry.WithRegistryTimeout(10*time.Second),
		registry.WithHTTPServer(&http.Server{
			Handler: router,
			Addr:    port,
		}),
	)

	// 注册服务
	if err := s.Register(); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// 启动 HTTP 服务器
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
