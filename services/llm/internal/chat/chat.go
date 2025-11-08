package chat

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"weave/services/llm/internal/config"

	"github.com/tmc/langchaingo/llms"
)

// 聊天会话核心组件
type Chat struct {
	pool    *LLMPool           // LLM连接池
	ctx     context.Context    // 上下文
	cancel  context.CancelFunc // 取消函数
	reader  *bufio.Reader      // 输入读取器
	history []string           // 对话历史记录
}

// 创建初始化Chat实例
func NewChat(pool *LLMPool) (*Chat, error) {
	// 加载配置(检查配置有效性)
	_, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %v", err)
	}

	// 创建带超时的上下文(30分钟)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	return &Chat{
		pool:    pool,
		ctx:     ctx,
		cancel:  cancel,
		reader:  bufio.NewReader(os.Stdin),     // 从标准输入读取
		history: make([]string, 0, maxHistory), // 初始化历史记录切片
	}, nil
}

// 从连接池获取LLM实例
func (c *Chat) GetLLM() llms.LLM {
	llm, err := c.pool.Get()
	if err != nil {
		log.Printf("Failed to get LLM from pool: %v", err)
		return nil
	}
	return llm
}

const maxHistory = 10 // 最大历史记录条数

// 聊天对话主循环
func (c *Chat) Start() error {
	fmt.Println("Welcome to AI PaiChat!")
	fmt.Println("输入 'help' 查看可用命令")
	for {
		// 获取用户输入
		fmt.Print("\nYou: ")
		input, err := c.reader.ReadString('\n')
		if err != nil {
			log.Printf("读取输入错误: %v", err)
			continue
		}
		input = strings.TrimSpace(input)

		// 处理特殊命令
		switch input {
		case "exit":
			fmt.Println("Thank you for using it, looking forward to our next encounter.")
			return nil
		case "clear":
			c.history = make([]string, 0)
			fmt.Println("Chat history cleared.")
			continue
		case "help":
			fmt.Println("可用命令:")
			fmt.Println("- exit: 退出程序")
			fmt.Println("- clear: 清空历史记录")
			fmt.Println("- history: 查看完整对话历史")
			continue
		case "history":
			fmt.Println("\n对话历史:")
			for _, msg := range c.history {
				fmt.Println(msg)
			}
			continue
		}

		// 构建提示词
		prompt := c.BuildPrompt(input)

		// 获取AI响应
		fmt.Print("\nPaiChat: ")
		var response strings.Builder
		llm := c.GetLLM()
		_, err = llms.GenerateFromSinglePrompt(c.ctx, llm, prompt,
			llms.WithTemperature(0.8),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				fmt.Print(string(chunk))
				response.Write(chunk)
				return nil
			}),
		)
		if err != nil {
			log.Printf("模型调用失败: %v\n提示词内容: %s", err, prompt)
			continue
		}

		// 保存对话历史
		c.addHistory(input, response.String())
		fmt.Println("\n---------------------------------------------")
	}
}

// 关闭聊天会话并释放资源
func (c *Chat) Close() {
	if c.cancel != nil {
		c.cancel() // 取消上下文
	}
}

// 构建完整提示词
func (c *Chat) BuildPrompt(input string) string {
	var prompt strings.Builder
	prompt.WriteString("PaiChat，回答应当:\n")
	prompt.WriteString("- 简洁明了\n- 逻辑清晰\n- 必要时提供示例\n\n")

	for _, msg := range c.history {
		prompt.WriteString(msg)
	}
	prompt.WriteString(fmt.Sprintf("You: %s\nPaiChat: ", input))
	return prompt.String()
}

// 添加对话到历史记录
func (c *Chat) addHistory(userInput, aiResponse string) {
	if len(c.history) >= maxHistory {
		c.history = c.history[1:]
	}
	timestamp := time.Now().Format("2025-05-01 16:05:05")
	c.history = append(c.history,
		fmt.Sprintf("[%s] You: %s\n[%s] PaiChat: %s\n",
			timestamp, userInput,
			timestamp, aiResponse),
	)
}
