package cache

import "time"

type CacheSystem interface {
	Get(key string) (interface{}, error)
	//GetWithTTL(key string) (interface{}, time.Duration, time.Time, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
}