package registry

import (
	"fmt"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Registry struct {
	apps  map[string]*Application
	lock  sync.RWMutex
	guard *Guard
}

// NewRegistry 初始化注册表
func NewRegistry() *Registry {
	registry := &Registry{
		apps:  make(map[string]*Application),
		guard: new(Guard),
	}

	// 删除过期实例
	go registry.evictTask()

	return registry
}

// Register 实例注册
func (r *Registry) Register(instance *Instance, latestTimestamp int64) (*Application, *utils.Error) {
	// 注册中心应用服务唯一标识
	key := fmt.Sprintf("%s-%s", instance.AppID, instance.Env)
	r.lock.RLock()
	app, ok := r.apps[key]
	r.lock.RUnlock()
	// 如果注册中心不存在该应用服务则添加
	if !ok {
		app = NewApplication(instance.AppID)
	}
	// 添加应用服务的实例
	_, isNew := app.AddInstance(instance, latestTimestamp)
	if isNew {
		// 重新计算一个周期需要的续约次数和保护模式阈值，2*0.85=1.6  取 1
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
	// 对当前节点进行续约
	in, ok := app.Renew(hostname)
	if !ok {
		return nil, utils.NotFound
	}
	// 如果续约成功，对续约总次数+1
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

// 获取全部应用服务
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

// 删除过期实例
func (r *Registry) evict() {
	now := time.Now().UnixNano()
	var expiredInstances []*Instance
	apps := r.getAllApplications()
	var registryLen int
	// 判断上次的更新次数和阈值的大小，true 表示上次更新次数小于阈值，需要启动自我保护机制
	protectStatus := r.guard.selfProtectStatus()
	for _, app := range apps {
		registryLen += app.GetInstanceLen()
		allInstances := app.GetAllInstances()
		for _, instance := range allInstances {
			// 当前时间和最后一次续约的时间间隔
			delta := now - instance.RenewTimestamp
			// 如果（当前时间和最后一次续约的时间间隔大于实例过期时间并且未开启保护模式）或者时间间隔大于实例最大过期时间，将实例加入到过期队列中等待被删除
			if !protectStatus && delta > int64(global.InstanceExpireDuration) ||
				delta > int64(global.InstanceMaxExpireDuration) {
				expiredInstances = append(expiredInstances, instance)
			}
		}
	}
	// 删除实例最大数量限制
	evictionLimit := registryLen - int(float64(registryLen)*global.SelfProtectThreshold)
	expiredLen := len(expiredInstances)
	// 即使过期实例数量大于删除实例最大数量限制，也不能超出该限制
	if expiredLen > evictionLimit {
		expiredLen = evictionLimit
	}
	if expiredLen == 0 {
		return
	}
	// 随机删除过期实例
	for i := 0; i < expiredLen; i++ {
		j := i + rand.Intn(len(expiredInstances)-i)
		expiredInstances[i], expiredInstances[j] = expiredInstances[j], expiredInstances[i]
		expiredInstance := expiredInstances[i]
		r.Cancel(expiredInstance.Env, expiredInstance.AppID, expiredInstance.Hostname, now)
		//todo 取消广播
		//global.Discovery.Nodes.Load().(*Nodes).Replicate(configs.Cancel, expiredInstance)
		log.Printf("### 删除实例 (%v, %v,%v)###\n", expiredInstance.Env, expiredInstance.AppID, expiredInstance.Hostname)

	}
}

// 定期删除过期节点
func (r *Registry) evictTask() {
	// 每 60 秒检查一次节点的可用性
	ticker := time.Tick(global.CheckEvictInterval)
	// 每 15 分钟重置一次注册表
	resetTicker := time.Tick(global.ResetGuardNeedCountInterval)
	for {
		select {
		case <-ticker:
			log.Println("### 每 60 秒随机删除过期实例 ###")
			r.guard.storeLastCount()
			r.evict()
		case <-resetTicker:
			log.Println("### 每 15 分钟重置一次注册表 ###")
			var count int64
			for _, app := range r.getAllApplications() {
				count += int64(app.GetInstanceLen())
			}
			r.guard.setNeed(count)
		}
	}
}
