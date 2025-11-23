package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"weave/config"
	"weave/pkg"
	"weave/pkg/metrics"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecordHTTPRequest 可替换的指标记录函数，默认为metrics.RecordHTTPRequest
// 公开此变量以便在测试中可以替换
var RecordHTTPRequest = metrics.RecordHTTPRequest

// LoadBalancerController 负载均衡管理控制器
type LoadBalancerController struct{}

// LoadBalancerStatus 负载均衡状态信息
type LoadBalancerStatus struct {
	InstanceID   string            `json:"instance_id"`
	Status       string            `json:"status"`
	LastCheck    time.Time         `json:"last_check"`
	ResponseTime int64             `json:"response_time"`
	Weight       int               `json:"weight"`
	Requests     int64             `json:"requests"`
	SuccessRate  float64           `json:"success_rate"`
	Metadata     map[string]string `json:"metadata"`
}

// LoadBalancerStats 负载均衡统计信息
type LoadBalancerStats struct {
	TotalRequests     int64                `json:"total_requests"`
	ActiveConnections int                  `json:"active_connections"`
	AverageResponse   float64              `json:"average_response_time"`
	Instances         []LoadBalancerStatus `json:"instances"`
	LastUpdated       time.Time            `json:"last_updated"`
}

// GetLoadBalancerStatus 获取负载均衡状态
func (lbc *LoadBalancerController) GetLoadBalancerStatus(c *gin.Context) {
	startTime := time.Now()

	// 获取当前实例信息
	currentInstance := LoadBalancerStatus{
		InstanceID:   config.Config.Server.InstanceID,
		Status:       "healthy",
		LastCheck:    time.Now(),
		ResponseTime: 10, // 模拟响应时间（毫秒）
		Weight:       1,
		Requests:     getMockRequestCount(),
		SuccessRate:  0.99, // 模拟成功率
		Metadata: map[string]string{
			"version":    "1.0.0",
			"started_at": time.Now().Add(-time.Hour).Format(time.RFC3339),
		},
	}

	// 模拟其他实例信息（在实际环境中，这些信息应该从服务发现或配置中心获取）
	instances := []LoadBalancerStatus{currentInstance}

	// 添加其他模拟实例
	for i := 1; i <= 2; i++ {
		instance := LoadBalancerStatus{
			InstanceID:   fmt.Sprintf("weave-%d", i),
			Status:       "healthy",
			LastCheck:    time.Now().Add(-time.Duration(i) * time.Second),
			ResponseTime: 10 + int64(i*2),
			Weight:       1,
			Requests:     getMockRequestCount() + int64(i*100),
			SuccessRate:  0.98 - float64(i)*0.01,
			Metadata: map[string]string{
				"version":    "1.0.0",
				"started_at": time.Now().Add(-time.Hour - time.Duration(i)*time.Minute).Format(time.RFC3339),
			},
		}
		instances = append(instances, instance)
	}

	// 按实例ID排序
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].InstanceID < instances[j].InstanceID
	})

	stats := LoadBalancerStats{
		TotalRequests:     getTotalRequests(instances),
		ActiveConnections: getActiveConnections(),
		AverageResponse:   getAverageResponseTime(instances),
		Instances:         instances,
		LastUpdated:       time.Now(),
	}

	// 记录请求指标
	duration := time.Since(startTime).Seconds()
	RecordHTTPRequest("GET", "/loadbalancer/status", "200", duration)

	pkg.Info("Load balancer status requested",
		zap.String("instance_id", config.Config.Server.InstanceID),
		zap.Float64("duration", duration))

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   stats,
	})
}

// GetInstanceHealth 获取特定实例健康状态
func (lbc *LoadBalancerController) GetInstanceHealth(c *gin.Context) {
	instanceID := c.Param("instanceId")
	startTime := time.Now()

	// 模拟实例健康检查
	health := gin.H{
		"instance_id":   instanceID,
		"status":        "healthy",
		"last_check":    time.Now(),
		"response_time": 15, // 毫秒
		"uptime":        "1h30m45s",
		"memory_usage":  "45%",
		"cpu_usage":     "12%",
		"requests":      1250,
		"errors":        3,
		"success_rate":  0.9976,
	}

	// 记录请求指标
	duration := time.Since(startTime).Seconds()
	RecordHTTPRequest("GET", "/loadbalancer/instance/"+instanceID+"/health", "200", duration)

	pkg.Info("Instance health check requested",
		zap.String("instance_id", instanceID),
		zap.String("requester", config.Config.Server.InstanceID),
		zap.Float64("duration", duration))

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   health,
	})
}

// UpdateInstanceWeight 更新实例权重
func (lbc *LoadBalancerController) UpdateInstanceWeight(c *gin.Context) {
	instanceID := c.Param("instanceId")
	startTime := time.Now()

	var request struct {
		Weight int `json:"weight" binding:"required,min=1,max=10"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		pkg.Error("Invalid weight update request",
			zap.String("instance_id", instanceID),
			zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid weight value. Must be between 1 and 10.",
		})
		return
	}

	// 在实际环境中，这里应该更新负载均衡器的配置
	// 目前只是模拟成功响应

	duration := time.Since(startTime).Seconds()
	RecordHTTPRequest("PUT", "/loadbalancer/instance/"+instanceID+"/weight", "200", duration)

	pkg.Info("Instance weight updated",
		zap.String("instance_id", instanceID),
		zap.Int("new_weight", request.Weight),
		zap.String("updated_by", config.Config.Server.InstanceID),
		zap.Float64("duration", duration))

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Weight updated successfully",
		"data": gin.H{
			"instance_id": instanceID,
			"old_weight":  1,
			"new_weight":  request.Weight,
			"updated_at":  time.Now(),
		},
	})
}

// DrainInstance 排干实例（停止接收新请求）
func (lbc *LoadBalancerController) DrainInstance(c *gin.Context) {
	instanceID := c.Param("instanceId")
	startTime := time.Now()

	// 在实际环境中，这里应该将实例标记为排干状态
	// 目前只是模拟成功响应

	duration := time.Since(startTime).Seconds()
	RecordHTTPRequest("POST", "/loadbalancer/instance/"+instanceID+"/drain", "200", duration)

	pkg.Info("Instance drained",
		zap.String("instance_id", instanceID),
		zap.String("drained_by", config.Config.Server.InstanceID),
		zap.Float64("duration", duration))

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Instance is being drained",
		"data": gin.H{
			"instance_id": instanceID,
			"status":      "draining",
			"drained_at":  time.Now(),
		},
	})
}

// EnableInstance 启用实例
func (lbc *LoadBalancerController) EnableInstance(c *gin.Context) {
	instanceID := c.Param("instanceId")
	startTime := time.Now()

	// 在实际环境中，这里应该将实例重新启用
	// 目前只是模拟成功响应

	duration := time.Since(startTime).Seconds()
	RecordHTTPRequest("POST", "/loadbalancer/instance/"+instanceID+"/enable", "200", duration)

	pkg.Info("Instance enabled",
		zap.String("instance_id", instanceID),
		zap.String("enabled_by", config.Config.Server.InstanceID),
		zap.Float64("duration", duration))

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Instance enabled successfully",
		"data": gin.H{
			"instance_id": instanceID,
			"status":      "enabled",
			"enabled_at":  time.Now(),
		},
	})
}

// 辅助函数（模拟数据）
func getMockRequestCount() int64 {
	return int64(time.Now().Unix() % 10000)
}

func getTotalRequests(instances []LoadBalancerStatus) int64 {
	var total int64
	for _, instance := range instances {
		total += instance.Requests
	}
	return total
}

func getActiveConnections() int {
	// 模拟活跃连接数
	return 25 + int(time.Now().Unix()%50)
}

func getAverageResponseTime(instances []LoadBalancerStatus) float64 {
	if len(instances) == 0 {
		return 0
	}

	var total int64
	for _, instance := range instances {
		total += instance.ResponseTime
	}
	return float64(total) / float64(len(instances))
}
