package registry

import (
	"fmt"
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/global"
)

type Node struct {
	config      *config.Config // 节点配置
	addr        string         // 节点地址
	status      int            // 节点状态
	registerUrl string         // 注册地址
	cancelUrl   string         // 注销地址
	renewUrl    string         // 续约地址
}

func NewNode(config *config.Config, addr string) *Node {
	return &Node{
		config:      config,
		addr:        addr,
		status:      global.NodeStatusDown,
		registerUrl: fmt.Sprintf("http://%s%s", addr, global.RegisterURL),
		cancelUrl:   fmt.Sprintf("http://%s%s", addr, global.CancelURL),
		renewUrl:    fmt.Sprintf("http://%s%s", addr, global.RenewURL),
	}
}
