package metrics

import (
	"fmt"
	"runtime"
	"strconv"
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

	// 插件性能指标
	PluginExecutionCount    *prometheus.CounterVec
	PluginExecutionDuration *prometheus.HistogramVec
	PluginMethodCalls       *prometheus.CounterVec
	PluginErrors            *prometheus.CounterVec
	PluginMemoryUsage       *prometheus.GaugeVec
	PluginReloads           *prometheus.CounterVec

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

// init 初始化所有指标
func init() {
	// 插件性能指标初始化
	PluginExecutionCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "plugin_execution_total",
			Help: "Total number of plugin executions",
		},
		[]string{"plugin_name", "success"},
	)

	PluginExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "plugin_execution_duration_seconds",
			Help:    "Plugin execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"plugin_name", "success"},
	)

	PluginMethodCalls = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "plugin_method_calls_total",
			Help: "Total number of plugin method calls",
		},
		[]string{"plugin_name", "method", "success"},
	)

	PluginErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "plugin_errors_total",
			Help: "Total number of plugin errors",
		},
		[]string{"plugin_name", "error_type"},
	)

	PluginMemoryUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "plugin_memory_usage_bytes",
			Help: "Memory usage per plugin in bytes",
		},
		[]string{"plugin_name"},
	)

	PluginReloads = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "plugin_reloads_total",
			Help: "Total number of plugin reloads",
		},
		[]string{"plugin_name", "success"},
	)
}

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
		statusCode := fmt.Sprintf("%d", status)

		httpRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
		httpRequestDuration.WithLabelValues(method, path, statusCode).Observe(duration)
	}
}

// RecordHTTPRequest 记录HTTP请求（直接调用方式）
func RecordHTTPRequest(method, endpoint, status string, duration float64) {
	// 确保status是状态码的数字字符串表示
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

// RecordPluginExecution 记录插件执行情况
func RecordPluginExecution(pluginName string, success bool, duration time.Duration) {
	successStr := strconv.FormatBool(success)
	PluginExecutionCount.WithLabelValues(pluginName, successStr).Inc()
	PluginExecutionDuration.WithLabelValues(pluginName, successStr).Observe(duration.Seconds())
}

// RecordPluginMethodCall 记录插件方法调用
func RecordPluginMethodCall(pluginName, method string, success bool) {
	successStr := strconv.FormatBool(success)
	PluginMethodCalls.WithLabelValues(pluginName, method, successStr).Inc()
}

// RecordPluginError 记录插件错误
func RecordPluginError(pluginName, errorType string) {
	PluginErrors.WithLabelValues(pluginName, errorType).Inc()
}

// UpdatePluginMemoryUsage 更新插件内存使用
func UpdatePluginMemoryUsage(pluginName string, memoryBytes int64) {
	PluginMemoryUsage.WithLabelValues(pluginName).Set(float64(memoryBytes))
}

// RecordPluginReload 记录插件重载
func RecordPluginReload(pluginName string, success bool) {
	successStr := strconv.FormatBool(success)
	PluginReloads.WithLabelValues(pluginName, successStr).Inc()
}

// UpdateSystemMetrics 更新系统指标
func UpdateSystemMetrics() {
	// 更新系统运行时间
	uptime := time.Since(startTime).Seconds()
	systemUptime.Set(uptime)

	// 更新内存使用监控
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryUsage.Set(float64(memStats.Alloc))
}

// RecordError 记录错误
func RecordError(errorType, component string) {
	errorCount.WithLabelValues(errorType, component).Inc()
}

// PluginMonitoringMiddleware 创建插件监控中间件
func PluginMonitoringMiddleware(pluginName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		success := true

		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			success = false
		}

		duration := time.Since(start).Seconds()
		handlerName := fmt.Sprintf("%s %s", method, path)

		// 记录方法调用
		RecordPluginMethodCall(pluginName, handlerName, success)
		// 记录执行时间
		RecordPluginExecution(pluginName, success, time.Duration(duration*float64(time.Second)))

		// 如果失败，记录错误
		if !success {
			RecordPluginError(pluginName, "handler_error")
		}
	}
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
