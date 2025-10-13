package pkg

import (
	"errors"
	"fmt"
	"time"
)

// ErrorCode 定义错误码类型
type ErrorCode string

// 错误码分类前缀
const (
	// 客户端错误 (4xx)
	ClientErrorPrefix = "CLIENT_"
	// 服务器错误 (5xx)
	ServerErrorPrefix = "SERVER_"
	// 数据库错误
	DatabaseErrorPrefix = "DB_"
	// 插件错误
	PluginErrorPrefix = "PLUGIN_"
	// 认证错误
	AuthErrorPrefix = "AUTH_"
	// 参数验证错误
	ValidationErrorPrefix = "VALIDATION_"
)

// 常用错误码定义
const (
	// 客户端错误
	ErrBadRequest           ErrorCode = "BAD_REQUEST"
	ErrUnauthorized         ErrorCode = "UNAUTHORIZED"
	ErrForbidden            ErrorCode = "FORBIDDEN"
	ErrNotFound             ErrorCode = "NOT_FOUND"
	ErrConflict             ErrorCode = "CONFLICT"
	ErrTooManyRequests      ErrorCode = "TOO_MANY_REQUESTS"
	ErrUnsupportedMediaType ErrorCode = "UNSUPPORTED_MEDIA_TYPE"

	// 服务器错误
	ErrInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrNotImplemented     ErrorCode = "NOT_IMPLEMENTED"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrGatewayTimeout     ErrorCode = "GATEWAY_TIMEOUT"

	// 数据库错误
	ErrDatabaseError       ErrorCode = "DATABASE_ERROR"
	ErrDatabaseConnection  ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrDatabaseQuery       ErrorCode = "DATABASE_QUERY_ERROR"
	ErrDatabaseTransaction ErrorCode = "DATABASE_TRANSACTION_ERROR"
	ErrDatabaseConstraint  ErrorCode = "DATABASE_CONSTRAINT_ERROR"

	// 插件错误
	ErrPluginError      ErrorCode = "PLUGIN_ERROR"
	ErrPluginNotFound   ErrorCode = "PLUGIN_NOT_FOUND"
	ErrPluginDisabled   ErrorCode = "PLUGIN_DISABLED"
	ErrPluginDependency ErrorCode = "PLUGIN_DEPENDENCY_ERROR"
	ErrPluginInit       ErrorCode = "PLUGIN_INIT_ERROR"
	ErrPluginExecution  ErrorCode = "PLUGIN_EXECUTION_ERROR"

	// 认证错误
	ErrAuthInvalidToken     ErrorCode = "AUTH_INVALID_TOKEN"
	ErrAuthExpiredToken     ErrorCode = "AUTH_EXPIRED_TOKEN"
	ErrAuthInsufficientRole ErrorCode = "AUTH_INSUFFICIENT_ROLE"
	ErrAuthRateLimited      ErrorCode = "AUTH_RATE_LIMITED"

	// 参数验证错误
	ErrValidationRequired ErrorCode = "VALIDATION_REQUIRED"
	ErrValidationFormat   ErrorCode = "VALIDATION_FORMAT_ERROR"
	ErrValidationRange    ErrorCode = "VALIDATION_RANGE_ERROR"
	ErrValidationUnique   ErrorCode = "VALIDATION_UNIQUE_ERROR"
	ErrValidationLength   ErrorCode = "VALIDATION_LENGTH_ERROR"
)

// 错误码对应的默认错误信息
var DefaultErrorMessages = map[ErrorCode]string{
	ErrBadRequest:           "请求参数错误",
	ErrUnauthorized:         "未授权访问",
	ErrForbidden:            "权限不足",
	ErrNotFound:             "请求的资源不存在",
	ErrConflict:             "请求冲突",
	ErrTooManyRequests:      "请求过于频繁",
	ErrUnsupportedMediaType: "不支持的媒体类型",
	ErrInternalError:        "服务器内部错误",
	ErrNotImplemented:       "功能尚未实现",
	ErrServiceUnavailable:   "服务不可用",
	ErrGatewayTimeout:       "网关超时",
	ErrDatabaseError:        "数据库错误",
	ErrDatabaseConnection:   "数据库连接失败",
	ErrDatabaseQuery:        "数据库查询错误",
	ErrDatabaseTransaction:  "数据库事务错误",
	ErrDatabaseConstraint:   "数据库约束违反",
	ErrPluginError:          "插件错误",
	ErrPluginNotFound:       "插件不存在",
	ErrPluginDisabled:       "插件已禁用",
	ErrPluginDependency:     "插件依赖错误",
	ErrPluginInit:           "插件初始化失败",
	ErrPluginExecution:      "插件执行错误",
	ErrAuthInvalidToken:     "无效的令牌",
	ErrAuthExpiredToken:     "令牌已过期",
	ErrAuthInsufficientRole: "角色权限不足",
	ErrAuthRateLimited:      "认证请求受限",
	ErrValidationRequired:   "缺少必要参数",
	ErrValidationFormat:     "参数格式错误",
	ErrValidationRange:      "参数值超出范围",
	ErrValidationUnique:     "值必须唯一",
	ErrValidationLength:     "参数长度不符合要求",
}

// HTTPStatusMap 错误码对应的HTTP状态码
var HTTPStatusMap = map[ErrorCode]int{
	// 客户端错误 (4xx)
	ErrBadRequest:           400,
	ErrUnauthorized:         401,
	ErrForbidden:            403,
	ErrNotFound:             404,
	ErrConflict:             409,
	ErrTooManyRequests:      429,
	ErrUnsupportedMediaType: 415,

	// 服务器错误 (5xx)
	ErrInternalError:      500,
	ErrNotImplemented:     501,
	ErrServiceUnavailable: 503,
	ErrGatewayTimeout:     504,

	// 数据库错误
	ErrDatabaseError:       500,
	ErrDatabaseConnection:  503,
	ErrDatabaseQuery:       500,
	ErrDatabaseTransaction: 500,
	ErrDatabaseConstraint:  400,

	// 插件错误
	ErrPluginError:      500,
	ErrPluginNotFound:   404,
	ErrPluginDisabled:   403,
	ErrPluginDependency: 500,
	ErrPluginInit:       500,
	ErrPluginExecution:  500,

	// 认证错误
	ErrAuthInvalidToken:     401,
	ErrAuthExpiredToken:     401,
	ErrAuthInsufficientRole: 403,
	ErrAuthRateLimited:      429,

	// 参数验证错误
	ErrValidationRequired: 400,
	ErrValidationFormat:   400,
	ErrValidationRange:    400,
	ErrValidationUnique:   409,
	ErrValidationLength:   400,
}

// AppError 应用错误结构
type AppError struct {
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Err       error       `json:"-"`                   // 不序列化到JSON
	Details   interface{} `json:"details,omitempty"`   // 可选的错误详情
	RequestID string      `json:"requestId,omitempty"` // 请求ID，用于追踪
	Timestamp int64       `json:"timestamp"`           // 错误发生时间戳
	Path      string      `json:"path,omitempty"`      // 请求路径
}

// Error 实现error接口
func (e *AppError) Error() string {
	baseMsg := fmt.Sprintf("%s: %s", e.Code, e.Message)
	if e.Err != nil {
		baseMsg += fmt.Sprintf(" (original: %v)", e.Err)
	}
	if e.RequestID != "" {
		baseMsg += fmt.Sprintf(" [requestID: %s]", e.RequestID)
	}
	return baseMsg
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

// WithDetails 添加错误详情
func (e *AppError) WithDetails(details interface{}) *AppError {
	e.Details = details
	return e
}

// WithRequestID 添加请求ID
func (e *AppError) WithRequestID(requestID string) *AppError {
	e.RequestID = requestID
	return e
}

// WithPath 添加请求路径
func (e *AppError) WithPath(path string) *AppError {
	e.Path = path
	return e
}

// New 创建一个新的AppError
func New(code ErrorCode, message string, err error) *AppError {
	// 如果消息为空，使用默认消息
	if message == "" {
		if defaultMsg, exists := DefaultErrorMessages[code]; exists {
			message = defaultMsg
		}
	}

	return &AppError{
		Code:      code,
		Message:   message,
		Err:       err,
		Timestamp: time.Now().Unix(),
	}
}

// 客户端错误辅助函数
func NewBadRequest(message string, err error) *AppError {
	return New(ErrBadRequest, message, err)
}

// 添加别名以保持兼容性
func NewBadRequestError(message string, err error) *AppError {
	return NewBadRequest(message, err)
}

func NewUnauthorized(message string, err error) *AppError {
	return New(ErrUnauthorized, message, err)
}

// 添加别名以保持兼容性
func NewUnauthorizedError(message string, err error) *AppError {
	return NewUnauthorized(message, err)
}

func NewForbidden(message string, err error) *AppError {
	return New(ErrForbidden, message, err)
}

// 添加别名以保持兼容性
func NewForbiddenError(message string, err error) *AppError {
	return NewForbidden(message, err)
}

func NewNotFound(message string, err error) *AppError {
	return New(ErrNotFound, message, err)
}

// 添加别名以保持兼容性
func NewNotFoundError(message string, err error) *AppError {
	return NewNotFound(message, err)
}

func NewConflict(message string, err error) *AppError {
	return New(ErrConflict, message, err)
}

// 添加别名以保持兼容性
func NewConflictError(message string, err error) *AppError {
	return NewConflict(message, err)
}

func NewTooManyRequests(message string, err error) *AppError {
	return New(ErrTooManyRequests, message, err)
}

func NewUnsupportedMediaType(message string, err error) *AppError {
	return New(ErrUnsupportedMediaType, message, err)
}

// 服务器错误辅助函数
func NewInternalError(message string, err error) *AppError {
	return New(ErrInternalError, message, err)
}

// 添加别名以保持兼容性
func NewInternalServerError(message string, err error) *AppError {
	return NewInternalError(message, err)
}

func NewNotImplemented(message string, err error) *AppError {
	return New(ErrNotImplemented, message, err)
}

func NewServiceUnavailable(message string, err error) *AppError {
	return New(ErrServiceUnavailable, message, err)
}

// 添加别名以保持兼容性
func NewServiceUnavailableError(message string, err error) *AppError {
	return NewServiceUnavailable(message, err)
}

func NewGatewayTimeout(message string, err error) *AppError {
	return New(ErrGatewayTimeout, message, err)
}

// 数据库错误辅助函数
func NewDatabaseError(message string, err error) *AppError {
	return New(ErrDatabaseError, message, err)
}

func NewDatabaseConnectionError(message string, err error) *AppError {
	return New(ErrDatabaseConnection, message, err)
}

func NewDatabaseQueryError(message string, err error) *AppError {
	return New(ErrDatabaseQuery, message, err)
}

func NewDatabaseTransactionError(message string, err error) *AppError {
	return New(ErrDatabaseTransaction, message, err)
}

func NewDatabaseConstraintError(message string, err error) *AppError {
	return New(ErrDatabaseConstraint, message, err)
}

// 插件错误辅助函数
func NewPluginError(message string, err error) *AppError {
	return New(ErrPluginError, message, err)
}

func NewPluginNotFoundError(message string, err error) *AppError {
	return New(ErrPluginNotFound, message, err)
}

func NewPluginDisabledError(message string, err error) *AppError {
	return New(ErrPluginDisabled, message, err)
}

func NewPluginDependencyError(message string, err error) *AppError {
	return New(ErrPluginDependency, message, err)
}

func NewPluginInitError(message string, err error) *AppError {
	return New(ErrPluginInit, message, err)
}

func NewPluginExecutionError(message string, err error) *AppError {
	return New(ErrPluginExecution, message, err)
}

// 认证错误辅助函数
func NewAuthInvalidTokenError(message string, err error) *AppError {
	return New(ErrAuthInvalidToken, message, err)
}

func NewAuthExpiredTokenError(message string, err error) *AppError {
	return New(ErrAuthExpiredToken, message, err)
}

func NewAuthInsufficientRoleError(message string, err error) *AppError {
	return New(ErrAuthInsufficientRole, message, err)
}

func NewAuthRateLimitedError(message string, err error) *AppError {
	return New(ErrAuthRateLimited, message, err)
}

// 添加通用认证错误函数作为别名
func NewAuthError(message string, err error) *AppError {
	return NewUnauthorizedError(message, err)
}

// 参数验证错误辅助函数
func NewValidationRequiredError(message string, err error) *AppError {
	return New(ErrValidationRequired, message, err)
}

func NewValidationFormatError(message string, err error) *AppError {
	return New(ErrValidationFormat, message, err)
}

func NewValidationRangeError(message string, err error) *AppError {
	return New(ErrValidationRange, message, err)
}

func NewValidationUniqueError(message string, err error) *AppError {
	return New(ErrValidationUnique, message, err)
}

func NewValidationLengthError(message string, err error) *AppError {
	return New(ErrValidationLength, message, err)
}

// 添加通用验证错误函数作为别名
func NewValidationError(message string, err error) *AppError {
	return NewValidationFormatError(message, err)
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

// 错误类型判断辅助函数
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrNotFound
	}
	return false
}

func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrUnauthorized
	}
	return false
}

func IsForbidden(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrForbidden
	}
	return false
}

func IsBadRequest(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrBadRequest
	}
	return false
}

func IsConflict(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrConflict
	}
	return false
}

func IsDatabaseError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrDatabaseError ||
			appErr.Code == ErrDatabaseConnection ||
			appErr.Code == ErrDatabaseQuery ||
			appErr.Code == ErrDatabaseTransaction ||
			appErr.Code == ErrDatabaseConstraint
	}
	return false
}

func IsPluginError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrPluginError ||
			appErr.Code == ErrPluginNotFound ||
			appErr.Code == ErrPluginDisabled ||
			appErr.Code == ErrPluginDependency ||
			appErr.Code == ErrPluginInit ||
			appErr.Code == ErrPluginExecution
	}
	return false
}

func IsAuthError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrUnauthorized ||
			appErr.Code == ErrAuthInvalidToken ||
			appErr.Code == ErrAuthExpiredToken ||
			appErr.Code == ErrAuthInsufficientRole ||
			appErr.Code == ErrAuthRateLimited
	}
	return false
}

func IsValidationError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrValidationRequired ||
			appErr.Code == ErrValidationFormat ||
			appErr.Code == ErrValidationRange ||
			appErr.Code == ErrValidationUnique ||
			appErr.Code == ErrValidationLength
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
