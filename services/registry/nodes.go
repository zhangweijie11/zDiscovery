package registry

import (
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/global"
	"log"
)

type Nodes struct {
	nodes    []*Node
	selfAddr string
}

// NewNodes 初始化所有注册中心服务实例
func NewNodes(conf *config.Config) *Nodes {
	nodes := make([]*Node, 0, len(conf.Nodes))
	// 将注册中心的节点加入节点列表
	for _, addr := range conf.Nodes {
		node := NewNode(conf, addr)
		nodes = append(nodes, node)
	}

	return &Nodes{
		nodes:    nodes,
		selfAddr: conf.HttpServer,
	}
}

func (nodes *Nodes) AllNodes() []*Node {
	result := make([]*Node, 0, len(nodes.nodes))
	for _, node := range nodes.nodes {
		n := &Node{
			addr:   node.addr,
			status: node.status,
		}
		result = append(result, n)
	}

	return result
}

// SetUp 节点上线
func (nodes *Nodes) SetUp() {
	for _, node := range nodes.nodes {
		// 节点上线
		if node.addr == nodes.selfAddr {
			node.status = global.NodeStatusUp
		}
	}
}

// Replicate 将当前节点同步至全部节点
func (nodes *Nodes) Replicate(action global.Action, instance *Instance) error {
	if len(nodes.nodes) == 0 {
		return nil
	}
	// 将实例依次注册到每个节点
	for _, node := range nodes.nodes {
		if node.addr != nodes.selfAddr {
			log.Printf("### 对节点地址 %v 进行操作,操作方法 %v ,当前主机名称 %v ###\n", node.addr, action, instance.Hostname)
			go nodes.action(node, action, instance)
		}
	}
	return nil
}

func (nodes *Nodes) action(node *Node, action global.Action, instance *Instance) {
	switch action {
	case global.Register:
		go node.Register(instance)
	case global.Renew:
		go node.Renew(instance)
	case global.Cancel:
		go node.Cancel(instance)
	}
}
