package plugins

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// SampleDependentPlugin 示例依赖插件
type SampleDependentPlugin struct {
	pluginManager *pluginManager
}

// NewSampleDependentPlugin 创建新的示例依赖插件
func NewSampleDependentPlugin() *SampleDependentPlugin {
	return &SampleDependentPlugin{}
}

// Name 返回插件名称
func (p *SampleDependentPlugin) Name() string {
	return "sample_dependent"
}

// Description 返回插件描述
func (p *SampleDependentPlugin) Description() string {
	return "一个依赖其他插件的示例插件"
}

// Version 返回插件版本
func (p *SampleDependentPlugin) Version() string {
	return "1.0.0"
}

// GetDependencies 返回依赖的插件
func (p *SampleDependentPlugin) GetDependencies() []string {
	return []string{"sample_optimized", "hello_plugin"} // 依赖 sample_optimized 和 hello_plugin 插件
}

// GetConflicts 返回冲突的插件
func (p *SampleDependentPlugin) GetConflicts() []string {
	return []string{} // 与其他插件无冲突
}

// SetPluginManager 设置插件管理器
func (p *SampleDependentPlugin) SetPluginManager(manager *pluginManager) {
	p.pluginManager = manager
}

// Init 初始化插件
func (p *SampleDependentPlugin) Init() error {
	// 在初始化时可以访问依赖的插件
	if optimizedPlugin, exists := p.pluginManager.GetPlugin("sample_optimized"); exists {
		fmt.Printf("成功访问依赖的插件: %s\n", optimizedPlugin.Name())
	}
	return nil
}

// Shutdown 关闭插件
func (p *SampleDependentPlugin) Shutdown() error {
	return nil
}

// GetRoutes 获取路由定义
func (p *SampleDependentPlugin) GetRoutes() []Route {
	return []Route{
		{
			Path:         "/",
			Method:       "GET",
			Handler:      p.handleIndex,
			Description:  "示例依赖插件首页",
			AuthRequired: false,
		},
		{
			Path:         "/use-dependency",
			Method:       "GET",
			Handler:      p.handleUseDependency,
			Description:  "使用依赖的插件功能",
			AuthRequired: false,
		},
		{
			Path:         "/dependencies",
			Method:       "GET",
			Handler:      p.handleGetDependencies,
			Description:  "获取插件依赖信息",
			AuthRequired: false,
		},
	}
}

// RegisterRoutes 注册路由（旧版接口）
func (p *SampleDependentPlugin) RegisterRoutes(router *gin.Engine) {
	// 为空实现，使用新版接口
}

// Execute 执行插件功能
func (p *SampleDependentPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	return gin.H{"result": "success", "plugin": p.Name()}, nil
}

// GetDefaultMiddlewares 获取默认中间件
func (p *SampleDependentPlugin) GetDefaultMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

// handleIndex 处理首页请求
func (p *SampleDependentPlugin) handleIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"plugin":      p.Name(),
		"version":     p.Version(),
		"description": p.Description(),
		"dependencies": p.GetDependencies(),
		"available_endpoints": []string{
			"GET /plugins/sample_dependent/ - 获取插件信息",
			"GET /plugins/sample_dependent/use-dependency - 使用依赖的插件功能",
			"GET /plugins/sample_dependent/dependencies - 获取插件依赖信息",
		},
	})
}

// handleUseDependency 处理使用依赖插件的请求
func (p *SampleDependentPlugin) handleUseDependency(c *gin.Context) {
	// 使用依赖的插件执行某些功能
	action := c.DefaultQuery("action", "greet")
	params := map[string]interface{}{"action": action}
	
	result, err := p.pluginManager.ExecutePlugin("sample_optimized", params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"plugin":               p.Name(),
		"used_dependency":      "sample_optimized",
		"action":               action,
		"result_from_dependency": result,
	})
}

// handleGetDependencies 处理获取依赖信息的请求
func (p *SampleDependentPlugin) handleGetDependencies(c *gin.Context) {
	// 获取所有已注册插件的依赖关系
	depGraph := p.pluginManager.GetDependencyGraph()
	
	// 获取当前插件的依赖状态
	var dependenciesStatus []map[string]interface{}
	for _, depName := range p.GetDependencies() {
		if depPlugin, exists := p.pluginManager.GetPlugin(depName); exists {
			dependenciesStatus = append(dependenciesStatus, map[string]interface{}{
				"name":      depName,
				"version":   depPlugin.Version(),
				"description": depPlugin.Description(),
				"status":    "available",
			})
		} else {
			dependenciesStatus = append(dependenciesStatus, map[string]interface{}{
				"name":   depName,
				"status": "missing",
			})
		}
	}

	c.JSON(200, gin.H{
		"plugin":       p.Name(),
		"dependencies": dependenciesStatus,
		"dependency_graph": depGraph,
	})
}