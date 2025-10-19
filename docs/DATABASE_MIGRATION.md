# ToolCat 数据库迁移工具

## 1. 概述

本文档提供 ToolCat 数据库迁移工具的详细使用说明。该工具基于 `golang-migrate` 库实现，用于管理数据库结构的版本化迁移，支持迁移的应用、回滚、状态查询等功能，确保数据库结构变更可以被安全、可追踪地管理。

## 2. 迁移工具简介

数据库迁移工具是一个独立的命令行应用程序，位于项目的 `pkg/migrate` 目录下。它提供了以下核心功能：

- 应用未执行的数据库迁移
- 回滚已应用的迁移
- 创建新的迁移文件
- 查询当前数据库迁移状态
- 生成初始数据库架构迁移

## 3. 迁移文件结构

### 3.1 存储路径

迁移文件存储在项目的 `pkg/migrate/data_sql` 目录中。

### 3.2 文件命名规则

每个迁移包含两个文件：

1. 应用迁移文件：`{version}_{name}.up.sql`
2. 回滚迁移文件：`{version}_{name}.down.sql`

其中：
- `{version}`：3位数字的版本号（如：001、002、003等）
- `{name}`：迁移的描述性名称，使用下划线分隔单词

### 3.3 文件内容

- **up.sql**：包含要应用的SQL语句，如创建表、添加字段等
- **down.sql**：包含对应的回滚SQL语句，如删除表、移除字段等

每个迁移文件应包含单一的SQL语句或相关的SQL语句集合，确保迁移可以原子性地应用和回滚。

### 3.4 迁移文件示例

**创建用户表示例（001_create_users_table.up.sql）**：
```sql
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**回滚创建用户表（001_create_users_table.down.sql）**：
```sql
DROP TABLE IF EXISTS users;
```

**添加用户状态字段示例（002_add_user_status.up.sql）**：
```sql
ALTER TABLE users ADD COLUMN status ENUM('active', 'inactive', 'suspended') DEFAULT 'active' AFTER password;
```

**回滚添加用户状态字段（002_add_user_status.down.sql）**：
```sql
ALTER TABLE users DROP COLUMN status;
```

## 4. 工具使用方法

### 4.1 基本命令格式

```bash
go run ./pkg/migrate/main.go <command> [arguments]
```

### 4.2 可用命令

#### 4.2.1 up - 应用所有未执行的迁移

```bash
go run ./pkg/migrate/main.go up
```

此命令会应用所有尚未应用到数据库的迁移。

**成功输出示例：**
```
Database connection established successfully host=localhost port=3306 database=toolcat
Migrations applied successfully
```

**操作步骤：**
1. 确保数据库服务正在运行
2. 执行上述命令
3. 观察输出是否有错误信息
4. 可以使用`status`命令验证迁移是否成功应用

#### 4.2.2 down - 回滚最后一个应用的迁移

```bash
go run ./pkg/migrate/main.go down
```

此命令会回滚最后一个应用的迁移。

**注意事项：**
- 回滚操作会删除数据，请谨慎使用，特别是在生产环境
- 每次执行只会回滚一个迁移版本
- 如果需要回滚多个版本，可以多次执行此命令

#### 4.2.3 status - 查看当前迁移状态

```bash
go run ./pkg/migrate/main.go status
```

此命令会显示当前数据库的迁移状态，包括当前版本、是否处于脏状态以及可用的迁移版本。

**输出示例：**
```
Database connection established successfully host=localhost port=3306 database=toolcat
Current version: 1
Available versions: 1 2 3
```

**状态说明：**
- `Current version`：当前已应用到数据库的最新迁移版本
- `Available versions`：所有可用的迁移版本文件
- 如果数据库处于脏状态，会显示警告信息

#### 4.2.4 create - 创建新的迁移文件

```bash
go run ./pkg/migrate/main.go create <migration_name>
```

此命令会创建两个空的迁移文件（up.sql和down.sql），使用自动生成的版本号。

**参数：**
- `migration_name`：迁移的描述性名称，将用于生成文件名

**输出示例：**
```
Created migration files: pkg/migrate/data_sql/002_add_user_table.up.sql, pkg/migrate/data_sql/002_add_user_table.down.sql
```

**使用示例：**
```bash
# 创建添加工具表的迁移文件
go run ./pkg/migrate/main.go create add_tools_table

# 创建添加索引的迁移文件
go run ./pkg/migrate/main.go create add_indexes_to_users_table
```

创建后，需要手动编辑生成的SQL文件，添加相应的SQL语句。

#### 4.2.5 init - 初始化迁移环境

```bash
go run ./pkg/migrate/main.go init
```

此命令会初始化迁移环境，确保迁移目录存在并与数据库连接。

## 5. 迁移版本管理

### 5.1 版本号规则

- 版本号从0开始递增
- 版本号格式为3位数字，不足3位前面补0（如：001、010、100）
- 版本号必须连续，不允许跳号

### 5.2 迁移表

迁移工具使用数据库中的 `schema_migrations` 表来跟踪已应用的迁移版本。该表包含以下字段：

- `version`：当前数据库版本号
- `dirty`：布尔值，表示数据库是否处于脏状态（上次迁移是否失败）

## 6. 常见问题处理

### 6.1 脏数据库状态

如果迁移执行失败，数据库可能会处于脏状态。在执行新的迁移前，需要先解决脏状态问题。

**错误信息示例：**
```
Failed to apply migrations: failed to apply migrations: Dirty database version 1. Fix and force version.
```

**解决方法：**
1. 检查并修复导致迁移失败的问题（通常是SQL语法错误或数据库约束冲突）
2. 手动修改数据库使其与迁移文件预期状态一致
3. 解决脏标记方法：
   ```bash
   # 连接到MySQL数据库
   mysql -h localhost -u root -p toolcat
   
   # 查看schema_migrations表
   SELECT * FROM schema_migrations;
   
   # 将dirty字段设置为0
   UPDATE schema_migrations SET dirty = 0;
   ```
4. 如果问题无法解决，可以考虑重置迁移状态，但请注意这将丢失迁移历史

### 6.2 迁移文件路径问题

如果迁移工具无法找到迁移文件，请检查：

1. 迁移文件是否位于 `pkg/migrate/data_sql` 目录下
2. 迁移文件命名是否符合 `{version}_{name}.up.sql` 和 `{version}_{name}.down.sql` 格式
3. 版本号是否为3位数字格式

### 6.3 SQL语法错误

**错误信息示例：**
```
migration failed in line 0: <SQL语句> (details: Error 1064 (42000): You have an error in your SQL syntax)
```

**解决方法：**
1. 检查迁移文件中的SQL语法
2. 确保每个迁移文件中的SQL语句都是有效的
3. 避免在单个迁移文件中包含多个不相关的SQL语句

## 7. 最佳实践

### 7.1 迁移文件编写

- 为每个逻辑变更创建单独的迁移文件
- 使用描述性的迁移名称
- 确保每个up迁移都有对应的down迁移
- 迁移文件应该是幂等的（可以安全地多次执行）
- 在生产环境应用迁移前，先在测试环境验证

### 7.2 迁移版本管理

- 版本号应由迁移工具自动生成，避免手动指定
- 团队协作时，确保迁移版本号不会冲突
- 记录迁移的应用历史，便于问题排查

### 7.3 测试和验证

- 在应用迁移前，备份数据库
- 迁移应用后，验证数据完整性和应用功能
- 定期执行状态检查，确保数据库结构与代码期望一致

## 8. 高级功能

### 8.1 生成初始数据库架构

迁移工具支持基于当前模型生成初始数据库架构的迁移文件。使用此功能可以快速创建初始数据库结构。

### 8.2 与CI/CD集成

数据库迁移可以集成到持续集成/持续部署流程中，确保每次代码部署时数据库结构也能自动更新。

## 9. 注意事项

- 生产环境的数据库迁移操作应当谨慎执行，建议在低峰期进行
- 大型数据库迁移可能需要较长时间，应合理安排维护窗口
- 迁移过程中可能需要锁定表，会影响应用的正常运行
- 复杂的数据库迁移可能需要编写自定义脚本来处理数据转换

## 10. 故障恢复

如果迁移操作失败，建议采取以下步骤：

1. 分析错误日志，确定失败原因
2. 根据需要修改迁移文件或修复数据库状态
3. 在开发环境验证修复方案
4. 再次尝试应用迁移
5. 如果问题严重，可以考虑回滚到上一个稳定版本

通过遵循本指南中的最佳实践和建议，可以确保数据库迁移过程的安全性和可维护性，减少因数据库变更带来的风险。