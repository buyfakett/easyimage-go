package handler

import (
	"runtime"
	"sync"

	"easyimage_go/utils/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	memAlloc      prometheus.Gauge
	numGoroutines prometheus.Gauge
	totalAlloc    prometheus.Gauge
	// config handlers counters
	configChangeCounter prometheus.Counter
	configReadCounter   prometheus.Counter

	metricsInitialized bool
	metricsMutex       sync.Mutex
)

// ensureMetricsInitialized 确保指标已初始化
func ensureMetricsInitialized() {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	if metricsInitialized {
		return
	}

	memAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "go_memory_alloc_bytes",
		Help:        "Current memory allocation in bytes.",
		ConstLabels: prometheus.Labels{"server_name": config.Cfg.Server.Name},
	})
	numGoroutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "go_num_goroutines",
		Help:        "Number of Goroutines.",
		ConstLabels: prometheus.Labels{"server_name": config.Cfg.Server.Name},
	})
	totalAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "go_memory_total_alloc_bytes",
		Help:        "Total memory allocated in bytes.",
		ConstLabels: prometheus.Labels{"server_name": config.Cfg.Server.Name},
	})

	prometheus.MustRegister(memAlloc, numGoroutines, totalAlloc, configChangeCounter, configReadCounter)

	metricsInitialized = true
}

// updateMetrics 更新运行时指标
func updateMetrics() {
	ensureMetricsInitialized()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memAlloc.Set(float64(m.Alloc))
	numGoroutines.Set(float64(runtime.NumGoroutine()))
	totalAlloc.Set(float64(m.TotalAlloc))
}

// Metrics gin 版本
//
//	@Tags		监控
//	@Summary	普罗米修斯监控
//	@Router		/api/metrics [get]
func Metrics(c *gin.Context) {
	updateMetrics()
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
