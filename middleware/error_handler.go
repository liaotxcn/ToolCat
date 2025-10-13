package middleware

import (
	"errors"
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
		requestID := c.GetString("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Set("X-Request-ID", requestID)
		}

		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误作为主要错误
			err := c.Errors.Last().Err
			var appErr *pkg.AppError
			var statusCode int

			// 检查是否为AppError类型
			if errors.As(err, &appErr) {
				// 设置请求ID和路径
				appErr.WithRequestID(requestID).WithPath(c.Request.URL.Path)
				// 获取对应的HTTP状态码
				statusCode = pkg.GetHTTPStatus(appErr)
				// 以统一格式返回错误
				c.JSON(statusCode, appErr)
			} else {
				// 对于非AppError类型的错误，创建一个内部错误
				appErr = pkg.NewInternalError("Internal server error", err)
				appErr.WithRequestID(requestID).WithPath(c.Request.URL.Path)
				statusCode = http.StatusInternalServerError
				c.JSON(statusCode, appErr)
			}

			// 记录错误日志
			duration := time.Since(reqStart)
			logFields := []zap.Field{
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("request_id", requestID),
				zap.Int("status_code", statusCode),
				zap.String("error_code", string(appErr.Code)),
				zap.String("error_message", appErr.Message),
				zap.Duration("duration", duration),
				zap.String("remote_addr", c.ClientIP()),
			}

			// 根据错误类型设置不同的日志级别
			if statusCode >= 500 {
				pkg.With(logFields...).Error("Request failed with server error")
			} else {
				pkg.With(logFields...).Warn("Request failed with client error")
			}

			// 确保响应已写入
			c.Abort()
			return
		}

		// 记录成功请求的信息
		duration := time.Since(reqStart)
		if c.Writer.Status() >= 400 {
			// 没有捕获到错误但状态码是4xx，记录警告日志
			pkg.With(
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("request_id", requestID),
				zap.Int("status_code", c.Writer.Status()),
				zap.Duration("duration", duration),
				zap.String("remote_addr", c.ClientIP()),
			).Warn("Request completed with non-success status")
		} else {
			// 记录成功请求的信息（调试级别）
			pkg.Debug("Request processed successfully",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("request_id", requestID),
				zap.Int("status_code", c.Writer.Status()),
				zap.Duration("duration", duration),
			)
		}
	}
}

// generateRequestID 生成一个简单的请求ID
func generateRequestID() string {
	// 在实际项目中，应该使用更安全的方式生成请求ID
	// 这里为了简化，使用时间戳和随机数的组合
	return time.Now().Format("20060102150405") + "-" + pkg.RandomString(8)
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
