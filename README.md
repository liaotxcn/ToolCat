# ToolCat - A high-performance, highly scalable, and easily extendable plugin-based tool integration service platform developed in Golang

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
</div>

## 📋 项目简介

ToolCat 基于 Golang 开发的高性能、高效率、插件化易扩展的工具服务平台，旨在提供灵活的工具集成和管理解决方案。

## 🌟 项目特点

### 🚀 高性能/效率
- 基于 Gin 框架构建，处理请求速度快，并发能力强
- 数据库连接池优化，支持高并发访问
- 模块化架构设计，代码结构清晰，易于维护和扩展
- 配置管理支持环境变量覆盖，便于不同环境部署
- 优雅关闭机制，确保服务安全退出

### 🔌 插件化易扩展
- 统一的插件接口设计，支持热插拔
- 插件管理器统一注册、管理和执行插件
- 示例插件（Hello、Note）展示了完整的插件开发流程
- 插件可独立注册路由，拥有独立命名空间

## 📂 项目架构

```
├── config/         # 配置文件管理
├── controllers/    # 控制器层
├── internal/       # 内部包
├── main.go         # 程序入口
├── middleware/     # 中间件
├── models/         # 数据模型
├── pkg/            # 公共包
├── plugins/        # 插件系统
├── routers/        # 路由管理
├── utils/          # 工具函数
└── web/            # 前端代码
```

## 🛠️ 核心组件

### 插件系统
ToolCat 的核心特色是其灵活高效的插件系统，允许开发者轻松扩展平台功能。

```go
// 插件接口定义
type Plugin interface {
    Name() string              // 插件名称
    Description() string       // 插件描述
    Version() string           // 插件版本
    Init() error               // 初始化插件
    Shutdown() error           // 关闭插件
    
    // 路由管理（新方式）- 推荐使用
    GetRoutes() []Route
    GetDefaultMiddlewares() []gin.HandlerFunc
    
    // 路由管理（旧方式）- 为兼容性保留
    RegisterRoutes(*gin.Engine) // 注册路由
    
    Execute(map[string]interface{}) (interface{}, error) // 执行功能
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

插件管理器负责插件的整个生命周期管理，包括注册、注销、查询和执行插件功能。

## 快速开始

1. 克隆代码库
```bash
git clone https://github.com/liaotxcn/toolcat.git
cd toolcat
```

2. 初始化数据库
创建数据库，并确保配置文件相关配置信息正确

3. 编译并运行
```bash
go mod tidy
go build -o toolcat
go run main.go
```

服务将在 http://localhost:8081 启动。

## 插件开发指南

### 创建新插件
1. 实现 `plugins.Plugin` 接口
2. 在 `main.go` 的 `registerPlugins` 函数中注册插件

### 插件示例（使用推荐的 GetRoutes 方法）
```go
// 示例插件结构
type MyPlugin struct{}

// 实现 Plugin 接口的方法
func (p *MyPlugin) Name() string { return "myplugin" }
func (p *MyPlugin) Description() string { return "我的自定义插件" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Init() error { /* 初始化逻辑 */ return nil }
func (p *MyPlugin) Shutdown() error { /* 关闭逻辑 */ return nil }

// 使用推荐的 GetRoutes 方法注册路由
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
    }
}

// 路由处理函数
func (p *MyPlugin) handleIndex(c *gin.Context) {
    c.JSON(200, gin.H{
        "plugin": p.Name(),
        "version": p.Version(),
    })
}

func (p *MyPlugin) handleGetData(c *gin.Context) {
    id := c.Query("id")
    c.JSON(200, gin.H{
        "id": id,
        "data": "示例数据",
    })
}

// 中间件示例
func (p *MyPlugin) logMiddleware(c *gin.Context) {
    // 记录请求日志
    c.Next()
}

// 为兼容性保留的 RegisterRoutes 方法
func (p *MyPlugin) RegisterRoutes(router *gin.Engine) {
    // 注意：推荐使用 GetRoutes 方法，此方法仅为兼容性保留
    // 这里可以保留空实现或添加日志提示
}

// 插件执行逻辑
func (p *MyPlugin) Execute(params map[string]interface{}) (interface{}, error) {
    // 实现插件功能
    return map[string]interface{}{"result": "success"}, nil
}
```

### 插件示例（旧的 RegisterRoutes 方法 - 仅为兼容性保留）
```go
// 注册插件路由（旧方式 - 不推荐）
func (p *MyPlugin) RegisterRoutes(router *gin.Engine) {
    group := router.Group(fmt.Sprintf("/plugins/%s", p.Name()))
    {
        group.GET("/", func(c *gin.Context) {
            c.JSON(200, gin.H{"plugin": p.Name()})
        })
        // 添加更多路由...
    }
}
```

### 两种路由注册方式的对比
| 特性 | GetRoutes 方法（推荐） | RegisterRoutes 方法（兼容性保留） |
|------|-----------------------|-----------------------------------|
| 路由定义 | 使用 Route 结构体数组 | 直接操作 gin.Engine 对象 |
| 元数据支持 | ✅ 完整支持 | ❌ 不支持 |
| 自动路由组 | ✅ 自动创建 | ❌ 需要手动创建 |
| 中间件管理 | ✅ 支持全局和路由级别 | ❌ 需要手动添加 |
| 文档生成 | ✅ 支持自动生成 API 文档 | ❌ 不支持 |

## 🤝 贡献指南

欢迎对项目进行贡献！感谢！

1. **Fork 仓库**并克隆到本地
2. **创建分支**进行开发（`git checkout -b feature/your-feature`）
3. **提交代码**并确保通过测试
4. **创建 Pull Request** 描述您的更改
5. 等待**代码审查**并根据反馈进行修改

---

### <div align="center"> <strong>✨ 持续更新完善中... ✨</strong> </div>



