package plugins

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// HelloPlugin 示例插件
type HelloPlugin struct{}

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

// RegisterRoutes 注册插件路由
func (p *HelloPlugin) RegisterRoutes(router *gin.Engine) {
	// 注册插件相关路由
	pluginGroup := router.Group(fmt.Sprintf("/plugins/%s", p.Name()))
	{
		pluginGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"plugin":      p.Name(),
				"description": p.Description(),
				"version":     p.Version(),
			})
		})

		pluginGroup.GET("/say", func(c *gin.Context) {
			name := c.DefaultQuery("name", "World")
			result, _ := p.Execute(map[string]interface{}{
				"name": name,
			})
			c.JSON(200, result)
		})
	}
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
