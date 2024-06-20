package metrices

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Metrics struct {
	totalRequests prometheus.Counter
	totalLatency  prometheus.Histogram
	successHits   prometheus.Counter
	failureHits   prometheus.Counter
}

func NewMetrics() *Metrics {
	m := &Metrics{
		totalRequests: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Total number of HTTP requests",
		}),
		totalLatency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "request_latency_seconds",
			Help:    "Histogram of latencies for HTTP requests",
			Buckets: prometheus.DefBuckets,
		}),
		successHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "success_hits",
			Help: "Number of successful HTTP requests",
		}),
		failureHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "failure_hits",
			Help: "Number of failed HTTP requests",
		}),
	}

	prometheus.MustRegister(m.totalRequests, m.totalLatency, m.successHits, m.failureHits)
	return m
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics collection for the /metrics endpoint
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		startTime := time.Now()
		m.totalRequests.Inc()

		c.Next() // process the request

		latency := time.Since(startTime)
		m.totalLatency.Observe(latency.Seconds())
		if c.Writer.Status() >= http.StatusOK && c.Writer.Status() < http.StatusMultipleChoices {
			m.successHits.Inc()
		} else if c.Writer.Status() == http.StatusInternalServerError {
			m.failureHits.Inc()
		}

		// Log the metrics (optional, for demonstration purposes)
		// log.Printf("Status: %d, Latency: %s", c.Writer.Status(), latency)

		// Log the metrics
		logrus.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"latency": latency,
		}).Info("Request processed")
	}
}
