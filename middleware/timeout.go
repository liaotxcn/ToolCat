package middleware

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"weave/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	// DefaultTimeout 默认超时时间
	DefaultTimeout time.Duration
	// PathTimeouts 特定路径的超时时间
	PathTimeouts map[string]time.Duration
	// OnTimeout 超时回调函数
	OnTimeout func(c *gin.Context, timeout time.Duration)
	// TimeoutHandler 超时处理函数
	TimeoutHandler gin.HandlerFunc
}

// DefaultTimeoutConfig 默认超时配置
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		DefaultTimeout: 30 * time.Second,
		PathTimeouts: map[string]time.Duration{
			"/api/v1/llm/chat":   60 * time.Second,  // LLM聊天接口需要更长时间
			"/api/v1/llm/stream": 120 * time.Second, // 流式接口需要更长时间
			"/api/v1/health":     5 * time.Second,   // 健康检查快速响应
			"/api/v1/metrics":    10 * time.Second,  // 监控指标
			"/auth/login":        10 * time.Second,  // 登录接口
			"/auth/register":     15 * time.Second,  // 注册接口
		},
		OnTimeout:      DefaultOnTimeout,
		TimeoutHandler: DefaultTimeoutHandler,
	}
}

// DefaultOnTimeout 默认超时回调函数
func DefaultOnTimeout(c *gin.Context, timeout time.Duration) {
	pkg.Warn("Request timeout",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Duration("timeout", timeout),
		zap.String("client_ip", c.ClientIP()))
}

// DefaultTimeoutHandler 默认超时处理函数
func DefaultTimeoutHandler(c *gin.Context) {
	c.JSON(http.StatusRequestTimeout, gin.H{
		"error": "Request timeout",
		"code":  "REQUEST_TIMEOUT",
	})
	c.Abort()
}

// TimeoutMiddleware 超时控制中间件
func TimeoutMiddleware(config TimeoutConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前路径的超时配置
		timeout := getTimeoutForPath(c.Request.URL.Path, config)

		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换请求的上下文
		c.Request = c.Request.WithContext(ctx)

		// 创建通道用于检测请求完成
		finished := make(chan struct{})

		// 启动goroutine执行请求处理
		go func() {
			defer close(finished)
			c.Next()
		}()

		// 等待请求完成或超时
		select {
		case <-finished:
			// 请求正常完成
			return
		case <-ctx.Done():
			// 请求超时
			if ctx.Err() == context.DeadlineExceeded {
				// 执行超时回调
				if config.OnTimeout != nil {
					config.OnTimeout(c, timeout)
				}

				// 执行超时处理
				if config.TimeoutHandler != nil {
					config.TimeoutHandler(c)
				} else {
					DefaultTimeoutHandler(c)
				}
			}
		}
	}
}

// getTimeoutForPath 获取特定路径的超时时间
func getTimeoutForPath(path string, config TimeoutConfig) time.Duration {
	// 精确匹配
	if timeout, exists := config.PathTimeouts[path]; exists {
		return timeout
	}

	// 前缀匹配
	for pattern, timeout := range config.PathTimeouts {
		if len(pattern) > 0 && pattern[len(pattern)-1] == '/' {
			if len(path) >= len(pattern) && path[:len(pattern)] == pattern {
				return timeout
			}
		}
	}

	// 返回默认超时时间
	return config.DefaultTimeout
}

// HTTPClientWithTimeout 创建带超时的HTTP客户端
func HTTPClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			ResponseHeaderTimeout: timeout,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
		},
	}
}

// TimeoutHTTPClient 带超时控制的HTTP客户端
type TimeoutHTTPClient struct {
	client  *http.Client
	timeout time.Duration
}

// NewTimeoutHTTPClient 创建带超时的HTTP客户端
func NewTimeoutHTTPClient(timeout time.Duration) *TimeoutHTTPClient {
	return &TimeoutHTTPClient{
		client:  HTTPClientWithTimeout(timeout),
		timeout: timeout,
	}
}

// Do 执行HTTP请求
func (t *TimeoutHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// 为请求设置超时
	ctx, cancel := context.WithTimeout(req.Context(), t.timeout)
	defer cancel()

	return t.client.Do(req.WithContext(ctx))
}

// Get 执行GET请求
func (t *TimeoutHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return t.Do(req)
}

// Post 执行POST请求
func (t *TimeoutHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return t.Do(req)
}

// TimeoutWrapper 超时包装器
type TimeoutWrapper struct {
	timeout time.Duration
}

// NewTimeoutWrapper 创建超时包装器
func NewTimeoutWrapper(timeout time.Duration) *TimeoutWrapper {
	return &TimeoutWrapper{
		timeout: timeout,
	}
}

// Wrap 包装函数，添加超时控制
func (t *TimeoutWrapper) Wrap(fn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// 创建带超时的上下文
		timeoutCtx, cancel := context.WithTimeout(ctx, t.timeout)
		defer cancel()

		// 创建通道用于接收函数结果
		result := make(chan error, 1)

		// 启动goroutine执行函数
		go func() {
			result <- fn(timeoutCtx)
		}()

		// 等待函数完成或超时
		select {
		case err := <-result:
			return err
		case <-timeoutCtx.Done():
			if timeoutCtx.Err() == context.DeadlineExceeded {
				pkg.Warn("Function execution timeout",
					zap.Duration("timeout", t.timeout))
				return timeoutCtx.Err()
			}
			return timeoutCtx.Err()
		}
	}
}

// WrapWithResult 包装带返回值的函数，添加超时控制
func WrapWithResult[T any](wrapper *TimeoutWrapper, fn func(ctx context.Context) (T, error)) func(ctx context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		var zero T

		// 创建带超时的上下文
		timeoutCtx, cancel := context.WithTimeout(ctx, wrapper.timeout)
		defer cancel()

		// 创建通道用于接收函数结果
		result := make(chan struct {
			value T
			err   error
		}, 1)

		// 启动goroutine执行函数
		go func() {
			value, err := fn(timeoutCtx)
			result <- struct {
				value T
				err   error
			}{value: value, err: err}
		}()

		// 等待函数完成或超时
		select {
		case res := <-result:
			return res.value, res.err
		case <-timeoutCtx.Done():
			if timeoutCtx.Err() == context.DeadlineExceeded {
				pkg.Warn("Function execution timeout",
					zap.Duration("timeout", wrapper.timeout))
				return zero, timeoutCtx.Err()
			}
			return zero, timeoutCtx.Err()
		}
	}
}
