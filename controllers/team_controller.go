package controllers

import (
	"net/http"

	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TeamController 团队控制器
type TeamController struct{}

// UpdateTeam 更新团队信息
func (tc *TeamController) UpdateTeam(c *gin.Context) {
	// 解析请求参数
	var req struct {
		Name        string `json:"name" binding:"omitempty,min=2,max=100"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		err := pkg.NewValidationError("Invalid team data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 获取团队ID
	teamID := c.Param("id")
	if teamID == "" {
		err := pkg.NewValidationError("Team ID is required", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 获取当前用户信息
	userID := c.GetUint("user_id")
	tenantID := c.GetUint("tenant_id")

	// 查找团队
	var team models.Team
	if err := pkg.DB.Where("id = ? AND tenant_id = ?", teamID, tenantID).First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err := pkg.NewNotFoundError("Team not found", nil)
			c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		} else {
			err := pkg.NewDatabaseError("Failed to query team", err)
			c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		}
		return
	}

	// 检查权限：只有团队所有者可以更新团队信息
	var teamMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ? AND role = 'owner'", team.ID, userID).First(&teamMember).Error; err != nil {
		err := pkg.NewForbiddenError("Only team owners can update team information", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 如果要更新名称，检查名称是否已存在
	if req.Name != "" && req.Name != team.Name {
		var existingTeam models.Team
		if err := pkg.DB.Where("name = ? AND tenant_id = ? AND id != ?", req.Name, tenantID, teamID).First(&existingTeam).Error; err == nil {
			err := pkg.NewConflictError("Team name already exists", nil)
			c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
			return
		}
		team.Name = req.Name
	}

	// 更新描述
	if req.Description != team.Description {
		team.Description = req.Description
	}

	// 保存更新
	if err := pkg.DB.Save(&team).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to update team", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录审计日志
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "update",
		ResourceType: "team",
		ResourceID:   team.Name,
		NewValue:     team,
	})

	c.JSON(http.StatusOK, team)
}

// CreateTeam 创建团队
func (tc *TeamController) CreateTeam(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required,min=2,max=100"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		err := pkg.NewValidationError("Invalid team data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	tenantID := c.GetUint("tenant_id")
	ownerID := c.GetUint("user_id")

	// 防止重复名称（租户内唯一）
	var existing models.Team
	if err := pkg.DB.Where("name = ? AND tenant_id = ?", req.Name, tenantID).First(&existing).Error; err == nil {
		err := pkg.NewConflictError("Team name already exists", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	team := models.Team{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		TenantID:    tenantID,
	}
	if err := pkg.DB.Create(&team).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to create team", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 将创建者加入团队成员，角色为owner
	_ = pkg.DB.Create(&models.TeamMember{TeamID: team.ID, UserID: ownerID, Role: "owner", TenantID: tenantID}).Error

	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "create",
		ResourceType: "team",
		ResourceID:   team.Name,
		NewValue:     team,
	})

	c.JSON(http.StatusCreated, team)
}
