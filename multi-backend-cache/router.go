package main

import (
	"fmt"
	"log"
	handler "multi-backend-cache/Internal/Handler"
	"multi-backend-cache/Internal/cache"
	"multi-backend-cache/Internal/config"
	"multi-backend-cache/Internal/metrices"
	_ "multi-backend-cache/docs"
	"time"

	"github.com/pbnjay/memory"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// var inmemorycache *cache.LRUCache
var tenantCaches *cache.FixedTenantsCaches


func main() {
	// Load configuration
	config.LoadConfig("./Internal/config/config.yaml")

	defaultTTL := config.AppConfig.DefaultTTL
	
	// Initialize Redis cache using config.AppConfig
    redisConfig := config.AppConfig.Redis
    redisCache := cache.NewRedisCache(redisConfig.Address, redisConfig.Password, redisConfig.Database, time.Duration(defaultTTL)*time.Second)

    // Initialize Memcache with TTL conversion
    memcacheConfig := config.AppConfig.Memcache
    memCache := cache.NewMemCache(memcacheConfig.Address, int32(memcacheConfig.DefaultTTL)) // Convert int to int32


	totalCacheMemory := int(float64(memory.TotalMemory()) * config.AppConfig.MemoryUsagePercentage)

	// Initialize tenant-specific in-memory LRUCaches
	// Each tenant has a cache with a fixed capacity
	isTenantBased := config.AppConfig.IsTenantBased

	tenantCaches = cache.NewFixedTenantsCaches(isTenantBased, totalCacheMemory, time.Duration(defaultTTL))
	cacheSystem := handler.NewServer(tenantCaches, redisCache, memCache)

	router := gin.Default()

	host := fmt.Sprintf("http://%s:8080/swagger/doc.json", config.AppConfig.IP)
	url := ginSwagger.URL(host)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))


	metrics := metrices.NewMetrics()

	router.Use(metrics.Middleware())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(handler.ValidateCacheSystem())

	// Middleware for "inmemory" system
	router.Use(func(c *gin.Context) {
		if system := c.Query("system"); system == "inmemory" && isTenantBased {
			handler.ValidateTenant()(c)
		}
		c.Next()
	})

	// Cache System routes
	router.GET("/cache/:key", cacheSystem.GetCacheHandler)
	//router.GET("/cache/TTL/:key", cacheSystem.GetCacheWithTTLHandler)
	router.POST("/cache", cacheSystem.SetCacheHandler)
	router.DELETE("/cache/:key", cacheSystem.DeleteCacheHandler)
	router.PUT("/cache/clear", cacheSystem.ClearCacheHandler)

	// Start the HTTP server
	addr := ":8080"
	log.Printf("Server started at %s\n", addr)
	log.Fatal(router.Run(addr))
}
