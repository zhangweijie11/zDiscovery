package registry

import (
	"github.com/zhangweijie11/zDiscovery/global"
	"sync"
	"sync/atomic"
)

// Guard 统计中心
type Guard struct {
	renewCount     int64 // 所有服务续约次数，执行一次+1
	lastRenewCount int64 // 上次检查周期服务续约次数
	needRenewCount int64 // 一个周期总计需要续约次数，按一次续约 30 秒，周期 60 秒，一个实例就需要2 次，所以服务注册时+2，服务取消时-2
	threshold      int64 // 保护模式阈值，通过 needRenewCount  和阈值比例 （0.85）确定触发自我保护的值，
	// 例如一个周期需要续约 100 次，阈值为 0.85,则保护模式的阈值为 85 次，如果小于 85 说明需要开启保护模式，防止把全部实例都删除
	lock sync.RWMutex // 读写锁
}

func (gd *Guard) incrNeed() {
	gd.lock.Lock()
	defer gd.lock.Unlock()

	gd.needRenewCount += int64(global.CheckEvictInterval / global.RenewInterval)
	gd.threshold = int64(float64(gd.needRenewCount) * global.SelfProtectThreshold)
}

// 服务续约次数+1
func (gd *Guard) incrCount() {
	atomic.AddInt64(&gd.renewCount, 1)
}

func (gd *Guard) decrNeed() {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	gd.needRenewCount -= int64(global.CheckEvictInterval / global.RenewInterval)
	gd.threshold = int64(float64(gd.needRenewCount) * global.SelfProtectThreshold)
}

// 将 renewCount 的值存储到 lastRenewCount，并且将 renewCount 置为 0
func (gd *Guard) storeLastCount() {
	atomic.StoreInt64(&gd.lastRenewCount, atomic.SwapInt64(&gd.renewCount, 0))
}

// 判断上次的更新次数和保护模式阈值的大小
func (gd *Guard) selfProtectStatus() bool {
	return atomic.LoadInt64(&gd.lastRenewCount) < atomic.LoadInt64(&gd.threshold)
}

// 重置周期需要续约次数和保护模式的阈值，降低脏数据风险
func (gd *Guard) setNeed(count int64) {
	gd.lock.Lock()
	defer gd.lock.Unlock()
	gd.needRenewCount = count * int64(global.CheckEvictInterval/global.RenewInterval)
	gd.threshold = int64(float64(gd.needRenewCount) * global.SelfProtectThreshold)
}
