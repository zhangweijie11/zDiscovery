package global

import "time"

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
	ProtectTimeInterval         = 60 * time.Second   //two renew cycle
	NodePerceptionInterval      = 5 * time.Second    //nodesprotect
)
