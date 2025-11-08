package llm

import (
	"context"
	"sync"
	"time"
	"weave/pkg"
	"weave/plugins/core"
	"weave/services/llm/internal/chat"
	"weave/services/llm/internal/config"
	"weave/services/llm/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// LLMChatPlugin 实现LLM聊天功能的插件
type LLMChatPlugin struct {
	pool    *chat.LLMPool
	manager *core.PluginManager
}

// NewLLMChatPlugin 创建LLM聊天插件实例
func NewLLMChatPlugin() *LLMChatPlugin {
	return &LLMChatPlugin{}
}

// 基础信息接口实现
func (p *LLMChatPlugin) Name() string {
	return "LLMChat"
}

func (p *LLMChatPlugin) Description() string {
	return "提供LLM聊天和对话历史管理功能"
}

func (p *LLMChatPlugin) Version() string {
	return "1.0.0"
}

func (p *LLMChatPlugin) GetDependencies() []string {
	return []string{}
}

func (p *LLMChatPlugin) GetConflicts() []string {
	return []string{}
}

// 生命周期接口实现
func (p *LLMChatPlugin) Init() error {
	pkg.Info("Initializing LLM Chat Plugin...")

	// 初始化LLM连接池
	p.pool = chat.NewLLMPool(5, func() (llms.LLM, error) {
		// 加载配置信息
		appConfig, err := config.LoadConfig()
		if err != nil {
			pkg.Warn("Failed to load LLM config, using defaults")
			// 使用默认配置
			appConfig = &config.AppConfig{
				ModelName:   "deepseek-r1",
				ServerURL:   "http://localhost:11434",
				Port:        11434,
				MaxHistory:  20,
				Temperature: 0.7,
			}
		}

		// 创建新的LLM实例
		return ollama.New(
			ollama.WithModel(appConfig.ModelName),
			ollama.WithServerURL(appConfig.ServerURL),
		)
	})

	pkg.Info("LLM Chat Plugin initialized successfully")
	return nil
}

func (p *LLMChatPlugin) Shutdown() error {
	pkg.Info("Shutting down LLM Chat Plugin...")
	// 这里可以添加清理资源的代码
	return nil
}

func (p *LLMChatPlugin) OnEnable() error {
	pkg.Info("LLM Chat Plugin enabled")
	return nil
}

func (p *LLMChatPlugin) OnDisable() error {
	pkg.Info("LLM Chat Plugin disabled")
	return nil
}

// 路由注册接口实现
func (p *LLMChatPlugin) GetRoutes() []core.Route {
	return []core.Route{
		{
			Path:         "api/chat",
			Method:       "POST",
			Handler:      p.handleChat,
			Description:  "发送聊天消息",
			AuthRequired: false,
			Tags:         []string{"LLM", "Chat"},
		},
		{
			Path:         "api/history",
			Method:       "GET",
			Handler:      p.handleGetHistory,
			Description:  "获取对话历史",
			AuthRequired: false,
			Tags:         []string{"LLM", "History"},
		},
		{
			Path:         "api/clear-history",
			Method:       "POST",
			Handler:      p.handleClearHistory,
			Description:  "清空对话历史",
			AuthRequired: false,
			Tags:         []string{"LLM", "History"},
		},
	}
}

func (p *LLMChatPlugin) RegisterRoutes(router *gin.Engine) {
	// 为了兼容性保留旧接口实现
	// 实际上路由会通过GetRoutes方法获取并由PluginManager注册
	routes := p.GetRoutes()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handler)
	}
}

// 执行功能接口实现
func (p *LLMChatPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	// 可以实现一些通用功能供其他插件调用
	return nil, nil
}

func (p *LLMChatPlugin) GetDefaultMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

func (p *LLMChatPlugin) SetPluginManager(manager *core.PluginManager) {
	p.manager = manager
}

// HTTP处理函数
func (p *LLMChatPlugin) handleChat(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	chat, err := chat.NewChat(p.pool)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer chat.Close()

	prompt := chat.BuildPrompt(req.Message)
	response, err := chat.GetLLM().Call(context.Background(), prompt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 保存对话历史
	historyRepo := NewMemoryRepository()
	history := &models.ChatHistory{
		UserInput:  req.Message,
		AIResponse: response,
		Timestamp:  time.Now(),
	}
	historyRepo.SaveHistory(history)

	c.JSON(200, gin.H{
		"response": response,
	})
}

func (p *LLMChatPlugin) handleGetHistory(c *gin.Context) {
	// 获取LLM实例
	llm, err := p.pool.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get LLM instance"})
		return
	}
	defer p.pool.Put(llm)

	// 创建chat服务实例
	repo := NewMemoryRepository()
	chatService := chat.NewChatService(repo, llm)

	// 获取历史记录
	histories, err := chatService.GetHistory(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{
		"success": true,
		"data":    histories,
		"error":   "",
	})
}

func (p *LLMChatPlugin) handleClearHistory(c *gin.Context) {
	// 获取LLM实例
	llm, err := p.pool.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get LLM instance"})
		return
	}
	defer p.pool.Put(llm)

	// 创建chat服务实例
	repo := NewMemoryRepository()
	chatService := chat.NewChatService(repo, llm)

	// 清空历史记录
	if err := chatService.ClearHistory(context.Background()); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{
		"success": true,
		"data":    nil,
		"error":   "",
	})
}

// 内存存储库实现
var (
	memoryHistory []*models.ChatHistory
	memoryMutex   sync.RWMutex
)

// NewMemoryRepository 创建内存存储库实例
func NewMemoryRepository() chat.ChatRepository {
	return &memoryRepo{}
}

type memoryRepo struct{}

func (r *memoryRepo) SaveHistory(history *models.ChatHistory) error {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	memoryHistory = append(memoryHistory, history)
	return nil
}

func (r *memoryRepo) GetHistories() ([]*models.ChatHistory, error) {
	memoryMutex.RLock()
	defer memoryMutex.RUnlock()
	// 返回历史记录的副本
	histories := make([]*models.ChatHistory, len(memoryHistory))
	copy(histories, memoryHistory)
	return histories, nil
}

func (r *memoryRepo) ClearHistories() error {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	memoryHistory = make([]*models.ChatHistory, 0)
	return nil
}
