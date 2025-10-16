package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 全局指标注册表
var (
	// HTTP请求指标
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	// 数据库指标
	dbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)

	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	dbConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections",
			Help: "Current number of database connections",
		},
	)

	// 插件指标
	pluginsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "plugins_total",
			Help: "Total number of registered plugins",
		},
	)

	pluginsEnabled = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "plugins_enabled",
			Help: "Number of enabled plugins",
		},
	)

	// 系统指标
	memoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	systemUptime = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_uptime_seconds",
			Help: "System uptime in seconds",
		},
	)

	// 错误指标
	errorCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "component"},
	)

	// 初始启动时间
	startTime = time.Now()
)

// MetricsManager 指标管理器
type MetricsManager struct{}

// NewMetricsManager 创建指标管理器实例
func NewMetricsManager() *MetricsManager {
	return &MetricsManager{}
}

// RegisterMetricsRouter 注册Prometheus指标导出路由
func (mm *MetricsManager) RegisterMetricsRouter(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// HTTPMonitoringMiddleware HTTP请求监控中间件
func (mm *MetricsManager) HTTPMonitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		httpRequestsTotal.WithLabelValues(method, path, http.StatusText(status)).Inc()
		httpRequestDuration.WithLabelValues(method, path, http.StatusText(status)).Observe(duration)
	}
}

// RecordHTTPRequest 记录HTTP请求（直接调用方式）
func RecordHTTPRequest(method, endpoint, status string, duration float64) {
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
}

// RecordDatabaseQuery 记录数据库查询
func RecordDatabaseQuery(operation, table string, duration float64) {
	dbQueriesTotal.WithLabelValues(operation, table).Inc()
	dbQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// UpdateDatabaseConnections 更新数据库连接数
func UpdateDatabaseConnections(count int) {
	dbConnections.Set(float64(count))
}

// UpdatePluginStats 更新插件统计信息
func UpdatePluginStats(total, enabled int) {
	pluginsTotal.Set(float64(total))
	pluginsEnabled.Set(float64(enabled))
}

// UpdateSystemMetrics 更新系统指标
func UpdateSystemMetrics() {
	// 更新系统运行时间
	uptime := time.Since(startTime).Seconds()
	systemUptime.Set(uptime)

	// TODO: 添加内存使用监控（需要使用runtime包）
	// var memStats runtime.MemStats
	// runtime.ReadMemStats(&memStats)
	// memoryUsage.Set(float64(memStats.Alloc))
}

// RecordError 记录错误
func RecordError(errorType, component string) {
	errorCount.WithLabelValues(errorType, component).Inc()
}

// StartMetricsUpdater 启动指标更新器
func (mm *MetricsManager) StartMetricsUpdater(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			UpdateSystemMetrics()
		}
	}()
}
