package plugins

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Note 表示一条事件记录
type Note struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}

// NotePlugin 记事本插件
type NotePlugin struct {
	// 使用MySQL数据库存储
	mutex sync.RWMutex // 读写锁用于并发控制
}

// Name 返回插件名称
func (p *NotePlugin) Name() string {
	return "note"
}

// Description 返回插件描述
func (p *NotePlugin) Description() string {
	return "一个记事本插件，可以实现事件记录的增删查改功能"
}

// Version 返回插件版本
func (p *NotePlugin) Version() string {
	return "1.0.0"
}

// Init 初始化插件
func (p *NotePlugin) Init() error {
	// 插件初始化
	fmt.Println("NotePlugin: 记事本插件已初始化")
	return nil
}

// Shutdown 关闭插件
func (p *NotePlugin) Shutdown() error {
	// 插件关闭逻辑
	fmt.Println("NotePlugin: 记事本插件已关闭")
	return nil
}

// RegisterRoutes 注册插件路由
func (p *NotePlugin) RegisterRoutes(router *gin.Engine) {
	// 注册插件相关路由，保持原始路径格式以兼容系统
	pluginGroup := router.Group(fmt.Sprintf("/plugins/%s", p.Name()))
	{
		// 获取插件信息
		pluginGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"plugin":      p.Name(),
				"name":        p.Name(),
				"description": p.Description(),
				"version":     p.Version(),
				"endpoints": []string{
					"GET /plugins/note/ - 获取插件信息",
					"GET /plugins/note/notes - 获取所有笔记（支持分页和用户关联）",
					"GET /plugins/note/notes/:id - 获取单个笔记（用户关联）",
					"POST /plugins/note/notes - 创建新笔记（用户关联）",
					"PUT /plugins/note/notes/:id - 更新笔记（用户关联）",
					"DELETE /plugins/note/notes/:id - 删除笔记（用户关联）",
					"GET /plugins/note/notes/search - 搜索笔记（支持分页和用户关联）",
				},
			})
		})

		// 获取所有笔记（支持分页和用户关联）
		pluginGroup.GET("/notes", func(c *gin.Context) {
			// 获取用户ID，这里简化处理，实际应从认证中获取
			userID := c.DefaultQuery("user_id", "1")
			userIDUint, _ := strconv.ParseUint(userID, 10, 32)

			// 获取分页参数
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

			result, err := p.listNotes(uint(userIDUint), page, pageSize)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, result)
		})

		// 获取单个笔记（用户关联）
		pluginGroup.GET("/notes/:id", func(c *gin.Context) {
			id := c.Param("id")
			// 获取用户ID，这里简化处理，实际应从认证中获取
			userID := c.DefaultQuery("user_id", "1")
			userIDUint, _ := strconv.ParseUint(userID, 10, 32)

			result, err := p.getNote(uint(userIDUint), id)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, result)
		})

		// 创建新笔记（用户关联）
		pluginGroup.POST("/notes", func(c *gin.Context) {
			// 获取用户ID，这里简化处理，实际应从认证中获取
			userID := c.DefaultQuery("user_id", "1")
			userIDUint, _ := strconv.ParseUint(userID, 10, 32)

			var noteData struct {
				Title   string `json:"title" binding:"required"`
				Content string `json:"content" binding:"required"`
			}

			if err := c.ShouldBindJSON(&noteData); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			result, err := p.createNote(uint(userIDUint), noteData.Title, noteData.Content)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(201, result)
		})

		// 更新笔记（用户关联）
		pluginGroup.PUT("/notes/:id", func(c *gin.Context) {
			id := c.Param("id")
			// 获取用户ID，这里简化处理，实际应从认证中获取
			userID := c.DefaultQuery("user_id", "1")
			userIDUint, err := strconv.ParseUint(userID, 10, 32)
			if err != nil {
				log.Printf("Invalid user_id: %v", userID)
				c.JSON(400, gin.H{"error": "无效的用户ID"})
				return
			}

			log.Printf("Update note request: id=%s, user_id=%d", id, userIDUint)

			var noteData struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			}

			if err := c.ShouldBindJSON(&noteData); err != nil {
				log.Printf("Invalid request body: %v", err)
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			result, err := p.updateNote(uint(userIDUint), id, noteData.Title, noteData.Content)
			if err != nil {
				log.Printf("Update note failed: id=%s, user_id=%d, error=%v", id, userIDUint, err)
				if strings.Contains(err.Error(), "未找到") {
					c.JSON(404, gin.H{"error": err.Error()})
				} else {
					c.JSON(500, gin.H{"error": err.Error()})
				}
				return
			}

			c.JSON(200, result)
		})

		// 删除笔记（用户关联）
		pluginGroup.DELETE("/notes/:id", func(c *gin.Context) {
			id := c.Param("id")
			// 获取用户ID，这里简化处理，实际应从认证中获取
			userID := c.DefaultQuery("user_id", "1")
			userIDUint, _ := strconv.ParseUint(userID, 10, 32)

			result, err := p.deleteNote(uint(userIDUint), id)
			if err != nil {
				if strings.Contains(err.Error(), "未找到") {
					c.JSON(404, gin.H{"error": err.Error()})
				} else {
					c.JSON(500, gin.H{"error": err.Error()})
				}
				return
			}

			c.JSON(200, result)
		})

		// 搜索笔记（支持分页和用户关联）
		pluginGroup.GET("/notes/search", func(c *gin.Context) {
			// 获取用户ID，这里简化处理，实际应从认证中获取
			userID := c.DefaultQuery("user_id", "1")
			userIDUint, _ := strconv.ParseUint(userID, 10, 32)

			// 获取搜索关键词
			keyword := c.DefaultQuery("keyword", "")
			if keyword == "" {
				c.JSON(400, gin.H{"error": "搜索关键词不能为空"})
				return
			}

			// 获取分页参数
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

			result, err := p.searchNotes(uint(userIDUint), keyword, page, pageSize)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, result)
		})
	}

	// 同时保留兼容旧API路径格式，确保向前兼容
	api := router.Group("/api")
	{
		api.GET("/plugins/note/notes", func(c *gin.Context) {
			c.Request.URL.Path = "/plugins/note/notes"
			router.HandleContext(c)
		})
		api.GET("/plugins/note/notes/:id", func(c *gin.Context) {
			c.Request.URL.Path = "/plugins/note/notes/" + c.Param("id")
			router.HandleContext(c)
		})
		api.POST("/plugins/note/notes", func(c *gin.Context) {
			c.Request.URL.Path = "/plugins/note/notes"
			router.HandleContext(c)
		})
		api.PUT("/plugins/note/notes/:id", func(c *gin.Context) {
			c.Request.URL.Path = "/plugins/note/notes/" + c.Param("id")
			router.HandleContext(c)
		})
		api.DELETE("/plugins/note/notes/:id", func(c *gin.Context) {
			c.Request.URL.Path = "/plugins/note/notes/" + c.Param("id")
			router.HandleContext(c)
		})
		api.GET("/plugins/note/notes/search", func(c *gin.Context) {
			c.Request.URL.Path = "/plugins/note/notes/search"
			router.HandleContext(c)
		})
	}
}

// Execute 执行插件功能
func (p *NotePlugin) Execute(params map[string]interface{}) (interface{}, error) {
	action, ok := params["action"].(string)
	if !ok || action == "" {
		return nil, errors.New("缺少必要的 action 参数")
	}

	// 获取用户ID
	userID, _ := params["user_id"].(uint)

	switch action {
	case "list":
		// 获取分页参数
		page, _ := params["page"].(int)
		pageSize, _ := params["page_size"].(int)
		return p.listNotes(userID, page, pageSize)
	case "get":
		id, ok := params["id"].(string)
		if !ok || id == "" {
			return nil, errors.New("缺少必要的 id 参数")
		}
		return p.getNote(userID, id)
	case "create":
		title, ok := params["title"].(string)
		if !ok || title == "" {
			return nil, errors.New("标题不能为空")
		}
		content, ok := params["content"].(string)
		if !ok || content == "" {
			return nil, errors.New("内容不能为空")
		}
		return p.createNote(userID, title, content)
	case "update":
		id, ok := params["id"].(string)
		if !ok || id == "" {
			return nil, errors.New("缺少必要的 id 参数")
		}
		title, _ := params["title"].(string)
		content, _ := params["content"].(string)
		return p.updateNote(userID, id, title, content)
	case "delete":
		id, ok := params["id"].(string)
		if !ok || id == "" {
			return nil, errors.New("缺少必要的 id 参数")
		}
		return p.deleteNote(userID, id)
	case "search":
		keyword, ok := params["keyword"].(string)
		if !ok || keyword == "" {
			return nil, errors.New("搜索关键词不能为空")
		}
		page, _ := params["page"].(int)
		pageSize, _ := params["page_size"].(int)
		return p.searchNotes(userID, keyword, page, pageSize)
	default:
		return nil, errors.New("不支持的操作: " + action)
	}
}

// listNotes 获取当前用户的所有笔记
func (p *NotePlugin) listNotes(userID uint, page, pageSize int) (interface{}, error) {
	// 获取读锁
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var notes []models.Note
	var total int64

	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	db := pkg.DB.Where("user_id = ?", userID)

	// 获取总数
	if err := db.Model(&models.Note{}).Count(&total).Error; err != nil {
		log.Printf("Database error when counting notes: %v", err)
		return nil, fmt.Errorf("获取笔记总数失败，请稍后重试")
	}

	// 获取分页数据，按创建时间倒序
	if err := db.Offset(offset).Limit(pageSize).Order("created_time DESC").Find(&notes).Error; err != nil {
		log.Printf("Database error when listing notes: %v", err)
		return nil, fmt.Errorf("获取笔记列表失败，请稍后重试")
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return gin.H{
			"notes":      notes,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
		nil
}

// getNote 获取当前用户的单个笔记
func (p *NotePlugin) getNote(userID uint, id string) (interface{}, error) {
	// 获取读锁
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var note models.Note
	if err := pkg.DB.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		log.Printf("Database error when getting note %s: %v", id, err)
		return nil, errors.New("未找到指定的笔记或您没有权限访问")
	}
	return note, nil
}

// createNote 为当前用户创建新笔记
func (p *NotePlugin) createNote(userID uint, title, content string) (interface{}, error) {
	// 获取写锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 使用UUID生成唯一ID，更可靠
	id := "note-" + uuid.New().String()
	currentTime := time.Now()

	// 开始事务
	tx := pkg.DB.Begin()
	if tx.Error != nil {
		log.Printf("Failed to begin transaction: %v", tx.Error)
		return nil, fmt.Errorf("创建笔记失败，请稍后重试")
	}

	newNote := models.Note{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Content:     content,
		CreatedTime: currentTime,
		UpdatedTime: currentTime,
	}

	if err := tx.Create(&newNote).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error when creating note: %v", err)
		return nil, fmt.Errorf("创建笔记失败，请稍后重试")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("创建笔记失败，请稍后重试")
	}

	return newNote, nil
}

// updateNote 更新当前用户的笔记
func (p *NotePlugin) updateNote(userID uint, id, title, content string) (interface{}, error) {
	// 获取写锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 开始事务
	tx := pkg.DB.Begin()
	if tx.Error != nil {
		log.Printf("Failed to begin transaction: %v", tx.Error)
		return nil, fmt.Errorf("更新笔记失败，请稍后重试")
	}

	var note models.Note
	if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error when finding note %s: %v", id, err)
		return nil, errors.New("未找到指定的笔记或您没有权限访问")
	}

	// 更新笔记内容
	if title != "" {
		note.Title = title
	}
	if content != "" {
		note.Content = content
	}
	note.UpdatedTime = time.Now()

	if err := tx.Save(&note).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error when updating note %s: %v", id, err)
		return nil, fmt.Errorf("更新笔记失败，请稍后重试")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("更新笔记失败，请稍后重试")
	}

	return note, nil
}

// deleteNote 删除当前用户的笔记
func (p *NotePlugin) deleteNote(userID uint, id string) (interface{}, error) {
	// 获取写锁
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 开始事务
	tx := pkg.DB.Begin()
	if tx.Error != nil {
		log.Printf("Failed to begin transaction: %v", tx.Error)
		return nil, fmt.Errorf("删除笔记失败，请稍后重试")
	}

	var note models.Note
	if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error when finding note %s: %v", id, err)
		return nil, errors.New("未找到指定的笔记或您没有权限访问")
	}

	if err := tx.Delete(&note).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error when deleting note %s: %v", id, err)
		return nil, fmt.Errorf("删除笔记失败，请稍后重试")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("删除笔记失败，请稍后重试")
	}

	return gin.H{"message": "删除成功"}, nil
}

// searchNotes 搜索当前用户的笔记
func (p *NotePlugin) searchNotes(userID uint, keyword string, page, pageSize int) (interface{}, error) {
	// 获取读锁
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	var notes []models.Note
	var total int64

	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建搜索条件
	query := "%" + keyword + "%"
	db := pkg.DB.Where("user_id = ? AND (title LIKE ? OR content LIKE ?)", userID, query, query)

	// 获取总数
	if err := db.Model(&models.Note{}).Count(&total).Error; err != nil {
		log.Printf("Database error when counting search results: %v", err)
		return nil, fmt.Errorf("搜索笔记失败，请稍后重试")
	}

	// 获取分页数据，按创建时间倒序
	if err := db.Offset(offset).Limit(pageSize).Order("created_time DESC").Find(&notes).Error; err != nil {
		log.Printf("Database error when searching notes: %v", err)
		return nil, fmt.Errorf("搜索笔记失败，请稍后重试")
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return gin.H{
			"notes":      notes,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
		nil
}
