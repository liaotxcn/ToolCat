package models

import "time"

// Team 团队模型
// 用于在租户内组织和管理成员
type Team struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;index:idx_tenant_team_name,unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	OwnerID     uint      `gorm:"index" json:"owner_id"`
	TenantID    uint      `gorm:"index:idx_tenant_team_name,unique" json:"tenant_id"`
	Members     string    `gorm:"type:text;default:null;comment:团队成员列表（用户名形式）" json:"members"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TeamMember 团队成员模型
// 记录用户在团队内的角色
type TeamMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TeamID    uint      `gorm:"index" json:"team_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Role      string    `gorm:"size:50;default:member" json:"role"`
	TenantID  uint      `gorm:"index" json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
}
