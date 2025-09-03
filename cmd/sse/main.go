package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/boyane126/go-common/logger"

	"gobase/bootstrap"
	"gobase/internal/sse"
)

func init() {
	bootstrap.Bootstrap()

	// 初始化Redis和SSE服务
	bootstrap.SetupRedis()
}

func main() {
	if err := sse.InitSSEService(); err != nil {
		panic("初始化SSE服务失败: " + err.Error())
	}

	// 处理退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.InfoString("main", "服务关闭", "接收到退出信号，SSE服务正在关闭...")
		os.Exit(0)
	}()

	// 启动SSE服务器（阻塞运行）
	sse.StartSSEServer()
}
