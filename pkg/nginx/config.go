package nginx

import (
	"fmt"
	"time"
)

// LoadBalanceMethod 负载均衡算法类型
type LoadBalanceMethod string

const (
	// RoundRobin 轮询
	RoundRobin LoadBalanceMethod = "round_robin"
	// LeastConn 最少连接
	LeastConn LoadBalanceMethod = "least_conn"
	// IPHash IP哈希
	IPHash LoadBalanceMethod = "ip_hash"
	// WeightedRoundRobin 加权轮询
	WeightedRoundRobin LoadBalanceMethod = "weighted_round_robin"
)

// ServerConfig 后端服务器配置
type ServerConfig struct {
	Host        string        `json:"host"`         // 服务器地址，如 localhost:8081
	Weight      int           `json:"weight"`       // 权重，用于加权轮询
	Backup      bool          `json:"backup"`       // 是否为备份服务器
	Down        bool          `json:"down"`         // 是否标记为宕机
	MaxFails    int           `json:"max_fails"`    // 最大失败次数
	FailTimeout time.Duration `json:"fail_timeout"` // 失败超时时间
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled  bool          `json:"enabled"`  // 是否启用健康检查
	Path     string        `json:"path"`     // 健康检查路径
	Interval time.Duration `json:"interval"` // 检查间隔
	Timeout  time.Duration `json:"timeout"`  // 检查超时
	Fall     int           `json:"fall"`     // 失败阈值
	Rise     int           `json:"rise"`     // 成功阈值
}

// UpstreamConfig 上游服务器组配置
type UpstreamConfig struct {
	Name              string            `json:"name"`                // 上游服务器组名称
	Method            LoadBalanceMethod `json:"method"`              // 负载均衡算法
	Servers           []ServerConfig    `json:"servers"`             // 后端服务器列表
	HealthCheck       HealthCheckConfig `json:"health_check"`        // 健康检查配置
	KeepAlive         int               `json:"keep_alive"`          // 保持连接数
	KeepAliveTimeout  time.Duration     `json:"keep_alive_timeout"`  // 保持连接超时
	KeepAliveRequests int               `json:"keep_alive_requests"` // 每个连接的最大请求数
}

// LocationConfig 位置配置
type LocationConfig struct {
	Path            string            `json:"path"`              // 路径匹配规则
	UpstreamName    string            `json:"upstream_name"`     // 上游服务器组名称
	ProxySetHeaders map[string]string `json:"proxy_set_headers"` // 代理头设置
}

// ServerConfigBlock HTTP服务器配置
type ServerConfigBlock struct {
	Listen     []string         `json:"listen"`      // 监听端口
	ServerName []string         `json:"server_name"` // 服务器名称
	Locations  []LocationConfig `json:"locations"`   // 位置配置
	AccessLog  string           `json:"access_log"`  // 访问日志路径
	ErrorLog   string           `json:"error_log"`   // 错误日志路径
}

// NginxConfig Nginx主配置
type NginxConfig struct {
	WorkerProcesses   int                    `json:"worker_processes"`   // 工作进程数
	WorkerConnections int                    `json:"worker_connections"` // 每个工作进程的最大连接数
	AccessLog         string                 `json:"access_log"`         // 全局访问日志
	ErrorLog          string                 `json:"error_log"`          // 全局错误日志
	PidFile           string                 `json:"pid_file"`           // PID文件路径
	Events            map[string]interface{} `json:"events"`             // events配置
	Http              struct {
		SendTimeout      time.Duration       `json:"send_timeout"`       // 发送超时
		KeepAliveTimeout time.Duration       `json:"keep_alive_timeout"` // 保持连接超时
		Gzip             bool                `json:"gzip"`               // 是否启用gzip
		GzipCompLevel    int                 `json:"gzip_comp_level"`    // gzip压缩级别
		GzipTypes        []string            `json:"gzip_types"`         // gzip压缩类型
		Include          []string            `json:"include"`            // 包含的配置文件
		Upstreams        []UpstreamConfig    `json:"upstreams"`          // 上游服务器组
		Servers          []ServerConfigBlock `json:"servers"`            // HTTP服务器配置
	} `json:"http"`
}

// NewDefaultConfig 创建默认Nginx配置
func NewDefaultConfig() *NginxConfig {
	return &NginxConfig{
		WorkerProcesses:   4,
		WorkerConnections: 1024,
		AccessLog:         "logs/access.log",
		ErrorLog:          "logs/error.log",
		PidFile:           "logs/nginx.pid",
		Events: map[string]interface{}{
			"worker_connections": 1024,
		},
		Http: struct {
			SendTimeout      time.Duration       `json:"send_timeout"`
			KeepAliveTimeout time.Duration       `json:"keep_alive_timeout"`
			Gzip             bool                `json:"gzip"`
			GzipCompLevel    int                 `json:"gzip_comp_level"`
			GzipTypes        []string            `json:"gzip_types"`
			Include          []string            `json:"include"`
			Upstreams        []UpstreamConfig    `json:"upstreams"`
			Servers          []ServerConfigBlock `json:"servers"`
		}{
			SendTimeout:      60 * time.Second,
			KeepAliveTimeout: 75 * time.Second,
			Gzip:             true,
			GzipCompLevel:    6,
			GzipTypes:        []string{"text/plain", "text/css", "application/json", "application/javascript", "text/xml", "application/xml", "application/xml+rss", "text/javascript"},
			Include:          []string{"conf.d/*.conf"},
			Upstreams:        []UpstreamConfig{},
			Servers:          []ServerConfigBlock{},
		},
	}
}

// AddUpstream 添加上游服务器组
func (config *NginxConfig) AddUpstream(upstream UpstreamConfig) {
	config.Http.Upstreams = append(config.Http.Upstreams, upstream)
}

// AddServer 添加HTTP服务器配置
func (config *NginxConfig) AddServer(server ServerConfigBlock) {
	config.Http.Servers = append(config.Http.Servers, server)
}

// Validate 验证配置
func (config *NginxConfig) Validate() error {
	if len(config.Http.Upstreams) == 0 {
		return fmt.Errorf("至少需要配置一个上游服务器组")
	}

	if len(config.Http.Servers) == 0 {
		return fmt.Errorf("至少需要配置一个HTTP服务器")
	}

	for _, upstream := range config.Http.Upstreams {
		if upstream.Name == "" {
			return fmt.Errorf("上游服务器组名称不能为空")
		}
		if len(upstream.Servers) == 0 {
			return fmt.Errorf("上游服务器组 '%s' 至少需要一个后端服务器", upstream.Name)
		}
		for _, server := range upstream.Servers {
			if server.Host == "" {
				return fmt.Errorf("服务器地址不能为空")
			}
		}
	}

	return nil
}

// GetUpstreamByName 根据名称获取上游服务器组
func (config *NginxConfig) GetUpstreamByName(name string) (*UpstreamConfig, error) {
	for _, upstream := range config.Http.Upstreams {
		if upstream.Name == name {
			return &upstream, nil
		}
	}
	return nil, fmt.Errorf("未找到名为 '%s' 的上游服务器组", name)
}

// RemoveUpstream 根据名称移除上游服务器组
func (config *NginxConfig) RemoveUpstream(name string) error {
	for i, upstream := range config.Http.Upstreams {
		if upstream.Name == name {
			config.Http.Upstreams = append(config.Http.Upstreams[:i], config.Http.Upstreams[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("未找到名为 '%s' 的上游服务器组", name)
}

// String 返回负载均衡方法的字符串表示
func (method LoadBalanceMethod) String() string {
	switch method {
	case RoundRobin:
		return ""
	case LeastConn:
		return "least_conn;"
	case IPHash:
		return "ip_hash;"
	case WeightedRoundRobin:
		return ""
	default:
		return ""
	}
}

// IsWeighted 检查是否为加权算法
func (method LoadBalanceMethod) IsWeighted() bool {
	return method == WeightedRoundRobin
}
