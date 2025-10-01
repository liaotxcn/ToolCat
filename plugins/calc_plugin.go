package plugins

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CalcPlugin 计算器插件
type CalcPlugin struct{}

// Name 返回插件名称
func (p *CalcPlugin) Name() string {
	return "calculator"
}

// Description 返回插件描述
func (p *CalcPlugin) Description() string {
	return "一个简单的计算器插件，可以进行基本的算术运算"
}

// Version 返回插件版本
func (p *CalcPlugin) Version() string {
	return "1.0.0"
}

// Init 初始化插件
func (p *CalcPlugin) Init() error {
	// 插件初始化逻辑
	fmt.Println("CalcPlugin: 计算器插件已初始化")
	return nil
}

// Shutdown 关闭插件
func (p *CalcPlugin) Shutdown() error {
	// 插件关闭逻辑
	fmt.Println("CalcPlugin: 计算器插件已关闭")
	return nil
}

// RegisterRoutes 注册插件路由
func (p *CalcPlugin) RegisterRoutes(router *gin.Engine) {
	// 注册插件相关路由
	pluginGroup := router.Group(fmt.Sprintf("/plugins/%s", p.Name()))
	{
		pluginGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"plugin":               p.Name(),
				"description":          p.Description(),
				"version":              p.Version(),
				"supported_operations": []string{"add", "subtract", "multiply", "divide"},
			})
		})

		pluginGroup.GET("/calculate", func(c *gin.Context) {
			aStr := c.DefaultQuery("a", "0")
			bStr := c.DefaultQuery("b", "0")
			operation := c.DefaultQuery("operation", "add")

			result, err := p.Execute(map[string]interface{}{
				"a":         aStr,
				"b":         bStr,
				"operation": operation,
			})

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, result)
		})
	}
}

// Execute 执行插件功能
func (p *CalcPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	// 获取参数
	aStr, ok := params["a"].(string)
	if !ok {
		return nil, fmt.Errorf("参数 a 无效")
	}

	bStr, ok := params["b"].(string)
	if !ok {
		return nil, fmt.Errorf("参数 b 无效")
	}

	operation, ok := params["operation"].(string)
	if !ok {
		operation = "add"
	}

	// 转换数字
	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return nil, fmt.Errorf("参数 a 不是有效数字: %v", err)
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return nil, fmt.Errorf("参数 b 不是有效数字: %v", err)
	}

	// 执行计算
	var result float64
	var operationName string

	switch operation {
	case "add":
		result = a + b
		operationName = "加法"
	case "subtract":
		result = a - b
		operationName = "减法"
	case "multiply":
		result = a * b
		operationName = "乘法"
	case "divide":
		if b == 0 {
			return nil, fmt.Errorf("除数不能为零")
		}
		result = a / b
		operationName = "除法"
	default:
		return nil, fmt.Errorf("不支持的操作: %s", operation)
	}

	return map[string]interface{}{
			"a":              a,
			"b":              b,
			"operation":      operation,
			"operation_name": operationName,
			"result":         result,
		},
		nil
}
