package registry

import (
	"github.com/zhangweijie11/zDiscovery/schemas"
	"time"
)

type Instance struct {
	Env             string   `json:"env"`              // 当前环境
	AppID           string   `json:"appid"`            // 应用服务标识
	Hostname        string   `json:"hostname"`         // 主机名称
	Addresses       []string `json:"addresses"`        // 实例地址
	Version         string   `json:"version"`          // 应用服务版本
	Status          uint32   `json:"status"`           // 应用服务状态
	RegTimestamp    int64    `json:"reg_timestamp"`    // 注册时间
	UpTimestamp     int64    `json:"up_timestamp"`     // 上线时间
	RenewTimestamp  int64    `json:"renew_timestamp"`  // 续约时间
	DirtyTimestamp  int64    `json:"dirty_timestamp"`  // 脏时间
	LatestTimestamp int64    `json:"latest_timestamp"` // 最后更新时间
}

func NewInstance(req *schemas.RequestRegister) *Instance {
	now := time.Now().UnixNano()
	instance := &Instance{
		Env:             req.Env,
		AppID:           req.AppID,
		Hostname:        req.Hostname,
		Addresses:       req.Addresses,
		Version:         req.Version,
		Status:          req.Status,
		RegTimestamp:    now,
		UpTimestamp:     now,
		RenewTimestamp:  now,
		DirtyTimestamp:  now,
		LatestTimestamp: now,
	}

	return instance
}
