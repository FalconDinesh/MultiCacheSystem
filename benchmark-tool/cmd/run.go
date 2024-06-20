/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"net/http"
	"sync"
	"time"

	utils "benchmark-tool/utils"

	"github.com/spf13/cobra"
)

var (
	hostname string
	port     int
	server   string
	testType string
	clients  int
	num      int
	size     int
)

var (
	cacheContent string
	baseSetUrl   string
	baseGetUrl   string
	latencies    []time.Duration
)

// getKey simulates getting a key from the cache
func getKey(client *http.Client, key string) {
	// Simulate getting a key by sending a GET request
	url := fmt.Sprintf(baseGetUrl, key)
	startTime := time.Now()
	resp, err := client.Get(url)
	latencies = append(latencies, time.Since(startTime))
	if err != nil {
		fmt.Printf("Get request error: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

// setKey simulates setting a key in the cache
func setKey(client *http.Client, key, value string) {
	// Simulate setting a key by sending a POST request
	startTime := time.Now()
	resp, err := client.Post(baseSetUrl, "application/json", bytes.NewBuffer(utils.PrepareCacheData(key, value)))
	latencies = append(latencies, time.Since(startTime))
	if err != nil {
		fmt.Printf("Set request error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Set request failed: %s\n", resp.Status)
	}
}

// deleteKey simulates deleting a key from the cache
func deleteKey(client *http.Client, key string) {
	// Simulate deleting a key by sending a DELETE request
	url := fmt.Sprintf(baseGetUrl, key)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Printf("Delete request error: %v\n", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Delete request error: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

/*Worker for each client to send n number of requests to the cache server
 */
func runWorker(clientID int, wg *sync.WaitGroup, keys []string, client *http.Client, totalLen int) {
	defer wg.Done()
	for j := 0; j < num; j++ {
		keyIndex := clientID*num + j
		if keyIndex >= totalLen {
			fmt.Printf("Index out of range: %d\n", keyIndex)
			continue
		}
		key := keys[keyIndex]
		switch testType {
		case "set":
			setKey(client, key, cacheContent)
		case "get":
			getKey(client, key)
		case "delete":
			// deleteKey(client, key)
		default:
			fmt.Printf("Unknown test type: %s\n", testType)
		}
	}
}

// runBenchmark executes the benchmark based on the specified test type
func runBenchmark(cmd *cobra.Command) {
	var wg sync.WaitGroup

	assignFlagVal(cmd)
	keys := utils.GenerateKeys(num * clients)
	totalLen := len(keys)

	client := &http.Client{}
	start := time.Now()
	for i := 0; i < clients; i++ {
		wg.Add(1)
		go runWorker(i, &wg, keys, client, totalLen)
	}
	wg.Wait()
	utils.SortLatencies(latencies)
	elapsed := time.Since(start)
	printSummary(elapsed)
}

func assignFlagVal(cmd *cobra.Command) {
	hostname = cmd.Flag("hostname").Value.String()
	port, _ = cmd.Flags().GetInt("port")
	server = cmd.Flag("server").Value.String()
	testType = cmd.Flag("test").Value.String()
	clients, _ = cmd.Flags().GetInt("clients")
	num, _ = cmd.Flags().GetInt("num")
	size, _ = cmd.Flags().GetInt("size")

	baseSetUrl = fmt.Sprintf("http://%s:%d/cache?system=%s", hostname, port, server)
	baseGetUrl = fmt.Sprintf("http://%s:%d/cache/%s?system=%s", hostname, port, "%s", server)
	cacheContent = utils.GenerateValue(size)
	latencies = nil
}

// printSummary outputs the results of the benchmark
func printSummary(elapsed time.Duration) {
	totalRequests := num * clients
	fmt.Printf("\n====== %s-%s ======\n\n", strings.ToUpper(testType), strings.ToUpper(server))
	fmt.Printf("%d requests completed in %v\n", totalRequests, elapsed)
	fmt.Printf("%d parallel clients\n", clients)
	fmt.Printf("%d bytes payload\n", size)
	fmt.Printf("keep alive: 1\n")

	fmt.Printf("\nThroughput: %.2f requests per second\n", float64(totalRequests)/float64(elapsed.Seconds()))

	// Placeholder for actual latency distribution calculation
	avgLatency := utils.AverageDuration(latencies)
	minLatency := latencies[0]
	p50 := utils.Percentile(latencies, 50)
	p95 := utils.Percentile(latencies, 95)
	p99 := utils.Percentile(latencies, 99)
	maxLatency := latencies[len(latencies)-1]

	fmt.Printf("\nLatency summary (msec):\n")
	fmt.Printf("avg     min     p50     p95     p99     max\n")
	fmt.Printf("%.3v     %.3v     %.3v     %.3v     %.3v     %.3v\n", avgLatency, minLatency, p50, p95, p99, maxLatency)

}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Command to execute the benchmark tool",
	Long: `This setup allows you to run a benchmark for different cache operations (set, get, delete) 
	with configurable parameters such as the number of clients, number of requests, payload size, 
	and whether to use multithreading. Adjust the HTTP request logic to match your actual server endpoints 
	and methods.`,
	Run: func(cmd *cobra.Command, args []string) {
		runBenchmark(cmd)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("hostname", "h", "127.0.0.1", "Hostname")
	runCmd.Flags().IntP("port", "p", 8080, "Server post")
	runCmd.Flags().StringP("server", "s", "redis", "Cache server to use - available options: redis, inmemory and memcache")
	runCmd.Flags().StringP("test", "t", "set", "Type of test to run")
	runCmd.Flags().IntP("clients", "c", 20, "Number of parallel connections")
	runCmd.Flags().IntP("num", "n", 1000, "Number of requests")
	runCmd.Flags().IntP("size", "d", 3, "Data size in bytes")
}
