package config

import (
	"encoding/json"
	"os"
)

// 配置结构
type AppConfig struct {
	ModelName   string  `json:"model_name"`  // 模型名称
	ServerURL   string  `json:"server_url"`  // 服务端URL
	Port        int     `json:"port"`        // 服务监听端口
	MaxHistory  int     `json:"max_history"` // 最大历史记录条数
	Temperature float64 `json:"temperature"` // 模型生成温度参数
}

// 加载应用程序配置
// 优先从config.json文件读取配置，如果不存在则使用默认值
// 环境变量可以覆盖配置文件中的值
func LoadConfig() (*AppConfig, error) {
	configFile := "config.json"
	if _, err := os.Stat(configFile); err == nil {
		// 读取配置文件内容
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		// 解析JSON配置
		var config AppConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, err
		}

		// 环境变量覆盖配置
		if model := os.Getenv("MODEL_NAME"); model != "" {
			config.ModelName = model
		}
		if url := os.Getenv("SERVER_URL"); url != "" {
			config.ServerURL = url
		}

		return &config, nil
	}

	// 默认配置
	return &AppConfig{
		ModelName:   "deepseek-r1", // 需与Ollama中模型名称一致
		ServerURL:   "http://localhost:11434",
		Port:        11434,
		MaxHistory:  20,
		Temperature: 0.7,
	}, nil
}
