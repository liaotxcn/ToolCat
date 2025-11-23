package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"weave/config"
	"weave/controllers"
	"weave/pkg/metrics"
)

var (
	// 模拟配置
	testInstanceID = "test-instance-123"
	// 记录被调用的请求信息
	recordedRequests []HTTPRequestInfo
)

// HTTPRequestInfo 记录HTTP请求的信息
type HTTPRequestInfo struct {
	Method   string
	Endpoint string
	Status   string
	Duration float64
}

// 设置模拟配置
func setupMockConfig() {
	config.Config.Server.InstanceID = testInstanceID
}

// mockRecordHTTPRequest 模拟记录HTTP请求的函数
func mockRecordHTTPRequest(method, endpoint, status string, duration float64) {
	recordedRequests = append(recordedRequests, HTTPRequestInfo{
		Method:   method,
		Endpoint: endpoint,
		Status:   status,
		Duration: duration,
	})
}

// 测试前设置
func setupTest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMockConfig()
	// 清空记录的请求
	recordedRequests = []HTTPRequestInfo{}
	// 替换controllers包中的指标记录函数为模拟函数
	controllers.RecordHTTPRequest = mockRecordHTTPRequest
}

// 测试后清理
func teardownTest() {
	// 重置为原始实现
	controllers.RecordHTTPRequest = metrics.RecordHTTPRequest
}

// TestGetLoadBalancerStatus 测试获取负载均衡状态
func TestGetLoadBalancerStatus(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	lbc := controllers.LoadBalancerController{}
	r := gin.New()
	r.GET("/loadbalancer/status", func(c *gin.Context) { lbc.GetLoadBalancerStatus(c) })

	// 发送GET请求
	req, _ := http.NewRequest(http.MethodGet, "/loadbalancer/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 验证响应结构
	var response struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	// 验证状态和数据存在
	if response.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", response.Status)
	}
	if response.Data == nil {
		t.Fatalf("expected non-nil data in response")
	}

	// 验证数据字段
	if _, ok := response.Data["total_requests"].(float64); !ok {
		t.Fatalf("expected total_requests in data")
	}
	if _, ok := response.Data["active_connections"].(float64); !ok {
		t.Fatalf("expected active_connections in data")
	}
	if _, ok := response.Data["average_response_time"].(float64); !ok {
		t.Fatalf("expected average_response_time in data")
	}
	if instances, ok := response.Data["instances"].([]interface{}); !ok || len(instances) == 0 {
		t.Fatalf("expected non-empty instances array in data")
	}
}

// TestGetInstanceHealth 测试获取特定实例健康状态
func TestGetInstanceHealth(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	lbc := controllers.LoadBalancerController{}
	r := gin.New()
	r.GET("/loadbalancer/instance/:instanceId/health", func(c *gin.Context) { lbc.GetInstanceHealth(c) })

	// 发送GET请求
	instanceID := "test-instance-1"
	req, _ := http.NewRequest(http.MethodGet, "/loadbalancer/instance/"+instanceID+"/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 验证响应结构
	var response struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	// 验证状态和数据
	if response.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", response.Status)
	}
	if response.Data == nil {
		t.Fatalf("expected non-nil data in response")
	}

	// 验证实例ID和健康状态
	if dataInstanceID, ok := response.Data["instance_id"].(string); !ok || dataInstanceID != instanceID {
		t.Fatalf("expected instance_id=%s, got %v", instanceID, response.Data["instance_id"])
	}
	if status, ok := response.Data["status"].(string); !ok || status != "healthy" {
		t.Fatalf("expected status=healthy, got %v", response.Data["status"])
	}
}

// TestUpdateInstanceWeight 测试更新实例权重成功
func TestUpdateInstanceWeight_Success(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	lbc := controllers.LoadBalancerController{}
	r := gin.New()
	r.PUT("/loadbalancer/instance/:instanceId/weight", func(c *gin.Context) { lbc.UpdateInstanceWeight(c) })

	// 准备请求体
	instanceID := "test-instance-1"
	newWeight := 5
	payload := `{"weight":5}`
	req, _ := http.NewRequest(http.MethodPut, "/loadbalancer/instance/"+instanceID+"/weight", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 验证响应结构
	var response struct {
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	// 验证状态和消息
	if response.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", response.Status)
	}
	if response.Message != "Weight updated successfully" {
		t.Fatalf("unexpected message: %q", response.Message)
	}

	// 验证数据
	if response.Data == nil {
		t.Fatalf("expected non-nil data in response")
	}
	if dataInstanceID, ok := response.Data["instance_id"].(string); !ok || dataInstanceID != instanceID {
		t.Fatalf("expected instance_id=%s, got %v", instanceID, response.Data["instance_id"])
	}
	if weight, ok := response.Data["new_weight"].(float64); !ok || int(weight) != newWeight {
		t.Fatalf("expected new_weight=%d, got %v", newWeight, response.Data["new_weight"])
	}
}

// TestUpdateInstanceWeight_InvalidInput 测试更新实例权重时的无效输入
func TestUpdateInstanceWeight_InvalidInput(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	lbc := controllers.LoadBalancerController{}
	r := gin.New()
	r.PUT("/loadbalancer/instance/:instanceId/weight", func(c *gin.Context) { lbc.UpdateInstanceWeight(c) })

	// 准备无效请求体（超出范围的权重值）
	instanceID := "test-instance-1"
	invalidPayload := `{"weight":15}` // 超出1-10的范围
	req, _ := http.NewRequest(http.MethodPut, "/loadbalancer/instance/"+instanceID+"/weight", strings.NewReader(invalidPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	// 验证响应结构
	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	// 验证状态和错误消息
	if response.Status != "error" {
		t.Fatalf("expected status=error, got %q", response.Status)
	}
	if response.Message != "Invalid weight value. Must be between 1 and 10." {
		t.Fatalf("unexpected error message: %q", response.Message)
	}
}

// TestDrainInstance 测试排干实例
func TestDrainInstance(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	lbc := controllers.LoadBalancerController{}
	r := gin.New()
	r.POST("/loadbalancer/instance/:instanceId/drain", func(c *gin.Context) { lbc.DrainInstance(c) })

	// 发送POST请求
	instanceID := "test-instance-1"
	req, _ := http.NewRequest(http.MethodPost, "/loadbalancer/instance/"+instanceID+"/drain", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 验证响应结构
	var response struct {
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	// 验证状态和消息
	if response.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", response.Status)
	}
	if response.Message != "Instance is being drained" {
		t.Fatalf("unexpected message: %q", response.Message)
	}

	// 验证数据
	if response.Data == nil {
		t.Fatalf("expected non-nil data in response")
	}
	if dataInstanceID, ok := response.Data["instance_id"].(string); !ok || dataInstanceID != instanceID {
		t.Fatalf("expected instance_id=%s, got %v", instanceID, response.Data["instance_id"])
	}
	if status, ok := response.Data["status"].(string); !ok || status != "draining" {
		t.Fatalf("expected status=draining, got %v", response.Data["status"])
	}
}

// TestEnableInstance 测试启用实例
func TestEnableInstance(t *testing.T) {
	setupTest(t)
	defer teardownTest()

	lbc := controllers.LoadBalancerController{}
	r := gin.New()
	r.POST("/loadbalancer/instance/:instanceId/enable", func(c *gin.Context) { lbc.EnableInstance(c) })

	// 发送POST请求
	instanceID := "test-instance-1"
	req, _ := http.NewRequest(http.MethodPost, "/loadbalancer/instance/"+instanceID+"/enable", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 验证响应结构
	var response struct {
		Status  string                 `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	// 验证状态和消息
	if response.Status != "ok" {
		t.Fatalf("expected status=ok, got %q", response.Status)
	}
	if response.Message != "Instance enabled successfully" {
		t.Fatalf("unexpected message: %q", response.Message)
	}

	// 验证数据
	if response.Data == nil {
		t.Fatalf("expected non-nil data in response")
	}
	if dataInstanceID, ok := response.Data["instance_id"].(string); !ok || dataInstanceID != instanceID {
		t.Fatalf("expected instance_id=%s, got %v", instanceID, response.Data["instance_id"])
	}
	if status, ok := response.Data["status"].(string); !ok || status != "enabled" {
		t.Fatalf("expected status=enabled, got %v", response.Data["status"])
	}
}
