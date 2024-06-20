package handler

import (
	"multi-backend-cache/Internal/cache"
	"multi-backend-cache/Internal/config"
	utils "multi-backend-cache/packageUtils/Utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/* Structure for multiple Cache System
 */
type Server struct {
	tenantCaches  *cache.FixedTenantsCaches // Use FixedTenantsCaches for multi-tenant support
	redisCache    cache.CacheSystem
	memCache      cache.CacheSystem
	// inmemoryCache cache.CacheSystem
	mu            sync.Mutex
}

// /* Structure for multiple Cache System
//  */
//  type Server struct {
// 	tenantCaches *cache.FixedTenantsCaches // Use FixedTenantsCaches for multi-tenant support
// 	redisCache   cache.CacheSystem
// 	memCache     cache.CacheSystem
// 	mu           sync.Mutex
// }

/* Creating a New server
 */
func NewServer(tenantCaches *cache.FixedTenantsCaches, redisCache cache.CacheSystem, memCache cache.CacheSystem) *Server {
	return &Server{
		tenantCaches:  tenantCaches,
		redisCache:    redisCache,
		memCache:      memCache,
		// inmemoryCache: inmemoryCache,
	}
}

// /* Creating a New server
//  */
//  func NewServer(tenantCaches *cache.FixedTenantsCaches, redisCache cache.CacheSystem, memCache cache.CacheSystem) *Server {
// 	return &Server{
// 		tenantCaches: tenantCaches,
// 		redisCache:   redisCache,
// 		memCache:     memCache,
// 	}
// }

/* Determine the cache Library Type based on URI Param.
 */
func (s *Server) determineCacheLibraryType(cacheType string, tenantID string) cache.CacheSystem {
	//cacheType := mux.Vars(r)["cacheType"]
	switch cacheType {
	case "redis":
		return s.redisCache
	case "memcache":
		return s.memCache
	case "inmemory":
		if config.AppConfig.IsTenantBased {
			return s.tenantCaches.GetCache(tenantID)
		} else {
			return s.tenantCaches.GetCache(cache.DefaultTenant)
		}
	default:
		return nil
	}
}

// @Summary Get value from cache by key
// @Description Retrieve a value from the cache using the provided key and cache type
// @ID get-cache-by-key
// @Accept  json
// @Produce  json
// @Param   key        path    string  true  "Cache Key"
// @Param   system      query   string  true  "Cache Type"
// @Success 200  "status: ok"
// @Failure 400  "Bas Request"
// @Failure 404  "Not Found"
// @Failure 500  "Internal Server Error"
// @Router /cache/{key} [get]
func (s *Server) GetCacheHandler(c *gin.Context) {
	key := c.Param("key")
	tenantID := c.Query("tenantID")
	CacheLibraryType := c.Query("system")
	cache := s.determineCacheLibraryType(CacheLibraryType, tenantID)

	if cache == nil {
		logrus.Error("Unsupported cache type, please provide supported cache System", cache)
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	value, err := cache.Get(key)
	if err != nil {

		if err.Error() == utils.NotFound.Error() {
			logrus.Errorf("Error: key %s: %v", key, err)
			utils.RespondError(c.Writer, http.StatusNotFound, err.Error())
			return
		}
		logrus.Errorf("Error while getting cache for key %s: %v", key, err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
	logrus.Infof("Cache retrieved for key %s: %v", key, value)
	utils.RespondJSON(c.Writer, http.StatusOK, value)
}

// func (s *Server) GetCacheWithTTLHandler(c *gin.Context) {
// 	key := c.Param("key")
// 	tenantID := c.Query("tenantID")
// 	CacheLibraryType := c.Query("system")
// 	cache := s.determineCacheLibraryType(CacheLibraryType, tenantID)

// 	if cache == nil {
// 		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
// 		return
// 	}

// 	value, ttl, expiryTime, err := cache.GetWithTTL(key)
// 	if err != nil {
// 		utils.LogError("Error while getting cache with TTL", err)
// 		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to get cache with TTL")
// 		return
// 	}

// 	utils.RespondJSON(c.Writer, http.StatusOK, map[string]interface{}{
// 		"value":       value,
// 		"ttl":         ttl.Seconds(),
// 		"expiry_time": expiryTime,
// 	})
// }

//	type SetCachePayload struct {
//		Key        string        `json:"key"`
//		Value      string        `json:"value"`
//		TTL        time.Duration `json:"ttl"`
//		ExpiryTime time.Time
//	}

// @Summary Set value in cache
// @Description Set a value in the cache with a specified key and TTL (Time-To-Live)
// @ID set-cache-value
// @Accept json
// @Produce json
// @Param system query string true "Cache Type"
// @Param payload body cache.CacheData true "Cache Payload"
// @Success 200  "status: ok"
// @Failure 400  "Bad Request"
// @Failure 404  "Not Found"
// @Failure 500  "Internal Server Error"
// @Router /cache [post]
func (s *Server) SetCacheHandler(c *gin.Context) {
	var payload cache.CacheData
	if err := c.ShouldBindJSON(&payload); err != nil {

		logrus.Error("Invalid request payload", err)
		utils.RespondError(c.Writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if payload.Key == "" {

		logrus.Error("Key must not be null")
		utils.RespondError(c.Writer, http.StatusNotFound, "Key must not be null")
		return
	}

	CacheLibraryType := c.Query("system")
	tenantID := c.Query("tenantID")
	cache := s.determineCacheLibraryType(CacheLibraryType, tenantID)
	if cache == nil {
		logrus.Error("Unsupported cache type, please provide supported cache System", cache)
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	logrus.Debugf("Setting cache for key %s with TTL %s", payload.Key, payload.TTL)
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Set(payload.Key, payload.Value, payload.TTL); err != nil {
		logrus.Errorf("Error while setting cache for key %s: %v", payload.Key, err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to set cache")
		return
	}

	logrus.Infof("Cache set for key %s", payload.Key)
	utils.RespondJSON(c.Writer, http.StatusOK, map[string]string{"status": "ok"})
}

// @Summary Delete value from cache by key
// @Description Delete a value from the cache using the provided key and cache type
// @ID delete-cache-by-key
// @Accept  json
// @Produce  json
// @Param   key        path    string  true  "Cache Key"
// @Param   system      query   string  true  "Cache Type"
// @Success 200  "status: ok"
// @Failure 400  "Bad Request"
// @Failure 404  "Not Found"
// @Failure 500  "Internal Server Error"
// @Router /cache/{key} [delete]
func (s *Server) DeleteCacheHandler(c *gin.Context) {
	key := c.Param("key")
	CacheLibraryType := c.Query("system")
	tenantID := c.Query("tenantID")

	cache := s.determineCacheLibraryType(CacheLibraryType, tenantID)

	if cache == nil {
		logrus.Error("Unsupported cache type")
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	logrus.Debugf("Deleting cache for key %s", key)
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Delete(key); err != nil {
		if err.Error() == utils.NotFound.Error() {
			logrus.Errorf("Error for key %s: %v", key, err)
			utils.RespondError(c.Writer, http.StatusNotFound, "Cache not Found - Failed to delete cache")
			return
		}
		logrus.Errorf("Error while deleting cache for key %s: %v", key, err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to delete cache")
		return
	}

	logrus.Infof("Cache deleted for key %s", key)
	utils.RespondJSON(c.Writer, http.StatusOK, map[string]string{"status": "ok"})
}

// @Summary Clear all caches
// @Description clear caches for the provided cache type
// @ID clear-cache
// @Accept  json
// @Produce  json
// @Param   system      query   string  true  "Cache Type"
// @Success 200  "status: ok"
// @Failure 400  "Bad Request"
// @Failure 404  "Not Found"
// @Failure 500  "Internal Server Error"
// @Router /cache/clear [put]
func (s *Server) ClearCacheHandler(c *gin.Context) {
	CacheLibraryType := c.Query("system")
	tenantID := c.Query("tenantID")
	cache := s.determineCacheLibraryType(CacheLibraryType, tenantID)

	if cache == nil {
		utils.RespondError(c.Writer, http.StatusBadRequest, "Unsupported cache type")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := cache.Clear(); err != nil {
		utils.LogError("Error while clearing cache", err)
		utils.RespondError(c.Writer, http.StatusInternalServerError, "Failed to clear cache")
		return
	}
	utils.RespondJSON(c.Writer, http.StatusOK, map[string]string{"status": "ok"})
}
