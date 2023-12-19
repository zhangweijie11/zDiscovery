package global

import (
	"github.com/zhangweijie11/zDiscovery/services/registry"
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
	RenewInterval               = 30 * time.Second   // 心跳监测间隔时间
	CheckEvictInterval          = 60 * time.Second   // 剔除实例间隔时间
	SelfProtectThreshold        = 0.85               //保护模式阈值
	ResetGuardNeedCountInterval = 15 * time.Minute   //ticker reset guard need count
	InstanceExpireDuration      = 90 * time.Second   //instance's renewTimestamp after this will be canceled
	InstanceMaxExpireDuration   = 3600 * time.Second //instance's renewTimestamp after this will be canceled
	ProtectTimeInterval         = 60 * time.Second   // 两次续约时间间隔
	NodePerceptionInterval      = 5 * time.Second    // 更新节点间隔
)

type Action int

const (
	Register Action = iota
	Renew
	Cancel
	Delete
)

var Discovery *registry.Discovery
