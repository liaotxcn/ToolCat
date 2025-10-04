package middleware

import (
	"net/http"
	"strings"
	"toolcat/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// 验证token有效性
		tokenString := parts[1]
		userID, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 将用户ID存储在上下文
		c.Set("userID", userID)

		// 继续处理请求
		c.Next()
	}
}

// LogMiddleware 日志中间件
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前的日志记录
		// 可以记录请求路径、方法、IP等信息

		// 继续处理请求
		c.Next()

		// 请求后的日志记录
		// 可以记录响应状态码、处理时间等信息
	}
}
