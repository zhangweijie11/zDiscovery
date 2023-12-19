package registry

import "sync"

type Application struct {
	appid           string
	instances       map[string]*Instance
	latestTimeStamp int64
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
	if latestTimestamp <= app.latestTimeStamp {
		latestTimestamp = app.latestTimeStamp + 1
	}
	app.latestTimeStamp = latestTimestamp
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
