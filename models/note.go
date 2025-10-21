package models

import (
	"time"
)

// Note 笔记模型
type Note struct {
	ID          string    `gorm:"primaryKey;size:100" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"` // 添加用户ID字段，建立索引提高查询效率
	TenantID    uint      `gorm:"index" json:"tenant_id"`
	Title       string    `gorm:"size:255;not null;index" json:"title"` // 添加索引
	Content     string    `gorm:"type:text;not null" json:"content"`
	CreatedTime time.Time `gorm:"index" json:"created_time"` // 添加索引
	UpdatedTime time.Time `json:"updated_time"`
}

// AddNoteToMigration 将Note模型添加到迁移函数
func init() {
	// 修改现有的MigrateTables函数来包含Note模型
}
