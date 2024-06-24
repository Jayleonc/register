package ginx

import (
	"context"
	"errors"
	"fmt"
	"github.com/Jayleonc/register/internal/pkg/netx"
	"github.com/Jayleonc/register/registry"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*gin.Engine
	Addr   string
	Client *clientv3.Client
}

func (s *Server) Start() error {
	if s.Engine == nil {
		s.Engine = gin.Default()
	}

	// 配置路由
	s.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	srv := &http.Server{
		Addr:    s.Addr,
		Handler: s.Engine,
	}

	server, err := registrySelf(s.Client, srv)
	defer server.Close()
	if err != nil {
		return err
	}

	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 打印启动信息
	tip(s)

	// 创建一个通道来接收系统信号，并优雅关停服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting...")
	return nil
}

func tip(s *Server) {
	// https://patorjk.com/software/taag/#p=display&f=Slant&t=Composer
	fmt.Println(" ___  _  ______ _____ ___  ____  / /__        ________  ____  ____ ")
	fmt.Println("/ _ \\| |/_/ __ `/ __ `__ \\/ __ \\/ / _ \\______/ ___/ _ \\/ __ \\/ __ \\")
	fmt.Println("/  __/>  </ /_/ / / / / / / /_/ / /  __/_____/ /  /  __/ /_/ / /_/ /")
	fmt.Println("\\___/_/|_|\\__,_/_/ /_/ /_/ .___/_/\\___/     /_/   \\___/ .___/\\____/ ")
	fmt.Println("                        /_/                          /_/           ")

	fmt.Printf("Listening and serving HTTP on http://%s%s\n", netx.GetOutboundIP(), s.Addr)
}

func registrySelf(client *clientv3.Client, srv *http.Server) (*registry.Server, error) {
	// 创建服务器实例
	server := registry.NewServer(
		"registry",
		registry.MustWithAddress("localhost:8080"),
		registry.WithRegistry(client),
		registry.WithRegistryTimeout(10*time.Second),
		registry.WithHTTPServer(srv),
	)

	// 注册服务
	if err := server.Register(); err != nil {
		return nil, fmt.Errorf("failed to register service: %v", err)
	}

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

	return server, nil
}
