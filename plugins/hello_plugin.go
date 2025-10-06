package plugins

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// HelloPlugin 示例插件
type HelloPlugin struct{
	pluginManager *pluginManager
}

// NewHelloPlugin 创建新的HelloPlugin实例
func NewHelloPlugin() *HelloPlugin {
	return &HelloPlugin{}
}

// Name 返回插件名称
func (p *HelloPlugin) Name() string {
	return "hello"
}

// Description 返回插件描述
func (p *HelloPlugin) Description() string {
	return "一个简单的示例插件"
}

// Version 返回插件版本
func (p *HelloPlugin) Version() string {
	return "1.0.0"
}

// GetDependencies 返回依赖的插件
func (p *HelloPlugin) GetDependencies() []string {
	return []string{} // 不依赖其他插件
}

// GetConflicts 返回冲突的插件
func (p *HelloPlugin) GetConflicts() []string {
	return []string{} // 与其他插件无冲突
}

// SetPluginManager 设置插件管理器
func (p *HelloPlugin) SetPluginManager(manager *pluginManager) {
	p.pluginManager = manager
}

// Init 初始化插件
func (p *HelloPlugin) Init() error {
	// 插件初始化逻辑
	fmt.Println("HelloPlugin: 插件已初始化")
	return nil
}

// Shutdown 关闭插件
func (p *HelloPlugin) Shutdown() error {
	// 插件关闭逻辑
	fmt.Println("HelloPlugin: 插件已关闭")
	return nil
}

// RegisterRoutes 保留旧的方法以确保兼容性
// 在使用新的GetRoutes方法后，这个方法实际上不会被调用
func (p *HelloPlugin) RegisterRoutes(router *gin.Engine) {
	// 这个方法在使用新的GetRoutes时不会被调用
	// 保留只是为了兼容性
	fmt.Printf("%s: 注意：使用了旧的RegisterRoutes方法，建议使用新的GetRoutes方法\n", p.Name())
}

// GetRoutes 返回插件的路由定义
func (p *HelloPlugin) GetRoutes() []Route {
	return []Route{
		{
			Path:         "/",
			Method:       "GET",
			Handler: func(c *gin.Context) {
				c.JSON(200, gin.H{
					"plugin":      p.Name(),
					"description": p.Description(),
					"version":     p.Version(),
				})
			},
			Description:  "获取插件信息",
			AuthRequired: false,
		},
		{
			Path:         "/say",
			Method:       "GET",
			Handler: func(c *gin.Context) {
				name := c.DefaultQuery("name", "World")
				result, _ := p.Execute(map[string]interface{}{
					"name": name,
				})
				c.JSON(200, result)
			},
			Description:  "Say Hello API",
			AuthRequired: false,
			Params: map[string]string{
				"name": "可选，问候的对象名称",
			},
		},
	}
}

// GetDefaultMiddlewares 返回插件的默认中间件
func (p *HelloPlugin) GetDefaultMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

// Execute 执行插件功能
func (p *HelloPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	name, ok := params["name"].(string)
	if !ok || name == "" {
		name = "World"
	}

	message := fmt.Sprintf("Hello, %s!", name)
	return map[string]interface{}{
			"message": message,
			"params":  params,
		},
		nil
}
