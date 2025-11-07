package main

import (
	"io"
	"log"
	"os"

	"toolcat/services/llm/internal/chat"
	"toolcat/services/llm/internal/config"
	"toolcat/services/llm/internal/server"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	// 初始化日志系统
	// 创建日志文件，权限设置为644
	logFile, err := os.OpenFile("llm.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close() // 确保程序退出前关闭日志文件

	// 配置日志输出格式和级别
	log.SetFlags(log.LstdFlags | log.Lshortfile) // 显示时间和文件名
	logLevel := os.Getenv("LOG_LEVEL")           // 从环境变量获取日志级别
	if logLevel == "debug" {
		// 调试模式下同时输出到控制台和文件
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	} else {
		// 非调试模式只输出到文件
		log.SetOutput(logFile)
	}

	// 初始化LLM连接池
	pool := chat.NewLLMPool(5, func() (llms.LLM, error) {
		// 加载配置信息
		config, _ := config.LoadConfig()
		// 创建新的LLM实例
		return ollama.New(
			ollama.WithModel(config.ModelName),     // 模型名称
			ollama.WithServerURL(config.ServerURL), // 服务器地址
		)
	})

	// 启动HTTP服务器（异步）
	go server.StartWebServer(pool)

	// 初始化聊天实例
	c, err := chat.NewChat(pool)
	if err != nil {
		log.Fatal("failed to create chat: ", err) // 如果初始化失败，记录错误并退出
	}
	defer c.Close() // 确保程序退出前关闭聊天实例

	// 启动聊天主循环
	if err := c.Start(); err != nil {
		log.Fatal("chat ended with error: ", err) // 如果聊天循环异常退出，记录错误
	}
}
