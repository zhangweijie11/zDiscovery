package registry

import (
	"encoding/json"
	"fmt"
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

// 同步全部节点的数据
func (dis *Discovery) initSync() {
	// 获取当前存储的所有节点
	nodes := dis.Nodes.Load().(*Nodes)
	for _, node := range nodes.AllNodes() {
		// 如果是本身节点，无需同步
		if node.addr == nodes.selfAddr {
			continue
		}
		// 请求其他节点的获取全部节点数据的接口
		url := fmt.Sprintf("http://%s%s", node.addr, global.FetchAllURL)
		response, err := utils.HttpPost(url, nil)
		if err != nil {
			log.Println("同步其余节点数据时出现异常", err)
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
		if resp.Code != global.StatusOK {
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

// 注册自身，也就是添加注册中心的节点
func (dis *Discovery) regSelf() *Instance {
	now := time.Now().UnixNano()
	instance := &Instance{
		Env:             dis.config.Env,
		AppID:           global.DiscoveryAppId,
		Hostname:        dis.config.Hostname,
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
			log.Println("### 每 30 秒续约一次节点 ###")
			_, err := dis.Registry.Renew(instance.Env, instance.AppID, instance.Hostname)
			// 续约失败，表示当前节点不存在，需要重新注册
			if err == utils.NotFound {
				dis.Registry.Register(instance, now)
				dis.Nodes.Load().(*Nodes).Replicate(global.Register, instance)
			} else {
				dis.Nodes.Load().(*Nodes).Replicate(global.Renew, instance)
			}
		}
	}
}

// 更新节点列表
func (dis *Discovery) nodesPerception() {
	var lastTimestamp int64
	ticker := time.NewTicker(global.NodePerceptionInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Printf("### 定期自发现子节点， 一共发现节点, %v 个 ###\n", len(dis.Nodes.Load().(*Nodes).AllNodes()))
			fetchData, err := dis.Registry.Fetch(dis.config.Env, global.DiscoveryAppId, global.NodeStatusUp, lastTimestamp)
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
			// 设置更新时间，后续只会获取上线的节点
			lastTimestamp = fetchData.LatestTimestamp

			// 更新节点列表
			config := new(config.Config)
			*config = *dis.config
			config.Nodes = nodes

			ns := NewNodes(config)
			ns.SetUp()
			dis.Nodes.Store(ns)
			log.Printf("### 更新节点列表,数量变为 %v 个 ###\n", len(dis.Nodes.Load().(*Nodes).AllNodes()))
		}
	}
}

// 退出保护模式
func (dis *Discovery) exitProtect() {
	time.Sleep(global.ProtectTimeInterval)
	dis.protected = false
	log.Println("### 每 60 秒检测保护模式 ###")
}

func (dis *Discovery) CancelSelf() {
	log.Println("### 退出时注销本身节点 ###")
	dis.Registry.Cancel(dis.config.Env, global.DiscoveryAppId, dis.config.Hostname, time.Now().UnixNano())
	instance := &Instance{
		Env:      dis.config.Env,
		Hostname: dis.config.Hostname,
		AppID:    global.DiscoveryAppId,
	}
	dis.Nodes.Load().(*Nodes).Replicate(global.Cancel, instance) //broadcast
}

// NewDiscovery 初始化注册中心服务
func NewDiscovery(conf *config.Config) *Discovery {
	discovery := &Discovery{
		config:    conf,
		protected: false,
		Registry:  NewRegistry(),
	}

	// 初始化节点，一个节点就是一个独立的注册中心服务，可以在多个 goroutine 中安全地共享和更新变量
	discovery.Nodes.Store(NewNodes(conf))
	// 从其他节点同步数据
	discovery.initSync()
	// 注册当前的注册中心节点
	instance := discovery.regSelf()
	// 续约注册中心
	go discovery.renewTask(instance)
	// 更新节点列表
	go discovery.nodesPerception()
	// 退出保护模式
	go discovery.exitProtect()

	return discovery
}
