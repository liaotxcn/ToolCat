// 热重载示例插件
package main

import (
	"fmt"
	"time"
	"toolcat/plugins/core"

	"github.com/gin-gonic/gin"
)

// 插件版本号，用于测试热重载是否生效
var pluginVersion = "1.0.0"

// HotReloadExample 热重载示例插件结构
type HotReloadExample struct {
	manager   *core.PluginManager
	startTime time.Time
}

// 实现Plugin接口的方法
func (p *HotReloadExample) Name() string {
	return "hot_reload_example"
}

func (p *HotReloadExample) Description() string {
	return "热重载功能示例插件，可在不重启应用的情况下动态更新"
}

func (p *HotReloadExample) Version() string {
	return pluginVersion
}

func (p *HotReloadExample) GetDependencies() []string {
	return []string{}
}

func (p *HotReloadExample) GetConflicts() []string {
	return []string{}
}

func (p *HotReloadExample) Init() error {
	p.startTime = time.Now()
	fmt.Println("热重载示例插件已初始化")
	return nil
}

func (p *HotReloadExample) Shutdown() error {
	fmt.Println("热重载示例插件已关闭")
	return nil
}

func (p *HotReloadExample) OnEnable() error {
	fmt.Println("热重载示例插件已启用")
	return nil
}

func (p *HotReloadExample) OnDisable() error {
	fmt.Println("热重载示例插件已禁用")
	return nil
}

func (p *HotReloadExample) GetRoutes() []core.Route {
	return []core.Route{
		{
			Path:        "/plugin/hotreload",
			Method:      "GET",
			Handler:     p.handleHotReloadTest,
			Description: "测试热重载功能的API",
		},
		{
			Path:        "/plugin/hotreload/info",
			Method:      "GET",
			Handler:     p.handlePluginInfo,
			Description: "获取插件信息",
		},
	}
}

func (p *HotReloadExample) RegisterRoutes(router *gin.Engine) {
	// 为了保持兼容，不做任何操作，使用GetRoutes替代
}

func (p *HotReloadExample) Execute(params map[string]interface{}) (interface{}, error) {
	return map[string]interface{}{
			"plugin":  "hot_reload_example",
			"version": pluginVersion,
			"status":  "running",
			"uptime":  time.Since(p.startTime).String(),
		},
		nil
}

func (p *HotReloadExample) GetDefaultMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

func (p *HotReloadExample) SetPluginManager(manager *core.PluginManager) {
	p.manager = manager
}

// 测试热重载功能的处理函数
func (p *HotReloadExample) handleHotReloadTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":   "热重载示例插件响应",
		"version":   pluginVersion,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"uptime":    time.Since(p.startTime).String(),
		"tip":       "修改代码并重新编译后，刷新此页面可看到版本变化",
	})
}

// 获取插件信息的处理函数
func (p *HotReloadExample) handlePluginInfo(c *gin.Context) {
	plugins := p.manager.GetAllPluginsInfo()
	c.JSON(200, gin.H{
		"plugin":             "hot_reload_example",
		"version":            pluginVersion,
		"all_plugins":        plugins,
		"hot_reload_enabled": "true",
	})
}

// 必须导出的插件构造函数，热加载机制通过此函数创建插件实例
func NewPlugin() core.Plugin {
	return &HotReloadExample{}
}

// main函数（编译为共享库时会被忽略，但保留以便测试）
func main() {
	fmt.Println("这是热重载示例插件的main函数，编译为共享库时会被忽略")
}
