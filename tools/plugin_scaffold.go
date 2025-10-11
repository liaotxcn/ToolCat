package main

import (
	"bufio"
	"flag"
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
	PluginType  string // 插件类型：api, task, event等
}

// 验证插件名称是否有效
func validatePluginName(name string) bool {
	// 插件名称应该是驼峰式命名，以Plugin结尾，只包含字母和数字
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

// generatePluginCode 生成插件代码
func generatePluginCode(info PluginInfo) string {
	var code strings.Builder

	// 包声明
	code.WriteString("package plugins\n\n")

	// 导入部分
	code.WriteString("import (\n")
	code.WriteString("	\"fmt\"\n")
	code.WriteString("	\"strings\"\n")
	code.WriteString("	\"time\"\n")
	code.WriteString("	\"sync\"\n")
	code.WriteString("	\"github.com/gin-gonic/gin\"\n")
	code.WriteString("	\"errors\"\n")
	code.WriteString("	\"toolcat/pkg\"\n")
	code.WriteString("	\"toolcat/plugins/core\"\n")
	code.WriteString("\n")
	code.WriteString(")\n\n")

	// 类型验证（确保实现了core.Plugin接口）
	code.WriteString("// 类型验证：确保")
	code.WriteString(info.Name)
	code.WriteString("实现了core.Plugin接口\n")
	code.WriteString("var _ core.Plugin = &")
	code.WriteString(info.Name)
	code.WriteString("{}\n\n")

	// 插件结构体定义
	code.WriteString("// ")
	code.WriteString(info.Name)
	code.WriteString(" 是一个")
	code.WriteString(info.Description)
	code.WriteString("\n")
	code.WriteString("type ")
	code.WriteString(info.Name)
	code.WriteString(" struct {\n")
	code.WriteString("	mu        sync.RWMutex\n")
	code.WriteString("	config    map[string]interface{}\n")
	code.WriteString("	resources map[string]interface{}\n")
	code.WriteString("	version   string\n")
	code.WriteString("}\n\n")

	// 配置结构体定义
	code.WriteString("// ")
	code.WriteString(info.Name)
	code.WriteString("Config 定义插件的配置结构\n")
	code.WriteString("type ")
	code.WriteString(info.Name)
	code.WriteString("Config struct {\n")
	code.WriteString("	// 可以在这里定义配置项\n")
	code.WriteString("\tDebug bool `json:\"debug\"`\n")
	code.WriteString("\tTimeout int `json:\"timeout\"`\n")
	code.WriteString("}\n\n")

	// loadConfig方法
	code.WriteString("// loadConfig 加载并解析配置\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") loadConfig() (*")
	code.WriteString(info.Name)
	code.WriteString("Config, error) {\n")
	code.WriteString("	var config ")
	code.WriteString(info.Name)
	code.WriteString("Config\n")
	code.WriteString("\t// 这里应该从配置中心或配置文件加载配置\n")
	code.WriteString("\t// 目前使用默认值\n")
	code.WriteString("\tconfig.Debug = false\n")
	code.WriteString("\tconfig.Timeout = 30\n")
	code.WriteString("\n")
	code.WriteString("\treturn &config, nil\n")
	code.WriteString("}\n\n")

	// GetResources方法
	code.WriteString("// GetResources 获取资源列表的示例方法\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") GetResources() (interface{}, error) {\n")
	code.WriteString("	p.mu.RLock()\n")
	code.WriteString("	defer p.mu.RUnlock()\n")
	code.WriteString("\n")
	code.WriteString("	// 这里是实际的资源获取逻辑\n")
	code.WriteString("	return []map[string]interface{}{{\"id\": \"1\", \"name\": \"示例资源1\"}, {\"id\": \"2\", \"name\": \"示例资源2\"}},",)
	code.WriteString("\t\tnil\n")
	code.WriteString("}\n\n")

	// CreateResource方法
	code.WriteString("// CreateResource 创建资源的示例方法\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") CreateResource(name, value string) (interface{}, error) {\n")
	code.WriteString("	p.mu.Lock()\n")
	code.WriteString("	defer p.mu.Unlock()\n")
	code.WriteString("\n")
	code.WriteString("	// 这里是实际的资源创建逻辑\n")
	code.WriteString("	return map[string]interface{}{\"id\": \"3\", \"name\": name, \"value\": value},",)
	code.WriteString("\t\tnil\n")
	code.WriteString("}\n\n")

	// Name方法
	code.WriteString("// Name 返回插件名称\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") Name() string {\n")
	code.WriteString("	return \"")
	code.WriteString(info.Name)
	code.WriteString("\"\n")
	code.WriteString("}\n\n")

	// Description方法
	code.WriteString("// Description 返回插件描述\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") Description() string {\n")
	code.WriteString("	return \"")
	code.WriteString(info.Description)
	code.WriteString("\"\n")
	code.WriteString("}\n\n")

	// Version方法
	code.WriteString("// Version 返回插件版本\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") Version() string {\n")
	code.WriteString("	p.mu.RLock()\n")
	code.WriteString("	defer p.mu.RUnlock()\n")
	code.WriteString("	return p.version\n")
	code.WriteString("}\n\n")

	// SetVersion方法
	code.WriteString("// SetVersion 设置插件版本\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") SetVersion(version string) {\n")
	code.WriteString("	p.mu.Lock()\n")
	code.WriteString("	p.version = version\n")
	code.WriteString("	p.mu.Unlock()\n")
	code.WriteString("}\n\n")

	// Init方法
	code.WriteString("// Init 初始化插件\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") Init(config map[string]interface{}) error {\n")
	code.WriteString("	p.mu.Lock()\n")
	code.WriteString("	p.config = config\n")
	code.WriteString("	p.resources = make(map[string]interface{})\n")
	code.WriteString("	p.version = \"1.0.0\"\n")
	code.WriteString("	p.mu.Unlock()\n")
	code.WriteString("\n")
	code.WriteString("	// 加载详细配置\n")
	code.WriteString("	pluginConfig, err := p.loadConfig()\n")
	code.WriteString("	if err != nil {\n")
	code.WriteString("		pkg.Error(fmt.Sprintf(\"%s: 加载配置失败: %%v\", p.Name(), err))\n")
	code.WriteString("		return err\n")
	code.WriteString("	}\n")
	code.WriteString("\n")
	code.WriteString("	// 根据配置初始化插件\n")
	code.WriteString("	if pluginConfig.Debug {\n")
	code.WriteString("		pkg.Info(fmt.Sprintf(\"%s: 调试模式已启用\", p.Name()))\n")
	code.WriteString("	}\n")
	code.WriteString("\n")
	code.WriteString("	pkg.Info(fmt.Sprintf(\"%s: 初始化完成\", p.Name()))\n")
	code.WriteString("	return nil\n")
	code.WriteString("}\n\n")

	// Shutdown方法
	code.WriteString("// Shutdown 关闭插件\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") Shutdown() error {\n")
	code.WriteString("	// 执行清理操作\n")
	code.WriteString("	p.mu.Lock()\n")
	code.WriteString("	defer p.mu.Unlock()\n")
	code.WriteString("\n")
	code.WriteString("	// 清理资源\n")
	code.WriteString("	p.resources = nil\n")
	code.WriteString("	p.config = nil\n")
	code.WriteString("\n")
	code.WriteString("	pkg.Info(fmt.Sprintf(\"%s: 已关闭\", p.Name()))\n")
	code.WriteString("	return nil\n")
	code.WriteString("}\n\n")

	// GetRoutes方法
	code.WriteString("// GetRoutes 返回插件的路由定义\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") GetRoutes() []core.Route {\n")
	code.WriteString("	return []core.Route{\n")
	code.WriteString("	\t{\n")
	code.WriteString("	\t\tPath:         \"/\",\n")
	code.WriteString("	\t\tMethod:       \"GET\",\n")
	code.WriteString("	\t\tHandler: func(c *gin.Context) {\n")
	code.WriteString("	\t\t\tc.JSON(200, gin.H{\n")
	code.WriteString("	\t\t\t\t\"plugin\":      p.Name(),\n")
	code.WriteString("	\t\t\t\t\"description\": p.Description(),\n")
	code.WriteString("	\t\t\t\t\"version\":     p.Version(),\n")
	code.WriteString("	\t\t\t\t// 可以在这里添加插件的其他元信息\n")
	code.WriteString("	\t\t\t})\n")
	code.WriteString("	\t\t},\n")
	code.WriteString("	\t\tDescription:  \"获取插件信息\",\n")
	code.WriteString("	\t\tAuthRequired: false,\n")
	code.WriteString("	\t\tTags:         []string{\"info\", \"metadata\"},\n")
	code.WriteString("	\t},\n")
	code.WriteString("	\t{\n")
	code.WriteString("	\t\tPath:         \"/resources\",\n")
	code.WriteString("	\t\tMethod:       \"GET\",\n")
	code.WriteString("	\t\tHandler: func(c *gin.Context) {\n")
	code.WriteString("	\t\t\t// 这里是处理逻辑\n")
	code.WriteString("	\t\t\tresult, err := p.GetResources()\n")
	code.WriteString("	\t\t\tif err != nil {\n")
	code.WriteString("	\t\t\t\tc.JSON(500, gin.H{\"error\": err.Error()})\n")
	code.WriteString("	\t\t\t\treturn\n")
	code.WriteString("	\t\t\t}\n")
	code.WriteString("	\t\t\tc.JSON(200, result)\n")
	code.WriteString("	\t\t},\n")
	code.WriteString("	\t\tDescription:  \"获取资源列表\",\n")
	code.WriteString("	\t\tAuthRequired: false,\n")
	code.WriteString("	\t\tTags:         []string{\"resources\", \"list\"},\n")
	code.WriteString("	\t\tParams: map[string]string{\n")
	code.WriteString("	\t\t\t// 可以在这里定义查询参数\n")
	code.WriteString("	\t\t},\n")
	code.WriteString("	\t},\n")
	code.WriteString("	\t{\n")
	code.WriteString("	\t\tPath:         \"/resources\",\n")
	code.WriteString("	\t\tMethod:       \"POST\",\n")
	code.WriteString("	\t\tHandler: func(c *gin.Context) {\n")
	code.WriteString("	\t\t\t// 这里是处理逻辑\n")
	code.WriteString("	\t\t\tvar request struct {\n")
	code.WriteString("	\t\t\t\tName  string `json:\"name\" binding:\"required\"`\n")
	code.WriteString("	\t\t\t\tValue string `json:\"value\"`\n")
	code.WriteString("	\t\t\t}\n\n")
	code.WriteString("	\t\t\tif err := c.ShouldBindJSON(&request); err != nil {\n")
	code.WriteString("	\t\t\t\tc.JSON(400, gin.H{\"error\": err.Error()})\n")
	code.WriteString("	\t\t\t\treturn\n")
	code.WriteString("	\t\t\t}\n\n")
	code.WriteString("	\t\t\tresult, err := p.CreateResource(request.Name, request.Value)\n")
	code.WriteString("	\t\t\tif err != nil {\n")
	code.WriteString("	\t\t\t\tc.JSON(500, gin.H{\"error\": err.Error()})\n")
	code.WriteString("	\t\t\t\treturn\n")
	code.WriteString("	\t\t\t}\n")
	code.WriteString("	\t\t\tc.JSON(201, result)\n")
	code.WriteString("	\t\t},\n")
	code.WriteString("	\t\tDescription:  \"创建资源\",\n")
	code.WriteString("	\t\tAuthRequired: false,\n")
	code.WriteString("	\t\tTags:         []string{\"resources\", \"create\"},\n")
	code.WriteString("	\t},\n")
	code.WriteString("	\t// 可以根据需要添加更多路由\n")
	code.WriteString("	}\n")
	code.WriteString("}\n\n")

	// RegisterRoutes方法
	code.WriteString("// RegisterRoutes 保留旧的方法以确保兼容性\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") RegisterRoutes(router *gin.Engine) {\n")
	code.WriteString("	// 这个方法在使用新的GetRoutes时不会被调用\n")
	code.WriteString("	// 保留只是为了兼容性\n")
	code.WriteString("	pkg.Info(fmt.Sprintf(\"%s: 注意：使用了旧的RegisterRoutes方法，建议使用新的GetRoutes方法\", p.Name()))\n")
	code.WriteString("}\n\n")

	// GetDefaultMiddlewares方法
	code.WriteString("// GetDefaultMiddlewares 返回插件的默认中间件\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") GetDefaultMiddlewares() []gin.HandlerFunc {\n")
	code.WriteString("	return []gin.HandlerFunc{}\n")
	code.WriteString("}\n\n")

	// Execute方法
	code.WriteString("// Execute 执行插件功能\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") Execute(params map[string]interface{}) (interface{}, error) {\n")
	code.WriteString("	// 插件的核心执行逻辑\n")
	code.WriteString("	pkg.Info(fmt.Sprintf(\"%s: 执行插件功能，参数: %%v\", p.Name(), params))\n\n")
	code.WriteString("	// 示例：根据参数执行不同的操作\n")
	code.WriteString("	action, ok := params[\"action\"].(string)\n")
	code.WriteString("	if !ok {\n")
	code.WriteString("		action = \"default\"\n")
	code.WriteString("	}\n\n")
	code.WriteString("	switch action {\n")
	code.WriteString("	case \"get_info\":\n")
	code.WriteString("		// 返回插件信息\n")
	code.WriteString("		return map[string]interface{}{\n")
	code.WriteString("		\t\"name\":        p.Name(),\n")
	code.WriteString("		\t\"description\": p.Description(),\n")
	code.WriteString("		\t\"version\":     p.Version(),\n")
	code.WriteString("		},",)
	code.WriteString("\t\tnil\n")
	code.WriteString("	case \"do_something\":\n")
	code.WriteString("		// 执行特定操作\n")
	code.WriteString("		return p.DoSomething(params)\n")
	code.WriteString("	default:\n")
	code.WriteString("		return map[string]interface{}{\"message\": \"未知的操作\"},",)
	code.WriteString("\t\tnil\n")
	code.WriteString("	}\n")
	code.WriteString("}\n\n")

	// DoSomething方法
	code.WriteString("// DoSomething 执行特定操作的示例方法\n")
	code.WriteString("func (p *")
	code.WriteString(info.Name)
	code.WriteString(") DoSomething(params map[string]interface{}) (interface{}, error) {\n")
	code.WriteString("	return map[string]interface{}{\"result\": \"操作执行成功\"},",)
	code.WriteString("\t\tnil\n")
	code.WriteString("}\n\n")

	// 插件注册变量
	code.WriteString("// 插件注册变量\n")
	code.WriteString("var ")
	code.WriteString(strings.ToLower(info.Name[:1]) + info.Name[1:])
	code.WriteString(" = &")
	code.WriteString(info.Name)
	code.WriteString("{}\n\n")

	// init函数
	code.WriteString("// 初始化函数\n")
	code.WriteString("func init() {\n")
	code.WriteString("	// 根据插件类型执行特定的初始化逻辑\n")
	code.WriteString("	switch \"")
	code.WriteString(info.PluginType)
	code.WriteString("\" {\n")
	code.WriteString("	case \"task\":\n")
	code.WriteString("	\tpkg.Info(fmt.Sprintf(\"%s: 任务类型插件初始化完成\", ",)
	code.WriteString(strings.ToLower(info.Name[:1]) + info.Name[1:])
	code.WriteString(".Name()))\n")
	code.WriteString("	case \"event\":\n")
	code.WriteString("	\t// 注册事件监听器的示例\n")
	code.WriteString("	default:\n")
	code.WriteString("	\t// API类型插件无需特殊初始化\n")
	code.WriteString("	}\n")
	code.WriteString("}")

	return code.String()
}

// 保存插件代码到文件
func savePluginCode(code string, outputDir string, pluginName string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	fileName := strings.ToLower(strings.TrimSuffix(pluginName, "Plugin")) + "_plugin.go"
	filePath := filepath.Join(outputDir, fileName)

	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("文件 '%s' 已存在，请选择其他名称或删除现有文件", filePath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(code); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return writer.Flush()
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
			fmt.Println("插件名称无效！请使用驼峰式命名并以Plugin结尾，只包含字母和数字。")
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

	// 获取插件类型
	fmt.Print("请输入插件类型（api, task, event，默认: api）: ")
	pluginType, err := reader.ReadString('\n')
	if err != nil {
		return info, fmt.Errorf("读取输入失败: %w", err)
	}
	info.PluginType = strings.TrimSpace(pluginType)

	// 如果用户没有输入插件类型，使用默认值
	if info.PluginType == "" {
		info.PluginType = "api"
	}

	// 验证插件类型
	validTypes := map[string]bool{"api": true, "task": true, "event": true}
	if !validTypes[info.PluginType] {
		fmt.Println("警告: 无效的插件类型，将使用默认类型 'api'")
		info.PluginType = "api"
	}

	return info, nil
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
	readme.WriteString("\"\n}\n```")

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
	readme.WriteString("\n\n## 使用示例\n\n```go\n// 在main.go中注册插件\nimport (\n\t\"toolcat/plugins\"\n\t\"toolcat/plugins/core\"\n)\n\nfunc registerPlugins() {\n\t// 注册插件\n\tcore.GlobalPluginManager.Register(plugins." + strings.ToLower(info.Name[:1]) + info.Name[1:] + ")\n}\n```")

	// 写入开发说明
	readme.WriteString("\n\n## 开发说明\n\n- 编辑 " + strings.ToLower(strings.TrimSuffix(info.Name, "Plugin")) + "_plugin.go 文件实现具体功能\n- 可以添加自定义字段到" + info.Name + "结构体中\n- 根据需要实现更多的API接口\n- 实现GetDefaultMiddlewares方法添加插件特定的中间件")

	return readme.String()
}

// 保存README文档
func saveReadme(content string, filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("警告: README文件 '%s' 已存在，跳过创建。\n", filePath)
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建README文件失败: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := writer.WriteString(content); err != nil {
		return fmt.Errorf("写入README文件失败: %w", err)
	}

	return writer.Flush()
}

func main() {
	fmt.Println("=== ToolCat 插件脚手架生成器 ===")
	fmt.Println("这个工具将帮助你快速生成符合规范的插件代码文件。")

	// 解析命令行参数
	namePtr := flag.String("name", "", "插件名称（驼峰式，以Plugin结尾，如TodoPlugin）")
	descPtr := flag.String("desc", "", "插件描述")
	versionPtr := flag.String("version", "1.0.0", "插件版本")
	typePtr := flag.String("type", "api", "插件类型（api, task, event）")
	dirPtr := flag.String("dir", "", "输出目录")
	nonInteractivePtr := flag.Bool("non-interactive", false, "非交互式模式")
	flag.Parse()

	var info PluginInfo
	var err error

	// 如果是非交互式模式且提供了必要的参数
	if *nonInteractivePtr && *namePtr != "" {
		info.Name = *namePtr
		if !validatePluginName(info.Name) {
			fmt.Println("错误: 插件名称无效！请使用驼峰式命名并以Plugin结尾，只包含字母和数字。")
			os.Exit(1)
		}

		info.Identifier = generateIdentifier(info.Name)
		info.Description = *descPtr
		info.Version = *versionPtr
		info.PluginType = *typePtr

		// 验证插件类型
		validTypes := map[string]bool{"api": true, "task": true, "event": true}
		if !validTypes[info.PluginType] {
			fmt.Println("警告: 无效的插件类型，将使用默认类型 'api'")
			info.PluginType = "api"
		}
	} else {
		// 获取插件信息（交互式）
		info, err = getPluginInfoFromInput()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			os.Exit(1)
		}
	}

	// 插件信息摘要:
	fmt.Println("\n插件信息摘要:")
	fmt.Printf("名称: %s\n", info.Name)
	fmt.Printf("标识符: %s\n", info.Identifier)
	fmt.Printf("描述: %s\n", info.Description)
	fmt.Printf("版本: %s\n", info.Version)
	fmt.Printf("类型: %s\n", info.PluginType)

	// 生成插件代码
	code := generatePluginCode(info)

	// 确定输出目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		os.Exit(1)
	}

	// 确定最终输出目录
	var outputDir string
	if *dirPtr != "" {
		// 如果用户通过命令行参数指定了输出目录
		outputDir = *dirPtr
	} else {
		// 尝试找到plugins目录
		projectRoot := currentDir
		pluginsFound := false
		for i := 0; i < 3; i++ {
			pluginsDir := filepath.Join(projectRoot, "plugins")
			if stat, err := os.Stat(pluginsDir); err == nil && stat.IsDir() {
				outputDir = pluginsDir
				pluginsFound = true
				break
			}

			// 向上一级目录查找
			projectRoot = filepath.Dir(projectRoot)
		}

		// 如果找不到plugins目录，使用当前目录
		if !pluginsFound {
			outputDir = currentDir
			fmt.Println("警告: 未找到plugins目录，将使用当前目录作为输出目录。")
		}
	}

	// 生成并保存README文件（仅当输出到plugins目录时）
	readmeGenerated := false
	if strings.Contains(outputDir, "plugins") {
		readmeContent := generatePluginReadme(info)
		readmePath := filepath.Join(outputDir, strings.ToLower(strings.TrimSuffix(info.Name, "Plugin"))+"_plugin.md")
		saveReadme(readmeContent, readmePath)
		fmt.Printf("✅ 插件文档已成功生成到: %s\n", readmePath)
		readmeGenerated = true
	}

	// 保存插件代码
	if err := savePluginCode(code, outputDir, info.Name); err != nil {
		fmt.Printf("保存插件代码失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ 插件代码已成功生成到: %s\n", filepath.Join(outputDir, strings.ToLower(strings.TrimSuffix(info.Name, "Plugin"))+"_plugin.go"))

	// 根据插件类型和输出目录显示不同的后续步骤
	if strings.Contains(outputDir, "plugins") {
		fmt.Println("\n下一步操作:")
		fmt.Println("1. 编辑生成的插件文件，实现具体功能")
		fmt.Println("2. 在main.go中注册你的插件")
		fmt.Println("3. 测试插件功能")
		if !readmeGenerated {
			fmt.Println("4. 编写插件文档")
		}
	} else {
		fmt.Println("\n注意: 请手动将生成的插件文件移动到项目的plugins目录中，并在main.go中注册。")
	}

	// 根据插件类型显示特定提示
	switch info.PluginType {
	case "task":
		fmt.Println("\n任务类型插件提示:")
		fmt.Println("- 实现定时任务逻辑，可使用time包或其他定时任务库")
		fmt.Println("- 在Init方法中启动任务调度器")
	case "event":
		fmt.Println("\n事件类型插件提示:")
		fmt.Println("- 实现事件监听器，订阅系统事件")
		fmt.Println("- 可以通过events.Publish发布自定义事件")
	default:
		// API类型插件无需特殊提示
	}
}
