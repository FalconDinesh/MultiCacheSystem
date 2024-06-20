package utils

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type CacheData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// generateValue creates a value of the specified size
func GenerateValue(size int) string {
	value := make([]byte, size)
	for i := range value {
		value[i] = 'a' // or any character you prefer
	}
	return string(value)
}

// generateKeys creates a slice of keys to be used in the benchmark
func GenerateKeys(n int) []string {
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("u-%d", i)
	}
	return keys
}

func PrepareCacheData(key string, value string) []byte {
	// Encode the key-value pair into JSON
	jsonData, err := json.Marshal(&CacheData{Key: key, Value: value})
	if err != nil {
		fmt.Printf("JSON encoding error: %v\n", err)
	}
	return jsonData
}

func SortLatencies(latencies []time.Duration){
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
}

func AverageDuration(latencies []time.Duration) time.Duration {
	var sum time.Duration
	for _, v := range latencies {
		sum += v
	}
	return time.Duration(int64(sum) / int64(len(latencies)))
}

func Percentile(latencies []time.Duration, p int) time.Duration {
    index := (p * len(latencies) / 100) - 1
    return latencies[index]
}
