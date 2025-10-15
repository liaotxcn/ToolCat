package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"toolcat/pkg"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// MigrationManager 迁移管理器
type MigrationManager struct {
	migrationsDir string
	migrate       *migrate.Migrate
}

// NewMigrationManager 创建新的迁移管理器
func NewMigrationManager() *MigrationManager {
	return &MigrationManager{
		migrationsDir: filepath.Join("pkg", "migrate", "data_sql"),
	}
}

// Init 初始化迁移管理器
func (mm *MigrationManager) Init() error {
	// 确保迁移目录存在
	if err := os.MkdirAll(mm.migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// 获取数据库连接
	sqlDB, err := pkg.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// 初始化数据库驱动
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// 创建迁移源
	migrationsFS := os.DirFS(mm.migrationsDir)
	source, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// 创建迁移实例
	mm.migrate, err = migrate.NewWithInstance("iofs", source, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return nil
}

// Up 执行所有未应用的迁移
func (mm *MigrationManager) Up() error {
	if err := mm.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

// Down 回滚一个迁移
func (mm *MigrationManager) Down() error {
	return mm.migrate.Down()
}

// CreateMigration 创建新的迁移文件
func (mm *MigrationManager) CreateMigration(name string) (string, error) {
	// 获取最新迁移版本号
	version, err := mm.getNextVersion()
	if err != nil {
		return "", err
	}

	// 创建迁移文件路径
	upFileName := fmt.Sprintf("%03d_%s.up.sql", version, name)
	downFileName := fmt.Sprintf("%03d_%s.down.sql", version, name)

	upFilePath := filepath.Join(mm.migrationsDir, upFileName)
	downFilePath := filepath.Join(mm.migrationsDir, downFileName)

	// 创建空的迁移文件
	if err := os.WriteFile(upFilePath, []byte("-- Write your migration here\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to create up migration file: %w", err)
	}

	if err := os.WriteFile(downFilePath, []byte("-- Write your rollback here\n"), 0644); err != nil {
		// 清理已创建的up文件
		os.Remove(upFilePath)
		return "", fmt.Errorf("failed to create down migration file: %w", err)
	}

	return upFilePath, nil
}

// getNextVersion 获取下一个迁移版本号
func (mm *MigrationManager) getNextVersion() (int, error) {
	availableVersions, err := mm.getAvailableVersions()
	if err != nil || len(availableVersions) == 0 {
		return 1, nil // 如果没有版本或出错，从1开始
	}

	// 找到最大版本号
	maxVersion := availableVersions[0]
	for _, v := range availableVersions {
		if v > maxVersion {
			maxVersion = v
		}
	}

	return maxVersion + 1, nil
}

// getAvailableVersions 从文件系统获取可用的迁移版本号
func (mm *MigrationManager) getAvailableVersions() ([]int, error) {
	entries, err := os.ReadDir(mm.migrationsDir)
	if err != nil {
		return []int{}, nil // 如果目录不存在或无法读取，返回空列表
	}

	versionMap := make(map[int]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.Contains(name, ".") {
			continue
		}

		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}

		var version int
		if _, err := fmt.Sscanf(parts[0], "%03d", &version); err != nil {
			continue
		}

		versionMap[version] = true
	}

	// 转换为有序切片
	var versions []int
	for v := range versionMap {
		versions = append(versions, v)
	}

	// 简单排序
	for i := 0; i < len(versions); i++ {
		for j := i + 1; j < len(versions); j++ {
			if versions[i] > versions[j] {
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}

	return versions, nil
}

// GetStatus 获取迁移状态
func (mm *MigrationManager) GetStatus() (string, error) {
	version, dirty, err := mm.migrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return "", fmt.Errorf("failed to get migration status: %w", err)
	}

	status := fmt.Sprintf("Current version: %d\n", version)
	if dirty {
		status += "Database is dirty, last migration failed\n"
	}

	// 从文件系统获取可用的迁移版本
	availableVersions, err := mm.getAvailableVersions()
	if err != nil {
		return status, err
	}

	status += "Available versions: "
	for _, v := range availableVersions {
		status += fmt.Sprintf("%d ", v)
	}

	return status, nil
}

// GenerateInitialMigrations 生成初始迁移文件（基于当前模型）
func (mm *MigrationManager) GenerateInitialMigrations() error {
	// 获取数据库连接
	sqlDB, err := pkg.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// 生成初始迁移文件
	upContent, err := mm.generateCreateTablesSQL(sqlDB)
	if err != nil {
		return err
	}

	downContent, err := mm.generateDropTablesSQL(sqlDB)
	if err != nil {
		return err
	}

	// 创建初始迁移文件
	upFilePath := filepath.Join(mm.migrationsDir, "001_initial_schema.up.sql")
	downFilePath := filepath.Join(mm.migrationsDir, "001_initial_schema.down.sql")

	if err := os.WriteFile(upFilePath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create initial up migration: %w", err)
	}

	if err := os.WriteFile(downFilePath, []byte(downContent), 0644); err != nil {
		// 清理已创建的up文件
		os.Remove(upFilePath)
		return fmt.Errorf("failed to create initial down migration: %w", err)
	}

	return nil
}

// generateCreateTablesSQL 生成创建表的SQL
func (mm *MigrationManager) generateCreateTablesSQL(db *sql.DB) (string, error) {
	// 这里可以基于当前模型生成创建表SQL
	// 或者从数据库中获取当前表结构生成SQL
	// 简单实现：返回当前模型的创建表SQL
	return `-- Initial schema creation

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    username varchar(50) NOT NULL,
    password varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY idx_username (username),
    UNIQUE KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Notes table
CREATE TABLE IF NOT EXISTS notes (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    title varchar(255) NOT NULL,
    content text,
    user_id bigint unsigned NOT NULL,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_user_id (user_id),
    CONSTRAINT fk_note_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Login history table
CREATE TABLE IF NOT EXISTS login_histories (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    username varchar(50) NOT NULL,
    ip_address varchar(45) NOT NULL,
    user_agent text,
    success tinyint(1) NOT NULL DEFAULT '0',
    error_message text,
    created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_username (username),
    KEY idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`, nil
}

// generateDropTablesSQL 生成删除表的SQL
func (mm *MigrationManager) generateDropTablesSQL(db *sql.DB) (string, error) {
	return `-- Rollback initial schema

DROP TABLE IF EXISTS login_histories;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS users;
`, nil
}
