package controllers

import (
	"net/http"

	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
)

// TeamController 团队控制器
type TeamController struct{}

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
