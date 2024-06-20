package test

import (
	"bytes"
	handler "multi-backend-cache/Internal/Handler"
	"multi-backend-cache/Internal/cache"
	"multi-backend-cache/Internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Function to set up Inmemory router, with capacity as 300 bytes and TTL to 10s
func setupInMemoryTenantRouter() *gin.Engine {
	config.LoadConfig("../Internal/config/config.yaml")
	config.AppConfig.IsTenantBased = true
	inmemorycache := cache.NewFixedTenantsCaches(true, 900, 10)
	cacheSystemType := handler.NewServer(inmemorycache, nil, nil)
	router := gin.Default()
	router.Use(handler.ValidateTenant())
	router.GET("/cache/:key", cacheSystemType.GetCacheHandler)
	router.POST("/cache", cacheSystemType.SetCacheHandler)
	router.DELETE("/cache/:key", cacheSystemType.DeleteCacheHandler)
	// router.GET("/cache/TTL/:key", cacheSystemType.GetCacheWithTTLHandler)
	router.PUT("/cache/clear", cacheSystemType.ClearCacheHandler)

	return router
}

// Test for set function
func TestInMemTenantPostCacheHandler(t *testing.T) {
	router := setupInMemoryTenantRouter()

	t.Run("Valid Data", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"key": "1", "value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Key not passed", func(t *testing.T) {
		invalidData := `{"key": "","value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", bytes.NewBufferString(invalidData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Invalid Data Format (key as int)", func(t *testing.T) {
		invalidData := `{"key": 1 ,"value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", bytes.NewBufferString(invalidData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("Invalid Tenant", func(t *testing.T) {
		invalidData := `{"key": 1 ,"value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant9", bytes.NewBufferString(invalidData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Tenant not provided", func(t *testing.T) {
		invalidData := `{"key": 1 ,"value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory", bytes.NewBufferString(invalidData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Test for get function
func TestInMemTenantGetCacheHandler(t *testing.T) {
	router := setupInMemoryTenantRouter()

	// First, post a cache entry
	w := httptest.NewRecorder()
	reqBody := `{"key": "2", "value": {"name":"session"},"ttl":200}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Run("Valid Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/2?system=inmemory&tenantID=tenant1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `{"name":"session"}`)
	})

	t.Run("InValid Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=inmemory&tenantID=tenant1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Invalid Tenant", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/2?system=inmemory&tenantID=tenant9", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Tenant Not provided", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=inmemory", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Test to check deleteAPI
func TestInMemTenantDeleteCacheHandler(t *testing.T) {
	router := setupInMemoryTenantRouter()

	// post a cache entry in tenant 1
	w := httptest.NewRecorder()
	reqBody := `{"key": "3", "value": "cache", "ttl": 300}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// post a cache entry in tenant 2
	w = httptest.NewRecorder()
	reqBody = `{"key": "3", "value": "cache", "ttl": 300}`
	req, _ = http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant2", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	//Delete the cache with key 3 in tenant1
	t.Run("Valid Key", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/cache/3?system=inmemory&tenantID=tenant1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	//check if the cache is deleted in tenant1
	t.Run("Check if deleted", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=inmemory&tenantID=tenant1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	//check if the cache with key3 in tenant2 is not deleted
	t.Run("Check if same key in differnt tenant not deleted", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=inmemory&tenantID=tenant2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}

// Test for data expiry with default TTL
func TestInMemoryTenantDataExpiryWithDefaultTTL(t *testing.T) {
	router := setupInMemoryTenantRouter()

	w := httptest.NewRecorder()
	reqBodyExpiry := `{"key": "10", "value": "sessions"}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant2", strings.NewReader(reqBodyExpiry))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	t.Run("Expiry check", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/cache/10?system=inmemory&tenantID=tenant2", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "sessions")

		time.Sleep(10 * time.Second)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/cache/10?system=inmemory&tenantID=tenant2", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

	})
}

// Test to check clear API
func TestInMemTenantClearCacheHandler(t *testing.T) {
	router := setupInMemoryTenantRouter()

	// First, post a cache entry in tenant1
	w := httptest.NewRecorder()
	reqBody := `{"key": "3", "value": "cache", "ttl": 300}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	//2nd cache entry in tenant1
	w = httptest.NewRecorder()
	reqBody = `{"key": "4", "value": "new var", "ttl": 300}`
	req, _ = http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	//1st entry in tenant2

	w = httptest.NewRecorder()
	reqBody = `{"key": "4", "value": "new var1", "ttl": 300}`
	req, _ = http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant2", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	//2nd entry in tenant2
	w = httptest.NewRecorder()
	reqBody = `{"key": "5", "value": "new var2", "ttl": 300}`
	req, _ = http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant2", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Then, clear the cache entry in tenant1
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/cache/clear?system=inmemory&tenantID=tenant1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Ensure the cache entry is deleted in tenant1
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/3?system=inmemory&tenantID=tenant1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/4?system=inmemory&tenantID=tenant1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	//Ensure that the cache is not deleted in tenant2
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/5?system=inmemory&tenantID=tenant2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/4?system=inmemory&tenantID=tenant2", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
