package registry

import "sync"

type Registry struct {
	apps  map[string]*Application
	lock  sync.RWMutex
	guard *Guard
}

type name struct {
}
