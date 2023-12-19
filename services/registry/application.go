package registry

import (
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"sync"
	"time"
)

type Application struct {
	appid           string
	instances       map[string]*Instance
	latestTimestamp int64
	lock            sync.RWMutex
}

// NewApplication 单个服务可能存在多个实例，多节点部署
func NewApplication(appid string) *Application {
	return &Application{
		appid:     appid,
		instances: make(map[string]*Instance),
	}
}

// 更新应用服务的最新更新时间
func (app *Application) upLatestTimestamp(latestTimestamp int64) {
	if latestTimestamp <= app.latestTimestamp {
		latestTimestamp = app.latestTimestamp + 1
	}
	app.latestTimestamp = latestTimestamp
}

// AddInstance 应用服务增加实例
func (app *Application) AddInstance(instance *Instance, latestTimestamp int64) (*Instance, bool) {
	app.lock.Lock()
	defer app.lock.Unlock()
	appIns, ok := app.instances[instance.Hostname]
	if ok {
		// 节点的上限时间同步为服务的上线时间
		instance.UpTimestamp = appIns.UpTimestamp
		if instance.DirtyTimestamp < appIns.DirtyTimestamp {
			instance = appIns
		}
	}
	app.instances[instance.Hostname] = instance
	// 根据实例的最后更新时间修改应用服务的最后更新时间
	app.upLatestTimestamp(latestTimestamp)
	// 如果ok 为 true 表示已经存在该实例
	return instance, !ok

}

// Renew 续约
func (app *Application) Renew(hostname string) (*Instance, bool) {
	app.lock.Lock()
	defer app.lock.Unlock()
	appIn, ok := app.instances[hostname]
	if !ok {
		return nil, ok
	}
	// 修改实例的最后续约时间
	appIn.RenewTimestamp = time.Now().UnixNano()
	// 拷贝节点地址
	return copyInstance(appIn), true
}

// deep copy
func copyInstance(src *Instance) *Instance {
	dst := new(Instance)
	*dst = *src
	// 拷贝节点地址
	dst.Addresses = make([]string, len(src.Addresses))
	for i, addr := range src.Addresses {
		dst.Addresses[i] = addr
	}
	return dst
}

// GetInstance 获取应用服务的全部实例
func (app *Application) GetInstance(status uint32, latestTime int64) (*FetchData, *utils.Error) {
	app.lock.RLock()
	defer app.lock.RUnlock()
	// 如果最后更新时间大于应用服务的最后更新时间，则无需更新
	if latestTime >= app.latestTimestamp {
		return nil, utils.NotModified
	}
	fetchData := FetchData{
		Instances:       make([]*Instance, 0),
		LatestTimestamp: app.latestTimestamp,
	}
	var exists bool
	for _, instance := range app.instances {
		if status&instance.Status > 0 {
			exists = true
			newInstance := copyInstance(instance)
			fetchData.Instances = append(fetchData.Instances, newInstance)
		}
	}
	if !exists {
		return nil, utils.NotFound
	}
	return &fetchData, nil
}

// Cancel 注销应用服务的节点
func (app *Application) Cancel(hostname string, latestTimestamp int64) (*Instance, bool, int) {
	newInstance := new(Instance)
	app.lock.Lock()
	defer app.lock.Unlock()
	appIn, ok := app.instances[hostname]
	if !ok {
		return nil, ok, 0
	}
	// 删除实例
	delete(app.instances, hostname)
	appIn.LatestTimestamp = latestTimestamp
	app.upLatestTimestamp(latestTimestamp)
	*newInstance = *appIn
	return newInstance, true, len(app.instances)
}

// GetAllInstances 获取所有实例
func (app *Application) GetAllInstances() []*Instance {
	app.lock.RLock()
	defer app.lock.RUnlock()
	rs := make([]*Instance, 0, len(app.instances))
	for _, instance := range app.instances {
		newInstance := new(Instance)
		*newInstance = *instance
		rs = append(rs, newInstance)
	}
	return rs
}
