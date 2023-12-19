package schemas

// RequestRegister 注册参数
type RequestRegister struct {
	Env             string   `json:"env"`
	AppID           string   `json:"appid"`
	Hostname        string   `json:"hostname"`
	Addresses       []string `json:"addresses"`
	Status          uint32   `json:"status"`
	Version         string   `json:"version"`
	LatestTimestamp int64    `json:"latest_timestamp"`
	DirtyTimestamp  int64    `json:"dirty_timestamp"` //other node send
	Replication     bool     `json:"replication"`     //other node send
}

// RequestRenew 心跳监测参数
type RequestRenew struct {
	Env            string `form:"env"`
	AppId          string `form:"appid"`
	Hostname       string `form:"hostname"`
	DirtyTimestamp int64  `form:"dirty_timestamp"` //other node send
	Replication    bool   `form:"replication"`     //other node send
}

// RequestCancel 注销参数
type RequestCancel struct {
	Env             string `form:"env"`
	AppId           string `form:"appid"`
	Hostname        string `form:"hostname"`
	LatestTimestamp int64  `form:"last_timestamp"` //other node send
	Replication     bool   `form:"replication"`    //other node send
}

// RequestFetch 同步参数
type RequestFetch struct {
	Env    string `from:"env"`
	AppId  string `form:"appid"`
	Status uint32 `form:"status"`
}

// RequestFetchs 全部同步参数
type RequestFetchs struct {
	Env    string   `form:"env"`
	AppId  []string `form:"appid"`
	Status uint32   `form:"status"`
}

// RequestNodes 获取节点参数
type RequestNodes struct {
	Env string `form:"env"`
}
