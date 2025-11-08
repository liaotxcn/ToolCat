# Weave 插件脚手架工具

Weave 提供了一个高效生成Weave插件代码的脚手架，帮助开发者快速创建符合规范的插件项目结构。

## 功能特性

- 交互式命令行界面，引导用户输入插件信息
- 自动生成符合Plugin接口规范的插件代码
- 支持自定义插件名称、标识符、描述和版本
- 支持三种插件类型：API、Task、Event
- 自动生成插件文档
- 自动检测项目结构，将生成的文件放置到正确位置
- 支持非交互式模式，可以通过命令行参数直接生成插件
- 向上递归查找项目结构（最多3层），确保文件放置到正确位置

## 使用方法

### 快速运行（推荐）

不需要编译，直接使用Go命令运行：

```bash
# 在项目根目录下运行
cd /path/to/Weave
# Windows/Linux/Mac
go run tools/plugin_scaffold.go
```

### 编译运行（可选）

1. 确保已安装Go语言环境（版本1.16或更高）
2. 打开终端，进入项目根目录
3. 运行以下命令编译代码：
   ```bash
   go build -o plugin_scaffold.exe tools/plugin_scaffold.go  # Windows
   go build -o plugin_scaffold tools/plugin_scaffold.go      # Linux/Mac
   ```
4. 运行编译后的程序：
   ```bash
   # Windows
   plugin_scaffold.exe
   
   # Linux/Mac
   ./plugin_scaffold
   ```
5. 按照命令行提示输入插件信息

### 命令行参数

工具支持以下命令行参数：

```bash
plugin_scaffold [参数]
  -name string
        插件名称（驼峰式，以Plugin结尾，如TodoPlugin）
  -desc string
        插件描述
  -version string
        插件版本 (默认 "1.0.0")
  -type string
        插件类型（api, task, event） (默认 "api")
  -dir string
        输出目录
  -non-interactive
        非交互式模式
```

### 非交互式模式示例

```bash
# 快速创建一个名为TodoPlugin的API类型插件（编译后）
plugin_scaffold -name TodoPlugin -desc "待办事项管理插件" -version "1.0.0" -type api -non-interactive

# 使用go run直接运行（推荐）
go run tools/plugin_scaffold.go -name TodoPlugin -desc "待办事项管理插件" -version "1.0.0" -type api -non-interactive

# 指定输出目录
plugin_scaffold -name CalendarPlugin -dir /path/to/plugins -non-interactive
```

## 交互式配置选项

运行工具后，您需要输入以下信息：

1. **插件名称**：必须使用驼峰式命名并以Plugin结尾（如`TodoPlugin`）
2. **插件标识符**：必须使用小写字母和下划线（默认会根据插件名称自动生成）
3. **插件描述**：简要描述插件的功能和用途
4. **插件版本**：遵循语义化版本规范（默认：`1.0.0`）
5. **插件类型**：支持`api`（API类型）、`task`（任务类型）、`event`（事件类型）（默认：`api`）

## 生成的文件

工具会生成以下文件：

1. **插件代码文件**：`plugins/{plugin_name}_plugin.go`
   - 包含完整的Plugin接口实现，包括最新添加的`GetDependencies()`、`GetConflicts()`、`OnEnable()`、`OnDisable()`和`SetPluginManager()`方法
   - 预设了基本的路由和处理函数示例，使用推荐的`GetRoutes()`方法
   - 包含辅助方法的示例实现
   - 根据插件类型生成特定的功能代码
   - 内置插件依赖管理的基础结构
   - 支持热重载的生命周期管理

2. **插件文档文件**：`plugins/{plugin_name}_plugin.md`
   - 包含插件基本信息
   - 列出所有API接口及其说明
   - 提供配置说明和使用示例
   - 包含插件注册示例代码

## 生成的插件结构

生成的插件代码包含以下主要部分：

- **插件结构体**：用于存储插件状态和依赖，包含资源跟踪器
- **配置结构体**：用于插件配置管理
- **基础信息方法**：`Name()`, `Description()`, `Version()`
- **生命周期方法**：`Init()`, `Shutdown()`
- **路由定义方法**：`GetRoutes()`
- **中间件方法**：`GetDefaultMiddlewares()`
- **执行方法**：`Execute()`（根据插件类型自动调整）
- **辅助方法**：示例业务逻辑实现
  - `loadConfig()`: 加载插件配置
  - `trackResource()`: 跟踪创建的资源
  - `releaseAllResources()`: 释放所有资源

## 后续操作

生成插件代码后，您需要：

1. 编辑生成的插件文件，实现具体的业务逻辑
2. 在`main.go`中注册您的插件
3. 根据需要添加更多路由和功能
4. 测试插件功能
5. 更新插件文档

## 插件注册示例

在`main.go`文件中，您需要添加以下代码来注册您的插件：

```go
import (
	"weave/plugins"
	"weave/plugins/core"
)

func registerPlugins() {
	// 注册您的插件
	core.GlobalPluginManager.Register(plugins.{pluginInstance})
}
```

其中`{pluginInstance}`是生成的插件实例变量名，通常是插件名称的首字母小写形式。

例如，如果您生成的是`TodoPlugin`，那么注册代码应该是：

```go
core.GlobalPluginManager.Register(plugins.TodoPlugin)
```

## 插件类型详解

工具支持以下三种插件类型：

### API类型插件

**API类型插件**是最常用的插件类型，用于提供HTTP API接口。

- 自动生成基本的CRUD API路由模板
- 适合实现RESTful服务
- 包含路由定义和处理函数示例
- 预设了常见的API响应格式

### Task类型插件

**Task类型插件**用于实现定时任务或后台任务。

- 自动导入time包
- 适合实现周期性执行的任务
- 示例代码包含任务调度器的基本结构
- 可以在Init方法中启动任务

### Event类型插件

**Event类型插件**用于实现事件监听和事件处理。

- 自动导入events包
- 适合实现基于事件的功能
- 示例代码包含事件监听器的基本结构
- 可以订阅和发布系统事件

## 注意事项

- 插件名称必须遵循驼峰式命名规范并以Plugin结尾
- 插件标识符必须只包含小写字母和下划线
- 生成的代码是一个模板，需要根据实际需求进行修改和扩展
- 请确保在注册插件前已完成所有必要的初始化操作
- 如果有特殊的中间件需求，请实现`GetDefaultMiddlewares()`方法

## 常见问题

### Q: 为什么找不到plugins目录？
**A:** 工具会自动向上递归查找最多3层目录来找到plugins目录。如果仍然找不到，请确保您在正确的项目结构中运行工具，或者使用`-dir`参数指定输出目录。

### Q: 生成的插件无法编译？
**A:** 请检查是否正确安装了Go语言环境，以及是否安装了所有必要的依赖包。

### Q: 如何自定义插件的中间件？
**A:** 编辑生成的插件文件，实现`GetDefaultMiddlewares()`方法，添加您需要的中间件函数。

## 扩展建议

### 通用扩展建议

- 根据实际需求添加数据库模型和操作
- 添加缓存机制提高性能
- 实现插件配置管理功能
- 添加单元测试和集成测试
- 实现资源监控和日志记录功能

### API类型插件特定建议

- 实现更复杂的业务逻辑和API接口
- 添加请求验证和错误处理
- 实现API版本控制
- 添加API文档自动生成功能
- 实现限流和权限控制

### Task类型插件特定建议

- 实现任务调度器和任务队列
- 添加任务状态监控和重试机制
- 实现任务优先级管理
- 添加任务执行日志和统计
- 实现任务依赖关系管理

### Event类型插件特定建议

- 实现事件总线和事件分发器
- 添加事件过滤和转换功能
- 实现事件持久化和重播机制
- 添加事件处理错误重试策略
- 实现事件订阅管理界面

## 版本历史

- **1.0.0**：初始版本，提供基本的插件代码生成功能
- **1.1.0**：增加了插件类型支持（API、Task、Event），增加了非交互式模式，改进了项目结构检测
- **1.2.0**：修复了字符串格式化错误，优化了代码生成逻辑，改进了命令行参数解析，修复了Windows环境下的兼容性问题