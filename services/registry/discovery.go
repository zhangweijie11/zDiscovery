package registry

import (
	"encoding/json"
	"fmt"
	"github.com/skyhackvip/service_discovery/configs"
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"log"
	"sync/atomic"
)

type Discovery struct {
	config    *config.Config
	protected bool
	Registry  *Registry
	Nodes     atomic.Value
}

func (dis *Discovery) initSync() {
	nodes := dis.Nodes.Load().(*Nodes)
	for _, node := range nodes.AllNodes() {
		// 本身节点，无需同步
		if node.addr == nodes.selfAddr {
			continue
		}
		url := fmt.Sprintf("http://%s%s", node.addr, global.FetchAllURL)
		response, err := utils.HttpPost(url, nil)
		if err != nil {
			log.Println(err)
			continue
		}
		var resp struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string][]*Instance `json:"data"`
		}
		err = json.Unmarshal([]byte(response), &resp)
		if err != nil {
			log.Printf("get from %v error : %v", url, err)
			continue
		}
		if resp.Code != configs.StatusOK {
			log.Printf("get from %v error : %v", url, resp.Message)
			continue
		}
		dis.protected = false
		for _, v := range resp.Data {
			for _, instance := range v {
				dis.Registry.Register(instance, instance.LatestTimestamp)
			}
		}
	}

	nodes.SetUp()
}

func NewDiscovery(config *config.Config) *Discovery {
	discovery := &Discovery{
		config:    config,
		protected: false,
		Registry:  NewRegistry(),
	}

	// 初始化节点，一个节点就是一个独立的注册中心服务
	discovery.Nodes.Store(NewNodes(config))
	// 从其他节点同步数据
	discovery.initSync()

	return discovery
}
