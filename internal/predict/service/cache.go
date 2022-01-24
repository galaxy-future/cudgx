package service

import "sync"

var cacheManager *ViewCache

type SEntries struct {
	entries map[string]*RedundancySeries
}

type ViewCache struct {
	timestamp int64
	services  map[string]*SEntries
	lock      sync.RWMutex
}

func (cache *ViewCache) getRedundancySeries(serviceName, metricName string, timestamp int64) *RedundancySeries {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	if timestamp > cache.timestamp {
		return nil
	}
	serviceEntries, exists := cache.services[serviceName]
	if !exists {
		serviceEntries = &SEntries{
			entries: make(map[string]*RedundancySeries),
		}
	}

	seriesEntries, _ := serviceEntries.entries[metricName]
	return seriesEntries
}

func (cache *ViewCache) setRedundancySeries(serviceName, metricName string, timestamp int64, series *RedundancySeries) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	//过期的结果
	if timestamp < cache.timestamp {
		return
	} else if timestamp > cache.timestamp { //更新
		cache.services = make(map[string]*SEntries)
		cache.timestamp = timestamp
	}
	serviceEntries, exists := cache.services[serviceName]
	if !exists {
		serviceEntries = &SEntries{
			entries: make(map[string]*RedundancySeries),
		}
	}
	serviceEntries.entries[metricName] = series
}

func init() {
	cacheManager = &ViewCache{
		timestamp: 0,
		services:  map[string]*SEntries{},
		lock:      sync.RWMutex{},
	}
}
