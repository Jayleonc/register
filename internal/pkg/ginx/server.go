package ginx

import (
	"context"
	"errors"
	"fmt"
	"github.com/Jayleonc/register/pkg/netx"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	*gin.Engine
	Addr string
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    s.Addr,
		Handler: s.Engine,
	}
	// ============================================================================
	// 启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// ============================================================================
	// 打印启动信息
	tip(s)

	// ============================================================================
	// 创建一个通道来接收系统信号，并优雅关停服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 创建一个超时上下文
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
	fmt.Println(" ___  _  ______ _____ ___  ____  / /__        ________  ____  ____ \n / _ \\| |/_/ __ `/ __ `__ \\/ __ \\/ / _ \\______/ ___/ _ \\/ __ \\/ __ \\\n/  __/>  </ /_/ / / / / / / /_/ / /  __/_____/ /  /  __/ /_/ / /_/ /\n\\___/_/|_|\\__,_/_/ /_/ /_/ .___/_/\\___/     /_/   \\___/ .___/\\____/ \n                        /_/                          /_/           ")

	fmt.Printf("Listening and serving HTTP on http://%s%s\n", netx.GetOutboundIP(), s.Addr)
}
