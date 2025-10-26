package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"size:50;not null;unique" json:"username"`
	Password       string    `gorm:"size:100;not null" json:"password,omitempty"`
	Email          string    `gorm:"size:100;unique" json:"email"`
	TenantID       uint      `gorm:"index" json:"tenant_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	// 添加关联关系
	Notes           []Note         `gorm:"foreignKey:UserID" json:"notes,omitempty"`
	LoginHistories  []LoginHistory `gorm:"foreignKey:Username;references:Username" json:"login_histories,omitempty"`
	AuditLogs       []AuditLog     `gorm:"foreignKey:UserID" json:"audit_logs,omitempty"`
}

// Tool 工具模型
type Tool struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"size:255" json:"icon"`
	PluginName  string    `gorm:"size:100;not null" json:"plugin_name"`
	IsEnabled   bool      `gorm:"default:true" json:"is_enabled"`
	TenantID    uint      `gorm:"index" json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToolHistory 工具使用历史模型
type ToolHistory struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	UserID   uint      `json:"user_id"`
	ToolID   uint      `json:"tool_id"`
	TenantID uint      `gorm:"index" json:"tenant_id"`
	UsedAt   time.Time `json:"used_at"`
	Params   string    `gorm:"type:text" json:"params"`
	Result   string    `gorm:"type:text" json:"result"`
}

// MigrateTables 执行数据库迁移
func MigrateTables(db *gorm.DB) error {
	// 自动迁移表结构
	if err := db.AutoMigrate(&User{}, &Tool{}, &ToolHistory{}, &Note{}, &LoginHistory{}, &AuditLog{}); err != nil {
		return err
	}

	return nil
}
