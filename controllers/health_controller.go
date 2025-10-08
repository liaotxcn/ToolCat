package controllers

import (
	"time"
	"toolcat/pkg"
	"toolcat/plugins"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthController 健康检查控制器
type HealthController struct{}

// GetHealth 全面健康检查
func (hc *HealthController) GetHealth(c *gin.Context) {
	// 初始化健康检查结果
	result := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
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
	if overallStatus == "ok" {
		c.JSON(200, result)
	} else {
		c.JSON(503, result)
	}
}

// checkDatabaseHealth 检查数据库连接健康状态
func checkDatabaseHealth() gin.H {
	startTime := time.Now()
	db := pkg.DB

	// 执行简单的SQL查询来测试连接
	err := db.Exec("SELECT 1").Error
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		pkg.Error("Database health check failed", zap.Error(err))
		return gin.H{
			"healthy":      false,
			"error":        err.Error(),
			"responseTime": duration,
		}
	}

	return gin.H{
		"healthy":      true,
		"responseTime": duration,
	}
}

// checkPluginHealth 检查插件系统健康状态
func checkPluginHealth() gin.H {
	pluginStatuses := []gin.H{}
	allPluginsInfo := plugins.PluginManager.GetAllPluginsInfo()

	for _, pluginInfo := range allPluginsInfo {
		plugin := pluginInfo.Plugin
		status, exists := plugins.PluginManager.GetPluginStatus(plugin.Name())
		if !exists {
			status = "not_registered"
		}
		pluginStatuses = append(pluginStatuses, gin.H{
			"name":    plugin.Name(),
			"version": plugin.Version(),
			"enabled": pluginInfo.IsEnabled,
			"status":  status,
			"healthy": pluginInfo.IsEnabled && status == "enabled",
		})
	}

	return gin.H{
		"pluginCount":    len(allPluginsInfo),
		"pluginStatuses": pluginStatuses,
	}
}
