package main

import "sync"

type GKVS struct {
	storage map[string]GKVSTypes
	mutex   sync.RWMutex
}

func NewGKVS() *GKVS {
	return &GKVS{
		storage: make(map[string]GKVSTypes),
	}
}

func (g *GKVS) Set(key string, value GKVSTypes) GKVSTypes {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.storage[key] = value
	return value
}

func (g *GKVS) Get(key string) GKVSTypes {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	value, exists := g.storage[key]
	if !exists {
		return NewGKVSNone()
	}
	return value
}

func (g *GKVS) Delete(key string) GKVSTypes {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	value, exists := g.storage[key]
	if !exists {
		return NewGKVSNone()
	}
	delete(g.storage, key)
	return value
}
