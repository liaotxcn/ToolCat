package config

import (
	"os"
	"strconv"
)

// Config 应用程序配置结构
var Config = struct {
	// 服务器配置
	Server struct {
		Port int
	}

	// 数据库配置
	Database struct {
		Driver   string
		Host     string
		Port     int
		Username string
		Password string
		DBName   string
		Charset  string
	}

	// 日志配置
	Log struct {
		Level string
		Path  string
	}
}{
	// 默认配置
	Server: struct {
		Port int
	}{
		Port: 8080,
	},

	// MySQL 数据库配置
	Database: struct {
		Driver   string
		Host     string
		Port     int
		Username string
		Password string
		DBName   string
		Charset  string
	}{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		DBName:   "toolcat",
		Charset:  "utf8mb4",
	},
	Log: struct {
		Level string
		Path  string
	}{
		Level: "info",
		Path:  "logs/",
	},
}

// LoadConfig 从环境变量加载配置
func LoadConfig() {
	// 服务器端口
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			Config.Server.Port = p
		}
	}

	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		Config.Database.Host = host
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			Config.Database.Port = p
		}
	}

	if username := os.Getenv("DB_USERNAME"); username != "" {
		Config.Database.Username = username
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		Config.Database.Password = password
	}

	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		Config.Database.DBName = dbname
	}

	// 日志配置
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		Config.Log.Level = logLevel
	}

	if logPath := os.Getenv("LOG_PATH"); logPath != "" {
		Config.Log.Path = logPath
	}
}
