package middleware

import (
	"net/http"
	"time"

	"toolcat/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 错误处理中间件
type ErrorHandler struct{}

// NewErrorHandler 创建一个新的错误处理中间件
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandlerFunc 错误处理中间件的处理函数
func (eh *ErrorHandler) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		reqStart := time.Now()

		// 处理请求
		c.Next()

		// 如果状态码是4xx或5xx，记录错误日志
		if c.Writer.Status() >= 400 {
			duration := time.Since(reqStart)
			pkg.With(
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status_code", c.Writer.Status()),
				zap.Duration("duration", duration),
				zap.String("remote_addr", c.ClientIP()),
			).Warn("Request failed")
		} else {
			// 记录成功请求的信息（调试级别）
			duration := time.Since(reqStart)
			pkg.Debug("Request processed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status_code", c.Writer.Status()),
				zap.Duration("duration", duration),
			)
		}
	}
}



// responseWriterWrapper 用于包装http.ResponseWriter，捕获状态码
// 这个结构体用于内部跟踪响应状态码，以便在中间件中记录日志
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 重写WriteHeader方法，捕获状态码
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
