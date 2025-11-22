package config

import (
	"net/http"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	// 最大重试次数
	MaxRetries int
	// 初始延迟时间
	InitialDelay time.Duration
	// 最大延迟时间
	MaxDelay time.Duration
	// 延迟倍数
	Multiplier float64
	// 随机化因子
	RandomizationFactor float64
	// 重试条件判断函数
	RetryableFunc func(error) bool
	// 重试回调函数
	OnRetry func(attempt int, err error)
}

// TimeoutConfig 超时配置
type TimeoutConfig struct {
	// 默认超时时间
	DefaultTimeout time.Duration
	// 路径特定超时时间
	PathTimeouts map[string]time.Duration
	// 超时回调函数
	OnTimeout func(*http.Request, time.Duration)
	// 超时处理函数
	TimeoutHandler func(*http.Request, time.Duration) error
}

// RetryAndTimeoutConfig 重试和超时配置
type RetryAndTimeoutConfig struct {
	// 全局重试配置
	GlobalRetry RetryConfig
	// 全局超时配置
	GlobalTimeout TimeoutConfig
	// 服务特定配置
	Services map[string]ServiceConfig
}

// ServiceConfig 服务特定配置
type ServiceConfig struct {
	// 重试配置
	Retry *RetryConfig
	// 超时配置
	Timeout *TimeoutConfig
	// 是否启用重试
	EnableRetry bool
	// 是否启用超时
	EnableTimeout bool
}

// DefaultRetryAndTimeoutConfig 默认重试和超时配置
func DefaultRetryAndTimeoutConfig() RetryAndTimeoutConfig {
	return RetryAndTimeoutConfig{
		GlobalRetry: RetryConfig{
			MaxRetries:          3,
			InitialDelay:        100 * time.Millisecond,
			MaxDelay:            5 * time.Second,
			Multiplier:          2.0,
			RandomizationFactor: 0.1,
			RetryableFunc:       defaultRetryableFunc,
			OnRetry:             defaultOnRetry,
		},
		GlobalTimeout: TimeoutConfig{
			DefaultTimeout: 30 * time.Second,
			PathTimeouts:    make(map[string]time.Duration),
			OnTimeout:       defaultOnTimeout,
			TimeoutHandler:  defaultTimeoutHandler,
		},
		Services: map[string]ServiceConfig{
			"llm": {
				Retry: &RetryConfig{
					MaxRetries:          5,
					InitialDelay:        200 * time.Millisecond,
					MaxDelay:            10 * time.Second,
					Multiplier:          1.5,
					RandomizationFactor: 0.2,
					RetryableFunc:       defaultRetryableFunc,
					OnRetry:             defaultOnRetry,
				},
				Timeout: &TimeoutConfig{
					DefaultTimeout: 120 * time.Second,
					PathTimeouts: map[string]time.Duration{
						"/api/v1/llm/chat":   120 * time.Second,
						"/api/v1/llm/stream": 300 * time.Second, // 流式响应需要更长时间
					},
					OnTimeout:      defaultOnTimeout,
					TimeoutHandler: defaultTimeoutHandler,
				},
				EnableRetry:   true,
				EnableTimeout: true,
			},
			"auth": {
				Retry: &RetryConfig{
					MaxRetries:          2,
					InitialDelay:        100 * time.Millisecond,
					MaxDelay:            2 * time.Second,
					Multiplier:          2.0,
					RandomizationFactor: 0.1,
					RetryableFunc:       defaultRetryableFunc,
					OnRetry:             defaultOnRetry,
				},
				Timeout: &TimeoutConfig{
					DefaultTimeout: 15 * time.Second,
					PathTimeouts: map[string]time.Duration{
						"/auth/login":    10 * time.Second,
						"/auth/register": 15 * time.Second,
						"/auth/refresh":  5 * time.Second,
					},
					OnTimeout:      defaultOnTimeout,
					TimeoutHandler: defaultTimeoutHandler,
				},
				EnableRetry:   true,
				EnableTimeout: true,
			},
			"api": {
				Retry: &RetryConfig{
					MaxRetries:          3,
					InitialDelay:        150 * time.Millisecond,
					MaxDelay:            5 * time.Second,
					Multiplier:          2.0,
					RandomizationFactor: 0.1,
					RetryableFunc:       defaultRetryableFunc,
					OnRetry:             defaultOnRetry,
				},
				Timeout: &TimeoutConfig{
					DefaultTimeout: 30 * time.Second,
					PathTimeouts: map[string]time.Duration{
						"/api/v1/health":  5 * time.Second,
						"/api/v1/metrics": 10 * time.Second,
						"/api/v1/status":  5 * time.Second,
					},
					OnTimeout:      defaultOnTimeout,
					TimeoutHandler: defaultTimeoutHandler,
				},
				EnableRetry:   true,
				EnableTimeout: true,
			},
		},
	}
}

// GetRetryConfig 获取指定服务的重试配置
func (c *RetryAndTimeoutConfig) GetRetryConfig(serviceName string) RetryConfig {
	if serviceConfig, exists := c.Services[serviceName]; exists && serviceConfig.Retry != nil {
		return *serviceConfig.Retry
	}
	return c.GlobalRetry
}

// GetTimeoutConfig 获取指定服务的超时配置
func (c *RetryAndTimeoutConfig) GetTimeoutConfig(serviceName string) TimeoutConfig {
	if serviceConfig, exists := c.Services[serviceName]; exists && serviceConfig.Timeout != nil {
		return *serviceConfig.Timeout
	}
	return c.GlobalTimeout
}

// IsRetryEnabled 检查服务是否启用重试
func (c *RetryAndTimeoutConfig) IsRetryEnabled(serviceName string) bool {
	if serviceConfig, exists := c.Services[serviceName]; exists {
		return serviceConfig.EnableRetry
	}
	return true // 默认启用
}

// IsTimeoutEnabled 检查服务是否启用超时
func (c *RetryAndTimeoutConfig) IsTimeoutEnabled(serviceName string) bool {
	if serviceConfig, exists := c.Services[serviceName]; exists {
		return serviceConfig.EnableTimeout
	}
	return true // 默认启用
}

// CustomRetryConfig 自定义重试配置
type CustomRetryConfig struct {
	ServiceName string
	Config      RetryConfig
}

// CustomTimeoutConfig 自定义超时配置
type CustomTimeoutConfig struct {
	ServiceName string
	Config      TimeoutConfig
}

// WithCustomRetry 添加自定义重试配置
func (c *RetryAndTimeoutConfig) WithCustomRetry(serviceName string, config RetryConfig) *RetryAndTimeoutConfig {
	if c.Services == nil {
		c.Services = make(map[string]ServiceConfig)
	}

	if serviceConfig, exists := c.Services[serviceName]; exists {
		serviceConfig.Retry = &config
		serviceConfig.EnableRetry = true
		c.Services[serviceName] = serviceConfig
	} else {
		c.Services[serviceName] = ServiceConfig{
			Retry:         &config,
			EnableRetry:   true,
			EnableTimeout: true,
		}
	}

	return c
}

// WithCustomTimeout 添加自定义超时配置
func (c *RetryAndTimeoutConfig) WithCustomTimeout(serviceName string, config TimeoutConfig) *RetryAndTimeoutConfig {
	if c.Services == nil {
		c.Services = make(map[string]ServiceConfig)
	}

	if serviceConfig, exists := c.Services[serviceName]; exists {
		serviceConfig.Timeout = &config
		serviceConfig.EnableTimeout = true
		c.Services[serviceName] = serviceConfig
	} else {
		c.Services[serviceName] = ServiceConfig{
			Timeout:       &config,
			EnableRetry:   true,
			EnableTimeout: true,
		}
	}

	return c
}

// DisableRetryForService 禁用指定服务的重试
func (c *RetryAndTimeoutConfig) DisableRetryForService(serviceName string) *RetryAndTimeoutConfig {
	if c.Services == nil {
		c.Services = make(map[string]ServiceConfig)
	}

	if serviceConfig, exists := c.Services[serviceName]; exists {
		serviceConfig.EnableRetry = false
		c.Services[serviceName] = serviceConfig
	} else {
		c.Services[serviceName] = ServiceConfig{
			EnableRetry:   false,
			EnableTimeout: true,
		}
	}

	return c
}



// 默认函数实现

// defaultRetryableFunc 默认重试条件判断函数
func defaultRetryableFunc(err error) bool {
	if err == nil {
		return false
	}
	
	// 这里可以根据实际需求定义哪些错误需要重试
	// 例如：网络错误、超时错误、临时性错误等
	return true
}

// defaultOnRetry 默认重试回调函数
func defaultOnRetry(attempt int, err error) {
	// 默认实现：可以在这里记录日志
	// 例如：log.Printf("第 %d 次重试，错误: %v", attempt, err)
}

// defaultOnTimeout 默认超时回调函数
func defaultOnTimeout(req *http.Request, duration time.Duration) {
	// 默认实现：可以在这里记录超时日志
	// 例如：log.Printf("请求超时: %s, 耗时: %v", req.URL.Path, duration)
}

// defaultTimeoutHandler 默认超时处理函数
func defaultTimeoutHandler(req *http.Request, duration time.Duration) error {
	// 默认实现：返回超时错误
	return &timeoutError{
		message:  "request timeout",
		path:     req.URL.Path,
		duration: duration,
	}
}

// timeoutError 超时错误
type timeoutError struct {
	message  string
	path     string
	duration time.Duration
}

func (e *timeoutError) Error() string {
	return e.message
}

// DisableTimeoutForService 禁用指定服务的超时
func (c *RetryAndTimeoutConfig) DisableTimeoutForService(serviceName string) *RetryAndTimeoutConfig {
	if c.Services == nil {
		c.Services = make(map[string]ServiceConfig)
	}

	if serviceConfig, exists := c.Services[serviceName]; exists {
		serviceConfig.EnableTimeout = false
		c.Services[serviceName] = serviceConfig
	} else {
		c.Services[serviceName] = ServiceConfig{
			EnableRetry:   true,
			EnableTimeout: false,
		}
	}

	return c
}
