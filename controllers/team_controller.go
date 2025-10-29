package controllers

import (
	"net/http"
	"strconv"

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

// GetTeamMembers 获取团队成员列表
func (tc *TeamController) GetTeamMembers(c *gin.Context) {
	// 获取团队ID
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		err := pkg.NewValidationError("Invalid team ID", err)
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

	// 检查用户是否为团队成员
	var teamMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ?", teamID, userID).First(&teamMember).Error; err != nil {
		err := pkg.NewForbiddenError("You are not a member of this team", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 查询团队成员列表
	var members []models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND tenant_id = ?", teamID, tenantID).Find(&members).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to query team members", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, members)
}

// AddTeamMember 添加团队成员
func (tc *TeamController) AddTeamMember(c *gin.Context) {
	// 获取团队ID
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		err := pkg.NewValidationError("Invalid team ID", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 解析请求参数
	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required,oneof=admin member"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		err := pkg.NewValidationError("Invalid member data", err)
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

	// 检查权限：只有团队所有者或管理员可以添加成员
	var currentMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ? AND role IN ('owner', 'admin')", teamID, userID).First(&currentMember).Error; err != nil {
		err := pkg.NewForbiddenError("Only team owners or admins can add members", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 检查用户是否已存在
	var existingMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ?", teamID, req.UserID).First(&existingMember).Error; err == nil {
		err := pkg.NewConflictError("User is already a member of the team", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 添加新成员
	newMember := models.TeamMember{
		TeamID:   uint(teamID),
		UserID:   req.UserID,
		Role:     req.Role,
		TenantID: tenantID,
	}

	if err := pkg.DB.Create(&newMember).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to add team member", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录审计日志
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "add_member",
		ResourceType: "team",
		ResourceID:   team.Name,
		NewValue:     newMember,
	})

	c.JSON(http.StatusCreated, newMember)
}

// RemoveTeamMember 移除团队成员
func (tc *TeamController) RemoveTeamMember(c *gin.Context) {
	// 获取团队ID和成员ID
	teamIDStr := c.Param("id")
	memberIDStr := c.Param("memberId")
	
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		err := pkg.NewValidationError("Invalid team ID", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}
	
	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		err := pkg.NewValidationError("Invalid member ID", err)
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

	// 查找要移除的成员
	var teamMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ? AND tenant_id = ?", teamID, memberID, tenantID).First(&teamMember).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err := pkg.NewNotFoundError("Team member not found", nil)
			c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		} else {
			err := pkg.NewDatabaseError("Failed to query team member", err)
			c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		}
		return
	}

	// 检查权限：只有团队所有者或管理员可以移除成员，且不能移除所有者
	if teamMember.Role == "owner" {
		err := pkg.NewForbiddenError("Cannot remove team owner", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}
	
	var currentMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ? AND role IN ('owner', 'admin')", teamID, userID).First(&currentMember).Error; err != nil {
		err := pkg.NewForbiddenError("Only team owners or admins can remove members", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 移除成员
	if err := pkg.DB.Delete(&teamMember).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to remove team member", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录审计日志
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "remove_member",
		ResourceType: "team",
		ResourceID:   team.Name,
		OldValue:     teamMember,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Team member removed successfully"})
}

// UpdateMemberRole 更新团队成员角色
func (tc *TeamController) UpdateMemberRole(c *gin.Context) {
	// 获取团队ID和成员ID
	teamIDStr := c.Param("id")
	memberIDStr := c.Param("memberId")
	
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		err := pkg.NewValidationError("Invalid team ID", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}
	
	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		err := pkg.NewValidationError("Invalid member ID", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 解析请求参数
	var req struct {
		Role string `json:"role" binding:"required,oneof=admin member"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		err := pkg.NewValidationError("Invalid role data", err)
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

	// 查找要更新的成员
	var teamMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ? AND tenant_id = ?", teamID, memberID, tenantID).First(&teamMember).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err := pkg.NewNotFoundError("Team member not found", nil)
			c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		} else {
			err := pkg.NewDatabaseError("Failed to query team member", err)
			c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		}
		return
	}

	// 检查权限：只有团队所有者可以更新成员角色
	var currentMember models.TeamMember
	if err := pkg.DB.Where("team_id = ? AND user_id = ? AND role = 'owner'", teamID, userID).First(&currentMember).Error; err != nil {
		err := pkg.NewForbiddenError("Only team owners can update member roles", nil)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 保存旧角色用于审计日志
	oldRole := teamMember.Role
	
	// 更新角色
	teamMember.Role = req.Role
	if err := pkg.DB.Save(&teamMember).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to update member role", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 记录审计日志
	_ = pkg.AuditLogFromContext(c, pkg.AuditLogOptions{
		Action:       "update_member_role",
		ResourceType: "team",
		ResourceID:   team.Name,
		OldValue:     map[string]interface{}{"user_id": teamMember.UserID, "role": oldRole},
		NewValue:     map[string]interface{}{"user_id": teamMember.UserID, "role": req.Role},
	})

	c.JSON(http.StatusOK, teamMember)
}
