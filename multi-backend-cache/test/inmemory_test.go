package test

import (
	"bytes"
	handler "multi-backend-cache/Internal/Handler"
	"multi-backend-cache/Internal/cache"
	"multi-backend-cache/Internal/config"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIsExpired(t *testing.T) {
	a, _ := time.Parse(time.RFC3339Nano, "2024-06-17T13:38:18.616332+05:30")
	got := cache.IsExpired(a)
	want := true
	assert.Equal(t, want, got)
}

func TestCalculateSize(t *testing.T) {
	// Create a CacheData instance
	cacheData := &cache.CacheData{
		Key:        "test_key",
		Value:      "test_value",
		TTL:        60,
		ExpiryTime: time.Now().Add(time.Second * 10),
	}

	// Call the calculateSize function
	size := cache.CalculateSize(cacheData)
	expectedSize := int(reflect.TypeOf(*cacheData).Size())
	if size != expectedSize {
		t.Errorf("Expected size: %d, got: %d", expectedSize, size)
	}
}

// Function to set up Inmemory router, with capacity as 300 bytes and TTL to 10s
func setupInMemoryRouter() *gin.Engine {
	// config.LoadConfig("../Internal/config/config.yaml")
	config.AppConfig.IsTenantBased = false
	inmemorycache := cache.NewFixedTenantsCaches(false, 300, 10)
	cacheSystemType := handler.NewServer(inmemorycache, nil, nil)
	router := gin.Default()

	router.GET("/cache/:key", cacheSystemType.GetCacheHandler)
	router.POST("/cache", cacheSystemType.SetCacheHandler)
	router.DELETE("/cache/:key", cacheSystemType.DeleteCacheHandler)
	// router.GET("/cache/TTL/:key", cacheSystemType.GetCacheWithTTLHandler)
	router.PUT("/cache/clear", cacheSystemType.ClearCacheHandler)

	return router
}

// Test for set function
func TestInMemPostCacheHandler(t *testing.T) {
	router := setupInMemoryRouter()

	t.Run("Valid Data", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"key": "1", "value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Key not passed", func(t *testing.T) {
		invalidData := `{"key": "","value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory", bytes.NewBufferString(invalidData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("Invalid Data Format (key as int)", func(t *testing.T) {
		invalidData := `{"key": 1 ,"value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
		req, _ := http.NewRequest("POST", "/cache?system=inmemory", bytes.NewBufferString(invalidData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

// Test for get function
func TestInMemGetCacheHandler(t *testing.T) {
	router := setupInMemoryRouter()

	// First, post a cache entry
	w := httptest.NewRecorder()
	reqBody := `{"key": "2", "value": {"name":"session"},"ttl":200}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Run("Valid Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/2?system=inmemory", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `{"name":"session"}`)
	})

	t.Run("InValid Key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=inmemory", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	//Ensure that passing tenants in URL, does not make a differnece
	// t.Run("Sending tenant id in url", func(t *testing.T) {
	// 	invalidData := `{"key": "1" ,"value": {"id":"12345","name":"Abcd"}, "ttl": 300}`
	// 	req, _ := http.NewRequest("POST", "/cache?system=inmemory&tenantID=tenant1", bytes.NewBufferString(invalidData))
	// 	req.Header.Set("Content-Type", "application/json")
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)
	// 	assert.Equal(t, http.StatusOK, w.Code)
	// })

	// t.Run("getting the valid cache with tenantId", func(t *testing.T) {
	// 	req, _ := http.NewRequest("GET", "/cache/1?system=inmemory&tenantID=tenant1", nil)
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Contains(t, w.Body.String(), `{"id":"12345","name":"Abcd"}`)
	// })
	// t.Run("getting a valid cache without tenantId", func(t *testing.T) {
	// 	req, _ := http.NewRequest("GET", "/cache/1?system=inmemory&tenantID=tenant1", nil)
	// 	w := httptest.NewRecorder()
	// 	router.ServeHTTP(w, req)

	// 	assert.Equal(t, http.StatusOK, w.Code)
	// 	assert.Contains(t, w.Body.String(), `{"id":"12345","name":"Abcd"}`)
	// })

}

// Test to check deleteAPI
func TestInMemDeleteCacheHandler(t *testing.T) {
	router := setupInMemoryRouter()

	// First, post a cache entry
	w := httptest.NewRecorder()
	reqBody := `{"key": "3", "value": "cache", "ttl": 300}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Run("Valid Key", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/cache/3?system=inmemory", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// assert.Contains(t, w.Body.String(), "deleted")
	})
	t.Run("Check if deleted", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/3?system=inmemory", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("InValid Key", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/cache/4?system=inmemory", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		// assert.Contains(t, w.Body.String(), "deleted")
	})
}

// Test for data expiry with default TTL
func TestInMemoryDataExpiryWithDefaultTTL(t *testing.T) {
	router := setupInMemoryRouter()

	w := httptest.NewRecorder()
	reqBodyExpiry := `{"key": "10", "value": "sessions"}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory", strings.NewReader(reqBodyExpiry))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	t.Run("Expiry check", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/cache/10?system=inmemory", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "sessions")

		time.Sleep(10 * time.Second)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/cache/10?system=inmemory", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

	})
}

// Test to check clear API
func TestInMemClearCacheHandler(t *testing.T) {
	router := setupInMemoryRouter()

	// First, post a cache entry
	w := httptest.NewRecorder()
	reqBody := `{"key": "3", "value": "cache", "ttl": 300}`
	req, _ := http.NewRequest("POST", "/cache?system=inmemory", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	//2nd cache entry
	w = httptest.NewRecorder()
	reqBody = `{"key": "4", "value": "new var", "ttl": 300}`
	req, _ = http.NewRequest("POST", "/cache?system=inmemory", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Then, delete the cache entry
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/cache/clear?system=inmemory", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Ensure the cache entry is deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/3?system=inmemory", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/cache/4?system=inmemory", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// unc TestGetCacheWithTTLHandler(t *testing.T) {
// 	router := setupInMemoryRouter()

// 	// First, post a cache entry
// 	w := httptest.NewRecorder()
// 	reqBody := `{"key": "2", "value": "session", "ttl": 300}`
// 	req, _ := http.NewRequest("POST", "/cache?system=inmemory", strings.NewReader(reqBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Then, get the cache entry
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/cache/TTL/2?system=inmemory", nil)
// 	router.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// Parse the response body
// 	type CacheDataTemp struct {
// 		Key        string      `json:"key"`
// 		Value      interface{} `json:"value"`
// 		TTL        float64     `json:"ttl"` // Changed TTL type to float64 for unmarshaling
// 		ExpiryTime time.Time   `json:"expiry_time"`
// 	}
// 	var response CacheDataTemp
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)

// 	expectedExpiryTime := time.Now().Add(300 * time.Second)

// 	// Allow a margin of error in the expiry time due to processing delays
// 	marginOfError := 2 * time.Second
// 	assert.WithinDuration(t, expectedExpiryTime, response.ExpiryTime, marginOfError)
// 	tolerance := 2.0
// 	assert.Equal(t, "session", response.Value)
// 	expectedTTL := expectedExpiryTime.Sub(time.Now()).Seconds()
// 	assert.InDelta(t, time.Duration(expectedTTL), time.Duration(response.TTL), tolerance) //added tolerance for delays
// }
