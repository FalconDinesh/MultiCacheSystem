package cache

// import (
// 	"sync"
// )

// type TenantLRUCacheManager struct {
// 	caches map[string]*LRUCache
// 	mu     sync.RWMutex
// }

// func NewTenantLRUCacheManager() *TenantLRUCacheManager {
// 	return &TenantLRUCacheManager{
// 		caches: make(map[string]*LRUCache),
// 	}
// }

// func (manager *TenantLRUCacheManager) GetCache(tenantID string) *LRUCache {
// 	manager.mu.RLock()
// 	cache, exists := manager.caches[tenantID]
// 	manager.mu.RUnlock()
// 	if exists {
// 		return cache
// 	}

// 	manager.mu.Lock()
// 	defer manager.mu.Unlock()
// 	cache = NewLRUCache(300, 60) // Customize size and TTL as needed
// 	manager.caches[tenantID] = cache
// 	return cache
// }
