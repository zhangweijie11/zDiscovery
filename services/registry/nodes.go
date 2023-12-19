package registry

import (
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/global"
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
