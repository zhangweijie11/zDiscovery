package main

import (
	"github.com/zhangweijie11/zDiscovery/api"
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/services/registry"
	"log"
	"net/http"
)

func main() {
	// 初始化配置
	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Println("load config error:", err)
		return
	}
	// 初始化注册中心
	global.Discovery = registry.NewDiscovery(conf)

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
}
