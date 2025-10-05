package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"toolcat/config"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CSRFMiddleware 跨站请求伪造防护中间件
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果CSRF防护未启用，直接跳过
		if !config.Config.CSRF.Enabled {
			c.Next()
			return
		}

		// 对于GET、HEAD、OPTIONS、TRACE请求，不做CSRF验证
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" || c.Request.Method == "TRACE" {
			// 确保CSRF令牌已设置
			ensureCSRFToken(c)
			c.Next()
			return
		}

		// 验证CSRF令牌
		if !validateCSRFToken(c) {
			pkg.Error("CSRF token validation failed", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token validation failed"})
			return
		}

		c.Next()
	}
}

// ensureCSRFToken 确保CSRF令牌已设置
func ensureCSRFToken(c *gin.Context) {
	// 检查Cookie中是否已有CSRF令牌
	token, err := c.Cookie(config.Config.CSRF.CookieName)
	if err != nil || token == "" {
		// 生成新的CSRF令牌
		token = generateCSRFToken(config.Config.CSRF.TokenLength)
		// 设置Cookie
		// 注意：当前Go版本不支持SameSite参数，在生产环境中应使用支持SameSite的较新版本
		c.SetCookie(
			config.Config.CSRF.CookieName,
			token,
			config.Config.CSRF.CookieMaxAge,
			config.Config.CSRF.CookiePath,
			config.Config.CSRF.CookieDomain,
			config.Config.CSRF.CookieSecure,
			config.Config.CSRF.CookieHttpOnly,
		)
	}

	// 将CSRF令牌添加到响应头中，以便前端可以获取
	c.Header(config.Config.CSRF.HeaderName, token)
}

// generateCSRFToken 生成随机的CSRF令牌
func generateCSRFToken(length int) string {
	token := make([]byte, length)
	if _, err := rand.Read(token); err != nil {
		pkg.Error("Failed to generate CSRF token", zap.Error(err))
		// 如果生成失败，返回一个备用令牌（不推荐，但作为最后的保障）
		return "fallback-csrf-token"
	}
	return hex.EncodeToString(token)
}

// validateCSRFToken 验证CSRF令牌
func validateCSRFToken(c *gin.Context) bool {
	// 从Cookie中获取CSRF令牌
	cookieToken, err := c.Cookie(config.Config.CSRF.CookieName)
	if err != nil || cookieToken == "" {
		return false
	}

	// 从请求头中获取CSRF令牌
	headerToken := c.GetHeader(config.Config.CSRF.HeaderName)
	if headerToken == "" {
		// 尝试从表单中获取CSRF令牌
		headerToken = c.PostForm(config.Config.CSRF.HeaderName)
	}

	// 验证令牌是否匹配
	return cookieToken == headerToken
}
