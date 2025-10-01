package plugins

import (
	"github.com/gin-gonic/gin"
)

// Plugin 插件接口定义
type Plugin interface {
	// Name 返回插件名称
	Name() string

	// Description 返回插件描述
	Description() string

	// Version 返回插件版本
	Version() string

	// Init 初始化插件
	Init() error

	// Shutdown 关闭插件
	Shutdown() error

	// RegisterRoutes 注册插件路由
	RegisterRoutes(router *gin.Engine)

	// Execute 执行插件功能
	Execute(params map[string]interface{}) (interface{}, error)
}

// PluginManager 插件管理器
var PluginManager = &pluginManager{
	plugins: make(map[string]Plugin),
}

type pluginManager struct {
	plugins map[string]Plugin
}

// Register 注册插件
func (pm *pluginManager) Register(plugin Plugin) error {
	name := plugin.Name()
	if _, exists := pm.plugins[name]; exists {
		return nil // 插件已存在，不重复注册
	}

	// 初始化插件
	if err := plugin.Init(); err != nil {
		return err
	}

	pm.plugins[name] = plugin
	return nil
}

// Unregister 注销插件
func (pm *pluginManager) Unregister(name string) error {
	plugin, exists := pm.plugins[name]
	if !exists {
		return nil // 插件不存在，无需注销
	}

	// 关闭插件
	if err := plugin.Shutdown(); err != nil {
		return err
	}

	delete(pm.plugins, name)
	return nil
}

// GetPlugin 获取插件
func (pm *pluginManager) GetPlugin(name string) (Plugin, bool) {
	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// ListPlugins 列出所有插件
func (pm *pluginManager) ListPlugins() []string {
	names := make([]string, 0, len(pm.plugins))
	for name := range pm.plugins {
		names = append(names, name)
	}
	return names
}

// ExecutePlugin 执行插件功能
func (pm *pluginManager) ExecutePlugin(name string, params map[string]interface{}) (interface{}, error) {
	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, nil // 插件不存在
	}

	return plugin.Execute(params)
}