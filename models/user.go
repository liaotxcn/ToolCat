package models

import (
	"time"
	"toolcat/pkg"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;not null;unique" json:"username"`
	Password  string    `gorm:"size:100;not null" json:"password,omitempty"`
	Email     string    `gorm:"size:100;unique" json:"email"`
	TenantID  uint      `gorm:"index" json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
func MigrateTables() error {
	// 自动迁移表结构
	if err := pkg.DB.AutoMigrate(&User{}, &Tool{}, &ToolHistory{}, &Note{}, &LoginHistory{}); err != nil {
		return err
	}

	return nil
}
