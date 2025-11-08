package models

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"weave/models"
	"weave/pkg"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// 初始化测试数据库（使用临时文件SQLite以确保持久化）
func setupTestDB(t *testing.T) {
	// 始终为每个测试创建独立的SQLite数据库，避免共享全局连接造成干扰
	// 如果之前已有连接，先关闭释放资源
	if pkg.DB != nil {
		_ = pkg.CloseDatabase()
		pkg.DB = nil
	}

	// 为当前测试创建临时数据库文件
	dbPath := filepath.Join(t.TempDir(), "weave_test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite test db: %v", err)
	}
	pkg.DB = db
	// 单连接以简化测试环境
	sqlDB, err := pkg.DB.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	if err := models.MigrateTables(pkg.DB); err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}
}

// 清理测试数据
func cleanupTestDB(t *testing.T) {
	// 删除测试数据但保留表结构（忽略可能的错误，确保清理不中断）
	if pkg.DB != nil {
		pkg.DB.Exec("DELETE FROM users")
		pkg.DB.Exec("DELETE FROM tools")
		pkg.DB.Exec("DELETE FROM tool_histories")
		pkg.DB.Exec("DELETE FROM notes")
		pkg.DB.Exec("DELETE FROM login_histories")
		// 关闭数据库以释放文件句柄，避免TempDir清理失败
		_ = pkg.CloseDatabase()
		pkg.DB = nil
	}
}

// TestUserCreate 测试创建用户
func TestUserCreate(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 设置测试模式
	gin.SetMode(gin.TestMode)

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

// TestMigrateTables 验证迁移后表是否存在
func TestMigrateTables(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 使用SQLite元数据检查表是否存在
	tables := []string{"users", "tools", "tool_histories", "notes", "login_histories", "audit_logs"}
	for _, tbl := range tables {
		var count int64
		// sqlite_master 查询检查表存在
		res := pkg.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tbl).Scan(&count)
		if res.Error != nil {
			t.Fatalf("error checking table %s existence: %v", tbl, res.Error)
		}
		if count == 0 {
			t.Fatalf("table %s does not exist after migration", tbl)
		}
	}
}

// TestUserUpdate 测试更新用户信息
func TestUserUpdate(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	// 设置测试模式
	gin.SetMode(gin.TestMode)

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

	// 更新用户信息（更新现有字段，避免使用不存在的列）
	newEmail := "updated@example.com"

	err = pkg.DB.Model(&user).Updates(map[string]interface{}{
		"email": newEmail,
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

	// 设置测试模式
	gin.SetMode(gin.TestMode)

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
