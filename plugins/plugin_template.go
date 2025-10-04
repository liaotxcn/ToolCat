package plugins

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// 插件开发模板
//
// 使用说明：
// 1. 复制此文件并重命名为 your_plugin_name.go
// 2. 将所有 "YourPlugin" 替换为你的插件名称（如 "TodoPlugin"）
// 3. 将所有 "your_plugin" 替换为你的插件标识符（如 "todo"）
// 4. 实现各个方法的具体逻辑
// 5. 在main.go中注册你的插件

// YourPlugin 你的插件结构体
// 可以在这里添加插件所需的字段和依赖

// 这里是示例字段，根据实际需求修改

type YourPlugin struct {
	// 添加插件所需的字段，例如数据库连接、缓存等
	// db *gorm.DB
	// cache *redis.Client
}

// Name 返回插件名称
func (p *YourPlugin) Name() string {
	// 返回插件的唯一标识符，将用于路由和插件管理
	// 建议使用小写字母和下划线
	return "your_plugin"
}

// Description 返回插件描述
func (p *YourPlugin) Description() string {
	// 返回插件的详细描述，说明插件的功能和用途
	return "这是一个插件模板，用于快速开发新插件"
}

// Version 返回插件版本
func (p *YourPlugin) Version() string {
	// 返回插件的版本号，遵循语义化版本规范（SemVer）
	return "1.0.0"
}

// Init 初始化插件
func (p *YourPlugin) Init() error {
	// 插件初始化逻辑
	// 可以在这里进行数据库连接、配置加载等操作
	log.Printf("%s: 插件已初始化", p.Name())

	// 如果初始化失败，返回错误信息
	// return errors.New("初始化失败的原因")
	return nil
}

// Shutdown 关闭插件
func (p *YourPlugin) Shutdown() error {
	// 插件关闭逻辑
	// 可以在这里释放资源、关闭连接等操作
	log.Printf("%s: 插件已关闭", p.Name())

	// 如果关闭失败，返回错误信息
	// return errors.New("关闭失败的原因")
	return nil
}

// RegisterRoutes 注册插件路由
func (p *YourPlugin) RegisterRoutes(router *gin.Engine) {
	// 注册插件相关路由
	// 路由前缀为 /plugins/插件名称
	pluginGroup := router.Group(fmt.Sprintf("/plugins/%s", p.Name()))
	{
		// 获取插件信息接口
		pluginGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"plugin":      p.Name(),
				"description": p.Description(),
				"version":     p.Version(),
				// 可以在这里添加插件的其他元信息
			})
		})

		// 示例API接口：获取资源列表
		pluginGroup.GET("/resources", func(c *gin.Context) {
			// 这里是处理逻辑
			// 例如获取查询参数、调用服务层、返回结果等
			result, err := p.GetResources()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, result)
		})

		// 示例API接口：创建资源
		pluginGroup.POST("/resources", func(c *gin.Context) {
			// 这里是处理逻辑
			// 例如绑定请求体、验证数据、调用服务层、返回结果等
			var request struct {
				Name  string `json:"name" binding:"required"`
				Value string `json:"value"`
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			result, err := p.CreateResource(request.Name, request.Value)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(201, result)
		})

		// 可以根据需要添加更多路由
	}
}

// Execute 执行插件功能
func (p *YourPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	// 插件的核心执行逻辑
	// params 包含调用插件时传递的参数
	log.Printf("%s: 执行插件功能，参数: %v", p.Name(), params)

	// 示例：根据参数执行不同的操作
	action, ok := params["action"].(string)
	if !ok {
		action = "default"
	}

	switch action {
	case "get_info":
		// 返回插件信息
		return map[string]interface{}{
				"name":        p.Name(),
				"description": p.Description(),
				"version":     p.Version(),
			},
			nil
	case "do_something":
		// 执行特定操作
		return p.DoSomething(params)
	default:
		return map[string]interface{}{
				"message": "未知的操作",
				"params":  params,
			},
			nil
	}
}

// 以下是插件的辅助方法，根据实际需求添加

// GetResources 获取资源列表的示例方法
func (p *YourPlugin) GetResources() (interface{}, error) {
	// 实际的业务逻辑
	return []map[string]interface{}{
			{
				"id":   "1",
				"name": "示例资源1",
			},
			{
				"id":   "2",
				"name": "示例资源2",
			},
		},
		nil
}

// CreateResource 创建资源的示例方法
func (p *YourPlugin) CreateResource(name, value string) (interface{}, error) {
	// 实际的业务逻辑
	return map[string]interface{}{
			"id":     "new-id",
			"name":   name,
			"value":  value,
			"status": "created",
		},
		nil
}

// DoSomething 执行特定操作的示例方法
func (p *YourPlugin) DoSomething(params map[string]interface{}) (interface{}, error) {
	// 实际的业务逻辑
	return map[string]interface{}{
			"result": "操作已完成",
			"params": params,
		},
		nil
}

// 插件使用示例
// 在main.go中注册插件：
// plugins.PluginManager.Register(&plugins.YourPlugin{})

// 前端可以通过以下方式调用插件API：
// GET /plugins/your_plugin/
// GET /plugins/your_plugin/resources
// POST /plugins/your_plugin/resources

// 系统可以通过PluginManager执行插件功能：
// result, err := plugins.PluginManager.ExecutePlugin("your_plugin", map[string]interface{}{
//	"action": "do_something",
//	"param1": "value1",
// })
