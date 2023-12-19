package registry

import "sync"

type Application struct {
	appid           string
	instance        map[string]*Instance
	latestTimeStamp int64
	lock            sync.RWMutex
}
