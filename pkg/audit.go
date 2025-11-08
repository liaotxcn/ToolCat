package pkg

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"weave/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuditLogger 审计日志记录器
type AuditLogger struct{}

// NewAuditLogger 创建新的审计日志记录器
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{}
}

// AuditLogOptions 审计日志选项
type AuditLogOptions struct {
	UserID       uint
	Username     string
	Action       string
	ResourceType string
	ResourceID   string
	OldValue     interface{}
	NewValue     interface{}
	IPAddress    string
	UserAgent    string
	TenantID     uint
}

// Log 记录审计日志
func (al *AuditLogger) Log(options AuditLogOptions) error {
	// 转换OldValue和NewValue为JSON字符串
	oldValueStr := ""
	if options.OldValue != nil {
		oldValue, err := json.Marshal(options.OldValue)
		if err != nil {
			return fmt.Errorf("failed to marshal old value: %v", err)
		}
		oldValueStr = string(oldValue)
	}

	newValueStr := ""
	if options.NewValue != nil {
		newValue, err := json.Marshal(options.NewValue)
		if err != nil {
			return fmt.Errorf("failed to marshal new value: %v", err)
		}
		newValueStr = string(newValue)
	}

	// 创建审计日志记录
	auditLog := models.AuditLog{
		UserID:       options.UserID,
		Username:     options.Username,
		Action:       options.Action,
		ResourceType: options.ResourceType,
		ResourceID:   options.ResourceID,
		OldValue:     oldValueStr,
		NewValue:     newValueStr,
		IPAddress:    options.IPAddress,
		UserAgent:    options.UserAgent,
		TenantID:     options.TenantID,
		CreatedAt:    time.Now(),
	}

	// 保存到数据库（异步保存，不阻塞主流程）
	go func() {
		if err := DB.Create(&auditLog).Error; err != nil {
			Error("Failed to save audit log",
				zap.Error(err),
				zap.String("action", options.Action),
				zap.String("resource_type", options.ResourceType),
			)
		}
	}()

	return nil
}

// FromContext 从Gin上下文中提取信息并记录审计日志
func (al *AuditLogger) FromContext(c *gin.Context, options AuditLogOptions) error {
	// 从上下文中获取IP地址和用户代理
	options.IPAddress = c.ClientIP()
	options.UserAgent = c.Request.UserAgent()

	// 从上下文中尝试获取用户信息（如果存在）
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			options.UserID = id
		}
	}

	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			options.Username = name
		}
	}

	if tenantID, exists := c.Get("tenant_id"); exists {
		if id, ok := tenantID.(uint); ok {
			options.TenantID = id
		}
	}

	return al.Log(options)
}

// 审计日志中间件
func AuditLogMiddleware() gin.HandlerFunc {
	auditLogger := NewAuditLogger()

	return func(c *gin.Context) {
		// 跳过不需要审计的路径
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/health") || strings.HasPrefix(path, "/metrics") {
			c.Next()
			return
		}

		// 记录请求开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 对于写操作（POST/PUT/DELETE）进行审计日志记录
		method := c.Request.Method
		if method == "POST" || method == "PUT" || method == "DELETE" {
			// 异步记录审计日志，避免影响响应时间
			go func() {
				action := strings.ToLower(method)
				resourceType := extractResourceType(path)
				resourceID := extractResourceID(path)

				auditLogger.FromContext(c, AuditLogOptions{
					Action:       action,
					ResourceType: resourceType,
					ResourceID:   resourceID,
				})
			}()
		}

		// 记录处理时间（调试用）
		Debug("Audit middleware processed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

// 从路径中提取资源类型
func extractResourceType(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == "api" {
		// Skip version segment like v1, v2 and return actual resource
		if len(parts) >= 3 && strings.HasPrefix(parts[1], "v") {
			return parts[2]
		}
		return parts[1]
	}
	if len(parts) >= 3 && parts[0] == "plugins" {
		return "plugin_" + parts[1]
	}
	return "unknown"
}

// 从路径中提取资源ID
func extractResourceID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		// 尝试提取最后一个路径段作为ID
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] != "" && !strings.Contains(parts[i], "=") {
				return parts[i]
			}
		}
	}
	return ""
}

// 便捷函数
var globalAuditLogger = NewAuditLogger()

// AuditLog 全局审计日志记录函数
func AuditLog(options AuditLogOptions) error {
	return globalAuditLogger.Log(options)
}

// AuditLogFromContext 从上下文记录审计日志
func AuditLogFromContext(c *gin.Context, options AuditLogOptions) error {
	return globalAuditLogger.FromContext(c, options)
}
