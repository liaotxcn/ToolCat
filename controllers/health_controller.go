package controllers

import (
	"fmt"
	"time"

	"weave/config"
	"weave/pkg"
	"weave/pkg/metrics"
	"weave/plugins"
	"weave/plugins/core"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthController 健康检查控制器
type HealthController struct{}

// GetHealth 全面健康检查
func (hc *HealthController) GetHealth(c *gin.Context) {
	// 开始时间
	startTime := time.Now()

	// 初始化健康检查结果
	result := gin.H{
		"status":      "ok",
		"timestamp":   time.Now().Unix(),
		"instance_id": config.Config.Server.InstanceID,
	}

	// 检查数据库连接健康状态
	dbHealth := checkDatabaseHealth()
	result["database"] = dbHealth

	// 检查插件系统健康状态
	pluginHealth := checkPluginHealth()
	result["plugins"] = pluginHealth

	// 检查整体系统健康状态
	overallStatus := "ok"
	if !dbHealth["healthy"].(bool) {
		overallStatus = "degraded"
	}

	for _, status := range pluginHealth["pluginStatuses"].([]gin.H) {
		if !status["healthy"].(bool) {
			overallStatus = "degraded"
			break
		}
	}

	result["status"] = overallStatus

	// 根据整体状态设置HTTP状态码
	statusCode := 200
	if overallStatus != "ok" {
		statusCode = 503
		// 使用统一错误码系统返回服务不可用错误
		serviceErr := pkg.NewServiceUnavailableError("System health is degraded", nil)
		serviceErr.WithDetails(map[string]interface{}{
			"database_healthy": dbHealth["healthy"].(bool),
			"plugin_count":     pluginHealth["pluginCount"].(int),
		})
		c.Error(serviceErr)
	}

	// 记录请求持续时间
	duration := time.Since(startTime).Seconds()
	pkg.Info("Health check completed",
		zap.Float64("duration", duration),
		zap.String("status", overallStatus))

	// 记录健康检查指标 - 使用数字字符串格式作为状态码标签
	metrics.RecordHTTPRequest("GET", "/health", fmt.Sprintf("%d", statusCode), duration)
	if !dbHealth["healthy"].(bool) {
		metrics.RecordError("database", "health_check")
	}

	// 更新插件统计指标
	totalPlugins := pluginHealth["pluginCount"].(int)
	enabledPlugins := 0
	for _, status := range pluginHealth["pluginStatuses"].([]gin.H) {
		if status["enabled"].(bool) {
			enabledPlugins++
		}
	}
	metrics.UpdatePluginStats(totalPlugins, enabledPlugins)

	c.JSON(statusCode, result)
}

// checkDatabaseHealth 检查数据库连接健康状态
func checkDatabaseHealth() gin.H {
	startTime := time.Now()
	db := pkg.DB

	// 执行简单的SQL查询来测试连接
	err := db.Exec("SELECT 1").Error
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		// 使用统一错误码系统创建数据库错误
		dbErr := pkg.NewDatabaseError("Database health check failed", err)
		dbErr.WithDetails(map[string]interface{}{
			"query": "SELECT 1",
		})
		pkg.Error("Database health check failed", zap.Error(dbErr))
		return gin.H{
			"healthy":      false,
			"error":        dbErr.Error(),
			"responseTime": duration,
		}
	}

	return gin.H{
		"healthy":      true,
		"responseTime": duration,
	}
}

// PluginHealthCheck 检查指定插件的健康状态
func (hc *HealthController) PluginHealthCheck(c *gin.Context) {
	pluginName := c.Param("name")
	startTime := time.Now()
	success := true

	// 查找插件信息
	allPluginsInfo := plugins.PluginManager.GetAllPluginsInfo()
	var targetPluginInfo *core.PluginInfo
	for _, info := range allPluginsInfo {
		if info.Plugin.Name() == pluginName {
			targetPluginInfo = &info
			break
		}
	}

	if targetPluginInfo == nil {
		metrics.RecordPluginError(pluginName, "health_check_not_found")
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "插件不存在",
		})
		return
	}

	// 检查插件状态
	status, exists := plugins.PluginManager.GetPluginStatus(pluginName)
	if !exists {
		status = "not_registered"
		success = false
		metrics.RecordPluginError(pluginName, "status_not_found")
	}

	healthy := targetPluginInfo.IsEnabled && status == "enabled"
	if !healthy && targetPluginInfo.IsEnabled {
		success = false
		metrics.RecordPluginError(pluginName, "health_check_failed")
	}

	// 记录执行时间和结果
	duration := time.Since(startTime)
	metrics.RecordPluginMethodCall(pluginName, "HealthCheck", success)
	metrics.RecordPluginExecution(pluginName, success, duration)

	c.JSON(200, gin.H{
		"name":    pluginName,
		"version": targetPluginInfo.Plugin.Version(),
		"enabled": targetPluginInfo.IsEnabled,
		"status":  status,
		"healthy": healthy,
	})
}

// checkPluginHealth 检查插件系统健康状态
func checkPluginHealth() gin.H {
	pluginStatuses := []gin.H{}
	allPluginsInfo := plugins.PluginManager.GetAllPluginsInfo()

	for _, pluginInfo := range allPluginsInfo {
		plugin := pluginInfo.Plugin
		startTime := time.Now()
		status, exists := plugins.PluginManager.GetPluginStatus(plugin.Name())
		if !exists {
			status = "not_registered"
			metrics.RecordPluginError(plugin.Name(), "status_not_found")
		}
		healthy := pluginInfo.IsEnabled && status == "enabled"
		if !healthy && pluginInfo.IsEnabled {
			metrics.RecordPluginError(plugin.Name(), "health_check_failed")
		}
		// 记录插件健康检查执行时间
		duration := time.Since(startTime)
		metrics.RecordPluginExecution(plugin.Name(), healthy, duration)

		pluginStatuses = append(pluginStatuses, gin.H{
			"name":    plugin.Name(),
			"version": plugin.Version(),
			"enabled": pluginInfo.IsEnabled,
			"status":  status,
			"healthy": healthy,
		})
	}

	return gin.H{
		"pluginCount":    len(allPluginsInfo),
		"pluginStatuses": pluginStatuses,
	}
}
