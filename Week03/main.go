package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// 启动 http 服务
	g.Go(func() error {
		fmt.Println("http server start")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("http server stop")
				return ctx.Err()
			}
		}
	})
	// 启动第一个子任务
	g.Go(func() error {
		fmt.Println("task 1 start")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("task 1 stop")
				return ctx.Err()
			}
		}
	})

	// 启动第二个子任务等待信号
	g.Go(func() error {
		fmt.Println("task 2 start")
		// 等待中断信号
		quit := make(chan os.Signal)
		// 接收 syscall.SIGINT 和 syscall.SIGTERM 信号
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		return errors.New("Shuting down server")
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}
