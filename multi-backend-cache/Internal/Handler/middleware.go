package handler

import (
	"multi-backend-cache/Internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ValidateTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Query("tenantID")
		logrus.Debugf("Received tenantID from parameter: %s", tenantID)
		if tenantID == "" || !isTenantValid(tenantID) {
			logrus.Warnf("Invalid or missing tenantID: %s", tenantID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Tenant Not Found"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isTenantValid(tenantID string) bool {
	for _, tenant := range config.AppConfig.TenantIDs {
		if tenant == tenantID {
			logrus.Debugf("Valid tenantID: %s", tenantID)
			return true
		}
	}
	logrus.Warnf("TenantID not found in configuration: %s", tenantID)
	return false
}

func ValidateCacheSystem() gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheSystem := c.Query("system")
		if cacheSystem == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cache system not provided"})
			c.Abort()
			return
		}

		cacheSystems := config.AppConfig.CacheSystems
		found := false
		for _, sys := range cacheSystems {
			if sys == cacheSystem {
				found = true
				break
			}
		}

		if !found {
			logrus.Warnf("Invalid cache system: %s", cacheSystem)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cache system"})
			c.Abort()
			return
		}
		logrus.Debugf("Valid cache system: %s", cacheSystem)
		c.Next()
	}
}
