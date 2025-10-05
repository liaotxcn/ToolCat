package config

import (
	"os"
	"strconv"
)

// Config 应用程序配置结构
var Config struct {
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
	Logger struct {
		Level       string
		OutputPath  string
		ErrorPath   string
		Development bool
	}

	// JWT配置
	JWT struct {
		Secret           string
		AccessTokenExpiry int // 访问令牌过期时间（分钟）
		RefreshTokenExpiry int // 刷新令牌过期时间（小时）
	}
}

// 初始化默认值
func init() {
	// 服务器配置
	Config.Server.Port = 8081

	// 数据库配置
	Config.Database.Driver = "mysql"
	Config.Database.Host = "localhost"
	Config.Database.Port = 3306
	Config.Database.Username = "root"
	Config.Database.Password = "123456"
	Config.Database.DBName = "toolcat"
	Config.Database.Charset = "utf8mb4"

	// 日志配置
	Config.Logger.Level = "info"
	Config.Logger.OutputPath = "stdout"
	Config.Logger.ErrorPath = "stderr"
	Config.Logger.Development = false

	// JWT配置
	Config.JWT.Secret = "your-secret-key"
	Config.JWT.AccessTokenExpiry = 60 // 60分钟
	Config.JWT.RefreshTokenExpiry = 24 * 7 // 7天
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
		Config.Logger.Level = logLevel
	}

	if logOutputPath := os.Getenv("LOG_OUTPUT_PATH"); logOutputPath != "" {
		Config.Logger.OutputPath = logOutputPath
	}

	if logErrorPath := os.Getenv("LOG_ERROR_PATH"); logErrorPath != "" {
		Config.Logger.ErrorPath = logErrorPath
	}

	// JWT配置
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		Config.JWT.Secret = jwtSecret
	}

	if accessTokenExpiry := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY"); accessTokenExpiry != "" {
		if expiry, err := strconv.Atoi(accessTokenExpiry); err == nil {
			Config.JWT.AccessTokenExpiry = expiry
		}
	}

	if refreshTokenExpiry := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"); refreshTokenExpiry != "" {
		if expiry, err := strconv.Atoi(refreshTokenExpiry); err == nil {
			Config.JWT.RefreshTokenExpiry = expiry
		}
	}

	if devMode := os.Getenv("DEV_MODE"); devMode != "" {
		if dev, err := strconv.ParseBool(devMode); err == nil {
			Config.Logger.Development = dev
		}
	}
}
