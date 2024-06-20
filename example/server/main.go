package main

import (
	"errors"
	"github.com/Jayleonc/register/registry"
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

	// 使用 Gin 框架创建 HTTP 服务器
	engine := gin.Default()

	// 构建全局的 App
	apiDescriptor := registry.NewApiDescriptor(engine)

	// 配置路由，使用简化的 Register 方法
	apiDescriptor.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}, nil, []registry.Return{{Name: "status", Type: "string"}})

	apiDescriptor.GET("/logs/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{
			"logs": []string{"log1 for " + id, "log2 for " + id},
		})
	}, []registry.Param{{Name: "id", Type: "string"}}, []registry.Return{{Name: "logs", Type: "array of strings"}})

	apiDescriptor.POST("/logs", func(c *gin.Context) {
		var req struct {
			Log string `json:"log"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "success"})
	}, nil, []registry.Return{{Name: "status", Type: "string"}})

	srv := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	s := registry.NewServer(
		"user_service",
		registry.WithRegistry(client),
		registry.WithRegistryTimeout(10*time.Second),
		registry.WithHTTPServer(srv),
	)

	if err := s.Register(apiDescriptor); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	defer s.Close()

	// 启动Gin服务器
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
