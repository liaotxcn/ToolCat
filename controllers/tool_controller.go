package controllers

import (
	"net/http"
	"time"

	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
)

// ToolController 工具控制器
type ToolController struct{}

// GetTools 获取所有工具
func (tc *ToolController) GetTools(c *gin.Context) {
	var tools []models.Tool
	result := pkg.DB.Where("is_enabled = ?", true).Find(&tools)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, tools)
}

// GetTool 获取单个工具
func (tc *ToolController) GetTool(c *gin.Context) {
	id := c.Param("id")

	var tool models.Tool
	result := pkg.DB.First(&tool, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}

	c.JSON(http.StatusOK, tool)
}

// CreateTool 创建工具
func (tc *ToolController) CreateTool(c *gin.Context) {
	var tool models.Tool
	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := pkg.DB.Create(&tool)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, tool)
}

// UpdateTool 更新工具
func (tc *ToolController) UpdateTool(c *gin.Context) {
	id := c.Param("id")

	var tool models.Tool
	result := pkg.DB.First(&tool, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}

	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result = pkg.DB.Save(&tool)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, tool)
}

// DeleteTool 删除工具
func (tc *ToolController) DeleteTool(c *gin.Context) {
	id := c.Param("id")

	result := pkg.DB.Delete(&models.Tool{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tool deleted successfully"})
}

// ExecuteTool 执行工具
func (tc *ToolController) ExecuteTool(c *gin.Context) {
	id := c.Param("id")

	var tool models.Tool
	result := pkg.DB.First(&tool, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tool not found"})
		return
	}

	if !tool.IsEnabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tool is disabled"})
		return
	}

	// 记录工具使用历史
	history := models.ToolHistory{
		ToolID: tool.ID,
		UsedAt: time.Now(),
		// 从请求中获取参数和用户信息
	}

	result = pkg.DB.Create(&history)
	if result.Error != nil {
		// 记录失败不应影响工具执行
		c.JSON(http.StatusOK, gin.H{
			"tool_id": tool.ID,
			"message": "Tool executed successfully, but history recording failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tool_id": tool.ID,
		"message": "Tool executed successfully",
	})
}
