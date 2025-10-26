package controllers

import (
	"net/http"
	"strconv"
	"time"
	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
)

// AuditController 审计日志控制器
type AuditController struct{}

// GetAuditLogs 获取审计日志列表
func (ac *AuditController) GetAuditLogs(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	action := c.Query("action")
	resourceType := c.Query("resource_type")
	username := c.Query("username")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	query := pkg.DB.Model(&models.AuditLog{})

	// 添加租户过滤（多租户隔离）
	tenantID := c.GetUint("tenant_id")
	query = query.Where("tenant_id = ?", tenantID)

	// 添加过滤条件
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if username != "" {
		query = query.Where("username = ?", username)
	}
	if startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}
	if endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			query = query.Where("created_at <= ?", endTime)
		}
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to count audit logs", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 获取分页数据
	var auditLogs []models.AuditLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&auditLogs).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to fetch audit logs", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 计算总页数
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// 返回分页结果
	c.JSON(http.StatusOK, gin.H{
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"logs":        auditLogs,
	})
}

// GetAuditLog 获取单个审计日志详情
func (ac *AuditController) GetAuditLog(c *gin.Context) {
	id := c.Param("id")
	tenantID := c.GetUint("tenant_id")

	var auditLog models.AuditLog
	result := pkg.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&auditLog)
	if result.Error != nil {
		err := pkg.NewNotFoundError("Audit log not found", result.Error)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	c.JSON(http.StatusOK, auditLog)
}

// GetAuditStats 获取审计日志统计信息
func (ac *AuditController) GetAuditStats(c *gin.Context) {
	tenantID := c.GetUint("tenant_id")

	// 按操作类型统计
	type ActionStat struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}

	var actionStats []ActionStat
	if err := pkg.DB.Model(&models.AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("action").
		Find(&actionStats).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to get action stats", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 按资源类型统计
	type ResourceStat struct {
		ResourceType string `json:"resource_type"`
		Count        int64  `json:"count"`
	}

	var resourceStats []ResourceStat
	if err := pkg.DB.Model(&models.AuditLog{}).
		Select("resource_type, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("resource_type").
		Find(&resourceStats).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to get resource stats", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 最近7天的操作统计
	type DailyStat struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	// 计算日期范围
	sevenDaysAgo := time.Now().AddDate(0, 0, -6).Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)

	// 创建日期映射表，初始设置所有日期的计数为0
	var dailyStats []DailyStat
	countMap := make(map[string]int64)
	for i := 0; i < 7; i++ {
		date := sevenDaysAgo.AddDate(0, 0, i).Format("2006-01-02")
		countMap[date] = 0
	}

	// 使用单个SQL查询获取所有日期的统计数据
	// 使用日期函数将created_at转换为日期字符串进行分组
	var results []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	if err := pkg.DB.Model(&models.AuditLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("tenant_id = ? AND created_at >= ? AND created_at < ?", tenantID, sevenDaysAgo, today.Add(24*time.Hour)).
		Group("DATE(created_at)").
		Find(&results).Error; err != nil {
		err := pkg.NewDatabaseError("Failed to get daily stats", err)
		c.JSON(pkg.GetHTTPStatus(err), gin.H{"code": string(err.Code), "message": err.Message})
		return
	}

	// 将结果填充到映射表中
	for _, result := range results {
		countMap[result.Date] = result.Count
	}

	// 将映射表转换为响应数据结构
	for i := 0; i < 7; i++ {
		date := sevenDaysAgo.AddDate(0, 0, i).Format("2006-01-02")
		dailyStats = append(dailyStats, DailyStat{
			Date:  date,
			Count: countMap[date],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"action_stats":   actionStats,
		"resource_stats": resourceStats,
		"daily_stats":    dailyStats,
	})
}
