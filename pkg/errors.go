package pkg

import (
	"errors"
	"fmt"
)

// ErrorCode 定义错误码类型
type ErrorCode string

// 常用错误码定义
const (
	// 客户端错误
	ErrBadRequest     ErrorCode = "BAD_REQUEST"
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrForbidden      ErrorCode = "FORBIDDEN"
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrConflict       ErrorCode = "CONFLICT"

	// 服务器错误
	ErrInternalError  ErrorCode = "INTERNAL_ERROR"
	ErrDatabaseError  ErrorCode = "DATABASE_ERROR"
	ErrCacheError     ErrorCode = "CACHE_ERROR"
	ErrNetworkError   ErrorCode = "NETWORK_ERROR"

	// 插件错误
	ErrPluginError    ErrorCode = "PLUGIN_ERROR"
	ErrPluginNotFound ErrorCode = "PLUGIN_NOT_FOUND"
	ErrPluginDisabled ErrorCode = "PLUGIN_DISABLED"
)

// HTTPStatusMap 错误码对应的HTTP状态码
var HTTPStatusMap = map[ErrorCode]int{
	ErrBadRequest:     400,
	ErrUnauthorized:   401,
	ErrForbidden:      403,
	ErrNotFound:       404,
	ErrConflict:       409,
	ErrInternalError:  500,
	ErrDatabaseError:  500,
	ErrCacheError:     500,
	ErrNetworkError:   500,
	ErrPluginError:    500,
	ErrPluginNotFound: 404,
	ErrPluginDisabled: 403,
}

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"` // 不序列化到JSON
	Details interface{} `json:"details,omitempty"` // 可选的错误详情
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (original: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is 实现errors.Is接口
func (e *AppError) Is(target error) bool {
	if targetAppErr, ok := target.(*AppError); ok {
		return e.Code == targetAppErr.Code
	}
	return false
}

// New 创建一个新的AppError
func New(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewBadRequest 创建一个400错误
func NewBadRequest(message string, err error) *AppError {
	return New(ErrBadRequest, message, err)
}

// NewUnauthorized 创建一个401错误
func NewUnauthorized(message string, err error) *AppError {
	return New(ErrUnauthorized, message, err)
}

// NewNotFound 创建一个404错误
func NewNotFound(message string, err error) *AppError {
	return New(ErrNotFound, message, err)
}

// NewInternalError 创建一个500错误
func NewInternalError(message string, err error) *AppError {
	return New(ErrInternalError, message, err)
}

// NewDatabaseError 创建一个数据库错误
func NewDatabaseError(message string, err error) *AppError {
	return New(ErrDatabaseError, message, err)
}

// NewPluginError 创建一个插件错误
func NewPluginError(message string, err error) *AppError {
	return New(ErrPluginError, message, err)
}

// Wrap 包装现有错误为AppError
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，直接返回
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsNotFound 判断错误是否为NotFound错误
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrNotFound
	}
	return false
}

// IsUnauthorized 判断错误是否为Unauthorized错误
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrUnauthorized
	}
	return false
}

// GetHTTPStatus 获取错误对应的HTTP状态码
func GetHTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		if status, exists := HTTPStatusMap[appErr.Code]; exists {
			return status
		}
	}
	return 500 // 默认返回500内部服务器错误
}