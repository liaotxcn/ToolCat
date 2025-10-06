# ToolCat 插件开发指南

## 1. 插件系统概述

ToolCat 插件系统允许开发者扩展应用功能，而无需修改核心代码。本指南将详细介绍如何开发、测试和部署 ToolCat 插件，推荐使用优化后的路由注册机制、插件依赖管理和热重载功能。

## 2. 插件接口定义

插件必须实现 `Plugin` 接口。ToolCat 提供了两种路由注册方式：传统的 `RegisterRoutes` 方法和优化后的 `GetRoutes` 方法。推荐使用新的 `GetRoutes` 方法进行路由注册。

```go
// Plugin 插件接口定义
type Plugin interface {
    // 基础信息接口
    Name() string        // 返回插件名称
    Description() string // 返回插件描述
    Version() string     // 返回插件版本

    // 生命周期接口
    Init() error     // 初始化插件
    Shutdown() error // 关闭插件

    // 路由注册接口
    // 新版接口：提供路由定义，由PluginManager统一注册
    GetRoutes() []Route // 获取插件路由定义
    // 旧版接口：为了兼容现有插件保留
    RegisterRoutes(router *gin.Engine) // 注册插件路由

    // 执行功能接口
    Execute(params map[string]interface{}) (interface{}, error) // 执行插件功能

    // 插件配置接口（可选）
    GetDefaultMiddlewares() []gin.HandlerFunc // 获取插件默认中间件
}

// Route 结构体定义了路由的元数据和处理函数
// 这是新的路由定义方式核心

type Route struct {
    Path         string            // 路由路径（不包含插件前缀）
    Method       string            // HTTP方法
    Handler      gin.HandlerFunc   // 处理函数
    Middlewares  []gin.HandlerFunc // 路由特定中间件
    Description  string            // 路由描述
    AuthRequired bool              // 是否需要认证
    Tags         []string          // 路由标签，用于文档生成
    Params       map[string]string // 参数说明，用于文档生成
}
```

## 3. 路由注册机制详解

### 3.1 优化后的路由注册机制

优化后的路由注册机制具有以下优势：

- **集中式路由管理**：PluginManager 负责统一注册和管理所有插件路由
- **路由元数据**：支持定义路由描述、参数说明、认证要求等元数据
- **统一的中间件机制**：支持全局和路由级别的中间件
- **自动路由组创建**：自动为插件创建路由组，格式为 `/plugins/{plugin_name}/`
- **类型安全**：通过结构体定义确保路由信息的完整性

### 3.2 两种路由注册方式的对比

| 特性 | GetRoutes 方法（推荐） | RegisterRoutes 方法（兼容性保留） |
|------|-----------------------|-----------------------------------|
| 路由定义 | 使用 Route 结构体数组 | 直接操作 gin.Engine 对象 |
| 元数据支持 | ✅ 完整支持 | ❌ 不支持 |
| 自动路由组 | ✅ 自动创建 | ❌ 需要手动创建 |
| 中间件管理 | ✅ 支持全局和路由级别 | ❌ 需要手动添加 |
| 文档生成 | ✅ 支持自动生成 API 文档 | ❌ 不支持 |

## 4. 开发插件的步骤

### 4.1 创建插件结构体

```go
package plugins

type MyPlugin struct{}
```

### 4.2 实现基础信息方法

```go
func (p *MyPlugin) Name() string {
    return "my_plugin"
}

func (p *MyPlugin) Description() string {
    return "我的自定义插件"
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}
```

### 4.3 实现生命周期方法

```go
func (p *MyPlugin) Init() error {
    // 初始化插件，加载配置、连接数据库等
    return nil
}

func (p *MyPlugin) Shutdown() error {
    // 关闭插件，清理资源
    return nil
}
```

### 4.4 实现路由注册（新方式）

推荐使用 `GetRoutes` 方法定义路由：

```go
func (p *MyPlugin) GetRoutes() []Route {
    return []Route{
        {
            Path:        "/",
            Method:      "GET",
            Handler:     p.handleIndex,
            Description: "插件主页",
            AuthRequired: false,
            Tags:        []string{"home"},
        },
        {
            Path:        "/api/data",
            Method:      "GET",
            Handler:     p.handleGetData,
            Middlewares: []gin.HandlerFunc{p.authMiddleware},
            Description: "获取数据API",
            AuthRequired: true,
            Tags:        []string{"data", "api"},
            Params: map[string]string{
                "id": "数据ID",
            },
        },
    }
}

// 定义插件的默认中间件
func (p *MyPlugin) GetDefaultMiddlewares() []gin.HandlerFunc {
    return []gin.HandlerFunc{
        p.logMiddleware,
        p.corsMiddleware,
    }
}
```

### 4.5 实现路由处理函数

```go
func (p *MyPlugin) handleIndex(c *gin.Context) {
    c.JSON(200, gin.H{
        "plugin": p.Name(),
        "version": p.Version(),
    })
}

func (p *MyPlugin) handleGetData(c *gin.Context) {
    id := c.Query("id")
    // 处理业务逻辑
    c.JSON(200, gin.H{
        "id": id,
        "data": "示例数据",
    })
}
```

### 4.6 实现中间件

```go
func (p *MyPlugin) logMiddleware(c *gin.Context) {
    // 记录请求日志
    c.Next()
}

func (p *MyPlugin) authMiddleware(c *gin.Context) {
    // 验证用户认证
    c.Next()
}

func (p *MyPlugin) corsMiddleware(c *gin.Context) {
    // 设置CORS头
    c.Next()
}
```

### 4.7 实现 Execute 方法

```go
func (p *MyPlugin) Execute(params map[string]interface{}) (interface{}, error) {
    // 实现插件的核心功能
    // params 包含调用插件时传递的参数
    return map[string]interface{}{"result": "success"}, nil
}
```

### 4.8 保留兼容性（建议）

为了确保与旧版系统的兼容性，建议保留 `RegisterRoutes` 方法的空实现或添加兼容性提示：

```go
func (p *MyPlugin) RegisterRoutes(router *gin.Engine) {
    // 注意：使用GetRoutes方法后，这个方法不会被调用
    // 保留只是为了兼容性
    log.Printf("%s: 注意：使用了旧的RegisterRoutes方法，建议使用新的GetRoutes方法", p.Name())
}
```

## 5. 在 main.go 中注册插件

```go
// 注册插件
func registerPlugins(router *gin.Engine) {
    // 设置路由引擎到PluginManager
    plugins.PluginManager.SetRouter(router)

    // 注册自定义插件
    myPlugin := &plugins.MyPlugin{}
    if err := plugins.PluginManager.Register(myPlugin); err != nil {
        log.Printf("Failed to register plugin %s: %v", myPlugin.Name(), err)
    } else {
        log.Printf("Successfully registered plugin: %s", myPlugin.Name())
    }

    // 注意：使用新的路由注册机制后，不需要手动调用插件的RegisterRoutes方法
    // PluginManager会自动处理路由注册
}
```

## 6. 访问插件API

使用新的路由注册机制后，插件的API路径将自动格式化为：`/plugins/{plugin_name}/{route_path}`

例如，对于上面的示例插件，API路径将是：
- `GET /plugins/my_plugin/` - 插件主页
- `GET /plugins/my_plugin/api/data?id=123` - 获取数据API

## 7. 路由设计最佳实践

1. **使用语义化路径**：设计清晰、有意义的API路径，例如 `/users` 而非 `/getUsers`
2. **统一的命名规范**：使用小写字母和连字符(-)或下划线(_)分隔单词，保持一致性
3. **版本化API**：考虑在路由中包含版本信息，如 `/v1/resource`
4. **使用中间件进行横切关注点**：将认证、日志、错误处理等逻辑抽象为中间件
5. **完整的路由元数据**：提供详细的路由描述、参数说明等元数据，便于文档生成
6. **RESTful API 设计**：遵循RESTful原则，使用正确的HTTP方法表示操作类型
7. **路径参数与查询参数分离**：使用路径参数表示资源标识，查询参数表示过滤、排序等操作
8. **请求验证中间件**：为需要验证输入的路由创建专门的验证中间件
9. **合理的错误处理**：为不同类型的错误提供明确的HTTP状态码和错误信息
10. **示例代码注释**：在路由处理函数中添加简短注释说明其功能和参数

## 8. 示例插件分析

`sample_optimized_plugin.go` 提供了一个完整的示例，展示如何使用新的路由注册机制开发插件。该示例包含：

- 完整的 `Plugin` 接口实现
- 使用 `GetRoutes` 方法定义多个路由，包含请求示例
- 实现路由级别的中间件和全局默认中间件
- 包含路由元数据、参数说明等信息
- 保留了兼容性的 `RegisterRoutes` 方法实现并添加了兼容性提示
- 提供了多种中间件示例：日志中间件、CORS中间件和请求验证中间件
- 展示了如何在验证中间件中将处理后的数据存储在上下文中供后续处理函数使用

## 9. 调试插件

开发插件时，可以使用以下方法进行调试：

1. **日志记录**：使用 `log` 包记录关键信息和错误
2. **断点调试**：使用GoLand或VSCode等IDE进行断点调试
3. **API测试**：使用Postman或curl测试插件API
4. **检查注册状态**：通过 `PluginManager.ListPlugins()` 检查插件是否正确注册

## 10. 常见问题解答

### 10.1 插件注册失败怎么办？

检查以下几点：
- 确保插件实现了所有必需的接口方法
- 检查 `Init()` 方法是否返回错误
- 确认插件名称不与已注册的插件冲突
- 查看日志输出获取详细错误信息

### 10.2 如何处理插件的依赖关系？

- 在 `Init()` 方法中初始化所有依赖
- 考虑使用Go模块管理外部依赖
- 对于插件间的依赖，确保依赖的插件先注册

### 10.3 如何保证插件的安全性？

- 实现适当的认证和授权机制
- 验证所有用户输入，防止注入攻击
- 限制插件的权限范围
- 记录敏感操作日志

### 10.4 新的路由注册机制与旧机制有何不同？

新机制使用 `GetRoutes()` 方法返回路由定义，PluginManager 负责统一注册路由；而旧机制需要插件自己通过 `RegisterRoutes()` 方法注册路由。新机制提供了更好的集中管理、元数据支持和中间件机制。

## 11. 插件依赖管理机制

ToolCat 插件系统支持插件间的依赖管理，允许一个插件声明其依赖于其他插件，并在运行时安全地访问这些依赖的插件。

### 11.1 依赖管理接口

Plugin 接口新增了以下依赖管理相关的方法：

```go
// Plugin 插件接口定义（包含依赖管理部分）
type Plugin interface {
    // 基础信息接口（现有）
    Name() string        // 返回插件名称
    Description() string // 返回插件描述
    Version() string     // 返回插件版本
    
    // 生命周期接口（现有）
    Init() error     // 初始化插件
    Shutdown() error // 关闭插件
    
    // 路由注册接口（现有）
    GetRoutes() []Route // 获取插件路由定义
    RegisterRoutes(router *gin.Engine) // 注册插件路由
    
    // 执行功能接口（现有）
    Execute(params map[string]interface{}) (interface{}, error) // 执行插件功能
    
    // 插件配置接口（现有，可选）
    GetDefaultMiddlewares() []gin.HandlerFunc // 获取插件默认中间件
    
    // 依赖管理接口（新增）
    GetDependencies() []string // 获取插件依赖的其他插件名称列表
    GetConflicts() []string    // 获取与当前插件冲突的插件名称列表
    SetPluginManager(manager *PluginManager) // 设置插件管理器引用
}
```

### 11.2 PluginInfo 结构体的扩展

```go
// PluginInfo 存储插件的元数据信息（包含依赖管理部分）
type PluginInfo struct {
    Name        string                 // 插件名称
    Description string                 // 插件描述
    Version     string                 // 插件版本
    Routes      []Route                // 插件的路由定义
    Middlewares []gin.HandlerFunc      // 插件的默认中间件
    Dependencies []string              // 插件依赖的其他插件名称列表
    Conflicts   []string               // 与当前插件冲突的插件名称列表
    Status      string                 // 插件状态（"enabled", "disabled", "error"）
    Error       error                  // 如果插件状态为error，存储错误信息
    Plugin      Plugin                 // 插件实例引用
}
```

### 11.3 声明插件依赖

在插件中实现 GetDependencies 和 GetConflicts 方法以声明依赖关系：

```go
// 声明插件依赖
func (p *MyPlugin) GetDependencies() []string {
    // 返回当前插件依赖的其他插件名称
    return []string{
        "sample_optimized",  // 依赖名为"sample_optimized"的插件
        "hello_plugin"       // 依赖名为"hello_plugin"的插件
    }
}

// 声明插件冲突
func (p *MyPlugin) GetConflicts() []string {
    // 返回与当前插件冲突的其他插件名称
    return []string{
        "legacy_plugin"      // 与名为"legacy_plugin"的插件冲突
    }
}

// 设置插件管理器引用
func (p *MyPlugin) SetPluginManager(manager *plugins.PluginManager) {
    p.pluginManager = manager
}
```

### 11.4 访问依赖的插件

在插件中，可以通过插件管理器访问依赖的插件：

```go
// 访问依赖的插件示例
func (p *MyPlugin) SomeFunction() {
    // 检查依赖的插件是否已注册并启用
    if p.pluginManager != nil {
        // 获取依赖的插件实例
        helloPlugin, err := p.pluginManager.GetPlugin("hello_plugin")
        if err == nil {
            // 调用依赖插件的方法
            result, err := helloPlugin.Execute(map[string]interface{}{
                "action": "greet",
                "name":   "ToolCat"
            })
            if err == nil {
                // 处理结果
                fmt.Printf("Hello plugin result: %v\n", result)
            }
        }
    }
}
```

### 11.5 依赖管理的工作原理

PluginManager 会在注册插件时进行以下检查：

1. **冲突检查**：验证当前注册的插件是否与已注册的插件冲突
2. **依赖检查**：验证当前注册的插件所需的所有依赖是否已注册
3. **循环依赖检测**：确保插件之间不会形成循环依赖

当使用批量注册方法 `RegisterPlugins()` 时，系统会自动根据依赖关系进行拓扑排序，确保依赖的插件先于依赖它的插件注册。

### 11.6 依赖管理最佳实践

1. **最小化依赖**：仅声明必要的依赖关系，减少插件间的耦合
2. **版本兼容性**：在插件文档中明确说明兼容的依赖插件版本
3. **处理依赖缺失**：在 Init() 方法中检查依赖插件是否存在并正确处理缺失情况
4. **避免循环依赖**：设计插件架构时避免形成循环依赖
5. **优雅降级**：在依赖插件不可用时，提供功能降级方案
6. **明确的错误信息**：当依赖问题导致插件初始化失败时，提供清晰的错误信息

### 11.7 批量注册插件

对于具有复杂依赖关系的插件集合，建议使用批量注册方法：

```go
// 批量注册插件，自动处理依赖关系
func registerMultiplePlugins(router *gin.Engine) {
    // 设置路由引擎到PluginManager
    plugins.PluginManager.SetRouter(router)
    
    // 创建插件实例
    helloPlugin := &plugins.HelloPlugin{}
    samplePlugin := &plugins.SampleOptimizedPlugin{}
    dependentPlugin := &plugins.SampleDependentPlugin{}
    
    // 批量注册，系统会自动处理依赖顺序
    pluginList := []plugins.Plugin{
        helloPlugin,
        samplePlugin,
        dependentPlugin,
    }
    
    if err := plugins.PluginManager.RegisterPlugins(pluginList); err != nil {
        log.Printf("Plugin registration failed: %v", err)
    }
}
```

## 12. 插件热重载支持

ToolCat 插件系统支持插件热重载功能，允许在不重启服务器的情况下启用、禁用和重载插件。这大大提高了插件开发和运维的效率。

### 12.1 热重载相关接口

Plugin 接口新增了以下热重载相关的生命周期方法：

```go
// Plugin 插件接口定义（包含热重载部分）
type Plugin interface {
    // 生命周期接口（新增热重载相关方法）
    Init() error     // 初始化插件
    OnEnable() error // 插件启用时调用（热重载相关）
    OnDisable() error // 插件禁用时调用（热重载相关）
    Shutdown() error // 关闭插件
    
    // 其他现有接口方法...
}
```

PluginInfo 结构体扩展了一个字段用于跟踪插件的启用状态：

```go
// PluginInfo 存储插件的元数据信息（包含热重载部分）
type PluginInfo struct {
    // 现有字段...
    IsEnabled bool    // 插件是否启用
}
```

### 12.2 实现热重载方法

插件开发者需要实现 `OnEnable()` 和 `OnDisable()` 方法以支持热重载功能：

```go
// OnEnable 插件启用时调用
func (p *MyPlugin) OnEnable() error {
    // 插件启用时执行的逻辑
    fmt.Printf("%s: 插件已启用\n", p.Name())
    
    // 可以在这里重新初始化资源、恢复服务等
    return nil
}

// OnDisable 插件禁用时调用
func (p *MyPlugin) OnDisable() error {
    // 插件禁用时执行的逻辑
    fmt.Printf("%s: 插件已禁用\n", p.Name())
    
    // 可以在这里释放临时资源、停止服务等
    return nil
}
```

### 12.3 热重载 API

ToolCat 提供了以下 API 用于管理插件的热重载状态：

```go
// 插件控制器中的热重载相关方法
// 获取所有插件信息
plugins.GET("/", pluginCtrl.GetAllPlugins)
// 获取插件状态
plugins.GET("/:name/status", pluginCtrl.GetPluginStatus)
// 启用插件
plugins.POST("/:name/enable", pluginCtrl.EnablePlugin)
// 禁用插件
plugins.POST("/:name/disable", pluginCtrl.DisablePlugin)
// 重载插件
plugins.POST("/:name/reload", pluginCtrl.ReloadPlugin)
// 获取插件依赖图
plugins.GET("/dependency-graph", pluginCtrl.GetDependencyGraph)
```

### 12.4 热重载工作原理

当使用热重载 API 时，PluginManager 会执行以下操作：

1. **启用插件**：设置插件状态为启用，调用 `OnEnable()` 方法，重新注册路由
2. **禁用插件**：设置插件状态为禁用，调用 `OnDisable()` 方法，移除路由
3. **重载插件**：先禁用插件，然后启用插件，相当于执行一个完整的循环

系统会保证在操作插件时的线程安全，并且在执行热重载操作时会考虑插件之间的依赖关系。

### 12.5 热重载最佳实践

1. **资源管理**：在 `OnEnable()` 中获取资源，在 `OnDisable()` 中释放资源
2. **状态保存**：重要状态应保存在持久化存储中，而不是内存中
3. **错误处理**：在热重载方法中提供适当的错误处理和日志记录
4. **依赖处理**：在启用插件时检查依赖插件的状态
5. **数据一致性**：确保热重载不会破坏系统数据的一致性

## 13. CSRF保护处理

ToolCat系统启用了CSRF（跨站请求伪造）保护机制，插件开发者需要了解如何在开发和测试插件时处理CSRF验证。

### 13.1 CSRF保护对插件的影响

对于插件中定义的非GET/HEAD/OPTIONS/TRACE请求，系统会自动应用CSRF保护机制，需要客户端同时在请求头和Cookie中提供有效的CSRF令牌。

### 13.2 插件开发中的CSRF处理

1. **前端插件集成**
   - 确保前端代码正确获取和使用CSRF令牌
   - 从Cookie(`XSRF-TOKEN`)或响应头(`X-CSRF-Token`)中获取令牌
   - 对所有非安全请求在请求头中添加`X-CSRF-Token`字段

2. **API测试**
   - 使用Apifox/Postman等工具测试API时，需要处理CSRF令牌
   - 先发送一个GET请求获取CSRF令牌，然后在后续请求中使用，或设定前置脚本操作进行测试

### 13.3 CSRF保护代码示例

以下是在插件中处理CSRF保护的示例代码：

```go
// 在插件的中间件中添加CSRF令牌到响应中
func (p *MyPlugin) addCSRFTokenMiddleware(c *gin.Context) {
    // 检查是否已有CSRF令牌
    token, err := c.Cookie(config.Config.CSRF.CookieName)
    if err != nil || token == "" {
        // 如果没有，生成一个新的令牌
        token := generateCSRFToken(config.Config.CSRF.TokenLength)
        // 设置Cookie
        c.SetCookie(
            config.Config.CSRF.CookieName,
            token,
            config.Config.CSRF.CookieMaxAge,
            config.Config.CSRF.CookiePath,
            config.Config.CSRF.CookieDomain,
            config.Config.CSRF.CookieSecure,
            config.Config.CSRF.CookieHttpOnly,
        )
        // 添加到响应头
        c.Header(config.Config.CSRF.HeaderName, token)
    }
    c.Next()
}

// 插件路由中使用
func (p *MyPlugin) GetRoutes() []Route {
    return []Route{
        {
            Path:        "/api/safe-operation",
            Method:      "POST",
            Handler:     p.handleSafeOperation,
            Middlewares: []gin.HandlerFunc{p.addCSRFTokenMiddleware},
            Description: "需要CSRF保护的操作",
            AuthRequired: true,
        },
    }
}
```

## 14. 结语

通过本指南，您应该能够理解 ToolCat 的插件系统，包括优化后的路由注册机制、插件依赖管理功能和热重载支持。使用这些功能可以使您的插件开发更加规范、高效和可维护，同时为构建复杂的插件生态系统提供坚实基础。

随着插件系统的不断完善，ToolCat 能够支持更加灵活和强大的插件扩展，为用户提供更加丰富的功能和更好的体验。