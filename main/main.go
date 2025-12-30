package main

import (
	"context"
	"crypto_alert/base"
	"crypto_alert/calculate"
	"crypto_alert/config"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	config.Init()

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	// 捕获 Ctrl+C 信号
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		fmt.Println("收到退出信息， 退出中...")
		cancle()
	}()
	if len(config.Cfg.Cycles) > 0 {
		fmt.Println("First cycle:", config.Cfg.Cycles[0])
	} else {
		fmt.Println("No cycles configured")
	}

	base.FetchBinanceSymbols()
	for _, c := range config.Cfg.Cycles {
		go calculate.MacdTicker(ctx, c.Cycle)
	}
	<-ctx.Done()
}
