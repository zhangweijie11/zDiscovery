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

func NewNodes(conf *config.Config) *Nodes {
	nodes := make([]*Node, 0, len(conf.Nodes))
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

func (nodes *Nodes) SetUp() {
	for _, node := range nodes.nodes {
		// 节点上线
		if node.addr == nodes.selfAddr {
			node.status = global.NodeStatusUp
		}
	}
}

// Replicate 节点复制（节点本身也作为注册中心服务的实例）
func (nodes *Nodes) Replicate(action global.Action, instance *Instance) error {
	if len(nodes.nodes) == 0 {
		return nil
	}
	// 将实例依次注册到每个节点
	for _, node := range nodes.nodes {
		if node.addr != nodes.selfAddr {
			log.Printf("### replicate node(%v),action(%v),hostname(%v) ###\n", node.addr, action, instance.Hostname)
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
