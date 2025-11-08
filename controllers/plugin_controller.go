package controllers

import (
	"net/http"
	"weave/plugins"

	"github.com/gin-gonic/gin"
)

// PluginController 插件控制器
// 用于处理插件相关的API请求
type PluginController struct{}

// GetAllPlugins 获取所有插件信息
// @Summary 获取所有插件信息
// @Description 获取系统中注册的所有插件信息，包括启用状态
// @Tags 插件管理
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/plugins [get]
func (pc *PluginController) GetAllPlugins(c *gin.Context) {
	pluginsInfo := plugins.PluginManager.GetAllPluginsInfo()

	// 准备响应数据
	response := make([]map[string]interface{}, 0, len(pluginsInfo))
	for _, info := range pluginsInfo {
		pluginData := map[string]interface{}{
			"name":         info.Plugin.Name(),
			"description":  info.Plugin.Description(),
			"version":      info.Plugin.Version(),
			"enabled":      info.IsEnabled,
			"dependencies": info.Dependencies,
			"conflicts":    info.Conflicts,
		}
		response = append(response, pluginData)
	}

	c.JSON(http.StatusOK, response)
}

// EnablePlugin 启用插件
// @Summary 启用插件
// @Description 启用指定的插件
// @Tags 插件管理
// @Security BearerAuth
// @Param name path string true "插件名称"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/plugins/{name}/enable [post]
func (pc *PluginController) EnablePlugin(c *gin.Context) {
	pluginName := c.Param("name")

	if err := plugins.PluginManager.EnablePlugin(pluginName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "插件启用成功", "plugin": pluginName})
}

// DisablePlugin 禁用插件
// @Summary 禁用插件
// @Description 禁用指定的插件
// @Tags 插件管理
// @Security BearerAuth
// @Param name path string true "插件名称"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/plugins/{name}/disable [post]
func (pc *PluginController) DisablePlugin(c *gin.Context) {
	pluginName := c.Param("name")

	if err := plugins.PluginManager.DisablePlugin(pluginName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "插件禁用成功", "plugin": pluginName})
}

// ReloadPlugin 重载插件
// @Summary 重载插件
// @Description 重载指定的插件（先禁用再启用）
// @Tags 插件管理
// @Security BearerAuth
// @Param name path string true "插件名称"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/plugins/{name}/reload [post]
func (pc *PluginController) ReloadPlugin(c *gin.Context) {
	pluginName := c.Param("name")

	if err := plugins.PluginManager.ReloadPlugin(pluginName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "插件重载成功", "plugin": pluginName})
}

// GetPluginStatus 获取插件状态
// @Summary 获取插件状态
// @Description 获取指定插件的详细状态信息
// @Tags 插件管理
// @Security BearerAuth
// @Param name path string true "插件名称"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/v1/plugins/{name}/status [get]
func (pc *PluginController) GetPluginStatus(c *gin.Context) {
	pluginName := c.Param("name")

	status, exists := plugins.PluginManager.GetPluginStatus(pluginName)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "插件不存在", "plugin": pluginName})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status, "plugin": pluginName})
}

// GetDependencyGraph 获取插件依赖图
// @Summary 获取插件依赖图
// @Description 获取所有插件的依赖关系图
// @Tags 插件管理
// @Security BearerAuth
// @Success 200 {object} map[string][]string
// @Router /api/v1/plugins/dependency-graph [get]
func (pc *PluginController) GetDependencyGraph(c *gin.Context) {
	dependencyGraph := plugins.PluginManager.GetDependencyGraph()
	c.JSON(http.StatusOK, dependencyGraph)
}
