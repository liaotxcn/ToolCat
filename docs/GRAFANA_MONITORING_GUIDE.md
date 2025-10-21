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
   - HTTP请求计数 (`http_requests_total`)：按方法、端点和状态码统计请求总数
   - HTTP请求延迟 (`http_request_duration_seconds`)：按方法、端点和状态码统计请求耗时分布

2. **数据库监控**
   - 数据库连接数 (`db_connections`)：显示当前活跃连接数
   - 数据库查询计数 (`db_queries_total`)：按操作类型和表名统计查询总数
   - 数据库查询延迟 (`db_query_duration_seconds`)：按操作类型和表名统计查询耗时分布

3. **插件监控**
   - 插件总数 (`plugins_total`)：显示所有已注册的插件数量
   - 已启用插件数 (`plugins_enabled`)：显示当前已启用的插件数量
   - 插件执行计数 (`plugin_execution_count`)：按插件名称和方法名统计执行次数
   - 插件执行延迟 (`plugin_execution_duration`)：按插件名称和方法名统计执行耗时
   - 插件方法调用 (`plugin_method_calls`)：按插件名称和方法名统计方法调用次数
   - 插件错误计数 (`plugin_errors`)：按插件名称和错误类型统计错误发生次数
   - 插件内存使用 (`plugin_memory_usage`)：按插件名称统计内存占用情况
   - 插件重载次数 (`plugin_reloads`)：按插件名称统计重载次数

4. **系统监控**
   - 内存使用 (`memory_usage_bytes`)：显示应用内存使用情况
   - 系统运行时间 (`system_uptime_seconds`)：显示应用已运行时长

5. **错误监控**
   - 错误计数 (`error_count`)：按错误类型和组件统计错误发生频率

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

### Grafana面板显示"无数据"或"No Data"

**症状**：Grafana面板显示"无数据"或"No Data"，但Prometheus服务看起来正常运行

**排查步骤**：

1. **验证Prometheus数据采集**：
   - 访问 http://localhost:9090 打开Prometheus界面
   - 点击顶部菜单的"Status" > "Targets"，检查ToolCat应用的抓取状态是否为"UP"
   - 如果状态显示为"DOWN"，检查错误消息并修复连接问题

2. **检查应用metrics端点**：
   - 确认ToolCat应用的/metrics端点是否正常工作
   - 可以通过以下命令检查：`curl http://localhost:8081/metrics`
   - 如果返回404或其他错误，需要检查应用的metrics初始化代码

3. **调整查询时间范围**：
   - 在Grafana界面右上角调整时间范围，尝试选择更长的时间段（如Last 6 hours）
   - 确认是否有历史数据显示

4. **检查Prometheus配置**：
   - 查看Prometheus配置文件是否正确配置了抓取间隔和目标
   - 检查配置文件路径：`pkg/metrics/prometheus.yml`
   - 确保抓取间隔设置合理（建议15秒或30秒）

5. **验证Grafana数据源连接**：
   - 在Grafana中进入数据源配置页面
   - 点击"Save & Test"按钮，确认连接状态
   - 如果显示连接成功但仍无数据，尝试重新保存数据源配置

6. **检查指标名称是否正确**：
   - 在Grafana面板中编辑查询，确认指标名称是否与Prometheus中实际采集的指标名称一致
   - 可以在Prometheus的Graph页面中使用自动补全功能确认正确的指标名称

### 监控指标不更新

**症状**：图表显示旧数据或数据长时间不变化

**排查步骤**：

1. 确认ToolCat应用正常运行且/metrics端点可访问
2. 检查Prometheus抓取配置，确认抓取间隔设置
3. 验证Grafana仪表盘的刷新间隔设置（默认10秒）
4. 检查系统时间同步，确保各容器间时间一致

### 仪表盘未自动加载

**症状**：Grafana中找不到ToolCat仪表盘

**排查步骤**：

1. 检查Docker卷挂载是否正确：`docker-compose ps grafana`
2. 检查provisioning配置文件格式是否正确：`pkg/grafana/provisioning/dashboards/`
3. 查看Grafana日志：`docker-compose logs grafana`
4. 尝试手动导入仪表盘JSON文件

### 常见指标查询问题

如果在使用指标查询时遇到问题，可以尝试以下基础查询验证系统是否正常：

```
# 检查基本的HTTP请求指标
sum(http_requests_total) by (method, endpoint)

# 检查系统指标（如果已配置）
process_resident_memory_bytes{job="toolcat"}

# 检查数据库连接
go_sql_stats_open_connections{database="toolcat"}
```

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