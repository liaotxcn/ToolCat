# ToolCat - A microkernel architecture plugin development/management platform developed with Golang, designed to provide high-performance, highly extensible, secure, and reliable plugin services

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Microkernel-Layered-6BA539?style=for-the-badge" alt="Architecture">
  <img src="https://img.shields.io/badge/Cloud_Native-3371E3?style=for-the-badge&logo=Docker&logoColor=white" alt="Cloud Native">
  <img src="https://img.shields.io/badge/Plugin System-FF6F00?style=for-the-badge&logo=ai&logoColor=white" alt="Plugin System">
</div>

## 📋 项目简介

ToolCat 基于 Golang 开发的高性能、高效率、插件化易扩展的工具服务平台。采用微内核+分层架构设计，允许开发者轻松集成并且高效管理各种工具服务，同时保持系统的高性能和可扩展性。

主要应用场景包括：
- 工具集成与管理平台
- 微服务组件聚合层
- 数据/服务流转中台
- API网关与服务编排
- 高效开发和原型验证平台

---

## 🏗️ 整体架构

<img width="1631" height="621" alt="image" src="https://github.com/user-attachments/assets/1eaa7b01-4afd-42d3-bb36-3f1dc3f63d23" />

ToolCat 采用**微内核+分层架构**设计模式，充分结合两种架构的优势，保障系统可用性，实现了高度的灵活性、可扩展性和良好的性能。

### 微内核与分层架构的融合

ToolCat 在微内核架构的基础上，融入了分层架构的设计思想，形成了一套完整且高效灵活的架构体系：

1. **微内核架构（插件体系）**：提供插件管理、生命周期控制和插件间通信机制
2. **分层架构（核心系统）**：将核心功能按关注点分离，形成清晰的层次结构

### 微内核架构组成

1. **核心内核（Core Kernel）**：提供基础运行时环境、插件管理、配置管理、日志服务、安全机制等基础功能
2. **插件系统（Plugin System）**：插件管理器负责插件的注册、生命周期管理、依赖解析和冲突检测
3. **扩展插件（Extensions）**：通过插件接口集成到核心系统，实现各种业务功能

### 分层架构组成

1. **接口层**：处理HTTP请求，包括路由管理和控制器
2. **业务层**：包含核心业务逻辑和插件系统
3. **数据层**：负责数据存储和访问
4. **基础设施层**：提供日志、配置、安全等服务

### 架构特点

**松耦合设计**：核心系统与插件之间通过定义良好的接口通信，降低模块间依赖

**热插拔能力**：插件可在运行时动态加载和卸载，无需重启系统

**功能隔离**：每个插件独立封装功能，拥有自己的命名空间和路由前缀

**依赖与冲突管理**：内置依赖解析和冲突检测机制，确保插件间和谐共存

**统一接口**：所有插件实现相同的`Plugin`接口，标准化开发流程

**可扩展性**：系统功能可按需扩展，无需修改内核代码

**层次清晰**：核心系统采用分层设计，代码组织合理，易于维护和扩展

**高性能**：分层设计优化了请求处理流程，提高系统响应速度

系统的核心是高效灵活的插件机制与分层结构的结合，允许功能模块以插件形式独立开发和部署，同时通过统一的接口进行交互。整体架构设计注重模块化、可扩展性和高性能。

---

## 🌟 项目特点

### 🏗️ 微内核+分层架构
- **核心稳定与层次清晰**：核心系统保持最小化，分层设计使代码组织合理，易于维护和扩展
- **功能扩展灵活**：通过插件机制按需扩展系统功能，无需修改内核代码
- **低耦合高内聚**：系统组件间松耦合，便于维护和升级
- **热插拔能力**：插件可在运行时动态加载和卸载，无需重启系统
- **功能隔离与统一管理**：每个插件独立封装功能，拥有自己的命名空间和路由前缀，同时核心服务通过分层架构统一管理
- **依赖与冲突管理**：内置依赖解析和冲突检测机制，确保插件间和谐共存
- **统一接口**：所有插件实现相同的`Plugin`接口，标准化开发流程
- **高性能**：分层设计优化了请求处理流程，提高系统响应速度

### 🚀 高性能/效率
- 基于 Gin 框架构建，处理请求速度快，并发能力强
- 数据库连接池优化，支持高并发访问
- 模块化架构设计，代码结构清晰，易于维护和扩展
- 支持环境变量覆盖，便于不同环境配置
- 高效路由管理，支持动态路由和参数绑定
- 分层架构优化了请求处理流程，提高系统响应速度

### 🔌 插件化易扩展
- 统一的插件接口设计，支持热插拔
- 插件管理器统一注册、管理和执行插件
- 插件可独立注册路由，拥有独立命名空间
- 插件依赖和冲突检测机制
- 脚手架工具便捷生成插件框架代码
- 示例插件（Hello、Note）展示了完整插件开发流程

### 🔒 安全可靠
- 基于 JWT 的认证授权系统
- 完善的 CSRF 保护机制
- 基于令牌桶算法的限流中间件
- 密码哈希存储与验证
- 详细的登录历史记录
- 统一的错误处理中间件
- 支持 HTTPS (可在配置中开启)
- 分层架构将安全机制统一封装在基础设施层，便于统一管理和维护

### 📊 可观测性
- 集成结构化日志系统 (zap)
- 健康检查接口，监控系统状态
- 详细的请求/响应日志
- 支持自定义监控指标
- 分层架构将监控功能独立封装，确保系统各层运行状态的可观测性

### 🚀 开发友好
- 完整的插件开发文档和示例
- 插件脚手架工具，快速生成插件模板
- 支持本地开发和 Docker 部署
- 清晰的项目结构和代码规范

---

## 📂 项目结构

ToolCat采用微内核+分层架构，项目结构清晰地反映了这一设计理念。核心系统采用分层组织，功能扩展则通过插件机制实现

```
├── .gitignore           # Git忽略文件配置
├── Dockerfile           # Docker构建文件
├── Makefile             # 构建脚本
├── config/              # 配置管理
├── controllers/         # API控制器[接口层]
│   ├── health_controller.go  # 健康检查控制器
│   ├── plugin_controller.go  # 插件管理控制器
│   ├── tool_controller.go    # 工具管理控制器
│   └── user_controller.go    # 用户管理控制器
├── docker-compose.yaml  # Docker Compose配置
├── docs/                # 项目文档
│   ├── API.md                # API文档
│   ├── PLUGIN_DEVELOPMENT_GUIDE.md  # 插件开发指南
│   └── PLUGIN_SCAFFOLD_USAGE.md     # 插件脚手架指南
├── go.mod               
├── main.go              
├── middleware/          # 中间件[接口层业务层中间]
│   ├── auth.go               # 认证
│   ├── buffer_request.go     # 请求缓冲
│   ├── cors.go               # CORS跨域
│   ├── csrf.go               # CSRF保护
│   ├── error_handler.go      # 错误处理
│   └── rate_limiter.go       # 限流
├── models/              # 数据模型[数据层]
├── pkg/                 # 公共包[基础设施层]
├── plugins/             # 插件系统[微内核架构核心]
│   ├── core/                 # 核心插件功能
│   ├── doc.go                # 插件包文档
│   ├── examples/             # 示例插件
│   │   ├── hello_plugin.go       # Hello插件示例
│   │   ├── sample_dependent_plugin.go # 依赖插件示例
│   │   └── sample_optimized_plugin.go # 优化插件示例
│   ├── features/             # 功能插件
│   ├── templates/            # 插件模板
│   └── watcher/              # 插件监控
├── routers/             # 路由定义注册
├── test/                # 单元/集成测试
├── tools/               # 开发工具
├── utils/               # 工具函数
└── web/                 # 前端代码
```

---

## 🧩 核心组件

ToolCat的核心组件设计充分体现了微内核+分层架构的思想，既保持了系统的灵活性和可扩展性，又保证了代码的结构清晰易维护

### 🔌 插件系统 - 微内核架构的核心实现
插件系统是ToolCat的核心组件，负责插件的注册、加载、卸载和生命周期管理。它实现了一套完整的插件机制，使系统能够以插件形式扩展功能。在微内核+分层架构中，插件系统连接了核心内核和各种业务扩展。

在微内核架构下，插件系统具备以下特性：
- **完整的生命周期管理**：从插件的初始化、注册、激活到关闭的全生命周期管理
- **自动依赖解析**：通过 `GetDependencies()` 方法自动解析插件间依赖关系
- **冲突检测机制**：通过 `GetConflicts()` 方法避免插件间功能冲突
- **路由自动注册**：支持两种路由注册方式，特别是推荐的 `GetRoutes()` 方法更符合微内核架构的设计理念
- **命名空间隔离**：每个插件拥有独立的命名空间，避免资源冲突
- **统一的中间件管理**：支持全局和插件级别的中间件配置

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

### 🔐 认证系统
认证系统位于分层架构的基础设施层，提供完善的身份认证和授权机制，支持多种认证方式，认证系统与插件系统紧密结合，确保插件的安全访问，同时通过分层设计实现了安全机制的统一管理
- 基于 JWT 的令牌认证
- 支持访问令牌和刷新令牌机制
- 密码哈希存储，增强安全性
- 登录历史记录，便于审计和追踪
- 基于角色的访问控制

### 🔄 中间件系统
中间件系统位于分层架构的接口层和业务层之间，支持全局中间件和插件级中间件，可用于日志记录、请求验证、性能监控等场景。中间件系统采用链式调用模式，灵活组合各种功能，体现了分层架构的请求处理优化
- 认证中间件：验证用户身份
- 限流中间件：防止API滥用
- CORS中间件：处理跨域请求
- CSRF保护中间件：防止跨站请求伪造
- 错误处理中间件：统一处理和记录错误

### 🩺 健康检查
健康检查功能覆盖了分层架构的各个层次，定期检查系统各组件的运行状态，确保系统稳定运行。支持自定义健康检查项，满足不同场景的需求。通过微内核+分层架构的设计，健康检查能够精确到每个插件和每个层次的运行状态
- 数据库连接健康检查
- 插件系统状态检查
- 整体系统健康评估
- 根据健康状态返回适当的HTTP状态码

---

## 快速开始

### 环境准备
- **Go 1.21+**（本地开发）
- **Docker** 和 **Docker Compose**（容器化部署）
- **Git**（用于克隆代码库）
- **MySQL 8.0+**（可选，如不使用Docker）

### 部署方式

#### 1. Docker Compose 部署（推荐）

1. 克隆代码库
```bash
git clone https://github.com/liaotxcn/toolcat.git
cd toolcat
```

2. 创建环境变量文件（可选但推荐）
创建`.env`文件，设置以下环境变量以增强安全性：
```bash
# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=your_secure_password
DB_NAME=toolcat

# JWT配置
JWT_SECRET=your_secure_jwt_secret
JWT_ACCESS_TOKEN_EXPIRY=60
JWT_REFRESH_TOKEN_EXPIRY=168

# MySQL配置
MYSQL_ROOT_PASSWORD=your_root_password

# 服务器配置
SERVER_PORT=8081

# 日志配置
LOG_LEVEL=info
DEV_MODE=false

# CSRF配置
CSRF_ENABLED=true
```

3. 启动服务
使用Docker Compose一键启动整个服务栈：
```bash
docker-compose up -d
```

   首次启动时，Docker Compose会自动：
   - 构建ToolCat应用的Docker镜像
   - 创建MySQL数据库容器
   - 配置网络和卷
   - 启动所有服务

4. 验证服务状态
查看所有服务是否正常运行：
```bash
docker-compose ps
```
正常情况下，`toolcat-app`和`toolcat-mysql`都应显示为`Up`状态。

5. 访问应用
服务启动后，可以通过以下URL访问ToolCat应用：
```
http://localhost:8081
```

### Docker Compose 命令

```bash
docker-compose down    // 停止服务
docker-compose logs -f toolcat-app   // 查看应用日志
docker-compose logs -f toolcat-mysql // 查看数据库日志
docker-compose exec toolcat-app /bin/sh             // 进入应用容器
docker-compose exec toolcat-mysql mysql -u root -p  // 进入数据库容器
docker-compose up --build -d        // 重新构建并启动服务
```

#### 2. 本地开发环境设置

1. 克隆代码库并进入项目目录
```bash
git clone https://github.com/liaotxcn/toolcat.git
cd toolcat
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
确保本地MySQL服务已启动，并创建数据库：
```sql
CREATE DATABASE toolcat;
```

4. 设置环境变量或修改`config/config.go`中的默认配置

5. 运行应用
```bash
go run main.go
```

6. 构建应用
```bash
go build
```

### 注意事项

1. **数据持久化**：MySQL数据存储在`mysql-data`卷中，确保数据不会丢失
2. **健康检查**：系统提供`/health`接口监控服务健康状态
3. **资源限制**：默认配置了CPU和内存限制，可根据实际需求在`docker-compose.yaml`中调整
4. **首次启动**：首次启动需要一些时间来构建镜像和初始化服务
5. **端口映射**：默认将容器的8081端口映射到主机的8081端口

服务将在 http://localhost:8081 启动。

---

## API文档

详细请阅读: [API文档](./docs/API.md)

## 插件开发指南 

详细请阅读: [插件开发指南](./docs/PLUGIN_DEVELOPMENT_GUIDE.md)

## 脚手架工具

详细请阅读: [插件脚手架工具](./docs/PLUGIN_SCAFFOLD_USAGE.md)

### 🔧 创建新插件

在ToolCat的微内核+分层架构下，创建新插件是扩展系统功能的主要方式。插件是一个实现了`Plugin`接口的Go结构体，通过这个接口，插件可以与核心系统进行交互。微内核架构提供了插件的灵活性，而分层架构则为插件内部的代码组织提供了良好的指导。

创建新插件非常简单，只需遵循以下步骤：
1. 实现 `plugins.Plugin` 接口，定义插件的基本信息、生命周期和功能
2. 在 `main.go` 的 `registerPlugins` 函数中注册插件

微内核+分层架构的插件开发优势：
- **低侵入性**：无需修改核心代码即可扩展系统功能
- **独立演进**：插件可以独立开发、测试和部署
- **标准化接口**：统一的插件接口简化了开发流程
- **灵活组合**：用户可以根据需求组合不同的插件
- **结构清晰**：分层架构思想指导插件内部代码组织，提高可维护性

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

---

## 🤝 贡献指南

欢迎对项目进行贡献！感谢！

1. **Fork 仓库**并克隆到本地
2. **创建分支**进行开发（`git checkout -b feature/your-feature`）
3. **提交代码**并确保通过测试
4. **创建 Pull Request** 描述您的更改
5. 等待**代码审查**并根据反馈进行修改

---

### <div align="center"> <strong>✨ 持续更新完善中... ✨</strong> </div>
