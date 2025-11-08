package chat

import (
	"context"
	"fmt"
	"strings"
	"time"
	"weave/services/llm/internal/models"

	"github.com/tmc/langchaingo/llms"
)

// 核心接口
type ChatService interface {
	// 发送消息并获取响应
	SendMessage(ctx context.Context, req *models.ChatRequest) (*models.ChatResponse, error)
	// 获取对话历史记录
	GetHistory(ctx context.Context) ([]*models.ChatHistory, error)
	// 清空对话历史记录
	ClearHistory(ctx context.Context) error
}

// 数据存储层的接口
type ChatRepository interface {
	// 保存单条对话记录
	SaveHistory(history *models.ChatHistory) error
	// 获取所有对话记录
	GetHistories() ([]*models.ChatHistory, error)
	// 清空所有对话记录
	ClearHistories() error
}

// ChatService接口
type chatService struct {
	repo ChatRepository // 数据存储层实例
	llm  llms.LLM       // 语言模型实例
}

// 创建新的聊天服务实例
func NewChatService(repo ChatRepository, llm llms.LLM) ChatService {
	return &chatService{repo: repo, llm: llm}
}

// 处理用户消息并返回AI响应
func (s *chatService) SendMessage(ctx context.Context, req *models.ChatRequest) (*models.ChatResponse, error) {
	// 验证消息内容非空
	if strings.TrimSpace(req.Message) == "" {
		return nil, fmt.Errorf("消息内容不能为空")
	}

	// 构建提示词
	prompt := fmt.Sprintf("You: %s\nPaiChat: ", req.Message)

	// 调用语言模型获取响应
	response, err := s.llm.Call(ctx, prompt)
	if err != nil {
		return &models.ChatResponse{
			Status: 500,
			Error:  err.Error(),
		}, err
	}

	// 保存对话记录
	history := &models.ChatHistory{
		Timestamp:  time.Now(),
		UserInput:  req.Message,
		AIResponse: response,
	}
	if err := s.repo.SaveHistory(history); err != nil {
		return nil, err
	}

	// 返回成功响应
	return &models.ChatResponse{
		Response: response,
		Status:   200,
	}, nil
}

// 获取所有对话记录
func (s *chatService) GetHistory(ctx context.Context) ([]*models.ChatHistory, error) {
	histories, err := s.repo.GetHistories()
	if err != nil {
		return nil, fmt.Errorf("获取历史记录失败: %v", err)
	}
	return histories, nil
}

// 清空所有对话记录
func (s *chatService) ClearHistory(ctx context.Context) error {
	if err := s.repo.ClearHistories(); err != nil {
		return fmt.Errorf("清空历史记录失败: %v", err)
	}
	return nil
}
