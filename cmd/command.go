package main

import (
	"fmt"
	"github.com/bincooo/chatgpt-adapter/internal/cache"
	"github.com/bincooo/chatgpt-adapter/internal/common"
	"github.com/bincooo/chatgpt-adapter/internal/gin.handler"
	"github.com/bincooo/chatgpt-adapter/logger"
	"github.com/bincooo/chatgpt-adapter/pkg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version  = "v2.1.0"
	proxies  string
	port     int
	logLevel = "info"
	vms      bool

	cmd = &cobra.Command{
		Use:   "ChatGPT-Adapter",
		Short: "GPT接口适配器",
		Long: "GPT接口适配器。统一适配接口规范，集成了bing、claude-2，gemini...\n" +
			"项目地址：https://github.com/bincooo/chatgpt-adapter",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			if vms {
				fmt.Println("模型可用列表:")
				for _, model := range handler.GlobalExtension.Models() {
					fmt.Println("- " + model.Id)
				}
				return
			}
			handler.Bind(port, version, proxies)
		},
	}
)

func main() {
	pkg.Init()
	cache.Init()
	common.Init()
	handler.InitExtensions()
	cmd.PersistentFlags().StringVar(&proxies, "proxies", "", "本地代理 proxies")
	cmd.PersistentFlags().IntVar(&port, "port", 8080, "服务端口 port")
	cmd.PersistentFlags().StringVar(&logLevel, "log", logLevel, "日志级别: debug|info|warn|error")
	cmd.PersistentFlags().BoolVar(&vms, "models", false, "查看所有模型")
	logger.Init(switchLogLevel())
	_ = cmd.Execute()
}

func switchLogLevel() logrus.Level {
	switch logLevel {
	case "debug":
		return logrus.DebugLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}
