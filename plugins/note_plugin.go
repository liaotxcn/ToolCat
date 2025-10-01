package plugins

import (
	"errors"
	"fmt"
	"time"

	"toolcat/models"
	"toolcat/pkg"

	"github.com/gin-gonic/gin"
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
	// 使用MySQL数据库存储，不需要内存存储字段
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
	// 注册插件相关路由
	pluginGroup := router.Group(fmt.Sprintf("/plugins/%s", p.Name()))
	{
		// 获取插件信息
		pluginGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"plugin":      p.Name(),
				"description": p.Description(),
				"version":     p.Version(),
				"endpoints": []string{
					"GET /plugins/note/ - 获取插件信息",
					"GET /plugins/note/notes - 获取所有笔记",
					"GET /plugins/note/notes/:id - 获取单个笔记",
					"POST /plugins/note/notes - 创建新笔记",
					"PUT /plugins/note/notes/:id - 更新笔记",
					"DELETE /plugins/note/notes/:id - 删除笔记",
				},
			})
		})

		// 获取所有笔记
		pluginGroup.GET("/notes", func(c *gin.Context) {
			result, _ := p.Execute(map[string]interface{}{
				"action": "list",
			})
			c.JSON(200, result)
		})

		// 获取单个笔记
		pluginGroup.GET("/notes/:id", func(c *gin.Context) {
			id := c.Param("id")
			result, err := p.Execute(map[string]interface{}{
				"action": "get",
				"id":     id,
			})
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, result)
		})

		// 创建新笔记
		pluginGroup.POST("/notes", func(c *gin.Context) {
			var noteData struct {
				Title   string `json:"title" binding:"required"`
				Content string `json:"content" binding:"required"`
			}

			if err := c.ShouldBindJSON(&noteData); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			result, err := p.Execute(map[string]interface{}{
				"action":  "create",
				"title":   noteData.Title,
				"content": noteData.Content,
			})

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.JSON(201, result)
		})

		// 更新笔记
		pluginGroup.PUT("/notes/:id", func(c *gin.Context) {
			id := c.Param("id")

			var noteData struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			}

			if err := c.ShouldBindJSON(&noteData); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			result, err := p.Execute(map[string]interface{}{
				"action":  "update",
				"id":      id,
				"title":   noteData.Title,
				"content": noteData.Content,
			})

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, result)
		})

		// 删除笔记
		pluginGroup.DELETE("/notes/:id", func(c *gin.Context) {
			id := c.Param("id")

			_, err := p.Execute(map[string]interface{}{
				"action": "delete",
				"id":     id,
			})

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, gin.H{"message": "笔记已删除"})
		})
	}
}

// Execute 执行插件功能
func (p *NotePlugin) Execute(params map[string]interface{}) (interface{}, error) {
	action, ok := params["action"].(string)
	if !ok || action == "" {
		return nil, errors.New("缺少必要的 action 参数")
	}

	switch action {
	case "list":
		return p.listNotes()
	case "get":
		id, ok := params["id"].(string)
		if !ok || id == "" {
			return nil, errors.New("缺少必要的 id 参数")
		}
		return p.getNote(id)
	case "create":
		title, ok := params["title"].(string)
		if !ok || title == "" {
			return nil, errors.New("标题不能为空")
		}
		content, ok := params["content"].(string)
		if !ok || content == "" {
			return nil, errors.New("内容不能为空")
		}
		return p.createNote(title, content)
	case "update":
		id, ok := params["id"].(string)
		if !ok || id == "" {
			return nil, errors.New("缺少必要的 id 参数")
		}
		title, _ := params["title"].(string)
		content, _ := params["content"].(string)
		return p.updateNote(id, title, content)
	case "delete":
		id, ok := params["id"].(string)
		if !ok || id == "" {
			return nil, errors.New("缺少必要的 id 参数")
		}
		return p.deleteNote(id)
	default:
		return nil, errors.New("不支持的操作: " + action)
	}
}

// listNotes 获取所有笔记
func (p *NotePlugin) listNotes() (interface{}, error) {
	var notes []models.Note
	if err := pkg.DB.Find(&notes).Error; err != nil {
		return nil, fmt.Errorf("查询笔记失败: %w", err)
	}
	return gin.H{"notes": notes, "total": len(notes)}, nil
}

// getNote 获取单个笔记
func (p *NotePlugin) getNote(id string) (interface{}, error) {
	var note models.Note
	if err := pkg.DB.Where("id = ?", id).First(&note).Error; err != nil {
		return nil, errors.New("未找到指定的笔记")
	}
	return note, nil
}

// createNote 创建新笔记
func (p *NotePlugin) createNote(title, content string) (interface{}, error) {
	// 生成唯一ID
	id := fmt.Sprintf("note-%d", time.Now().UnixNano())
	currentTime := time.Now()

	newNote := models.Note{
		ID:          id,
		Title:       title,
		Content:     content,
		CreatedTime: currentTime,
		UpdatedTime: currentTime,
	}

	if err := pkg.DB.Create(&newNote).Error; err != nil {
		return nil, fmt.Errorf("创建笔记失败: %w", err)
	}

	return newNote, nil
}

// updateNote 更新笔记
func (p *NotePlugin) updateNote(id, title, content string) (interface{}, error) {
	var note models.Note
	if err := pkg.DB.Where("id = ?", id).First(&note).Error; err != nil {
		return nil, errors.New("未找到指定的笔记")
	}

	// 更新笔记内容
	if title != "" {
		note.Title = title
	}
	if content != "" {
		note.Content = content
	}
	note.UpdatedTime = time.Now()

	if err := pkg.DB.Save(&note).Error; err != nil {
		return nil, fmt.Errorf("更新笔记失败: %w", err)
	}

	return note, nil
}

// deleteNote 删除笔记
func (p *NotePlugin) deleteNote(id string) (interface{}, error) {
	var note models.Note
	if err := pkg.DB.Where("id = ?", id).First(&note).Error; err != nil {
		return nil, errors.New("未找到指定的笔记")
	}

	if err := pkg.DB.Delete(&note).Error; err != nil {
		return nil, fmt.Errorf("删除笔记失败: %w", err)
	}

	return gin.H{"message": "删除成功"}, nil
}
