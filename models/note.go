package models

import (
	"time"
)

// Note 笔记模型
type Note struct {
	ID          string    `gorm:"primaryKey;size:100" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

// AddNoteToMigration 将Note模型添加到迁移函数
func init() {
	// 修改现有的MigrateTables函数来包含Note模型
}
