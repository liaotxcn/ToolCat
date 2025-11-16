package nginx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager 配置管理器
type Manager struct {
	config          *NginxConfig
	configFile      string
	outputDir       string
	generator       *Generator
	healthChecker   *MultiHealthChecker
	templates       map[string]string
}

// NewManager 创建配置管理器
func NewManager(configFile, outputDir string) *Manager {
	config := NewDefaultConfig()
	
	return &Manager{
		config:        config,
		configFile:    configFile,
		outputDir:     outputDir,
		generator:     NewGenerator(config),
		healthChecker: NewMultiHealthChecker(),
		templates:     make(map[string]string),
	}
}

// LoadConfig 加载配置文件
func (m *Manager) LoadConfig() error {
	if _, err := os.Stat(m.configFile); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置
		return m.SaveConfig()
	}

	data, err := ioutil.ReadFile(m.configFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config NginxConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	m.config = &config
	m.generator = NewGenerator(&config)

	return nil
}

// SaveConfig 保存配置文件
func (m *Manager) SaveConfig() error {
	// 确保目录存在
	dir := filepath.Dir(m.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	data, err := json.MarshalIndent(m.config, "", "    ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := ioutil.WriteFile(m.configFile, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// GenerateConfigs 生成Nginx配置文件
func (m *Manager) GenerateConfigs() error {
	// 确保输出目录存在
	if err := os.MkdirAll(m.outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 生成主配置文件
	mainConfig, err := m.generator.Generate()
	if err != nil {
		return fmt.Errorf("生成主配置失败: %v", err)
	}

	mainConfigFile := filepath.Join(m.outputDir, "nginx.conf")
	if err := ioutil.WriteFile(mainConfigFile, []byte(mainConfig), 0644); err != nil {
		return fmt.Errorf("写入主配置文件失败: %v", err)
	}

	// 生成upstream配置文件
	upstreamConfig, err := m.generator.GenerateUpstreamConfig()
	if err != nil {
		return fmt.Errorf("生成upstream配置失败: %v", err)
	}

	upstreamConfigFile := filepath.Join(m.outputDir, "upstream.conf")
	if err := ioutil.WriteFile(upstreamConfigFile, []byte(upstreamConfig), 0644); err != nil {
		return fmt.Errorf("写入upstream配置文件失败: %v", err)
	}

	// 生成server配置文件
	serverConfig, err := m.generator.GenerateServerConfig()
	if err != nil {
		return fmt.Errorf("生成server配置失败: %v", err)
	}

	serverConfigFile := filepath.Join(m.outputDir, "server.conf")
	if err := ioutil.WriteFile(serverConfigFile, []byte(serverConfig), 0644); err != nil {
		return fmt.Errorf("写入server配置文件失败: %v", err)
	}

	return nil
}

// GetConfig 获取配置
func (m *Manager) GetConfig() *NginxConfig {
	return m.config
}

// UpdateConfig 更新配置
func (m *Manager) UpdateConfig(config *NginxConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	m.config = config
	m.generator = NewGenerator(config)
	
	return m.SaveConfig()
}

// AddUpstream 添加上游服务器组
func (m *Manager) AddUpstream(upstream UpstreamConfig) error {
	m.config.AddUpstream(upstream)
	
	// 如果启用了健康检查，创建健康检查器
	if upstream.HealthCheck.Enabled {
		checker := NewHealthChecker(upstream.Name, upstream.HealthCheck, upstream.Servers)
		m.healthChecker.AddChecker(upstream.Name, checker)
	}
	
	return m.SaveConfig()
}

// RemoveUpstream 移除上游服务器组
func (m *Manager) RemoveUpstream(name string) error {
	if err := m.config.RemoveUpstream(name); err != nil {
		return err
	}
	
	// 停止健康检查器
	if checker, exists := m.healthChecker.checkers[name]; exists {
		checker.Stop()
		delete(m.healthChecker.checkers, name)
	}
	
	return m.SaveConfig()
}

// AddServer 添加HTTP服务器配置
func (m *Manager) AddServer(server ServerConfigBlock) error {
	m.config.AddServer(server)
	return m.SaveConfig()
}

// StartHealthCheck 启动健康检查
func (m *Manager) StartHealthCheck() {
	m.healthChecker.StartAll()
}

// StopHealthCheck 停止健康检查
func (m *Manager) StopHealthCheck() {
	m.healthChecker.StopAll()
}

// GetHealthStats 获取健康检查统计
func (m *Manager) GetHealthStats() map[string]HealthStats {
	return m.healthChecker.GetStats()
}

// GetHealthyServers 获取健康服务器列表
func (m *Manager) GetHealthyServers() map[string][]ServerConfig {
	return m.healthChecker.GetAllHealthyServers()
}

// CreateWeaveUpstream 创建Weave项目的上游配置
func (m *Manager) CreateWeaveUpstream(name string, servers []string, method LoadBalanceMethod) error {
	var serverConfigs []ServerConfig
	
	for i, server := range servers {
		weight := 1
		if method == WeightedRoundRobin {
			weight = i + 1 // 简单的权重分配
		}
		
		serverConfigs = append(serverConfigs, ServerConfig{
			Host:       server,
			Weight:     weight,
			MaxFails:   3,
			FailTimeout: 30 * time.Second,
		})
	}
	
	upstream := UpstreamConfig{
		Name:     name,
		Method:   method,
		Servers:  serverConfigs,
		HealthCheck: HealthCheckConfig{
			Enabled:  true,
			Path:     "/health",
			Interval: 10 * time.Second,
			Timeout:  5 * time.Second,
			Fall:     3,
			Rise:     2,
		},
		KeepAlive:         32,
		KeepAliveTimeout:  60 * time.Second,
		KeepAliveRequests: 100,
	}
	
	return m.AddUpstream(upstream)
}

// CreateWeaveServer 创建Weave项目的服务器配置
func (m *Manager) CreateWeaveServer(upstreamName string, listenPort int, serverNames []string) error {
	location := LocationConfig{
		Path:         "/",
		UpstreamName: upstreamName,
		ProxySetHeaders: map[string]string{
			"Host":             "$host",
			"X-Real-IP":        "$remote_addr",
			"X-Forwarded-For":  "$proxy_add_x_forwarded_for",
			"X-Forwarded-Proto": "$scheme",
			"Connection":        "keep-alive",
		},
	}
	
	server := ServerConfigBlock{
		Listen:     []string{fmt.Sprintf("0.0.0.0:%d", listenPort)},
		ServerName: serverNames,
		Locations:  []LocationConfig{location},
		AccessLog:  "logs/weave_access.log",
		ErrorLog:   "logs/weave_error.log",
	}
	
	return m.AddServer(server)
}

// ReloadNginx 重新加载Nginx配置
func (m *Manager) ReloadNginx() error {
	// 首先生成配置文件
	if err := m.GenerateConfigs(); err != nil {
		return fmt.Errorf("生成配置文件失败: %v", err)
	}

	// 测试配置文件
	nginxCmd := "nginx"
	if _, err := os.Stat("/usr/sbin/nginx"); err == nil {
		nginxCmd = "/usr/sbin/nginx"
	}
	
	// 测试配置
	if err := runCommand(nginxCmd, "-t"); err != nil {
		return fmt.Errorf("Nginx配置测试失败: %v", err)
	}
	
	// 重新加载配置
	if err := runCommand(nginxCmd, "-s", "reload"); err != nil {
		return fmt.Errorf("Nginx重新加载失败: %v", err)
	}
	
	return nil
}

// BackupConfig 备份配置
func (m *Manager) BackupConfig() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(m.outputDir, "backups", timestamp)
	
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("创建备份目录失败: %v", err)
	}
	
	// 备份配置文件
	backupFile := filepath.Join(backupDir, "config.json")
	data, err := json.MarshalIndent(m.config, "", "    ")
	if err != nil {
		return "", fmt.Errorf("序列化配置失败: %v", err)
	}
	
	if err := ioutil.WriteFile(backupFile, data, 0644); err != nil {
		return "", fmt.Errorf("写入备份文件失败: %v", err)
	}
	
	// 备份生成的配置文件
	configFiles := []string{"nginx.conf", "upstream.conf", "server.conf"}
	for _, file := range configFiles {
		src := filepath.Join(m.outputDir, file)
		dst := filepath.Join(backupDir, file)
		
		if _, err := os.Stat(src); err == nil {
			data, err := ioutil.ReadFile(src)
			if err != nil {
				return "", fmt.Errorf("读取配置文件失败: %v", err)
			}
			
			if err := ioutil.WriteFile(dst, data, 0644); err != nil {
				return "", fmt.Errorf("写入备份文件失败: %v", err)
			}
		}
	}
	
	return backupDir, nil
}

// RestoreConfig 恢复配置
func (m *Manager) RestoreConfig(backupDir string) error {
	configFile := filepath.Join(backupDir, "config.json")
	
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("读取备份配置文件失败: %v", err)
	}
	
	var config NginxConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析备份配置文件失败: %v", err)
	}
	
	return m.UpdateConfig(&config)
}

// ListBackups 列出所有备份
func (m *Manager) ListBackups() ([]string, error) {
	backupDir := filepath.Join(m.outputDir, "backups")
	
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return []string{}, nil
	}
	
	files, err := ioutil.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %v", err)
	}
	
	var backups []string
	for _, file := range files {
		if file.IsDir() {
			backups = append(backups, file.Name())
		}
	}
	
	return backups, nil
}

// runCommand 执行命令
func runCommand(name string, args ...string) error {
	cmd := fmt.Sprintf("%s %s", name, strings.Join(args, " "))
	return fmt.Errorf("命令执行功能需要实现: %s", cmd)
}

// ValidateConfig 验证配置
func (m *Manager) ValidateConfig() error {
	return m.config.Validate()
}

// GetConfigSummary 获取配置摘要
func (m *Manager) GetConfigSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"upstreams_count": len(m.config.Http.Upstreams),
		"servers_count":   len(m.config.Http.Servers),
		"worker_processes": m.config.WorkerProcesses,
		"worker_connections": m.config.WorkerConnections,
	}
	
	var upstreams []map[string]interface{}
	for _, upstream := range m.config.Http.Upstreams {
		upstreams = append(upstreams, map[string]interface{}{
			"name":    upstream.Name,
			"method":  upstream.Method,
			"servers": len(upstream.Servers),
			"health_check_enabled": upstream.HealthCheck.Enabled,
		})
	}
	summary["upstreams"] = upstreams
	
	var servers []map[string]interface{}
	for _, server := range m.config.Http.Servers {
		servers = append(servers, map[string]interface{}{
			"listen":       server.Listen,
			"server_name":  server.ServerName,
			"locations":    len(server.Locations),
		})
	}
	summary["servers"] = servers
	
	return summary
}