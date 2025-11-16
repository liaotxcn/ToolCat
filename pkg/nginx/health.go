package nginx

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthStatus 健康状态
type HealthStatus string

const (
	// StatusHealthy 健康
	StatusHealthy HealthStatus = "healthy"
	// StatusUnhealthy 不健康
	StatusUnhealthy HealthStatus = "unhealthy"
	// StatusUnknown 未知
	StatusUnknown HealthStatus = "unknown"
)

// ServerHealth 服务器健康状态
type ServerHealth struct {
	Server       ServerConfig  `json:"server"`
	Status       HealthStatus  `json:"status"`
	LastCheck    time.Time     `json:"last_check"`
	FailCount    int           `json:"fail_count"`
	SuccessCount int           `json:"success_count"`
	ResponseTime time.Duration `json:"response_time"`
	Error        string        `json:"error,omitempty"`
}

// HealthChecker 健康检查器
type HealthChecker struct {
	upstreamName   string
	config         HealthCheckConfig
	servers        []ServerConfig
	healthStatus   map[string]*ServerHealth
	client         *http.Client
	stopCh         chan struct{}
	mu             sync.RWMutex
	onStatusChange func(server string, status HealthStatus)
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(upstreamName string, config HealthCheckConfig, servers []ServerConfig) *HealthChecker {
	return &HealthChecker{
		upstreamName: upstreamName,
		config:       config,
		servers:      servers,
		healthStatus: make(map[string]*ServerHealth),
		client: &http.Client{
			Timeout: config.Timeout,
		},
		stopCh: make(chan struct{}),
	}
}

// SetStatusChangeCallback 设置状态变化回调
func (hc *HealthChecker) SetStatusChangeCallback(callback func(server string, status HealthStatus)) {
	hc.onStatusChange = callback
}

// Start 开始健康检查
func (hc *HealthChecker) Start() {
	if !hc.config.Enabled {
		return
	}

	// 初始化健康状态
	hc.initializeHealthStatus()

	// 启动检查协程
	go hc.runHealthCheck()
}

// Stop 停止健康检查
func (hc *HealthChecker) Stop() {
	close(hc.stopCh)
}

// GetHealthStatus 获取健康状态
func (hc *HealthChecker) GetHealthStatus() map[string]*ServerHealth {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	result := make(map[string]*ServerHealth)
	for k, v := range hc.healthStatus {
		// 复制对象避免并发问题
		health := *v
		result[k] = &health
	}
	return result
}

// GetHealthyServers 获取健康的服务器列表
func (hc *HealthChecker) GetHealthyServers() []ServerConfig {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	var healthyServers []ServerConfig
	for _, server := range hc.servers {
		if health, exists := hc.healthStatus[server.Host]; exists {
			if health.Status == StatusHealthy {
				healthyServers = append(healthyServers, server)
			}
		}
	}
	return healthyServers
}

// initializeHealthStatus 初始化健康状态
func (hc *HealthChecker) initializeHealthStatus() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	for _, server := range hc.servers {
		hc.healthStatus[server.Host] = &ServerHealth{
			Server:       server,
			Status:       StatusUnknown,
			LastCheck:    time.Now(),
			FailCount:    0,
			SuccessCount: 0,
		}
	}
}

// runHealthCheck 运行健康检查
func (hc *HealthChecker) runHealthCheck() {
	ticker := time.NewTicker(hc.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-hc.stopCh:
			return
		case <-ticker.C:
			hc.checkAllServers()
		}
	}
}

// checkAllServers 检查所有服务器
func (hc *HealthChecker) checkAllServers() {
	var wg sync.WaitGroup

	for _, server := range hc.servers {
		wg.Add(1)
		go func(s ServerConfig) {
			defer wg.Done()
			hc.checkServer(s)
		}(server)
	}

	wg.Wait()
}

// checkServer 检查单个服务器
func (hc *HealthChecker) checkServer(server ServerConfig) {
	start := time.Now()

	// 构建健康检查URL
	url := fmt.Sprintf("http://%s%s", server.Host, hc.config.Path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		hc.updateServerHealth(server.Host, StatusUnhealthy, 0, fmt.Sprintf("创建请求失败: %v", err))
		return
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Weave-HealthChecker/1.0")
	req.Header.Set("X-Health-Check", "true")

	resp, err := hc.client.Do(req)
	responseTime := time.Since(start)

	if err != nil {
		hc.updateServerHealth(server.Host, StatusUnhealthy, responseTime, fmt.Sprintf("请求失败: %v", err))
		return
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		hc.updateServerHealth(server.Host, StatusHealthy, responseTime, "")
	} else {
		hc.updateServerHealth(server.Host, StatusUnhealthy, responseTime, fmt.Sprintf("HTTP状态码: %d", resp.StatusCode))
	}
}

// updateServerHealth 更新服务器健康状态
func (hc *HealthChecker) updateServerHealth(serverHost string, status HealthStatus, responseTime time.Duration, errorMsg string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	health, exists := hc.healthStatus[serverHost]
	if !exists {
		return
	}

	oldStatus := health.Status
	health.LastCheck = time.Now()
	health.ResponseTime = responseTime
	health.Error = errorMsg

	switch status {
	case StatusHealthy:
		health.SuccessCount++
		health.FailCount = 0

		// 检查是否达到成功阈值
		if health.SuccessCount >= hc.config.Rise {
			health.Status = StatusHealthy
		}

	case StatusUnhealthy:
		health.FailCount++
		health.SuccessCount = 0

		// 检查是否达到失败阈值
		if health.FailCount >= hc.config.Fall {
			health.Status = StatusUnhealthy
		}
	}

	// 状态变化时触发回调
	if oldStatus != health.Status && hc.onStatusChange != nil {
		hc.onStatusChange(serverHost, health.Status)
	}
}

// HealthStats 健康统计信息
type HealthStats struct {
	TotalServers        int           `json:"total_servers"`
	HealthyServers      int           `json:"healthy_servers"`
	UnhealthyServers    int           `json:"unhealthy_servers"`
	UnknownServers      int           `json:"unknown_servers"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	LastCheck           time.Time     `json:"last_check"`
}

// GetHealthStats 获取健康统计信息
func (hc *HealthChecker) GetHealthStats() HealthStats {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	stats := HealthStats{
		TotalServers: len(hc.servers),
		LastCheck:    time.Now(),
	}

	var totalResponseTime time.Duration
	var responseTimeCount int

	for _, health := range hc.healthStatus {
		switch health.Status {
		case StatusHealthy:
			stats.HealthyServers++
		case StatusUnhealthy:
			stats.UnhealthyServers++
		case StatusUnknown:
			stats.UnknownServers++
		}

		if health.ResponseTime > 0 {
			totalResponseTime += health.ResponseTime
			responseTimeCount++
		}

		if health.LastCheck.After(stats.LastCheck) {
			stats.LastCheck = health.LastCheck
		}
	}

	if responseTimeCount > 0 {
		stats.AverageResponseTime = totalResponseTime / time.Duration(responseTimeCount)
	}

	return stats
}

// MultiHealthChecker 多上游健康检查器
type MultiHealthChecker struct {
	checkers map[string]*HealthChecker
	mu       sync.RWMutex
}

// NewMultiHealthChecker 创建多上游健康检查器
func NewMultiHealthChecker() *MultiHealthChecker {
	return &MultiHealthChecker{
		checkers: make(map[string]*HealthChecker),
	}
}

// AddChecker 添加健康检查器
func (mhc *MultiHealthChecker) AddChecker(upstreamName string, checker *HealthChecker) {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()

	mhc.checkers[upstreamName] = checker
}

// StartAll 启动所有健康检查器
func (mhc *MultiHealthChecker) StartAll() {
	mhc.mu.RLock()
	defer mhc.mu.RUnlock()

	for _, checker := range mhc.checkers {
		checker.Start()
	}
}

// StopAll 停止所有健康检查器
func (mhc *MultiHealthChecker) StopAll() {
	mhc.mu.RLock()
	defer mhc.mu.RUnlock()

	for _, checker := range mhc.checkers {
		checker.Stop()
	}
}

// GetStats 获取所有上游的统计信息
func (mhc *MultiHealthChecker) GetStats() map[string]HealthStats {
	mhc.mu.RLock()
	defer mhc.mu.RUnlock()

	stats := make(map[string]HealthStats)
	for upstreamName, checker := range mhc.checkers {
		stats[upstreamName] = checker.GetHealthStats()
	}
	return stats
}

// GetAllHealthyServers 获取所有上游的健康服务器
func (mhc *MultiHealthChecker) GetAllHealthyServers() map[string][]ServerConfig {
	mhc.mu.RLock()
	defer mhc.mu.RUnlock()

	result := make(map[string][]ServerConfig)
	for upstreamName, checker := range mhc.checkers {
		result[upstreamName] = checker.GetHealthyServers()
	}
	return result
}
