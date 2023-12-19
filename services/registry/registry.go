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
func (r *Registry) getApplication(appid, env string) (*Application, bool) {
	key := fmt.Sprintf("%s-%s", appid, env)
	r.lock.RLock()
	app, ok := r.apps[key]
	r.lock.RUnlock()
	return app, ok
}
func (r *Registry) Renew(env, appid, hostname string) (*Instance, *utils.Error) {
	// 获取应用服务
	app, ok := r.getApplication(appid, env)
	if !ok {
		return nil, utils.NotFound
	}
	// 修改实例的最后更新时间
	in, ok := app.Renew(hostname)
	if !ok {
		return nil, utils.NotFound
	}
	r.guard.incrCount()
	return in, nil
}

// Fetch 通过筛选条件获取应用服务的全部实例
func (r *Registry) Fetch(env, appid string, status uint32, latestTime int64) (*FetchData, *utils.Error) {
	app, ok := r.getApplication(appid, env)
	if !ok {
		return nil, utils.NotFound
	}
	return app.GetInstance(status, latestTime) //err = not modify
}

// Cancel 注销应用服务
func (r *Registry) Cancel(env, appid, hostname string, latestTimestamp int64) (*Instance, *utils.Error) {
	// 获取应用服务
	app, ok := r.getApplication(appid, env)
	if !ok {
		return nil, utils.NotFound
	}
	instance, ok, insLen := app.Cancel(hostname, latestTimestamp)
	if !ok {
		return nil, utils.NotFound
	}
	// 如果实例列表为空，则删除该应用服务
	if insLen == 0 {
		r.lock.Lock()
		delete(r.apps, fmt.Sprintf("%s-%s", appid, env))
		r.lock.Unlock()
	}
	r.guard.decrNeed()
	return instance, nil
}

func (r *Registry) getAllApplications() []*Application {
	r.lock.RLock()
	defer r.lock.RUnlock()
	apps := make([]*Application, 0, len(r.apps))
	for _, app := range r.apps {
		apps = append(apps, app)
	}
	return apps
}

// FetchAll 获取全部应用服务及其实例
func (r *Registry) FetchAll() map[string][]*Instance {
	apps := r.getAllApplications()
	rs := make(map[string][]*Instance)
	for _, app := range apps {
		rs[app.appid] = append(rs[app.appid], app.GetAllInstances()...)
	}
	return rs
}
