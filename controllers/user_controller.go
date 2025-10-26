package controllers

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"toolcat/models"
	"toolcat/pkg"
	"toolcat/utils"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct{}

// Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	// 定义注册请求结构体
	var registerRequest struct {
		Username        string `json:"username" binding:"required,min=3,max=50"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
		Email           string `json:"email" binding:"required,email"`
	}

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		err := pkg.NewValidationError("Invalid registration data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 检查两次输入的密码是否一致
	if registerRequest.Password != registerRequest.ConfirmPassword {
		err := pkg.NewValidationError("Passwords do not match", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := pkg.DB.Where("username = ?", registerRequest.Username).First(&existingUser)
	if result.Error == nil {
		err := pkg.NewConflictError("Username already exists", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 检查邮箱是否已存在
	result = pkg.DB.Where("email = ?", registerRequest.Email).First(&existingUser)
	if result.Error == nil {
		err := pkg.NewConflictError("Email already registered", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 对密码进行哈希处理
	passwordHash, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		err := pkg.NewInternalError("Failed to encrypt password", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 创建新用户
	newUser := models.User{
		Username: registerRequest.Username,
		Password: passwordHash,
		Email:    registerRequest.Email,
	}

	result = pkg.DB.Create(&newUser)
	if result.Error != nil {
		dbErr := pkg.NewDatabaseError("Failed to register user", result.Error)
		dbErr.WithDetails(map[string]interface{}{
			"username": registerRequest.Username,
			"email":    registerRequest.Email,
		})
		c.JSON(pkg.GetHTTPStatus(dbErr), gin.H{"code": string(dbErr.Code), "message": dbErr.Message})
		return
	}

	// 不返回密码信息
	newUser.Password = ""
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "user": newUser})
}

// Login 用户登录
func (uc *UserController) Login(c *gin.Context) {
	// 定义登录请求结构体
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		// 记录绑定失败的登录尝试
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "请求参数验证失败: "+err.Error(), 0)
		err := pkg.NewValidationError("Missing required login fields", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 查找用户
	var user models.User
	result := pkg.DB.Where("username = ?", loginRequest.Username).First(&user)
	if result.Error != nil {
		// 记录用户不存在的登录尝试
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "用户名或密码错误", 0)
		err := pkg.NewAuthError("Invalid username or password", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		// 记录密码错误的登录尝试
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "用户名或密码错误", user.TenantID)
		err := pkg.NewAuthError("Invalid username or password", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 生成访问令牌和刷新令牌（包含tenant_id）
	accessToken, err := utils.GenerateToken(user.ID, user.TenantID)
	if err != nil {
		// 记录生成token失败的情况
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "生成访问令牌失败: "+err.Error(), user.TenantID)
		err := pkg.NewInternalError("Failed to generate access token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.TenantID)
	if err != nil {
		// 记录生成刷新令牌失败的情况
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "生成刷新令牌失败: "+err.Error(), user.TenantID)
		err := pkg.NewInternalError("Failed to generate refresh token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录登录成功
	recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), true, "登录成功", user.TenantID)

	// 记录登录操作的审计日志
	loginUser := user
	loginUser.Password = "[REDACTED]"
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "login",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", user.ID),
		OldValue:     nil,
		NewValue: map[string]interface{}{
			"username":   user.Username,
			"ip_address": c.ClientIP(),
			"success":    true,
		},
	})

	// 不返回密码信息
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "access_token": accessToken, "refresh_token": refreshToken, "user": user})
}

// RefreshToken 刷新访问令牌
func (uc *UserController) RefreshToken(c *gin.Context) {
	// 定义刷新令牌请求结构体
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// 绑定JSON请求体
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		err := pkg.NewValidationError("Refresh token is required", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 验证刷新令牌（获取userID与tenantID）
	userID, tenantID, err := utils.VerifyRefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		err := pkg.NewAuthError("Invalid refresh token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 查找用户
	var user models.User
	result := pkg.DB.First(&user, userID)
	if result.Error != nil {
		err := pkg.NewNotFoundError("User not found", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 生成新的访问令牌（保持相同tenant_id）
	accessToken, err := utils.GenerateToken(userID, tenantID)
	if err != nil {
		err := pkg.NewInternalError("Failed to generate access token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 生成新的刷新令牌（保持相同tenant_id）
	refreshToken, err := utils.GenerateRefreshToken(userID, tenantID)
	if err != nil {
		err := pkg.NewInternalError("Failed to generate refresh token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 不返回密码信息
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "令牌刷新成功", "access_token": accessToken, "refresh_token": refreshToken, "user": user})
}

// recordLoginHistory 记录登录历史
func recordLoginHistory(username, ipAddress, userAgent string, success bool, message string, tenantID uint) {
	loginHistory := models.LoginHistory{
		Username:  username,
		IPAddress: ipAddress,
		Success:   success,
		Message:   message,
		UserAgent: userAgent,
		TenantID:  tenantID,
		LoginTime: time.Now(),
	}

	// 异步记录登录历史，不阻塞主流程
	go func() {
		if err := pkg.DB.Create(&loginHistory).Error; err != nil {
			// 记录失败不应影响主流程，可以记录到日志中
			fmt.Printf("Failed to record login history: %v\n", err)
		}
	}()
}

// GetUsers 获取所有用户
func (uc *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	tenantID := c.GetUint("tenant_id")
	// 根据需要预加载关联数据，避免N+1查询问题
	result := pkg.DB.Where("tenant_id = ?", tenantID).Find(&users)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to fetch users", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser 获取单个用户
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	var user models.User
	// 根据API需求预加载关联数据，这里根据常见使用场景选择预加载审计日志
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).
		Preload("AuditLogs", func(db *gorm.DB) *gorm.DB {
			// 只预加载最近30天的审计日志
			return db.Where("created_at > ?", time.Now().AddDate(0, 0, -30)).Order("created_at DESC").Limit(100)
		}).First(&user)
	if result.Error != nil {
		err := pkg.NewNotFoundError("User not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		err := pkg.NewValidationError("Invalid user data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 绑定租户ID，防止跨租户创建
	user.TenantID = c.GetUint("tenant_id")

	// 创建用户前先记录审计日志（不包含密码）
	logUser := user
	logUser.Password = "[REDACTED]"

	result := pkg.DB.Create(&user)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to create user", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录创建用户的审计日志
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "create",
		ResourceType: "user",
		ResourceID:   fmt.Sprintf("%d", user.ID),
		OldValue:     nil,
		NewValue:     logUser,
	})

	// 返回用户信息（不包含密码）
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// UpdateUser 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	// 获取原始用户信息
	var oldUser models.User
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&oldUser)
	if result.Error != nil {
		err := pkg.NewNotFoundError("User not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录原始值（不包含密码）
	auditOldUser := oldUser
	auditOldUser.Password = "[REDACTED]"

	// 绑定新的用户信息
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		err := pkg.NewValidationError("Invalid user data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 防止跨租户变更
	newUser.TenantID = tenantID
	newUser.ID = oldUser.ID // 确保ID不变

	// 如果没有更新密码，则保留原密码
	if newUser.Password == "" {
		newUser.Password = oldUser.Password
	}

	result = pkg.DB.Save(&newUser)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to update user", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录更新用户的审计日志
	auditNewUser := newUser
	auditNewUser.Password = "[REDACTED]"
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "update",
		ResourceType: "user",
		ResourceID:   id,
		OldValue:     auditOldUser,
		NewValue:     auditNewUser,
	})

	// 返回更新后的用户信息（不包含密码）
	newUser.Password = ""
	c.JSON(http.StatusOK, newUser)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	// 先获取要删除的用户信息，用于审计日志
	var user models.User
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user)
	if result.Error != nil {
		err := pkg.NewNotFoundError("User not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录要删除的用户信息（不包含密码）
	auditUser := user
	auditUser.Password = "[REDACTED]"

	// 执行删除操作
	result = pkg.DB.Delete(&user)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to delete user", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录删除用户的审计日志
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "delete",
		ResourceType: "user",
		ResourceID:   id,
		OldValue:     auditUser,
		NewValue:     nil,
	})

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
