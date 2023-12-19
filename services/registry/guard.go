package registry

import "sync"

type Guard struct {
	renewCount     int64        // 续约数量
	lastRenewCount int64        // 最后续约数量
	needRenewCount int64        // 需要续约数量
	threshold      int64        // 保护模式阈值
	lock           sync.RWMutex // 读写锁
}
