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
		// 使用统一错误码系统返回数据库错误
		dbErr := pkg.NewDatabaseError("Failed to fetch tools", result.Error)
		dbErr.WithDetails(map[string]interface{}{
			"filter": "is_enabled = true",
		})
		c.Error(dbErr)
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
		// 使用统一错误码系统返回未找到错误
		c.Error(pkg.NewNotFoundError("Tool not found", nil))
		return
	}

	c.JSON(http.StatusOK, tool)
}

// CreateTool 创建工具
func (tc *ToolController) CreateTool(c *gin.Context) {
	var tool models.Tool
	if err := c.ShouldBindJSON(&tool); err != nil {
		// 使用统一错误码系统返回参数验证错误
		c.Error(pkg.NewValidationError("Invalid tool data", err))
		return
	}

	result := pkg.DB.Create(&tool)
	if result.Error != nil {
		// 使用统一错误码系统返回数据库错误
		dbErr := pkg.NewDatabaseError("Failed to create tool", result.Error)
		dbErr.WithDetails(map[string]interface{}{
			"tool_name": tool.Name,
		})
		c.Error(dbErr)
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
		// 使用统一错误码系统返回未找到错误
		c.Error(pkg.NewNotFoundError("Tool not found", nil))
		return
	}

	if err := c.ShouldBindJSON(&tool); err != nil {
		// 使用统一错误码系统返回参数验证错误
		c.Error(pkg.NewValidationError("Invalid tool data", err))
		return
	}

	result = pkg.DB.Save(&tool)
	if result.Error != nil {
		// 使用统一错误码系统返回数据库错误
		dbErr := pkg.NewDatabaseError("Failed to update tool", result.Error)
		dbErr.WithDetails(map[string]interface{}{
			"tool_id": id,
		})
		c.Error(dbErr)
		return
	}

	c.JSON(http.StatusOK, tool)
}

// DeleteTool 删除工具
func (tc *ToolController) DeleteTool(c *gin.Context) {
	id := c.Param("id")

	result := pkg.DB.Delete(&models.Tool{}, id)
	if result.Error != nil {
		// 使用统一错误码系统返回数据库错误
		dbErr := pkg.NewDatabaseError("Failed to delete tool", result.Error)
		dbErr.WithDetails(map[string]interface{}{
			"tool_id": id,
		})
		c.Error(dbErr)
		return
	}

	if result.RowsAffected == 0 {
		// 使用统一错误码系统返回未找到错误
		c.Error(pkg.NewNotFoundError("Tool not found", nil))
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
		// 使用统一错误码系统返回未找到错误
		c.Error(pkg.NewNotFoundError("Tool not found", nil))
		return
	}

	if !tool.IsEnabled {
		// 使用统一错误码系统返回禁止访问错误
		c.Error(pkg.NewForbiddenError("Tool is disabled", nil))
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
		// 记录失败不应影响工具执行，但记录错误日志
		dbErr := pkg.NewDatabaseError("Failed to record tool usage history", result.Error)
		dbErr.WithDetails(map[string]interface{}{
			"tool_id": tool.ID,
		})
		// 记录但不中断执行
		c.Error(dbErr)
	}

	c.JSON(http.StatusOK, gin.H{
		"tool_id": tool.ID,
		"message": "Tool executed successfully",
	})
}
