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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查两次输入的密码是否一致
	if registerRequest.Password != registerRequest.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "两次输入的密码不一致"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := pkg.DB.Where("username = ?", registerRequest.Username).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	result = pkg.DB.Where("email = ?", registerRequest.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 对密码进行哈希处理
	passwordHash, err := utils.HashPassword(registerRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
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
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "请求参数验证失败: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var user models.User
	result := pkg.DB.Where("username = ?", loginRequest.Username).First(&user)
	if result.Error != nil {
		// 记录用户不存在的登录尝试
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "用户名或密码错误")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		// 记录密码错误的登录尝试
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "用户名或密码错误")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		// 记录生成token失败的情况
		recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), false, "生成token失败: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 记录登录成功
	recordLoginHistory(loginRequest.Username, c.ClientIP(), c.Request.UserAgent(), true, "登录成功")

	// 不返回密码信息
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token, "user": user})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := pkg.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result = pkg.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result := pkg.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
