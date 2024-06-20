package test

import (
	"bytes"
	handler "multi-backend-cache/Internal/Handler"
	"multi-backend-cache/Internal/cache"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRedisRouter() *gin.Engine {
	redisCache := cache.NewRedisCache("localhost:6379", "", 0, 10*time.Second)
	cacheSystemType := handler.NewServer(nil, redisCache, nil)
	router := gin.Default()

	router.GET("/cache/:key", cacheSystemType.GetCacheHandler)
	router.POST("/cache", cacheSystemType.SetCacheHandler)
	router.DELETE("/cache/:key", cacheSystemType.DeleteCacheHandler)
	// router.GET("/cache/TTL/:key", cacheSystemType.GetCacheWithTTLHandler)
	router.PUT("/cache/clear", cacheSystemType.ClearCacheHandler)

	return router
}

func TestRedisPostCacheHandler(t *testing.T) {
	router := setupRedisRouter()

	t.Run("Valid Data", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"key": "1", "value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=redis", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("No key passed", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"key": "","value": {"id":"12","name":"Abcddds"}, "ttl": 300}` //no key is passed
		req, _ := http.NewRequest("POST", "/cache?system=redis", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid Data Format (key as int)", func(t *testing.T) {
		invalidData := `{"key": 1 ,"value": {"id":"12345","name":"Abcd"}, "ttl": 300}` 
		req, _ := http.NewRequest("POST", "/cache?system=redis", bytes.NewBufferString(invalidData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

// Test for get function
func TestRedisGetCacheHandler(t *testing.T) {
	router := setupRedisRouter()

	// First, post a cache entry
	w := httptest.NewRecorder()
	reqBody := `{"key": "2", "value": "session", "ttl": 300}`
	req, _ := http.NewRequest("POST", "/cache?system=redis", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Then, get the cache entry

	t.Run("Valid Key With Cache Hit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/2?system=redis", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `session`)
	})

	t.Run("InValid Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=redis", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Test for data expiry with default TTL
func TestRedisDataExpiryWithDefaultTTL(t *testing.T) {
	router := setupRedisRouter()

	w := httptest.NewRecorder()
	reqBody := `{"key": "2", "value": "session"}`
	req, _ := http.NewRequest("POST", "/cache?system=redis", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Run("Expiry check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/2?system=redis", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "session")

		time.Sleep(10 * time.Second)

		req, _ = http.NewRequest("GET", "/cache/2?system=redis", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

	})
}

// Test to check deleteAPI
func TestRedisDeleteCacheHandler(t *testing.T) {
	router := setupRedisRouter()

	// First, post a cache entry
	w := httptest.NewRecorder()
	reqBody := `{"key": "3", "value": "cache", "ttl": 300}`
	req, _ := http.NewRequest("POST", "/cache?system=redis", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Then, delete the cache entry

	t.Run("Valid Key", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/cache/3?system=redis", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Check if deleted", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=redis", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InValid Key", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/cache/4?system=redis", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		// fmt.Println("body", w.Body.String())
		assert.Contains(t, w.Body.String(), `{"error":"Cache not Found - Failed to delete cache"}`)
	})

}

// Test to check clear API
func TestClearCacheHandler(t *testing.T) {
	router := setupRedisRouter()

	// First, post a cache entry
	w := httptest.NewRecorder()
	reqBody := `{"key": "3", "value": "cache", "ttl": 300}`
	req, _ := http.NewRequest("POST", "/cache?system=redis", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	//2nd cache entry
	w = httptest.NewRecorder()
	reqBody = `{"key": "4", "value": "new var", "ttl": 300}`
	req, _ = http.NewRequest("POST", "/cache?system=redis", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Then, delete the cache entry
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/cache/clear?system=redis", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Ensure the cache entry is deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/3?system=redis", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/4?system=redis", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
