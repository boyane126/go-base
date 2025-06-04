package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boyane126/go-common/logger"

	"gobase/bootstrap"
)

func init() {
	bootstrap.Bootstrap()
}

func main() {
	// 处理退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.InfoString("main", "服务关闭", "接收到退出信号，服务正在关闭...")
		os.Exit(0)
	}()

	go func() {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println("hello world")
		}
	}()

	select {}
}
