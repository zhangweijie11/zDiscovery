package registry

import (
	"encoding/json"
	"fmt"
	"github.com/skyhackvip/service_discovery/configs"
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"log"
	"net/url"
	"sync/atomic"
	"time"
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

func (dis *Discovery) regSelf() *Instance {
	now := time.Now().UnixNano()
	instance := &Instance{
		Env:             dis.config.Env,
		AppID:           global.DiscoveryAppId,
		Hostname:        dis.config.HostName,
		Addresses:       []string{"http://" + dis.config.HttpServer},
		Version:         "",
		Status:          global.NodeStatusUp,
		RegTimestamp:    now,
		UpTimestamp:     now,
		RenewTimestamp:  now,
		DirtyTimestamp:  now,
		LatestTimestamp: now,
	}

	dis.Registry.Register(instance, now)
	dis.Nodes.Load().(*Nodes).Replicate(global.Register, instance)

	return instance
}

func (dis *Discovery) renewTask(instance *Instance) {
	now := time.Now().UnixNano()
	ticker := time.NewTicker(global.RenewInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println("### discovery node renew every 30s ###")
			_, err := dis.Registry.Renew(instance.Env, instance.AppID, instance.Hostname)
			if err == utils.NotFound {
				dis.Registry.Register(instance, now)
				dis.Nodes.Load().(*Nodes).Replicate(global.Register, instance)
			} else {
				dis.Nodes.Load().(*Nodes).Replicate(global.Renew, instance)
			}

		}
	}
}

// update discovery nodes list
func (dis *Discovery) nodesPerception() {
	var lastTimestamp int64
	ticker := time.NewTicker(global.NodePerceptionInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println("### discovery node protect tick ###")
			log.Printf("### discovery nodes,len (%v) ###\n", len(dis.Nodes.Load().(*Nodes).AllNodes()))
			fetchData, err := dis.Registry.Fetch(dis.config.Env, configs.DiscoveryAppId, configs.NodeStatusUp, lastTimestamp)
			if err != nil || fetchData == nil {
				continue
			}
			var nodes []string
			for _, instance := range fetchData.Instances {
				for _, addr := range instance.Addresses {
					u, err := url.Parse(addr)
					if err == nil {
						nodes = append(nodes, u.Host)
					}
				}
			}
			lastTimestamp = fetchData.LatestTimestamp

			// 更新节点列表
			config := new(config.Config)
			*config = *dis.config
			config.Nodes = nodes

			ns := NewNodes(config)
			ns.SetUp()
			dis.Nodes.Store(ns)
			log.Printf("### discovery protect change nodes,len (%v) ###\n", len(dis.Nodes.Load().(*Nodes).AllNodes()))
		}
	}
}

func NewDiscovery(conf *config.Config) *Discovery {
	discovery := &Discovery{
		config:    conf,
		protected: false,
		Registry:  NewRegistry(),
	}

	// 初始化节点，一个节点就是一个独立的注册中心服务
	discovery.Nodes.Store(NewNodes(conf))
	// 从其他节点同步数据
	discovery.initSync()
	// 注册当前的注册中心节点
	instance := discovery.regSelf()
	// 续约注册中心
	go discovery.renewTask(instance)
	// 更新节点列表
	go discovery.nodesPerception()

	return discovery
}
