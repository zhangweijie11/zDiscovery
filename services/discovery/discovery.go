package discovery

import (
	"github.com/skyhackvip/service_discovery/model"
	"github.com/zhangweijie11/zDiscovery/config"
	"sync/atomic"
)

type Discovery struct {
	config    *config.Config
	protected bool
	Registry  *model.Registry
	Nodes     atomic.Value
}
