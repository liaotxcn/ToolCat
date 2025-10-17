package controllers

import (
	"fmt"
	"net/http"
	"time"

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
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "请求参数验证失败: "+err.Error())
		err := pkg.NewValidationError("Missing required login fields", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 查找用户
	var user models.User
	result := pkg.DB.Where("username = ?", loginRequest.Username).First(&user)
	if result.Error != nil {
		// 记录用户不存在的登录尝试
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "用户名或密码错误")
		err := pkg.NewAuthError("Invalid username or password", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		// 记录密码错误的登录尝试
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "用户名或密码错误")
		err := pkg.NewAuthError("Invalid username or password", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 生成访问令牌和刷新令牌
	accessToken, err := utils.GenerateToken(user.ID)
	if err != nil {
		// 记录生成token失败的情况
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "生成访问令牌失败: "+err.Error())
		err := pkg.NewInternalError("Failed to generate access token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		// 记录生成刷新令牌失败的情况
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "生成刷新令牌失败: "+err.Error())
		err := pkg.NewInternalError("Failed to generate refresh token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录登录成功
	recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), true, "登录成功")

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

	// 验证刷新令牌
	userID, err := utils.VerifyRefreshToken(refreshRequest.RefreshToken)
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

	// 生成新的访问令牌
	accessToken, err := utils.GenerateToken(userID)
	if err != nil {
		err := pkg.NewInternalError("Failed to generate access token", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 生成新的刷新令牌（可选：也可以继续使用原有的刷新令牌）
	refreshToken, err := utils.GenerateRefreshToken(userID)
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
func recordLoginHistory(username, ipAddress, userAgent string, success bool, message string) {
	loginHistory := models.LoginHistory{
		Username:  username,
		IPAddress: ipAddress,
		Success:   success,
		Message:   message,
		UserAgent: userAgent,
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
	result := pkg.DB.Find(&users)
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

	var user models.User
	result := pkg.DB.First(&user, id)
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

	result := pkg.DB.Create(&user)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to create user", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser 更新用户
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	result := pkg.DB.First(&user, id)
	if result.Error != nil {
		err := pkg.NewNotFoundError("User not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		err := pkg.NewValidationError("Invalid user data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	result = pkg.DB.Save(&user)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to update user", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result := pkg.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to delete user", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	if result.RowsAffected == 0 {
		err := pkg.NewNotFoundError("User not found", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
