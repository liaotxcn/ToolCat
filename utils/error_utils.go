package utils

import (
	"encoding/json"
	"net/http"
	"weave/pkg"

	"go.uber.org/zap"
)

// JSONErrorResponse 发送JSON格式的错误响应
func JSONErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	// 确保状态码是4xx或5xx
	if statusCode < 400 {
		statusCode = http.StatusInternalServerError
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// 准备错误响应数据
	var errResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	// 根据错误类型设置响应数据
	if appErr, ok := err.(*pkg.AppError); ok {
		errResponse.Code = string(appErr.Code)
		errResponse.Message = appErr.Message
	} else {
		errResponse.Code = "INTERNAL_ERROR"
		errResponse.Message = "Internal server error"
	}

	// 写入响应
	json.NewEncoder(w).Encode(errResponse)
}

// HandleAPIError 处理API错误并返回标准响应
func HandleAPIError(w http.ResponseWriter, r *http.Request, err error, defaultMessage string) {
	// 获取HTTP状态码
	statusCode := pkg.GetHTTPStatus(err)

	// 记录错误日志
	pkg.With(
		zap.Error(err),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	).Error(defaultMessage)

	// 返回JSON错误响应
	JSONErrorResponse(w, err, statusCode)
}

// HandleDatabaseError 处理数据库错误
func HandleDatabaseError(w http.ResponseWriter, r *http.Request, err error, operation string) {
	// 包装数据库错误
	appErr := pkg.NewDatabaseError(operation+" failed", err)

	// 处理API错误
	HandleAPIError(w, r, appErr, "Database operation failed")
}

// HandleValidationError 处理参数验证错误
func HandleValidationError(w http.ResponseWriter, r *http.Request, err error, field string) {
	// 包装验证错误
	appErr := pkg.NewBadRequest(field+" validation failed", err)

	// 处理API错误
	HandleAPIError(w, r, appErr, "Validation failed")
}

// HandleNotFoundError 处理资源未找到错误
func HandleNotFoundError(w http.ResponseWriter, r *http.Request, resource string, id string) {
	// 创建未找到错误
	appErr := pkg.NewNotFound(resource+" with id '"+id+"' not found", nil)

	// 处理API错误
	HandleAPIError(w, r, appErr, "Resource not found")
}

// HandleUnauthorizedError 处理未授权错误
func HandleUnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	// 包装未授权错误
	appErr := pkg.NewUnauthorized("Unauthorized access", err)

	// 处理API错误
	HandleAPIError(w, r, appErr, "Unauthorized access")
}

// HandlePluginError 处理插件错误
func HandlePluginError(w http.ResponseWriter, r *http.Request, err error, pluginName string) {
	// 包装插件错误
	appErr := pkg.NewPluginError("Plugin error in "+pluginName, err)

	// 处理API错误
	HandleAPIError(w, r, appErr, "Plugin execution failed")
}
