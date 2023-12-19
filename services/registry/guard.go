package registry

import (
	"github.com/zhangweijie11/zDiscovery/global"
	"sync"
)

// Guard 统计中心
type Guard struct {
	renewCount     int64        // 所有服务续约次数，执行一次+1
	lastRenewCount int64        // 上次检查周期服务续约次数
	needRenewCount int64        // 一个周期总计需要续约次数，按一次续约 30 秒，周期 60 秒，一个实例就需要2 次，所以服务注册时+2，服务取消时-2
	threshold      int64        // 保护模式阈值，通过 needRenewCount  和阈值比例 （0.85）确定触发自我保护的值
	lock           sync.RWMutex // 读写锁
}

func (gd *Guard) incrNeed() {
	gd.lock.Lock()
	defer gd.lock.Unlock()

	gd.needRenewCount += int64(global.CheckEvictInterval / global.RenewInterval)
	gd.threshold = int64(float64(gd.needRenewCount) + global.SelfProtectThreshold)
}
