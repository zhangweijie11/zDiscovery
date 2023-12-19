package main

import (
	"fmt"
	"github.com/skyhackvip/service_discovery/configs"
	"log"
)

func main() {
	// 初始化配置
	config, err := configs.LoadConfig("config.yaml")
	if err != nil {
		log.Println("load config error:", err)
		return
	}
	fmt.Println("------------>", config)
}
