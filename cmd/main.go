package main

import (
	"context"
	"github.com/zhangweijie11/zDiscovery/api"
	"github.com/zhangweijie11/zDiscovery/common"
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/services/registry"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化配置
	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Println("load config error:", err)
		return
	}
	// 初始化注册中心
	common.Discovery = registry.NewDiscovery(conf)
	router := api.InitRouter()
	srv := &http.Server{
		Addr:    conf.HttpServer,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%s\n", err)
		}
	}()

	// 重启
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit
	log.Println("shutdown discovery server...")

	// 注销
	common.Discovery.CancelSelf()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds")
	}
	log.Println("server exiting")
}
