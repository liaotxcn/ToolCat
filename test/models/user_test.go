package models

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"toolcat/models"
	"toolcat/pkg"

	"gorm.io/gorm"
)

// 初始化测试数据库
func setupTestDB(t *testing.T) {
	// 设置测试环境变量
	testVars := map[string]string{
		"DB_NAME":     "toolcat_test_model",
		"DB_HOST":     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		"DB_PORT":     getEnvOrDefault("TEST_DB_PORT", "3306"),
		"DB_USERNAME": getEnvOrDefault("TEST_DB_USERNAME", "root"),
		"DB_PASSWORD": getEnvOrDefault("TEST_DB_PASSWORD", "123456"),
	}

	// 设置环境变量
	for key, value := range testVars {
		os.Setenv(key, value)
	}

	// 初始化数据库连接
	err := pkg.InitDatabase()
	if err != nil {
		t.Skipf("Skipping database tests: %v\n请确保MySQL服务已启动，并且配置了正确的连接参数", err)
		return
	}

	// 迁移表结构
	models.MigrateTables()
}

// 清理测试数据
func cleanupTestDB(t *testing.T) {
	// 删除测试数据但保留表结构
	if pkg.DB != nil {
		pkg.DB.Exec("DELETE FROM users")
		pkg.DB.Exec("DELETE FROM tools")
		pkg.DB.Exec("DELETE FROM tool_histories")
		pkg.DB.Exec("DELETE FROM notes")
		pkg.DB.Exec("DELETE FROM login_histories")

		// 关闭数据库连接
		pkg.CloseDatabase()
	}
}

// getEnvOrDefault 获取环境变量，如果为空则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func TestMain(m *testing.M) {
	t := &testing.T{}
	var isSkipped bool

	setupTestDB(t)

	if isSkipped {
		os.Exit(0)
	}

	// 运行测试
	code := m.Run()

	// 清理
	cleanupTestDB(t)

	os.Exit(code)
}

// TestUserCreate 测试创建用户
func TestUserCreate(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 创建测试用户
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存用户到数据库
	err := pkg.DB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 验证用户是否创建成功（ID不为0）
	if user.ID == 0 {
		t.Fatal("User ID is 0, creation failed")
	}

	// 清理
	defer pkg.DB.Delete(&user)
}

// TestUserFind 测试查找用户
func TestUserFind(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 先创建一个测试用户
	user := models.User{
		Username: "paipai",
		Email:    "666666@qq.com",
		Password: "123456",
	}

	// 保存用户到数据库
	err := pkg.DB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user for find test: %v", err)
	}

	// 清理
	defer pkg.DB.Delete(&user)

	// 通过ID查找用户
	var foundUser models.User
	err = pkg.DB.First(&foundUser, user.ID).Error
	if err != nil {
		t.Fatalf("Failed to find user by ID: %v", err)
	}

	// 验证找到的用户是否正确
	if foundUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, foundUser.Username)
	}

	// 通过用户名查找用户
	var foundByUsername models.User
	err = pkg.DB.Where("username = ?", user.Username).First(&foundByUsername).Error
	if err != nil {
		t.Fatalf("Failed to find user by username: %v", err)
	}

	// 验证找到的用户是否正确
	if foundByUsername.ID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, foundByUsername.ID)
	}

	// 测试查找不存在的用户
	var notFoundUser models.User
	err = pkg.DB.First(&notFoundUser, 999999).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected record not found error, got %v", err)
	}
}

// TestUserUpdate 测试更新用户信息
func TestUserUpdate(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 先创建一个测试用户
	user := models.User{
		Username: "updateuser",
		Email:    "update@example.com",
		Password: "password123",
	}

	// 保存用户到数据库
	err := pkg.DB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user for update test: %v", err)
	}

	// 清理
	defer pkg.DB.Delete(&user)

	// 更新用户信息
	newNickname := "Updated Nickname"
	newEmail := "updated@example.com"

	err = pkg.DB.Model(&user).Updates(map[string]interface{}{
		"nickname": newNickname,
		"email":    newEmail,
	}).Error
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// 重新查询用户以验证更新
	var updatedUser models.User
	err = pkg.DB.First(&updatedUser, user.ID).Error
	if err != nil {
		t.Fatalf("Failed to find user after update: %v", err)
	}

	if updatedUser.Email != newEmail {
		t.Errorf("Expected email %s, got %s", newEmail, updatedUser.Email)
	}
}

// TestUserDelete 测试删除用户
func TestUserDelete(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 先创建一个测试用户
	user := models.User{
		Username: "deleteuser",
		Email:    "delete@example.com",
		Password: "password123",
	}

	// 保存用户到数据库
	err := pkg.DB.Create(&user).Error
	if err != nil {
		t.Fatalf("Failed to create user for delete test: %v", err)
	}

	// 记录用户ID用于后续验证
	userId := user.ID

	// 删除用户
	err = pkg.DB.Delete(&user).Error
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// 验证用户是否已被删除
	var deletedUser models.User
	err = pkg.DB.First(&deletedUser, userId).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected record not found error after delete, got %v", err)
	}
}

// TestMigrateTables 测试数据库迁移功能
func TestMigrateTables(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 调用迁移函数
	models.MigrateTables()

	// 验证表是否存在
	tableNames := []string{"users", "tools", "tool_histories", "notes", "login_histories"}

	for _, tableName := range tableNames {
		exists := false
		sql := fmt.Sprintf("SHOW TABLES LIKE '%s'", tableName)
		var result []string
		err := pkg.DB.Raw(sql).Scan(&result).Error
		if err != nil {
			t.Errorf("Failed to check if table %s exists: %v", tableName, err)
			continue
		}

		if len(result) > 0 {
			exists = true
		}

		if !exists {
			t.Errorf("Table %s does not exist after migration", tableName)
		}
	}
}
