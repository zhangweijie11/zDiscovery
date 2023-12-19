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
