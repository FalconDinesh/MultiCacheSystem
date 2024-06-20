package cache

import (
	"encoding/json"
	utils "multi-backend-cache/packageUtils/Utils"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/sirupsen/logrus"
)

type MemCache struct {
	client *memcache.Client
	ttl    int32
}

func NewMemCache(server string, ttl int32) *MemCache {
	client := memcache.New(server)
	logrus.Infof("Memcache initialized with server: %s", server)
	return &MemCache{client: client, ttl: ttl}
}

// Get retrieves a value from the cache by key
func (m *MemCache) Get(key string) (interface{}, error) {
	item, err := m.client.Get(key)
	// fmt.Println(item.Expiration)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, utils.NotFound
		}
		logrus.Errorf("Get: error getting key %s: %v", key, err)
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal(item.Value, &data)
	if err != nil {
		logrus.Errorf("Get: error unmarshaling value for key %s: %v", key, err)
		return nil, err
	}
	return data, nil
}

// GetWithTTL retrieves the value and its TTL for the specified key
// func (m *MemCache) GetWithTTL(key string) (interface{}, time.Duration, time.Time, error) {
// 	item, err := m.client.Get(key)
// 	fmt.Println("item Expiration:: ", item.Expiration)
// 	fmt.Println("item :: ", item)

// 	if err != nil {
// 		fmt.Println("item error :: ", err)
// 		return nil, 0, time.Time{}, err
// 	}
// 	return item.Value, time.Duration(item.Expiration) * time.Second, time.Time{}, nil
// }

func (m *MemCache) Set(key string, value interface{}, ttl time.Duration) error {
	ttlDuration := time.Duration(ttl) * time.Second
	val, err := json.Marshal(value)
	if err != nil {
		logrus.Errorf("Set: error marshaling value for key %s: %v", key, err)
		return err
	}
	actualTTL := int32(ttlDuration.Seconds())
	if ttl <= 0 {
		actualTTL = m.ttl
	}
	logrus.Infof("Setting KEY: %s with VALUE: %v and TTL: %d seconds", key, value, actualTTL)
	err = m.client.Set(&memcache.Item{Key: key, Value: val, Expiration: actualTTL})
	if err != nil {
		logrus.Errorf("Set: error setting key %s: %v", key, err)
	}
	return err
}

// Delete removes a value from the cache by key
func (m *MemCache) Delete(key string) error {
	err := m.client.Delete(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			logrus.Debugf("Delete: key %s does not exist", key)
			// return err
			return utils.NotFound
		}
		logrus.Errorf("Delete: error deleting key %s: %v", key, err)
	}
	return err
}

func (m *MemCache) Clear() error {
	logrus.Infof("Clearing all cache entries")
	err := m.client.FlushAll()
	if err != nil {
		logrus.Errorf("Error while clearing cache: %v", err)
		return err
	}
	logrus.Infof("Cache cleared successfully")
	return nil
}
