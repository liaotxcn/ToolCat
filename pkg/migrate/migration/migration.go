package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"toolcat/config"
	"toolcat/pkg"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// MigrationManager 迁移管理器
type MigrationManager struct {
	migrationsDir string
	migrate       *migrate.Migrate
	driverName    string // 数据库驱动名称
}

// NewMigrationManager 创建新的迁移管理器
func NewMigrationManager() *MigrationManager {
	// 获取数据库驱动名称
	driverName := config.Config.Database.Driver
	if driverName == "" {
		driverName = "mysql" // 默认MySQL
	}

	return &MigrationManager{
		migrationsDir: filepath.Join("pkg", "migrate", "data_sql"),
		driverName:    driverName,
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

	// 确保driverName被设置
	if mm.driverName == "" {
		mm.driverName = config.Config.Database.Driver
		if mm.driverName == "" {
			mm.driverName = "mysql" // 默认MySQL
		}
	}

	// 根据数据库驱动类型初始化迁移驱动
	var driver database.Driver

	switch mm.driverName {
	case "postgres":
		driver, err = postgres.WithInstance(sqlDB, &postgres.Config{})
	case "mysql":
		fallthrough
	default:
		driver, err = mysql.WithInstance(sqlDB, &mysql.Config{})
	}

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
	mm.migrate, err = migrate.NewWithInstance("iofs", source, mm.driverName, driver)
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
	// 获取数据库驱动类型
	driverName := "mysql" // 默认MySQL
	if mm.driverName != "" {
		driverName = mm.driverName
	}

	if driverName == "postgres" {
		// PostgreSQL版本的SQL
		return `-- Initial schema creation (PostgreSQL)

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE (username),
    UNIQUE (email)
);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为users表添加更新时间触发器
CREATE TRIGGER update_users_timestamp
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Notes table
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_note_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- 为notes表添加更新时间触发器
CREATE TRIGGER update_notes_timestamp
BEFORE UPDATE ON notes
FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Login history table
CREATE TABLE IF NOT EXISTS login_histories (
    id SERIAL NOT NULL,
    username VARCHAR(50) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    error_message TEXT,
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_login_histories_username ON login_histories(username);
CREATE INDEX idx_login_histories_created_at ON login_histories(created_at);
`, nil
	}

	// MySQL版本的SQL (默认)
	return `-- Initial schema creation (MySQL)

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
