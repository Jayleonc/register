package command

import (
	"github.com/Jayleonc/register/cmd/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	Flags = GlobalFlags{}
)

type GlobalFlags struct {
	configPath string
}

func NewWebCommand() *cobra.Command {
	w := &cobra.Command{
		Use:   "web",
		Short: "web server start.",
		Run:   run,
	}
	w.PersistentFlags().StringVarP(&Flags.configPath, "config", "c", "config/dev.yaml", "config file")
	return w
}

func run(cmd *cobra.Command, args []string) {
	runApp()
}

func runApp() {
	initConfig()

	app := wire.InitWebServer()

	// 启动 Web 服务器并等待其退出
	if err := app.Web.Start(); err != nil {
		log.Fatalf("Web server failed to start: %v", err)
	}

	// 服务器已优雅退出，现在关闭其他服务
	log.Println("Shutting down other services...")

	log.Println("All services stopped, exiting application.")
}

// todo 未来使用 etcd 配置中心
func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(Flags.configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
