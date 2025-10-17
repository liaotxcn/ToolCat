# ToolCat 监控系统指南

本文档详细介绍了如何使用Prometheus和Grafana进行ToolCat应用的监控可视化配置与使用。

## 监控架构概述

ToolCat监控系统采用以下架构：

1. **ToolCat应用**：暴露Prometheus格式的监控指标（`/metrics`端点）
2. **Prometheus**：负责抓取、存储和处理监控指标
3. **Grafana**：提供直观的可视化界面，连接Prometheus数据源

## 快速开始

### 启动监控系统

监控系统已集成到Docker Compose配置中，只需启动完整的服务栈：

```bash
# 启动所有服务（ToolCat应用、MySQL、Prometheus、Grafana）
docker-compose up -d
```

### 访问监控界面

启动后，可以访问以下服务：

- **ToolCat应用**：http://localhost:8081
- **Prometheus界面**：http://localhost:9090
- **Grafana界面**：http://localhost:3000

## Grafana使用指南

### 登录Grafana

1. 打开浏览器访问 http://localhost:3000
2. 使用默认凭据登录：
   - 用户名：admin
   - 密码：admin

首次登录时，系统会提示您修改默认密码。

### 数据源配置

数据源已通过provisioning自动配置，指向Prometheus服务：

- 名称：Prometheus
- URL：http://prometheus:9090
- 默认数据源：是

如需手动检查或修改数据源配置：

1. 点击左侧菜单的 **Configuration（配置）> Data sources（数据源）**
2. 找到并点击Prometheus数据源进行查看或编辑

### 监控仪表盘

系统已自动配置并加载了ToolCat监控仪表盘。

#### 访问仪表盘

1. 点击左侧菜单的 **Dashboards（仪表盘）> Manage（管理）**
2. 在 **ToolCat** 文件夹下找到 **ToolCat 应用监控** 仪表盘
3. 点击打开仪表盘

#### 仪表盘内容

仪表盘包含以下主要面板：

1. **HTTP请求监控**
   - HTTP请求速率：按方法和端点显示请求频率
   - HTTP请求延迟（95th percentile）：显示各端点的请求延迟
   - HTTP错误率（5xx）：显示服务器错误率

2. **数据库监控**
   - 数据库连接数：显示当前活跃连接数
   - 数据库查询延迟（95th percentile）：按操作类型显示查询延迟
   - 数据库查询速率：按操作类型和表名显示查询频率

3. **插件监控**
   - 插件状态：显示总插件数和已启用插件数

4. **错误监控**
   - 错误率：按错误类型和组件显示错误发生频率

5. **系统监控**
   - 内存使用：显示应用内存使用情况
   - 系统运行时间：显示应用已运行时长

### 仪表盘使用技巧

1. **时间范围调整**：
   - 使用右上角的时间选择器调整显示的时间范围
   - 默认显示最近1小时的数据

2. **刷新间隔**：
   - 默认每10秒自动刷新一次数据
   - 可以在右上角手动调整刷新频率或暂停自动刷新

3. **面板操作**：
   - 点击任何面板右上角的下拉菜单可以进行放大、编辑、导出等操作
   - 悬停在图表上可以查看详细的数据点信息

4. **自定义面板**：
   - 可以克隆现有面板并修改查询或样式
   - 可以添加新面板监控自定义指标

## 常见问题排查

### Grafana无法连接到Prometheus

**症状**：面板显示"No data"或连接错误

**排查步骤**：

1. 确认Prometheus服务正在运行：`docker-compose ps prometheus`
2. 检查Grafana数据源配置是否正确
3. 检查Prometheus是否成功抓取到ToolCat指标：访问 http://localhost:9090/targets

### 监控指标不更新

**症状**：图表显示旧数据或不变化

**排查步骤**：

1. 确认ToolCat应用正常运行且/metrics端点可访问
2. 检查Prometheus抓取配置是否正确
3. 验证Grafana仪表盘的刷新间隔设置

### 仪表盘未自动加载

**症状**：Grafana中找不到ToolCat仪表盘

**排查步骤**：

1. 检查Docker卷挂载是否正确
2. 检查provisioning配置文件格式是否正确
3. 查看Grafana日志：`docker-compose logs grafana`

## 高级配置

### 创建自定义告警

可以在Grafana中配置告警规则，当监控指标达到特定阈值时触发通知：

1. 打开要设置告警的面板
2. 点击面板标题 > Edit > Alert
3. 配置告警条件、通知渠道和消息模板

### 导出和分享仪表盘

可以将仪表盘导出为JSON文件或通过链接分享：

1. 打开仪表盘
2. 点击右上角的Share按钮
3. 选择Share Link或Export选项

## 自定义仪表盘

如需创建自定义仪表盘，请参考以下步骤：

1. 点击左侧菜单的 **+ > Create > Dashboard**
2. 点击 **Add new panel**
3. 在查询编辑器中编写Prometheus查询
4. 配置面板选项、图例和显示样式
5. 保存面板和仪表盘

## 常用Prometheus查询示例

### HTTP指标查询

- 总请求数：`sum(http_requests_total)`
- 按状态码统计请求：`sum by (status) (http_requests_total)`
- 平均请求延迟：`avg(rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])) by (endpoint)`

### 数据库指标查询

- 数据库查询总量：`sum(db_queries_total)`
- 数据库连接使用率：`db_connections / 100 * 100`（假设最大连接数为100）

### 插件指标查询

- 插件启用率：`plugins_enabled / plugins_total * 100`

### 错误指标查询

- 错误总数：`sum(errors_total)`
- 按组件分类的错误：`sum by (component) (errors_total)`

## 监控最佳实践

1. **定期检查监控仪表盘**：建立日常监控习惯，关注系统性能趋势
2. **设置适当的告警阈值**：根据应用特性和用户体验要求设置合理的告警条件
3. **保留足够的历史数据**：配置适当的存储保留期，以便进行趋势分析
4. **关注关键指标**：重点监控影响用户体验的关键指标，如响应时间、错误率等
5. **定期优化查询**：确保监控查询高效，避免对系统性能造成额外负担

## 升级指南

当ToolCat应用升级时，监控系统可能需要相应更新：

1. 检查是否有新增的监控指标
2. 更新仪表盘以包含新指标
3. 调整告警规则以适应新的性能特征

---

如有任何问题或需要进一步的监控配置支持，请参考Prometheus和Grafana的官方文档。