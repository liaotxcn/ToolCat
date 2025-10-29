# ToolCat API 接口文档

## 1. 概述

本文档提供 ToolCat 服务的API接口说明，包括认证接口、用户管理接口、工具管理接口和插件管理接口等。所有接口基于HTTP协议，使用JSON格式进行数据交换。

## 2. 基础信息

- 服务基础URL: `http://localhost:8081`
- API版本: v1
- API基础路径: `/api/v1`
- 认证方式: JWT (JSON Web Token)
- 数据格式: JSON
- CSRF保护: 启用（对于非GET/HEAD/OPTIONS/TRACE请求）

## 3. 认证机制

系统使用JWT (JSON Web Token)进行认证。用户登录成功后，服务器会返回一个JWT令牌，该令牌需要在后续的API请求中通过Authorization头传递。

JWT令牌包含用户的身份信息，有效期等。当令牌过期或无效时，API请求会返回401 Unauthorized错误。

## 4. 错误处理

所有API接口都使用标准的HTTP状态码来表示请求的结果：
- 200 OK: 请求成功
- 201 Created: 创建成功
- 400 Bad Request: 请求参数错误
- 401 Unauthorized: 未授权
- 403 Forbidden: 禁止访问（包含CSRF令牌验证失败）
- 404 Not Found: 资源不存在
- 429 Too Many Requests: 请求过于频繁，超出限流限制
- 500 Internal Server Error: 服务器错误

错误响应通常包含一个error字段，描述具体的错误信息。

## 5. CSRF保护机制

ToolCat服务启用了CSRF（跨站请求伪造）保护机制，对于非GET/HEAD/OPTIONS/TRACE的请求，需要进行CSRF令牌验证。

### 5.1 CSRF令牌获取

CSRF令牌会通过两种方式提供：

1. **Cookie**：服务器会在响应中设置名为`XSRF-TOKEN`的Cookie
2. **响应头**：服务器会在响应头中添加`X-CSRF-Token`字段

### 5.2 CSRF令牌使用

对于需要验证CSRF的请求（非GET/HEAD/OPTIONS/TRACE），需要同时满足以下条件：

1. 请求头中包含`X-CSRF-Token`字段，值为获取到的CSRF令牌
2. 请求中携带包含相同令牌值的`XSRF-TOKEN`Cookie

## 6. 认证接口

### 6.1 用户注册

**请求URL**: `/auth/register`
**请求方法**: POST
**请求体**: 
```json
{
  "username": "string",    // 用户名(必填，3-50个字符)
  "password": "string",    // 密码(必填，至少6个字符)
  "confirm_password": "string", // 确认密码(必填，必须与password一致)
  "email": "string"         // 邮箱(必填，有效的邮箱格式)
}
```

**成功响应**: 
```json
{
  "message": "注册成功",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2025-10-01T10:00:00Z",
    "updated_at": "2025-10-01T10:00:00Z"
  }
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败或用户名/邮箱已存在
```json
{
  "error": "错误信息"
}
```

### 6.2 用户登录

**请求URL**: `/auth/login`
**请求方法**: POST
**请求体**: 
```json
{
  "username": "string",    // 用户名(必填)
  "password": "string"     // 密码(必填)
}
```

**成功响应**: 
```json
{
  "message": "登录成功",
  "token": "JWT_TOKEN_HERE",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2025-10-01T10:00:00Z",
    "updated_at": "2025-10-01T10:00:00Z"
  }
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败
- 401 Unauthorized: 用户名或密码错误
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

## 7. API 接口 (需要认证)

所有API接口需要在请求头中包含JWT认证令牌：
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
```

### 7.1 用户管理接口

#### 7.1.1 获取所有用户

**请求URL**: `/api/v1/users`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}

**成功响应**: 
```json
[
  {
    "id": 1,
    "username": "testuser1",
    "email": "test1@example.com",
    "created_at": "2025-10-01T10:00:00Z",
    "updated_at": "2025-10-01T10:00:00Z"
  },
  {
    "id": 2,
    "username": "testuser2",
    "email": "test2@example.com",
    "created_at": "2025-10-02T11:00:00Z",
    "updated_at": "2025-10-02T11:00:00Z"
  }
]
```

**失败响应**: 
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.1.2 获取单个用户

**请求URL**: `/api/v1/users/:id`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- id: 用户ID

**成功响应**: 
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "created_at": "2025-10-01T10:00:00Z",
  "updated_at": "2025-10-01T10:00:00Z"
}
```

**失败响应**: 
- 404 Not Found: 用户不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.1.3 创建用户

**请求URL**: `/api/v1/users`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**请求体**: 
```json
{
  "username": "string",    // 用户名(必填，唯一)
  "password": "string",    // 密码(必填)
  "email": "string"         // 邮箱(唯一)
}
```

**成功响应**: 
```json
{
  "id": 3,
  "username": "newuser",
  "password": "hashed_password",
  "email": "new@example.com",
  "created_at": "2025-10-03T12:00:00Z",
  "updated_at": "2025-10-03T12:00:00Z"
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.1.4 更新用户

**请求URL**: `/api/v1/users/:id`
**请求方法**: PUT
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- id: 用户ID
**请求体**: 
```json
{
  "username": "string",    // 用户名(唯一)
  "password": "string",    // 密码
  "email": "string"         // 邮箱(唯一)
}
```

**成功响应**: 
```json
{
  "id": 1,
  "username": "updateduser",
  "password": "updated_hashed_password",
  "email": "updated@example.com",
  "created_at": "2025-10-01T10:00:00Z",
  "updated_at": "2025-10-04T13:00:00Z"
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败
- 404 Not Found: 用户不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.1.5 删除用户

**请求URL**: `/api/v1/users/:id`
**请求方法**: DELETE
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- id: 用户ID

**成功响应**: 
```json
{
  "message": "User deleted successfully"
}
```

**失败响应**: 
- 404 Not Found: 用户不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

### 7.2 工具管理接口

#### 7.2.1 获取所有工具

**请求URL**: `/api/v1/tools`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}

**成功响应**: 
```json
[
  {
    "id": 1,
    "name": "tool1",
    "description": "Description of tool 1",
    "icon": "tool1.png",
    "plugin_name": "plugin1",
    "is_enabled": true,
    "created_at": "2025-10-01T10:00:00Z",
    "updated_at": "2025-10-01T10:00:00Z"
  },
  {
    "id": 2,
    "name": "tool2",
    "description": "Description of tool 2",
    "icon": "tool2.png",
    "plugin_name": "plugin2",
    "is_enabled": true,
    "created_at": "2025-10-02T11:00:00Z",
    "updated_at": "2025-10-02T11:00:00Z"
  }
]
```

**失败响应**: 
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.2.2 获取单个工具

**请求URL**: `/api/v1/tools/:id`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- id: 工具ID

**成功响应**: 
```json
{
  "id": 1,
  "name": "tool1",
  "description": "Description of tool 1",
  "icon": "tool1.png",
  "plugin_name": "plugin1",
  "is_enabled": true,
  "created_at": "2025-10-01T10:00:00Z",
  "updated_at": "2025-10-01T10:00:00Z"
}
```

**失败响应**: 
- 404 Not Found: 工具不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.2.3 创建工具

**请求URL**: `/api/v1/tools`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**请求体**: 
```json
{
  "name": "string",         // 工具名称(必填，唯一)
  "description": "string",  // 工具描述
  "icon": "string",         // 工具图标路径
  "plugin_name": "string",  // 插件名称(必填)
  "is_enabled": true/false   // 是否启用
}
```

**成功响应**: 
```json
{
  "id": 3,
  "name": "newtool",
  "description": "Description of new tool",
  "icon": "newtool.png",
  "plugin_name": "plugin3",
  "is_enabled": true,
  "created_at": "2025-10-03T12:00:00Z",
  "updated_at": "2025-10-03T12:00:00Z"
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.2.4 更新工具

**请求URL**: `/api/v1/tools/:id`
**请求方法**: PUT
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- id: 工具ID
**请求体**: 
```json
{
  "name": "string",         // 工具名称(唯一)
  "description": "string",  // 工具描述
  "icon": "string",         // 工具图标路径
  "plugin_name": "string",  // 插件名称
  "is_enabled": true/false   // 是否启用
}
```

**成功响应**: 
```json
{
  "id": 1,
  "name": "updatedtool",
  "description": "Updated description",
  "icon": "updatedtool.png",
  "plugin_name": "plugin1",
  "is_enabled": false,
  "created_at": "2025-10-01T10:00:00Z",
  "updated_at": "2025-10-04T13:00:00Z"
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败
- 404 Not Found: 工具不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.2.5 删除工具

**请求URL**: `/api/v1/tools/:id`
**请求方法**: DELETE
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- id: 工具ID

**成功响应**: 
```json
{
  "message": "Tool deleted successfully"
}
```

**失败响应**: 
- 404 Not Found: 工具不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.2.6 执行工具

**请求URL**: `/api/v1/tools/:id/execute`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- id: 工具ID
**请求体**: 
```json
{
  // 工具执行所需的参数
}
```

**成功响应**: 
```json
{
  "tool_id": 1,
  "message": "Tool executed successfully"
}
```

**失败响应**: 
- 404 Not Found: 工具不存在
- 403 Forbidden: 工具已禁用
```json
{
  "error": "错误信息"
}
```

### 7.3 审计日志接口

#### 7.3.1 获取审计日志列表

**请求URL**: `/api/v1/audit/logs`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}
**查询参数**:
- page: 页码(可选，默认1)
- page_size: 每页数量(可选，默认10)
- start_time: 开始时间(可选，格式：2025-10-01T10:00:00Z)
- end_time: 结束时间(可选，格式：2025-10-02T10:00:00Z)
- user_id: 用户ID(可选)
- action: 操作类型(可选)

**成功响应**:
```json
{
  "logs": [
    {
      "id": 1,
      "user_id": 1,
      "username": "testuser",
      "action": "login",
      "resource_type": "auth",
      "resource_id": "1",
      "details": "用户登录成功",
      "ip": "127.0.0.1",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2025-10-01T10:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 10
}
```

**失败响应**:
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.3.2 获取单个审计日志详情

**请求URL**: `/api/v1/audit/logs/:id`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}
**URL参数**:
- id: 审计日志ID

**成功响应**:
```json
{
  "id": 1,
  "user_id": 1,
  "username": "testuser",
  "action": "login",
  "resource_type": "auth",
  "resource_id": "1",
  "details": "用户登录成功",
  "ip": "127.0.0.1",
  "user_agent": "Mozilla/5.0...",
  "created_at": "2025-10-01T10:00:00Z"
}
```

**失败响应**:
- 404 Not Found: 审计日志不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.3.3 获取审计日志统计信息

**请求URL**: `/api/v1/audit/stats`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}
**查询参数**:
- start_time: 开始时间(可选，格式：2025-10-01T10:00:00Z)
- end_time: 结束时间(可选，格式：2025-10-02T10:00:00Z)

**成功响应**:
```json
{
  "total_logs": 1000,
  "logs_per_day": [
    { "date": "2025-10-01", "count": 120 },
    { "date": "2025-10-02", "count": 150 }
  ],
  "actions_count": {
    "login": 300,
    "create": 200,
    "update": 150,
    "delete": 50
  }
}
```

**失败响应**:
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

### 7.4 插件管理接口

#### 7.4.1 获取所有插件

**请求URL**: `/api/v1/plugins`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}

**成功响应**:
```json
{
  "plugins": [
    {
      "name": "demo_plugin",
      "version": "1.0.0",
      "description": "示例插件",
      "enabled": true,
      "routes": [
        {
          "path": "/api/v1/demo",
          "method": "GET",
          "handler": "DemoHandler"
        }
      ],
      "dependencies": ["core_plugin"],
      "conflicts": ["conflicting_plugin"]
    }
  ]
}
```

**失败响应**:
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.4.2 获取插件状态

**请求URL**: `/api/v1/plugins/:name/status`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}
**URL参数**:
- name: 插件名称

**成功响应**:
```json
{
  "name": "demo_plugin",
  "enabled": true,
  "status": "running",
  "version": "1.0.0",
  "load_time": "2025-10-01T10:00:00Z"
}
```

**失败响应**:
- 404 Not Found: 插件不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.4.3 启用插件

**请求URL**: `/api/v1/plugins/:name/enable`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**URL参数**:
- name: 插件名称

**成功响应**:
```json
{
  "message": "插件启用成功",
  "plugin": {
    "name": "demo_plugin",
    "enabled": true,
    "version": "1.0.0",
    "status": "running"
  }
}
```

**失败响应**:
- 404 Not Found: 插件不存在
- 409 Conflict: 插件依赖冲突
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.4.4 禁用插件

**请求URL**: `/api/v1/plugins/:name/disable`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**URL参数**:
- name: 插件名称

**成功响应**:
```json
{
  "message": "插件禁用成功",
  "plugin": {
    "name": "demo_plugin",
    "enabled": false,
    "version": "1.0.0",
    "status": "disabled"
  }
}
```

**失败响应**:
- 404 Not Found: 插件不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.4.5 重载插件

**请求URL**: `/api/v1/plugins/:name/reload`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**URL参数**:
- name: 插件名称

**成功响应**:
```json
{
  "message": "插件重载成功",
  "plugin": {
    "name": "demo_plugin",
    "enabled": true,
    "version": "1.0.0",
    "reload_time": "2025-10-01T10:00:00Z"
  }
}
```

**失败响应**:
- 404 Not Found: 插件不存在
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.4.6 获取插件依赖图

**请求URL**: `/api/v1/plugins/dependency-graph`
**请求方法**: GET
**请求头**: Authorization: Bearer {token}

**成功响应**:
```json
{
  "nodes": [
    { "id": "core_plugin", "name": "核心插件", "enabled": true },
    { "id": "demo_plugin", "name": "示例插件", "enabled": true }
  ],
  "edges": [
    { "source": "demo_plugin", "target": "core_plugin", "type": "dependency" }
  ]
}
```

**失败响应**:
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

#### 7.4.7 加载插件

**请求URL**: `/api/v1/plugins/load`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**请求体**: 
```json
{
  // 插件加载所需的参数
}
```

**成功响应**: 
```json
{
  "message": "插件加载成功",
  "plugin": {
    "name": "demo_plugin",
    "version": "1.0.0",
    "status": "loaded",
    "enabled": false
  },
  "load_time": "2023-10-01T10:00:00Z"
}
```

#### 7.4.8 卸载插件

**请求URL**: `/api/v1/plugins/unload/:name`
**请求方法**: POST
**请求头**: Authorization: Bearer {token}
**URL参数**: 
- name: 插件名称

**成功响应**: 
```json
{
  "message": "插件卸载成功",
  "plugin": {
    "name": "demo_plugin",
    "status": "unloaded"
  },
  "unload_time": "2023-10-01T10:00:00Z"
}
```

## 8. 其他接口

### 8.1 根路径

**请求URL**: `/`
**请求方法**: GET

**响应**: 
```json
{
  "message": "欢迎使用工具猫(ToolCat)服务！",
  "version": "1.0.0",
  "api_base": "/api/v1",
  "health_check": "/health",
  "available_endpoints": ["/api/v1/users", "/api/v1/tools", "/api/v1/plugins", "/health", "/api/v1/audit/logs"],
  "timestamp": "2023-10-01T10:00:00Z"
}
```

### 8.2 健康检查

**请求URL**: `/health`
**请求方法**: GET

**响应**: 
```json
{
  "status": "ok",
  "timestamp": "2023-10-01T10:00:00Z",
  "version": "1.0.0"
}
```

## 9. 数据模型

### 9.1 用户模型(User)
```go
type User struct {
  ID        uint      `gorm:"primaryKey" json:"id"`
  Username  string    `gorm:"size:50;not null;unique" json:"username"`
  Password  string    `gorm:"size:100;not null" json:"password,omitempty"`
  Email     string    `gorm:"size:100;unique" json:"email"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}
```

### 9.2 工具模型(Tool)
```go
type Tool struct {
  ID          uint      `gorm:"primaryKey" json:"id"`
  Name        string    `gorm:"size:100;not null;unique" json:"name"`
  Description string    `gorm:"type:text" json:"description"`
  Icon        string    `gorm:"size:255" json:"icon"`
  PluginName  string    `gorm:"size:100;not null" json:"plugin_name"`
  IsEnabled   bool      `gorm:"default:true" json:"is_enabled"`
  CreatedAt   time.Time `json:"created_at"`
  UpdatedAt   time.Time `json:"updated_at"`
}
```

### 9.3 工具使用历史模型(ToolHistory)
```go
type ToolHistory struct {
  ID     uint      `gorm:"primaryKey" json:"id"`
  UserID uint      `json:"user_id"`
  ToolID uint      `json:"tool_id"`
  UsedAt time.Time `json:"used_at"`
  Params string    `gorm:"type:text" json:"params"`
  Result string    `gorm:"type:text" json:"result"`
}
```

### 9.4 登录历史模型(LoginHistory)
```go
type LoginHistory struct {
  ID        uint      `gorm:"primaryKey" json:"id"`
  Username  string    `gorm:"size:50;not null" json:"username"`
  IPAddress string    `gorm:"size:50" json:"ip_address"`
  Success   bool      `gorm:"not null" json:"success"`
  Message   string    `gorm:"size:255" json:"message"`
  UserAgent string    `gorm:"type:text" json:"user_agent"`
  LoginTime time.Time `json:"login_time"`
}
```

### 9.5 笔记模型(Note)
```go
type Note struct {
  ID          string    `gorm:"primaryKey;size:100" json:"id"`
  UserID      uint      `gorm:"not null;index" json:"user_id"`
  TenantID    uint      `gorm:"index" json:"tenant_id"`
  Title       string    `gorm:"size:255;not null;index" json:"title"`
  Content     string    `gorm:"type:text;not null" json:"content"`
  CreatedTime time.Time `gorm:"index" json:"created_time"`
  UpdatedTime time.Time `json:"updated_time"`
}
```

## 10. Note插件接口

Note插件是一个记事本插件，可以实现事件记录的增删查改功能。所有Note插件接口位于`/plugins/note`路径下。

### 10.1.1 获取插件信息

**请求URL**: `/plugins/note/`
**请求方法**: GET

**成功响应**: 
```json
{
  "plugin": "note",
  "name": "note",
  "description": "一个记事本插件，可以实现事件记录的增删查改功能",
  "version": "1.0.0",
  "endpoints": [
    "GET /plugins/note/ - 获取插件信息",
    "GET /plugins/note/notes - 获取所有笔记（需认证；按租户与用户隔离）",
    "GET /plugins/note/notes/:id - 获取单个笔记（需认证；按租户与用户隔离）",
    "POST /plugins/note/notes - 创建新笔记（需认证；按租户与用户隔离）",
    "PUT /plugins/note/notes/:id - 更新笔记（需认证；按租户与用户隔离）",
    "DELETE /plugins/note/notes/:id - 删除笔记（需认证；按租户与用户隔离）",
    "GET /plugins/note/notes/search - 搜索笔记（需认证；按租户与用户隔离）"
  ]
}
```

### 10.1.2 获取所有笔记

**请求URL**: `/plugins/note/notes`
**请求方法**: GET
**认证**: 需要携带 `Authorization: Bearer <token>`
**查询参数**: 
- page: 页码 (可选，默认1)
- page_size: 每页数量 (可选，默认10)

**成功响应**: 
```json
{
  "total": 100,
  "page": 1,
  "pageSize": 10,
  "totalPages": 10,
  "notes": [
    {
      "id": "note-12345678-1234-1234-1234-1234567890ab",
      "title": "测试笔记标题",
      "content": "测试笔记内容",
      "created_time": "2025-10-01T10:00:00Z",
      "updated_time": "2025-10-01T10:00:00Z"
    }
  ]
}
```

**失败响应**: 
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

### 10.1.3 获取单个笔记

**请求URL**: `/plugins/note/notes/:id`
**请求方法**: GET
**认证**: 需要携带 `Authorization: Bearer <token>`
**URL参数**: 
- id: 笔记ID

**成功响应**: 
```json
{
  "id": "note-12345678-1234-1234-1234-1234567890ab",
  "title": "测试笔记标题",
  "content": "测试笔记内容",
  "created_time": "2025-10-01T10:00:00Z",
  "updated_time": "2025-10-01T10:00:00Z"
}
```

**失败响应**: 
- 404 Not Found: 笔记不存在或无权限访问
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

### 10.1.4 创建新笔记

**请求URL**: `/plugins/note/notes`
**请求方法**: POST
**认证**: 需要携带 `Authorization: Bearer <token>`
**请求体**: 
```json
{
  "title": "字符串",  // 标题(必填)
  "content": "字符串" // 内容(必填)
}
```

**成功响应**: 
```json
{
  "id": "note-12345678-1234-1234-1234-1234567890ab",
  "title": "测试笔记标题",
  "content": "测试笔记内容",
  "created_time": "2025-10-01T10:00:00Z",
  "updated_time": "2025-10-01T10:00:00Z"
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

### 10.1.5 更新笔记

**请求URL**: `/plugins/note/notes/:id`
**请求方法**: PUT
**认证**: 需要携带 `Authorization: Bearer <token>`
**URL参数**: 
- id: 笔记ID
**请求体**: 
```json
{
  "title": "字符串",  // 标题(可选)
  "content": "字符串" // 内容(可选)
}
```

**成功响应**: 
```json
{
  "id": "note-12345678-1234-1234-1234-1234567890ab",
  "title": "更新后的笔记标题",
  "content": "更新后的笔记内容",
  "created_time": "2025-10-01T10:00:00Z",
  "updated_time": "2025-10-04T13:00:00Z"
}
```

**失败响应**: 
- 400 Bad Request: 请求参数验证失败
- 404 Not Found: 笔记不存在或无权限访问
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

### 10.1.6 删除笔记

**请求URL**: `/plugins/note/notes/:id`
**请求方法**: DELETE
**认证**: 需要携带 `Authorization: Bearer <token>`
**URL参数**: 
- id: 笔记ID

**成功响应**: 
```json
{
  "message": "删除成功"
}
```

**失败响应**: 
- 404 Not Found: 笔记不存在或无权限访问
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

### 10.1.7 搜索笔记

**请求URL**: `/plugins/note/notes/search`
**请求方法**: GET
**认证**: 需要携带 `Authorization: Bearer <token>`
**查询参数**: 
- keyword: 搜索关键词 (必填)
- page: 页码 (可选，默认1)
- page_size: 每页数量 (可选，默认10)

**成功响应**: 
```json
{
  "total": 10,
  "page": 1,
  "pageSize": 10,
  "totalPages": 1,
  "notes": [
    {
      "id": "note-12345678-1234-1234-1234-1234567890ab",
      "title": "包含关键词的笔记标题",
      "content": "包含关键词的笔记内容",
      "created_time": "2025-10-01T10:00:00Z",
      "updated_time": "2025-10-01T10:00:00Z"
    }
  ]
}
```

**失败响应**: 
- 500 Internal Server Error: 服务器错误
```json
{
  "error": "错误信息"
}
```

## 11. 安全提醒

1. 不要在客户端存储用户密码
2. 妥善保管JWT令牌，避免泄露
3. 定期更换密码和刷新令牌
4. 敏感操作前进行二次验证