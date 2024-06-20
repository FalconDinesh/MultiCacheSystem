package config

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	IsTenantBased         bool     `mapstructure:"IsTenantBased"`
	NumberOfTenants       string   `mapstructure:"NumberOfTenants"`
	TenantIDs             []string `mapstructure:"TenantIDs"`
	DefaultTTL            int      `mapstructure:"defaultTTL"`
	CacheSystems          []string `mapstructure:"CacheSystems"`
	MemoryUsagePercentage float64  `mapstructure:"MemoryUsagePercentage"`
	IP                    string   `mapstructure:"IP"`
	Redis      RedisConfig
    Memcache   MemcacheConfig
}

type RedisConfig struct {
    Address  string `mapstructure:"address"`
    Password string `mapstructure:"password"`
    Database int    `mapstructure:"database"`
}

type MemcacheConfig struct {
    Address    string `mapstructure:"address"`
    DefaultTTL int    `mapstructure:"defaultTTL"`
}

var AppConfig Config

func LoadConfig(configFile string) {
	absPath, err := filepath.Abs(configFile)
	if err != nil {
		logrus.Fatalf("Failed to get absolute path: %v", err)
	}
	logrus.Infof("Using config file: %s", absPath)

	viper.SetConfigFile(absPath)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Failed to read config file: %v", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		logrus.Fatalf("Failed to unmarshal config file: %v", err)
	}
	logrus.Infof("Config file loaded successfully: %+v", AppConfig)
}
