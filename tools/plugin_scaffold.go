package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 插件信息结构体

type PluginInfo struct {
	Name        string // 插件名称，如 "TodoPlugin"
	Identifier  string // 插件标识符，如 "todo"
	Description string // 插件描述
	Version     string // 插件版本
}

// 验证插件名称是否有效

func validatePluginName(name string) bool {
	// 插件名称应该是驼峰式命名，只包含字母
	matched, _ := regexp.MatchString("^[A-Z][a-zA-Z0-9]*Plugin$", name)
	return matched
}

// 验证插件标识符是否有效

func validatePluginIdentifier(identifier string) bool {
	// 插件标识符应该只包含小写字母和下划线
	matched, _ := regexp.MatchString("^[a-z][a-z0-9_]*$", identifier)
	return matched
}

// 从插件名称生成标识符

func generateIdentifier(name string) string {
	// 移除Plugin后缀
	nameWithoutPlugin := strings.TrimSuffix(name, "Plugin")
	if nameWithoutPlugin == "" {
		return ""
	}

	// 将驼峰式转换为下划线分隔的小写形式
	var result strings.Builder
	for i, char := range nameWithoutPlugin {
		if i > 0 && 'A' <= char && char <= 'Z' {
			result.WriteString("_")
		}
		result.WriteRune(char)
	}

	return strings.ToLower(result.String())
}

// 生成插件代码

func generatePluginCode(info PluginInfo) string {
	var code strings.Builder

	// 写入包声明和导入
	code.WriteString(`package plugins

import (
	"log"

	"github.com/gin-gonic/gin"
)

// `)
	code.WriteString(info.Name)
	code.WriteString(` 插件结构体
// 可以在这里添加插件所需的字段和依赖

type `)
	code.WriteString(info.Name)
	code.WriteString(` struct {
	// 添加插件所需的字段，例如数据库连接、缓存等
	// db *gorm.DB
	// cache *redis.Client
}

// Name 返回插件名称
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) Name() string {
	// 返回插件的唯一标识符，将用于路由和插件管理
	return "`)
	code.WriteString(info.Identifier)
	code.WriteString(`"
}

// Description 返回插件描述
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) Description() string {
	// 返回插件的详细描述，说明插件的功能和用途
	return "`)
	code.WriteString(info.Description)
	code.WriteString(`"
}

// Version 返回插件版本
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) Version() string {
	// 返回插件的版本号，遵循语义化版本规范（SemVer）
	return "`)
	code.WriteString(info.Version)
	code.WriteString(`"
}

// Init 初始化插件
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) Init() error {
	// 插件初始化逻辑
	// 可以在这里进行数据库连接、配置加载等操作
	log.Printf("%%s: 插件已初始化", p.Name())

	// 如果初始化失败，返回错误信息
	// return errors.New("初始化失败的原因")
	return nil
}

// Shutdown 关闭插件
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) Shutdown() error {
	// 插件关闭逻辑
	// 可以在这里释放资源、关闭连接等操作
	log.Printf("%%s: 插件已关闭", p.Name())

	// 如果关闭失败，返回错误信息
	// return errors.New("关闭失败的原因")
	return nil
}

// GetRoutes 返回插件的路由定义
// 推荐使用此方法代替 RegisterRoutes
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) GetRoutes() []Route {
	return []Route{
		{
			Path:         "/",
			Method:       "GET",
			Handler: func(c *gin.Context) {
				c.JSON(200, gin.H{
					"plugin":      p.Name(),
					"description": p.Description(),
					"version":     p.Version(),
					// 可以在这里添加插件的其他元信息
				})
			},
			Description:  "获取插件信息",
			AuthRequired: false,
			Tags:         []string{"info", "metadata"},
		},
		{
			Path:         "/resources",
			Method:       "GET",
			Handler: func(c *gin.Context) {
				// 这里是处理逻辑
				// 例如获取查询参数、调用服务层、返回结果等
				result, err := p.GetResources()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, result)
			},
			Description:  "获取资源列表",
			AuthRequired: false,
			Tags:         []string{"resources", "list"},
			Params: map[string]string{
				// 可以在这里定义查询参数
			},
		},
		{
			Path:         "/resources",
			Method:       "POST",
			Handler: func (c *gin.Context) {
		// 这里是处理逻辑
		// 例如绑定请求体、验证数据、调用服务层、返回结果等
		var request struct {
			Name  string ` + "`json:\"name\" binding:\"required\"`" + `
			Value string ` + "`json:\"value\"`" + `
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
	},
		Description:  "创建资源",
		AuthRequired: false,
		Tags:         []string{"resources", "create"},
		},
		// 可以根据需要添加更多路由
	}
}

// RegisterRoutes 保留旧的方法以确保兼容性
// 在使用新的GetRoutes方法后，这个方法实际上不会被调用
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) RegisterRoutes(router *gin.Engine) {
	// 这个方法在使用新的GetRoutes时不会被调用
	// 保留只是为了兼容性
	log.Printf("%%s: 注意：使用了旧的RegisterRoutes方法，建议使用新的GetRoutes方法", p.Name())
}

// GetDefaultMiddlewares 返回插件的默认中间件
// 这些中间件会应用到插件的所有路由上
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) GetDefaultMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 可以在这里添加认证中间件、日志中间件等
		// auth.Middleware(),
		// log.Middleware(),
	}
}

// Execute 执行插件功能
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) Execute(params map[string]interface{}) (interface{}, error) {
	// 插件的核心执行逻辑
	// params 包含调用插件时传递的参数
	log.Printf("%%s: 执行插件功能，参数: %%v", p.Name(), params)

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
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) GetResources() (interface{}, error) {
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
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) CreateResource(name, value string) (interface{}, error) {
	// 实际的业务逻辑
	return map[string]interface{}{
			"id":    "3", // 这里应该返回实际生成的ID
			"name":  name,
			"value": value,
		},
		nil
}

// DoSomething 执行特定操作的示例方法
func (p *`)
	code.WriteString(info.Name)
	code.WriteString(`) DoSomething(params map[string]interface{}) (interface{}, error) {
	// 实际的业务逻辑
	return map[string]interface{}{
			"result": "操作执行成功",
			"params": params,
		},
		nil
}

// 插件注册函数
var `)
	code.WriteString(strings.ToLower(info.Name[:1]) + info.Name[1:])
	code.WriteString(` = &`)
	code.WriteString(info.Name)
	code.WriteString(`{}

// 导出符号，便于外部使用
func init() {
	// 这个初始化函数会在包被导入时执行
	// 实际的插件注册应该在main.go中完成
}
`)

	return code.String()
}

// 保存插件代码到文件

func savePluginCode(code string, outputDir string, pluginName string) error {
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 生成文件名
	fileName := strings.ToLower(strings.TrimSuffix(pluginName, "Plugin")) + "_plugin.go"
	filePath := filepath.Join(outputDir, fileName)

	// 检查文件是否已存在
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("文件 '%s' 已存在，请选择其他名称或删除现有文件", filePath)
	}

	// 写入文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(code); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新文件缓冲区失败: %w", err)
	}

	return nil
}

// 从用户输入获取插件信息

func getPluginInfoFromInput() (PluginInfo, error) {
	reader := bufio.NewReader(os.Stdin)
	var info PluginInfo

	// 获取插件名称
	for {
		fmt.Print("请输入插件名称（驼峰式，以Plugin结尾，如TodoPlugin）: ")
		name, err := reader.ReadString('\n')
		if err != nil {
			return info, fmt.Errorf("读取输入失败: %w", err)
		}
		info.Name = strings.TrimSpace(name)

		if validatePluginName(info.Name) {
			break
		} else {
			fmt.Println("插件名称无效！请使用驼峰式命名并以Plugin结尾，只包含字母。")
		}
	}

	// 生成默认标识符
	defaultIdentifier := generateIdentifier(info.Name)

	// 获取插件标识符
	for {
		fmt.Printf("请输入插件标识符（小写字母和下划线，默认: %s）: ", defaultIdentifier)
		identifier, err := reader.ReadString('\n')
		if err != nil {
			return info, fmt.Errorf("读取输入失败: %w", err)
		}
		info.Identifier = strings.TrimSpace(identifier)

		// 如果用户没有输入，使用默认值
		if info.Identifier == "" {
			info.Identifier = defaultIdentifier
		}

		if validatePluginIdentifier(info.Identifier) {
			break
		} else {
			fmt.Println("插件标识符无效！请使用小写字母和下划线。")
		}
	}

	// 获取插件描述
	fmt.Print("请输入插件描述: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return info, fmt.Errorf("读取输入失败: %w", err)
	}
	info.Description = strings.TrimSpace(description)

	// 获取插件版本
	fmt.Print("请输入插件版本（默认: 1.0.0）: ")
	version, err := reader.ReadString('\n')
	if err != nil {
		return info, fmt.Errorf("读取输入失败: %w", err)
	}
	info.Version = strings.TrimSpace(version)

	// 如果用户没有输入版本，使用默认值
	if info.Version == "" {
		info.Version = "1.0.0"
	}

	return info, nil
}

func main() {
	fmt.Println("=== ToolCat 插件脚手架生成器 ===")
	fmt.Println("这个工具将帮助你快速生成符合规范的插件代码文件。")

	// 获取插件信息
	info, err := getPluginInfoFromInput()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n插件信息摘要:")
	fmt.Printf("名称: %s\n", info.Name)
	fmt.Printf("标识符: %s\n", info.Identifier)
	fmt.Printf("描述: %s\n", info.Description)
	fmt.Printf("版本: %s\n", info.Version)

	// 生成插件代码
	code := generatePluginCode(info)

	// 确定输出目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	// 尝试找到plugins目录
	projectRoot := currentDir
	for i := 0; i < 3; i++ {
		pluginsDir := filepath.Join(projectRoot, "plugins")
		if stat, err := os.Stat(pluginsDir); err == nil && stat.IsDir() {
			// 生成README文件
			readmeContent := generatePluginReadme(info)
			readmePath := filepath.Join(pluginsDir, strings.ToLower(strings.TrimSuffix(info.Name, "Plugin"))+"_plugin.md")
			saveReadme(readmeContent, readmePath)
			
			// 保存插件代码
			if err := savePluginCode(code, pluginsDir, info.Name); err != nil {
				fmt.Printf("保存插件代码失败: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("\n✅ 插件代码已成功生成到: %s\n", filepath.Join(pluginsDir, strings.ToLower(strings.TrimSuffix(info.Name, "Plugin"))+"_plugin.go"))
			fmt.Printf("✅ 插件文档已成功生成到: %s\n", readmePath)
			fmt.Println("\n下一步操作:")
			fmt.Println("1. 编辑生成的插件文件，实现具体功能")
		fmt.Println("2. 在main.go中注册你的插件")
		fmt.Println("3. 测试插件功能")
		fmt.Println("4. 编写更详细的文档")
		return
		}

		// 向上一级目录查找
		projectRoot = filepath.Dir(projectRoot)
	}

	// 如果找不到plugins目录，使用当前目录
	fmt.Println("警告: 未找到plugins目录，将使用当前目录作为输出目录。")
	if err := savePluginCode(code, currentDir, info.Name); err != nil {
		fmt.Printf("保存插件代码失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ 插件代码已成功生成到当前目录: %s\n", strings.ToLower(strings.TrimSuffix(info.Name, "Plugin"))+"_plugin.go")
	fmt.Println("\n注意: 请手动将生成的插件文件移动到项目的plugins目录中，并在main.go中注册。")
}

// 生成插件README文档
func generatePluginReadme(info PluginInfo) string {
	var readme strings.Builder

	// 写入标题和插件信息
	readme.WriteString("# ")
	readme.WriteString(info.Name)
	readme.WriteString(" 插件文档\n\n## 插件信息\n\n- **名称**: ")
	readme.WriteString(info.Name)
	readme.WriteString("\n- **标识符**: ")
	readme.WriteString(info.Identifier)
	readme.WriteString("\n- **描述**: ")
	readme.WriteString(info.Description)
	readme.WriteString("\n- **版本**: ")
	readme.WriteString(info.Version)

	// 写入功能说明
	readme.WriteString("\n\n## 功能说明\n\n这里是插件的详细功能说明，包括插件提供的主要功能和用途。")

	// 写入API接口
	readme.WriteString("\n\n## API 接口\n\n### GET /plugin/")
	readme.WriteString(info.Identifier)
	readme.WriteString("/\n\n**描述**: 获取插件信息\n**认证**: 不需要\n**参数**: 无\n**返回**: \n```json\n{\n  \"plugin\": \"")
	readme.WriteString(info.Identifier)
	readme.WriteString(",\n  \"description\": \"")
	readme.WriteString(info.Description)
	readme.WriteString(",\n  \"version\": \"")
	readme.WriteString(info.Version)
	readme.WriteString("\"")
	readme.WriteString("\n}\n```")

	// 写入获取资源列表接口
	readme.WriteString("\n\n### GET /plugin/")
	readme.WriteString(info.Identifier)
	readme.WriteString("/resources\n\n**描述**: 获取资源列表\n**认证**: 不需要\n**参数**: 无\n**返回**: 资源列表")

	// 写入创建资源接口
	readme.WriteString("\n\n### POST /plugin/")
	readme.WriteString(info.Identifier)
	readme.WriteString("/resources\n\n**描述**: 创建资源\n**认证**: 不需要\n**参数**: \n```json\n{\n  \"name\": \"资源名称\",\n  \"value\": \"资源值\"\n}\n```\n**返回**: 创建的资源信息")

	// 写入配置说明和依赖说明
	readme.WriteString("\n\n## 配置说明\n\n如果插件需要配置，请在此处说明配置项和配置方法。\n\n## 依赖说明\n\n如果插件依赖其他库或服务，请在此处说明。")

	// 写入使用示例
	readme.WriteString("\n\n## 使用示例\n\n```go\n// 在main.go中注册插件\nimport (\n\t\"toolcat/plugins\"\n)\n\nfunc registerPlugins() {\n\t// 注册插件\n\tplugins.PluginManager.Register(plugins.")
	readme.WriteString(strings.ToLower(info.Name[:1]) + info.Name[1:])
	readme.WriteString(")\n}\n```")

	// 写入开发说明
	readme.WriteString("\n\n## 开发说明\n\n- 编辑 ")
	readme.WriteString(strings.ToLower(strings.TrimSuffix(info.Name, "Plugin")))
	readme.WriteString("_plugin.go 文件实现具体功能\n- 可以添加自定义字段到")
	readme.WriteString(info.Name)
	readme.WriteString("结构体中\n- 根据需要实现更多的API接口\n- 实现GetDefaultMiddlewares方法添加插件特定的中间件")

	return readme.String()
}

// 保存README文档
func saveReadme(content string, filePath string) error {
	// 检查文件是否已存在
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("警告: README文件 '%s' 已存在，跳过创建。\n", filePath)
		return nil
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建README文件失败: %w", err)
	}
	defer file.Close()

	// 写入内容
	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(content); err != nil {
		return fmt.Errorf("写入README文件失败: %w", err)
	}

	// 刷新缓冲区
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新README文件缓冲区失败: %w", err)
	}

	return nil
}