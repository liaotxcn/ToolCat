package models

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAuditLogModel 测试审计日志模型
func TestAuditLogModel(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个测试用户
	user := models.User{
		Username: "testuser",
		Password: "hashedpassword",
		Email:    "test@example.com",
		TenantID: 1,
	}

	if err := pkg.DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer pkg.DB.Delete(&user)

	// 创建审计日志选项
	auditOptions := pkg.AuditLogOptions{
		UserID:       user.ID,
		Username:     user.Username,
		Action:       "test_action",
		ResourceType: "test_resource",
		ResourceID:   "test_id",
		OldValue: map[string]interface{}{
			"field1": "old_value",
		},
		NewValue: map[string]interface{}{
			"field1": "new_value",
		},
		IPAddress: "127.0.0.1",
		UserAgent: "Test Agent",
		TenantID:  user.TenantID,
	}

	// 记录审计日志
	err := pkg.AuditLog(auditOptions)
	assert.NoError(t, err, "Failed to record audit log")

	// 等待异步操作完成
	time.Sleep(100 * time.Millisecond)

	// 验证审计日志是否被正确保存
	var auditLogs []models.AuditLog
	result := pkg.DB.Where("resource_type = ? AND resource_id = ?", "test_resource", "test_id").Find(&auditLogs)
	assert.NoError(t, result.Error)
	assert.Equal(t, 1, len(auditLogs), "Audit log should be created")

	// 验证审计日志内容
	auditLog := auditLogs[0]
	assert.Equal(t, user.ID, auditLog.UserID)
	assert.Equal(t, user.Username, auditLog.Username)
	assert.Equal(t, "test_action", auditLog.Action)
	assert.Equal(t, "test_resource", auditLog.ResourceType)
	assert.Equal(t, "test_id", auditLog.ResourceID)
	assert.Equal(t, user.TenantID, auditLog.TenantID)

	// 验证OldValue和NewValue是否正确
	var oldValue map[string]interface{}
	var newValue map[string]interface{}

	err = json.Unmarshal([]byte(auditLog.OldValue), &oldValue)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(auditLog.NewValue), &newValue)
	assert.NoError(t, err)

	assert.Equal(t, "old_value", oldValue["field1"])
	assert.Equal(t, "new_value", newValue["field1"])

	// 清理测试数据
	pkg.DB.Where("resource_type = ? AND resource_id = ?", "test_resource", "test_id").Delete(&models.AuditLog{})
}

// TestAuditLogFromContext 测试从Gin上下文记录审计日志
func TestAuditLogFromContext(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个测试用户
	user := models.User{
		Username: "testuser2",
		Password: "hashedpassword",
		Email:    "test2@example.com",
		TenantID: 1,
	}

	if err := pkg.DB.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	defer pkg.DB.Delete(&user)

	// 创建一个Gin上下文
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/api/v1/test", strings.NewReader("{}"))
	c.Request.Header.Set("User-Agent", "Test Context Agent")
	c.Set("user_id", user.ID)
	c.Set("username", user.Username)
	c.Set("tenant_id", user.TenantID)

	// 创建审计日志选项
	auditOptions := pkg.AuditLogOptions{
		Action:       "context_test",
		ResourceType: "context_resource",
		ResourceID:   "context_id",
		OldValue:     nil,
		NewValue: map[string]interface{}{
			"test": "value",
		},
	}

	// 从上下文记录审计日志
	err := pkg.AuditLogFromContext(c, auditOptions)
	assert.NoError(t, err, "Failed to record audit log from context")

	// 等待异步操作完成
	time.Sleep(100 * time.Millisecond)

	// 验证审计日志是否被正确保存
	var auditLogs []models.AuditLog
	result := pkg.DB.Where("action = ? AND resource_type = ?", "context_test", "context_resource").Find(&auditLogs)
	assert.NoError(t, result.Error)
	assert.Equal(t, 1, len(auditLogs), "Audit log from context should be created")

	// 验证审计日志内容
	auditLog := auditLogs[0]
	assert.Equal(t, user.ID, auditLog.UserID)
	assert.Equal(t, user.Username, auditLog.Username)
	assert.Equal(t, "context_test", auditLog.Action)
	assert.Equal(t, "context_resource", auditLog.ResourceType)
	assert.Equal(t, "context_id", auditLog.ResourceID)
	assert.Equal(t, user.TenantID, auditLog.TenantID)
	assert.Contains(t, auditLog.UserAgent, "Test Context Agent")

	// 清理测试数据
	pkg.DB.Where("action = ? AND resource_type = ?", "context_test", "context_resource").Delete(&models.AuditLog{})
}

// TestAuditLogMiddleware 测试审计日志中间件
func TestAuditLogMiddleware(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个路由引擎并添加审计日志中间件
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(pkg.AuditLogMiddleware())

	// 创建一个测试API端点
	router.POST("/api/v1/test-resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Test successful"})
	})

	// 创建测试请求
	req, _ := http.NewRequest("POST", "/api/v1/test-resource", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 发送请求
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 等待异步操作完成
	time.Sleep(100 * time.Millisecond)

	// 验证审计日志是否被自动记录
	var auditLogs []models.AuditLog
	result := pkg.DB.Where("action = ? AND resource_type = ?", "post", "test-resource").Find(&auditLogs)
	assert.NoError(t, result.Error)
	assert.Equal(t, 1, len(auditLogs), "Audit log should be created by middleware")

	// 清理测试数据
	pkg.DB.Where("action = ? AND resource_type = ?", "post", "test-resource").Delete(&models.AuditLog{})
}