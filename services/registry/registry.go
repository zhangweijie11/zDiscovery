package registry

import (
	"fmt"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"sync"
)

type Registry struct {
	apps  map[string]*Application
	lock  sync.RWMutex
	guard *Guard
}

func NewRegistry() *Registry {
	registry := &Registry{
		apps:  make(map[string]*Application),
		lock:  sync.RWMutex{},
		guard: new(Guard),
	}

	return registry
}

// Register 实例注册
func (r *Registry) Register(instance *Instance, latestTimestamp int64) (*Application, *utils.Error) {
	// 注册中心应用服务唯一标识
	key := fmt.Sprintf("%s-%s", instance.AppID, instance.Env)
	r.lock.RLock()
	app, ok := r.apps[key]
	// 如果注册中心不存在该应用服务则添加
	if !ok {
		app = NewApplication(instance.AppID)
	}
	// 添加应用服务的实例
	_, isNew := app.AddInstance(instance, latestTimestamp)
	if isNew {
		// 重新计算一个周期需要的续约次数和保护模式阈值
		r.guard.incrNeed()
	}

	r.lock.Lock()
	r.apps[key] = app
	r.lock.Unlock()

	return app, nil
}
