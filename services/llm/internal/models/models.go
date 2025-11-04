package models

import "time"

// API请求结构
type ChatRequest struct {
	Message string `json:"message"`
}

// API响应结构
type ChatResponse struct {
	Response string `json:"response"`
	Status   int    `json:"status"`
	Error    string `json:"error,omitempty"`
}

// 对话历史记录结构
type ChatHistory struct {
	Timestamp  time.Time `json:"timestamp"`
	UserInput  string    `json:"user_input"`
	AIResponse string    `json:"ai_response"`
}

// Config 可以替代原有的AppConfig
type Config struct {
	ModelName   string  `json:"model_name"`
	ServerURL   string  `json:"server_url"`
	Port        int     `json:"port"`
	MaxHistory  int     `json:"max_history"`
	Temperature float64 `json:"temperature"`
}
