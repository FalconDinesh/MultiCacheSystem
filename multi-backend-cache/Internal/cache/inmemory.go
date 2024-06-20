package cache

import (
	"container/list"
	"encoding/json"
	utils "multi-backend-cache/packageUtils/Utils"
	"reflect"
	"sync"
	"time"

	"multi-backend-cache/Internal/config"

	"github.com/pbnjay/memory"
	"github.com/sirupsen/logrus"
)

// type MyDuration = time.Duration

type CacheData struct { //node

	Key        string        `json:"key" example:"1"`
	Value      interface{}   `json:"value" `
	TTL        time.Duration `json:"ttl" example:"100"`
	ExpiryTime time.Time     `json:"expirytime" example:"2021-05-25T00:53:16.535668Z" format:"date-time" swaggerignore:"true"`
}

// LRUCache represents the LRU cache, that consists of capacity, linkedlist as list, hashmap as index and lock
type LRUCache struct {
	capacity   int // in bytes
	used       int // in bytes
	list       *list.List
	index      map[string]*list.Element //key-> sring, value -> pointer to the list(*list.Element)
	lock       sync.Mutex
	defaultTTL time.Duration
}

type FixedTenantsCaches struct {
	caches map[string]*LRUCache
}

// GetCache retrieves the cache for the specified tenant.
func (ftc *FixedTenantsCaches) GetCache(tenantID string) *LRUCache {
	if cache, exists := ftc.caches[tenantID]; exists {
		return cache
	}
	return nil // Optionally handle the case where tenantID is not recognized.
}

// Initializes fixed tenant caches with predefined capacities.
const DefaultTenant string = "defaultTenant"

func NewFixedTenantsCaches(isTenantBased bool, totalCacheMemory int, defaultTTL time.Duration) *FixedTenantsCaches {
	tenantCaches := make(map[string]*LRUCache)

	if isTenantBased {
		tenantIDs := config.AppConfig.TenantIDs
		logrus.Infof("Tenant IDs: %v", tenantIDs)

		for _, id := range tenantIDs {
			tenantCaches[id] = NewLRUCache(int(totalCacheMemory)/len(tenantIDs), defaultTTL) // Define capacity per tenant here.
		}
	} else {
		tenantCaches[DefaultTenant] = NewLRUCache(totalCacheMemory, defaultTTL)
	}

	go checkMemoryForTenants(totalCacheMemory, tenantCaches)

	return &FixedTenantsCaches{
		caches: tenantCaches,
	}
}

// Checks the cache is expired or not
func IsExpired(expiryTime time.Time) bool {
	return time.Now().After(expiryTime)
}

// NewLRUCache creates a new LRU cache with the given capacity and ttl
func NewLRUCache(capacity int, defaultTTL time.Duration) *LRUCache {
	lru := &LRUCache{ // The "&" operator returns a pointer to the newly created LRUCache instance.
		capacity:   capacity,
		list:       list.New(),
		index:      make(map[string]*list.Element),
		defaultTTL: defaultTTL,
	}
	go DeleteExpiredCache(lru)

	return lru
}

// Go-routine that runs concurrently and for each 5 seconds scans the memory and deletes the
// expired ones
func DeleteExpiredCache(lru *LRUCache) {
	for range time.Tick(5 * time.Second) {
		lru.lock.Lock()
		// Iterate over the cache items and delete expired ones.
		for key, element := range lru.index {
			node := element.Value.(*CacheData)
			if IsExpired(node.ExpiryTime) {
				removeAndResize(lru, node, element)
				logrus.Infof("Deleted cache key %s with expiry time %v", key, node.ExpiryTime)
			}
		}
		lru.lock.Unlock()
	}
}

// Go-routine that concurrently checks whether the memory size increases and allocates cache memory accordingly
func checkMemoryForTenants(totalCacheMemory int, tenantCaches map[string]*LRUCache) {
	for range time.Tick(1 * time.Second) {
		numberOfTenants := len(tenantCaches)
		if totalCacheMemory < int(memory.TotalMemory()) {
			cacheMemory := float64(memory.TotalMemory()) * config.AppConfig.MemoryUsagePercentage
			capacity := int(cacheMemory) / numberOfTenants
			logrus.Infof("increased capacity for tenant :: %d", int(cacheMemory))
			for _, cache := range tenantCaches {
				cache.lock.Lock()
				cache.capacity = capacity
				cache.lock.Unlock()
			}
		}
	}
}

// GetAllCache retrieves all values from the cache
func (c *LRUCache) GetAllCache() []*CacheData {
	c.lock.Lock()
	defer c.lock.Unlock()

	var allCacheData []*CacheData
	for element := c.list.Front(); element != nil; element = element.Next() {
		node := element.Value.(*CacheData)
		nodeJSON, err := json.Marshal(node)
		if err != nil {
			logrus.Error("Error marshalling node to JSON:", err)
		} else {
			logrus.Debugf("All cached data without Expired Cache: %s", string(nodeJSON))
		}
		if IsExpired(node.ExpiryTime) {
			removeAndResize(c, node, element) // Entry has expired, remove it
		} else {
			allCacheData = append(allCacheData, node)
		}
	}
	return allCacheData
}

// GetCache returns the cache value for a specified key if exists
func (c *LRUCache) Get(key string) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var Err error
	if element, found := c.index[key]; found {
		logrus.Debugf("Existing cache found for key %s: %v", key, element.Value)
		node := element.Value.(*CacheData)
		nodeJSON, err := json.Marshal(node)
		if err != nil {
			logrus.Error("Error marshalling node to JSON:", err)
		} else {
			logrus.Infof("Cache data for key %s: %s", key, string(nodeJSON))
		}
		if IsExpired(node.ExpiryTime) { // Check if the entry has expired
			removeAndResize(c, node, element)
			return nil, utils.NotFound
		}
		c.list.MoveToFront(element)
		return node.Value, Err
	} else {
		logrus.Infof("Cache miss for key %s", key)
		return nil, utils.NotFound
	}
}

// Get cache along with data
// func (c *LRUCache) GetWithTTL(key string) (interface{}, time.Duration, time.Time, error) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	var Err error
// 	if element, found := c.index[key]; found {
// 		log.Println("Exsiting cache found", element.Value)
// 		node := element.Value.(*CacheData)
// 		nodeJSON, err := json.Marshal(node)
// 		if err != nil {
// 			log.Println("Error marshalling node to JSON:", err)
// 		} else {
// 			log.Printf("Cache data for key %s is: %s", key, string(nodeJSON))
// 		}
// 		size := CalculateSize(node)
// 		if IsExpired(node.ExpiryTime) { // Check if the entry has expired
// 			c.list.Remove(element)
// 			delete(c.index, key)
// 			c.used -= size
// 			// return nil, err
// 		}
// 		log.Println("stp-1 node.TTL:", node.TTL.Seconds())
// 		// expiryTime := CalculateExpiryTime(node.TTL) //extend expiry
// 		log.Println("stp-2 node.TTL:", node.TTL)
// 		//node.TTL = time.Duration((node.TTL * time.Nanosecond).Seconds())
// 		log.Println("stp-3 after node.TTL:", node.TTL)
// 		// node.ExpiryTime = expiryTime
// 		c.list.MoveToFront(element)
// 		ttl := time.Until(node.ExpiryTime)
// 		return node.Value, ttl, node.ExpiryTime, Err
// 	} else {
// 		Err := errors.New("key not found")
// 		return Err, 0, time.Time{}, nil
// 	}
// }

// Replaces the existing cache value with new value, along with resizing the cache.
func updateAndResize(c *LRUCache, node *CacheData, value interface{}, ttl time.Duration, expiryTime time.Time) {
	updateCacheUsed(c, node, false) // reduce the size of the node that is replaced
	node.Value = value
	node.TTL = ttl
	node.ExpiryTime = expiryTime
	updateCacheUsed(c, node, true) // Add the size of the new node back to cache
}

// setCache adds a value to the cache or updates the exisiting value
func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) error {
	logrus.Debugf("Setting key %s", key)
	c.lock.Lock()
	defer c.lock.Unlock()
	if ttl <= 0 {
		ttl = c.defaultTTL
	}
	logrus.Debugf("TTL for key %s: %s", key, ttl)
	expiryTime := CalculateExpiryTime(ttl)
	if element, found := c.index[key]; found {
		logrus.Infof("Updating existing cache for key %s", key)
		c.list.MoveToFront(element)

		node := element.Value.(*CacheData)
		updateAndResize(c, node, value, ttl, expiryTime)
	} else {
		logrus.Infof("Creating new cache node for key %s", key)
		newNode := &CacheData{Key: key, Value: value, TTL: ttl, ExpiryTime: expiryTime}
		newNodeSize := CalculateSize(newNode)

		for c.used+newNodeSize > c.capacity { // Recursively checked for freeing the last element to store the new one.
			logrus.Warn("Capacity Exceeded. Removing least recently used items.")
			backElement := c.list.Back()
			if backElement != nil {
				backNode := backElement.Value.(*CacheData)
				removeAndResize(c, backNode, backElement)
			}
		}
		element := c.list.PushFront(newNode)
		c.index[key] = element
		c.used += newNodeSize
	}
	return nil
}

// DeleteCache deletes a value from the cache
func (c *LRUCache) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, found := c.index[key]; found {
		node := element.Value.(*CacheData)
		removeAndResize(c, node, element)
		logrus.Infof("Deleted cache for key %s", key)
		return nil
	} else {
		logrus.Infof("Cache miss for key %s during deletion", key)
		return utils.NotFound
	}
}

// Returns expiry time of a cache based on it's TTL value
func CalculateExpiryTime(ttl time.Duration) time.Time {
	logrus.Debug("Calculating Expiry Time...")
	return time.Now().Add(ttl * time.Second)
}

// Returns the size of the node that holds the data, key and ttl
func CalculateSize(node *CacheData) int {
	size := reflect.TypeOf(*node).Size()
	return int(size)
}

// Function to clear the cache
func (c *LRUCache) Clear() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.list.Init()
	c.index = make(map[string]*list.Element)
	c.used = 0
	logrus.Infof("Cache cleared")
	return nil
}

// Function to update the cache used size 
func updateCacheUsed(c *LRUCache, node *CacheData, isAddition bool) {
	if isAddition { 				// Boolean to mention whether to add or remove size
		c.used += CalculateSize(node)
	} else {
		c.used -= CalculateSize(node)
	}
}

// deletes the data physically from node and map
func removeAndResize(c *LRUCache, node *CacheData, element *list.Element) {
	updateCacheUsed(c, node, false) // Reduce the size of the node that is replaced
	c.list.Remove(element)          // removes the node from list
	delete(c.index, node.Key)       // Deletes record from Map
}
