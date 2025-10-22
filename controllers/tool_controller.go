package controllers

import (
	"net/http"

	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
)

// ToolController 工具控制器
type ToolController struct{}

// GetTools 获取所有工具
func (tc *ToolController) GetTools(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")
	var tools []models.Tool
	result := pkg.DB.Where("tenant_id = ?", tenantID).Find(&tools)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to fetch tools", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// GetTool 获取单个工具
func (tc *ToolController) GetTool(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	var tool models.Tool
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&tool)
	if result.Error != nil {
		err := pkg.NewNotFoundError("Tool not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, tool)
}

// CreateTool 创建工具
func (tc *ToolController) CreateTool(c *gin.Context) {
	var tool models.Tool
	if err := c.ShouldBindJSON(&tool); err != nil {
		err := pkg.NewValidationError("Invalid tool data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	tool.TenantID = c.GetUint("tenant_id")

	result := pkg.DB.Create(&tool)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to create tool", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusCreated, tool)
}

// UpdateTool 更新工具
func (tc *ToolController) UpdateTool(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	var oldTool models.Tool
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&oldTool)
	if result.Error != nil {
		err := pkg.NewNotFoundError("Tool not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	var newTool models.Tool
	if err := c.ShouldBindJSON(&newTool); err != nil {
		err := pkg.NewValidationError("Invalid tool data", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	newTool.ID = oldTool.ID
	newTool.TenantID = tenantID

	result = pkg.DB.Save(&newTool)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to update tool", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, newTool)
}

// DeleteTool 删除工具
func (tc *ToolController) DeleteTool(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	var tool models.Tool
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&tool)
	if result.Error != nil {
		err := pkg.NewNotFoundError("Tool not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	result = pkg.DB.Delete(&tool)
	if result.Error != nil {
		err := pkg.NewDatabaseError("Failed to delete tool", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tool deleted successfully"})
}

// ExecuteTool 执行工具
func (tc *ToolController) ExecuteTool(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	var tool models.Tool
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&tool)
	if result.Error != nil {
		err := pkg.NewNotFoundError("Tool not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 执行逻辑保持不变
	c.JSON(http.StatusOK, gin.H{"message": "Tool execution started", "tool": tool})
}
