package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
)

// Generator 定义生成接口
type Generator interface {
	Generate(ctx context.Context, query string, documents []*schema.Document) (string, error)
}

// ArkGenerator 实现基于ARK的生成器
type ArkGenerator struct {
	baseURL   string
	apiKey    string
	modelName string
}

// NewArkGenerator 创建ARK生成器
func NewArkGenerator(baseURL, apiKey, modelName string) *ArkGenerator {
	return &ArkGenerator{
		baseURL:   baseURL,
		apiKey:    apiKey,
		modelName: modelName,
	}
}

// 请求结构体
type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 响应结构体
type chatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

// API错误响应结构体
type apiErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Generate 生成回答
func (g *ArkGenerator) Generate(ctx context.Context, query string, documents []*schema.Document) (string, error) {
	fmt.Printf("[%s] 开始处理查询: %s\n", time.Now().Format("2006-01-02 15:04:05"), query)
	fmt.Printf("[%s] 使用模型: %s, 基础URL: %s\n", time.Now().Format("2006-01-02 15:04:05"), g.modelName, g.baseURL)

	// 验证API密钥
	if g.apiKey == "" || strings.TrimSpace(g.apiKey) == "" {
		errMsg := "API密钥为空或无效"
		fmt.Printf("[%s] 错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 组合上下文信息
	context := ""
	if len(documents) > 0 {
		contextParts := make([]string, len(documents))
		for i, doc := range documents {
			// 如果元数据中有标题，添加标题信息
			titleInfo := ""
			if title, ok := doc.MetaData["title"].(string); ok && title != "" {
				titleInfo = fmt.Sprintf("标题: %s\n", title)
			}
			contextParts[i] = fmt.Sprintf("文档片段[%d]:\n%s%s\n", i+1, titleInfo, doc.Content)
		}
		context = strings.Join(contextParts, "\n---\n")
		fmt.Printf("[%s] 检索到 %d 个文档用于上下文\n", time.Now().Format("2006-01-02 15:04:05"), len(documents))
	} else {
		fmt.Printf("[%s] 未检索到相关文档\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	// 构建提示
	systemPrompt := "你是一个知识助手。基于提供的文档回答用户问题。如果文档中没有相关信息，请诚实地表明你不知道，不要编造答案。"
	userPrompt := query

	if context != "" {
		userPrompt = fmt.Sprintf("基于以下信息回答我的问题：\n\n%s\n\n问题：%s", context, query)
	}

	// 构建请求
	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	reqBody := chatRequest{
		Model:    g.modelName,
		Messages: messages,
	}

	// 序列化请求体
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		errMsg := fmt.Sprintf("序列化请求失败: %v", err)
		fmt.Printf("[%s] 错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 创建HTTP请求
	endpoint := fmt.Sprintf("%s/chat/completions", g.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		errMsg := fmt.Sprintf("创建HTTP请求失败: %v", err)
		fmt.Printf("[%s] 错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 添加头信息
	req.Header.Set("Content-Type", "application/json")
	// 隐藏API密钥的中间部分，只显示前后各8个字符
	maskedAPIKey := maskAPIKey(g.apiKey)
	fmt.Printf("[%s] 使用API密钥: %s\n", time.Now().Format("2006-01-02 15:04:05"), maskedAPIKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.apiKey))

	// 设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	fmt.Printf("[%s] 发送请求到: %s\n", time.Now().Format("2006-01-02 15:04:05"), endpoint)
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("发送请求失败: %v", err)
		fmt.Printf("[%s] 错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
		return "", fmt.Errorf(errMsg)
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)
	fmt.Printf("[%s] 收到响应，状态码: %d, 耗时: %v\n", time.Now().Format("2006-01-02 15:04:05"), resp.StatusCode, duration)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Sprintf("读取响应失败: %v", err)
		fmt.Printf("[%s] 错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		// 尝试解析错误响应
		var errorResp apiErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			// 成功解析到错误信息
			errorCode := errorResp.Error.Code
			errorMessage := errorResp.Error.Message

			// 针对认证错误提供更详细的提示
			if strings.Contains(strings.ToLower(errorCode), "auth") || strings.Contains(strings.ToLower(errorMessage), "auth") {
				errMsg := fmt.Sprintf("API认证失败: %s (错误码: %s) - 请检查ARK_API_KEY是否正确且有效", errorMessage, errorCode)
				fmt.Printf("[%s] 认证错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
				return "", fmt.Errorf(errMsg)
			}

			errMsg := fmt.Sprintf("API返回错误: %s (错误码: %s), 状态码: %d", errorMessage, errorCode, resp.StatusCode)
			fmt.Printf("[%s] API错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
			return "", fmt.Errorf(errMsg)
		} else {
			// 无法解析错误响应，返回原始错误
			errMsg := fmt.Sprintf("API返回错误, 状态码: %d, 响应内容: %s", resp.StatusCode, string(body))
			fmt.Printf("[%s] API错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
			return "", fmt.Errorf(errMsg)
		}
	}

	// 解析响应
	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		errMsg := fmt.Sprintf("解析响应失败: %v, 响应内容: %s", err, string(body))
		fmt.Printf("[%s] 错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 提取回答
	if len(chatResp.Choices) > 0 {
		answer := chatResp.Choices[0].Message.Content
		fmt.Printf("[%s] 成功生成回答 (长度: %d 字符)\n", time.Now().Format("2006-01-02 15:04:05"), len(answer))
		return answer, nil
	}

	errMsg := "API没有返回有效回答"
	fmt.Printf("[%s] 错误: %s\n", time.Now().Format("2006-01-02 15:04:05"), errMsg)
	return "", fmt.Errorf(errMsg)
}

// maskAPIKey 隐藏API密钥的中间部分
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 16 {
		return "***"
	}
	return apiKey[:8] + "..." + apiKey[len(apiKey)-8:]
}
