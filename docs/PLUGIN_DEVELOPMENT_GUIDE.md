# ToolCat 插件开发指南

## 1. 插件系统概述

ToolCat 插件系统允许开发者扩展应用功能，而无需修改核心代码。本指南将详细介绍如何开发、测试和部署 ToolCat 插件，特别关注优化后的路由注册机制。

## 2. 插件接口定义

插件必须实现 `Plugin` 接口。ToolCat 提供了两种路由注册方式：传统的 `RegisterRoutes` 方法和优化后的 `GetRoutes` 方法。推荐使用新的 `GetRoutes` 方法进行路由注册。

```go
// Plugin 接口定义了所有插件必须实现的方法
// 注意：新的 GetRoutes 方法已加入，推荐使用此方法代替 RegisterRoutes

type Plugin interface {
    // 基础信息
    Name() string
    Description() string
    Version() string

    // 生命周期管理
    Init() error
    Shutdown() error

    // 路由管理（新方式）- 推荐使用
    GetRoutes() []Route
    GetDefaultMiddlewares() []gin.HandlerFunc

    // 路由管理（旧方式）- 为兼容性保留
    RegisterRoutes(router *gin.Engine)

    // 功能执行
    Execute(params map[string]interface{}) (interface{}, error)
}

// Route 结构体定义了路由的元数据和处理函数
// 这是新的路由定义方式核心

type Route struct {
    Path         string                 // 路由路径
    Method       string                 // HTTP 方法（GET, POST, PUT, DELETE 等）
    Handler      gin.HandlerFunc        // 请求处理函数
    Middlewares  []gin.HandlerFunc      // 路由特定的中间件
    Description  string                 // 路由描述
    AuthRequired bool                   // 是否需要认证
    Tags         []string               // 路由标签，用于文档生成
    Params       map[string]string      // 参数说明，用于文档生成
    Metadata     map[string]interface{} // 自定义元数据
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

### 4.8 保留兼容性（可选）

为了确保与旧版系统的兼容性，可以保留 `RegisterRoutes` 方法的空实现：

```go
func (p *MyPlugin) RegisterRoutes(router *gin.Engine) {
    // 注意：使用GetRoutes方法后，这个方法不会被调用
    // 保留只是为了兼容性
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

1. **使用语义化路径**：设计清晰、有意义的API路径
2. **统一的命名规范**：使用小写字母和连字符(-)或下划线(_)分隔单词
3. **版本化API**：考虑在路由中包含版本信息，如 `/v1/resource`
4. **使用中间件进行横切关注点**：将认证、日志、错误处理等逻辑抽象为中间件
5. **完整的路由元数据**：提供详细的路由描述、参数说明等元数据，便于文档生成
6. **RESTful API 设计**：遵循RESTful原则，使用正确的HTTP方法表示操作类型

## 8. 示例插件分析

`sample_optimized_plugin.go` 提供了一个完整的示例，展示如何使用新的路由注册机制开发插件。该示例包含：

- 完整的 `Plugin` 接口实现
- 使用 `GetRoutes` 方法定义多个路由
- 实现路由级别的中间件和全局默认中间件
- 包含路由元数据、参数说明等信息
- 保留了兼容性的 `RegisterRoutes` 方法实现

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

## 11. 结语

通过本指南，您应该能够理解 ToolCat 的插件系统，特别是优化后的路由注册机制。使用新的 `GetRoutes` 方法可以使您的插件开发更加规范、高效和可维护。