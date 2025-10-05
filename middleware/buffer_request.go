package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestBufferMiddleware 请求缓冲中间件
// 预先读取请求体到内存，防止慢速客户端问题
func RequestBufferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于非GET、HEAD、OPTIONS请求，需要缓冲请求体
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead && c.Request.Method != http.MethodOptions {
			// 检查Content-Length，防止过大的请求
			contentLength := c.Request.ContentLength
			const maxBodySize = 10 * 1024 * 1024 // 10MB，可根据需求调整

			if contentLength > maxBodySize {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": "Request body too large",
				})
				c.Abort()
				return
			}

			// 读取整个请求体
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to read request body",
				})
				c.Abort()
				return
			}

			// 将请求体重置为已读取的内容，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// 将请求体存储在上下文中，方便控制器获取
			c.Set("requestBody", body)
		}

		// 继续处理请求
		c.Next()
	}
}
