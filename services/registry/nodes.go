package registry

import "github.com/zhangweijie11/zDiscovery/config"

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
