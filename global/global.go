package global

import (
	"time"
)

const (
	NodeStatusUp = iota + 1
	NodeStatusDown
)

const (
	StatusOK = 200
)

const (
	StatusReceive = iota + 1
	StatusNotReceive
)

const (
	RegisterURL = "/api/register"
	CancelURL   = "/api/cancel"
	RenewURL    = "/api/renew"
	FetchAllURL = "/api/fetchall"
)

const (
	DiscoveryAppId = "zDiscovery"
)

const (
	RenewInterval               = 30 * time.Second   // 心跳监测时间间隔
	CheckEvictInterval          = 60 * time.Second   // 删除实例时间间隔
	SelfProtectThreshold        = 0.85               //保护模式阈值
	ResetGuardNeedCountInterval = 15 * time.Minute   // 注册表重置时间间隔
	InstanceExpireDuration      = 90 * time.Second   // 实例的过期时间
	InstanceMaxExpireDuration   = 3600 * time.Second //实例的最大过期时间
	ProtectTimeInterval         = 60 * time.Second   // 两次续约时间间隔
	NodePerceptionInterval      = 5 * time.Second    // 更新节点时间间隔
)

type Action int

const (
	Register Action = iota
	Renew
	Cancel
	Delete
)
