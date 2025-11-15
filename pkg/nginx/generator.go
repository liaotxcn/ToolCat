package nginx

import (
	"fmt"
	"strings"
)

// Generator Nginx配置生成器
type Generator struct {
	config *NginxConfig
}

// NewGenerator 创建配置生成器
func NewGenerator(config *NginxConfig) *Generator {
	return &Generator{
		config: config,
	}
}

// Generate 生成完整的Nginx配置
func (g *Generator) Generate() (string, error) {
	if err := g.config.Validate(); err != nil {
		return "", fmt.Errorf("配置验证失败: %v", err)
	}

	var builder strings.Builder

	// 生成全局配置
	g.generateGlobal(&builder)

	// 生成events配置
	g.generateEvents(&builder)

	// 生成http配置
	g.generateHttp(&builder)

	return builder.String(), nil
}

// GenerateUpstream 生成upstream配置
func (g *Generator) GenerateUpstream(upstream UpstreamConfig) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("upstream %s {\n", upstream.Name))

	// 添加负载均衡算法
	if method := upstream.Method.String(); method != "" {
		builder.WriteString(fmt.Sprintf("    %s\n", method))
	}

	// 添加keepalive配置
	if upstream.KeepAlive > 0 {
		builder.WriteString(fmt.Sprintf("    keepalive %d;\n", upstream.KeepAlive))
	}

	// 添加服务器配置
	for _, server := range upstream.Servers {
		serverConfig := fmt.Sprintf("    server %s", server.Host)

		// 添加权重
		if upstream.Method.IsWeighted() && server.Weight > 0 {
			serverConfig += fmt.Sprintf(" weight=%d", server.Weight)
		}

		// 添加备份标记
		if server.Backup {
			serverConfig += " backup"
		}

		// 添加宕机标记
		if server.Down {
			serverConfig += " down"
		}

		// 添加失败配置
		if server.MaxFails > 0 {
			serverConfig += fmt.Sprintf(" max_fails=%d", server.MaxFails)
		}

		if server.FailTimeout > 0 {
			serverConfig += fmt.Sprintf(" fail_timeout=%ds", int(server.FailTimeout.Seconds()))
		}

		serverConfig += ";\n"
		builder.WriteString(serverConfig)
	}

	// 添加健康检查配置
	if upstream.HealthCheck.Enabled {
		g.generateHealthCheck(&builder, upstream.HealthCheck)
	}

	builder.WriteString("}\n\n")

	return builder.String()
}

// generateGlobal 生成全局配置
func (g *Generator) generateGlobal(builder *strings.Builder) {
	builder.WriteString("# Nginx配置文件 - 由Weave项目自动生成\n\n")

	// worker_processes
	if g.config.WorkerProcesses > 0 {
		builder.WriteString(fmt.Sprintf("worker_processes %d;\n", g.config.WorkerProcesses))
	} else {
		builder.WriteString("worker_processes auto;\n")
	}

	// 全局日志
	if g.config.ErrorLog != "" {
		builder.WriteString(fmt.Sprintf("error_log %s;\n", g.config.ErrorLog))
	}

	if g.config.PidFile != "" {
		builder.WriteString(fmt.Sprintf("pid %s;\n", g.config.PidFile))
	}

	builder.WriteString("\n")
}

// generateEvents 生成events配置
func (g *Generator) generateEvents(builder *strings.Builder) {
	builder.WriteString("events {\n")

	// worker_connections
	if g.config.WorkerConnections > 0 {
		builder.WriteString(fmt.Sprintf("    worker_connections %d;\n", g.config.WorkerConnections))
	}

	// 其他events配置
	if g.config.Events != nil {
		for key, value := range g.config.Events {
			if key != "worker_connections" {
				builder.WriteString(fmt.Sprintf("    %s %v;\n", key, value))
			}
		}
	}

	builder.WriteString("}\n\n")
}

// generateHttp 生成http配置
func (g *Generator) generateHttp(builder *strings.Builder) {
	builder.WriteString("http {\n")

	// 基础配置
	g.generateHttpBasic(builder)

	// 日志格式
	g.generateLogFormats(builder)

	// 全局日志
	if g.config.AccessLog != "" {
		builder.WriteString(fmt.Sprintf("    access_log %s main;\n", g.config.AccessLog))
	}

	// 超时配置
	if g.config.Http.SendTimeout > 0 {
		builder.WriteString(fmt.Sprintf("    send_timeout %ds;\n", int(g.config.Http.SendTimeout.Seconds())))
	}

	if g.config.Http.KeepAliveTimeout > 0 {
		builder.WriteString(fmt.Sprintf("    keepalive_timeout %ds;\n", int(g.config.Http.KeepAliveTimeout.Seconds())))
	}

	// Gzip配置
	g.generateGzip(builder)

	// 包含其他配置文件
	for _, include := range g.config.Http.Include {
		builder.WriteString(fmt.Sprintf("    include %s;\n", include))
	}

	builder.WriteString("\n")

	// 生成upstream配置
	for _, upstream := range g.config.Http.Upstreams {
		upstreamConfig := g.GenerateUpstream(upstream)
		builder.WriteString(strings.ReplaceAll(upstreamConfig, "\n", "\n    "))
	}

	// 生成server配置
	for _, server := range g.config.Http.Servers {
		g.generateServer(builder, server)
	}

	builder.WriteString("}\n")
}

// generateHttpBasic 生成http基础配置
func (g *Generator) generateHttpBasic(builder *strings.Builder) {
	builder.WriteString("    include       mime.types;\n")
	builder.WriteString("    default_type  application/octet-stream;\n\n")

	builder.WriteString("    sendfile        on;\n")
	builder.WriteString("    tcp_nopush      on;\n")
	builder.WriteString("    tcp_nodelay     on;\n\n")
}

// generateLogFormats 生成日志格式
func (g *Generator) generateLogFormats(builder *strings.Builder) {
	builder.WriteString("    log_format main '$remote_addr - $remote_user [$time_local] \"$request\" '\n")
	builder.WriteString("                    '$status $body_bytes_sent \"$http_referer\" '\n")
	builder.WriteString("                    '\"$http_user_agent\" \"$http_x_forwarded_for\"';\n\n")
}

// generateGzip 生成gzip配置
func (g *Generator) generateGzip(builder *strings.Builder) {
	if g.config.Http.Gzip {
		builder.WriteString("    gzip on;\n")

		if g.config.Http.GzipCompLevel > 0 {
			builder.WriteString(fmt.Sprintf("    gzip_comp_level %d;\n", g.config.Http.GzipCompLevel))
		}

		if len(g.config.Http.GzipTypes) > 0 {
			builder.WriteString(fmt.Sprintf("    gzip_types %s;\n", strings.Join(g.config.Http.GzipTypes, " ")))
		}

		builder.WriteString("    gzip_vary on;\n")
		builder.WriteString("    gzip_min_length 1024;\n\n")
	}
}

// generateServer 生成server配置
func (g *Generator) generateServer(builder *strings.Builder, server ServerConfigBlock) {
	builder.WriteString("    server {\n")

	// 监听端口
	for _, listen := range server.Listen {
		builder.WriteString(fmt.Sprintf("        listen %s;\n", listen))
	}

	// 服务器名称
	for _, name := range server.ServerName {
		builder.WriteString(fmt.Sprintf("        server_name %s;\n", name))
	}

	// 日志配置
	if server.AccessLog != "" {
		builder.WriteString(fmt.Sprintf("        access_log %s main;\n", server.AccessLog))
	}

	if server.ErrorLog != "" {
		builder.WriteString(fmt.Sprintf("        error_log %s;\n", server.ErrorLog))
	}

	builder.WriteString("\n")

	// 生成location配置
	for _, location := range server.Locations {
		g.generateLocation(builder, location)
	}

	builder.WriteString("    }\n\n")
}

// generateLocation 生成location配置
func (g *Generator) generateLocation(builder *strings.Builder, location LocationConfig) {
	builder.WriteString(fmt.Sprintf("        location %s {\n", location.Path))

	// 代理到上游服务器
	if location.UpstreamName != "" {
		builder.WriteString(fmt.Sprintf("            proxy_pass http://%s;\n", location.UpstreamName))

		// 默认代理头设置
		defaultHeaders := map[string]string{
			"Host":              "$host",
			"X-Real-IP":         "$remote_addr",
			"X-Forwarded-For":   "$proxy_add_x_forwarded_for",
			"X-Forwarded-Proto": "$scheme",
		}

		// 合并自定义头
		allHeaders := make(map[string]string)
		for k, v := range defaultHeaders {
			allHeaders[k] = v
		}
		for k, v := range location.ProxySetHeaders {
			allHeaders[k] = v
		}

		// 生成proxy_set_header
		for key, value := range allHeaders {
			builder.WriteString(fmt.Sprintf("            proxy_set_header %s %s;\n", key, value))
		}

		// 代理超时设置
		builder.WriteString("            proxy_connect_timeout 30s;\n")
		builder.WriteString("            proxy_send_timeout 60s;\n")
		builder.WriteString("            proxy_read_timeout 60s;\n")
	}

	builder.WriteString("        }\n\n")
}

// generateHealthCheck 生成健康检查配置
func (g *Generator) generateHealthCheck(builder *strings.Builder, healthCheck HealthCheckConfig) {
	// 注意：标准Nginx不支持原生健康检查，这里使用第三方模块的配置格式
	// 实际使用时可能需要安装nginx_upstream_check_module等第三方模块

	builder.WriteString(fmt.Sprintf("    check interval=%ds fall=%d rise=%d timeout=%ds",
		int(healthCheck.Interval.Seconds()),
		healthCheck.Fall,
		healthCheck.Rise,
		int(healthCheck.Timeout.Seconds())))

	if healthCheck.Path != "" {
		builder.WriteString(fmt.Sprintf(" type=http uri=%s", healthCheck.Path))
	}

	builder.WriteString(";\n")
}

// GenerateUpstreamConfig 生成单独的upstream配置文件
func (g *Generator) GenerateUpstreamConfig() (string, error) {
	if len(g.config.Http.Upstreams) == 0 {
		return "", fmt.Errorf("没有配置上游服务器组")
	}

	var builder strings.Builder
	builder.WriteString("# Upstream配置 - 由Weave项目自动生成\n\n")

	for _, upstream := range g.config.Http.Upstreams {
		config := g.GenerateUpstream(upstream)
		builder.WriteString(config)
	}

	return builder.String(), nil
}

// GenerateServerConfig 生成单独的server配置文件
func (g *Generator) GenerateServerConfig() (string, error) {
	if len(g.config.Http.Servers) == 0 {
		return "", fmt.Errorf("没有配置HTTP服务器")
	}

	var builder strings.Builder
	builder.WriteString("# Server配置 - 由Weave项目自动生成\n\n")

	for _, server := range g.config.Http.Servers {
		g.generateServer(&builder, server)
	}

	return builder.String(), nil
}
