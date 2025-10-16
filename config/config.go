package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 应用程序配置结构
var Config struct {
	// 配置文件设置
	ConfigFiles struct {
		Path string
		Type string // yaml/json
	}

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
		Secret             string
		AccessTokenExpiry  int // 访问令牌过期时间（分钟）
		RefreshTokenExpiry int // 刷新令牌过期时间（小时）
	}

	// CSRF配置
	CSRF struct {
		Enabled        bool
		CookieName     string
		HeaderName     string
		TokenLength    int
		CookieMaxAge   int // 秒
		CookiePath     string
		CookieDomain   string
		CookieSecure   bool
		CookieHttpOnly bool
		CookieSameSite string
	}

	// 数据库迁移配置
	AutoMigrate bool

	// 插件配置
	Plugins struct {
		Dir            string
		WatcherEnabled bool
		ScanInterval   int // 秒
		HotReload      bool
	}
}

// 初始化默认值
func init() {
	// 配置文件设置
	Config.ConfigFiles.Path = "./config/config.yaml"
	Config.ConfigFiles.Type = "yaml"

	// 服务器配置
	Config.Server.Port = 8081

	// 数据库配置（非敏感字段默认值）
	Config.Database.Driver = "mysql"
	Config.Database.Host = "localhost"
	Config.Database.Port = 3306
	Config.Database.DBName = "toolcat"
	Config.Database.Charset = "utf8mb4"
	// 敏感字段（数据库用户名和密码）将通过环境变量或配置文件设置
	Config.Database.Username = ""
	Config.Database.Password = ""

	// 日志配置
	Config.Logger.Level = "info"
	Config.Logger.OutputPath = "stdout"
	Config.Logger.ErrorPath = "stderr"
	Config.Logger.Development = false

	// JWT配置
	Config.JWT.Secret = "" // 敏感信息，将通过环境变量或配置文件设置
	Config.JWT.AccessTokenExpiry = 60      // 60分钟
	Config.JWT.RefreshTokenExpiry = 24 * 7 // 7天

	// CSRF配置
	Config.CSRF.Enabled = true
	Config.CSRF.CookieName = "XSRF-TOKEN"
	Config.CSRF.HeaderName = "X-CSRF-Token"
	Config.CSRF.TokenLength = 32
	Config.CSRF.CookieMaxAge = 3600 * 24 * 7 // 7天
	Config.CSRF.CookiePath = "/"
	Config.CSRF.CookieDomain = ""
	Config.CSRF.CookieSecure = false   // 开发环境下为false
	Config.CSRF.CookieHttpOnly = false // 必须为false以便前端可以读取
	Config.CSRF.CookieSameSite = "Lax"

	// 数据库迁移配置
	Config.AutoMigrate = true

	// 插件配置
	Config.Plugins.Dir = "./plugins"
	Config.Plugins.WatcherEnabled = true
	Config.Plugins.ScanInterval = 5 // 5秒
	Config.Plugins.HotReload = true
}

// ValidateConfig 验证配置的有效性，特别是敏感信息
func ValidateConfig() error {
	// 检查必要的敏感配置项
	if Config.Database.Username == "" {
		return fmt.Errorf("数据库用户名未配置，请设置DB_USERNAME环境变量或在配置文件中指定")
	}
	
	if Config.Database.Password == "" {
		return fmt.Errorf("数据库密码未配置，请设置DB_PASSWORD环境变量或在配置文件中指定")
	}
	
	if Config.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥未配置，请设置JWT_SECRET环境变量或在配置文件中指定")
	}
	
	return nil
}

// SanitizeConfig 清理配置中的敏感信息，用于日志输出
func SanitizeConfig() map[string]interface{} {
	// 创建配置的安全副本用于日志输出
	sanitized := map[string]interface{}{
		"Server": map[string]interface{}{
			"Port": Config.Server.Port,
		},
		"Database": map[string]interface{}{
			"Driver":  Config.Database.Driver,
			"Host":    Config.Database.Host,
			"Port":    Config.Database.Port,
			"Username": Config.Database.Username,
			"Password": "***", // 隐藏密码
			"DBName":   Config.Database.DBName,
			"Charset":  Config.Database.Charset,
		},
		"Logger": map[string]interface{}{
			"Level":       Config.Logger.Level,
			"OutputPath":  Config.Logger.OutputPath,
			"ErrorPath":   Config.Logger.ErrorPath,
			"Development": Config.Logger.Development,
		},
		"JWT": map[string]interface{}{
			"Secret":             "***", // 隐藏密钥
			"AccessTokenExpiry":  Config.JWT.AccessTokenExpiry,
			"RefreshTokenExpiry": Config.JWT.RefreshTokenExpiry,
		},
		"CSRF": map[string]interface{}{
			"Enabled":        Config.CSRF.Enabled,
			"CookieName":     Config.CSRF.CookieName,
			"HeaderName":     Config.CSRF.HeaderName,
			"TokenLength":    Config.CSRF.TokenLength,
			"CookieMaxAge":   Config.CSRF.CookieMaxAge,
			"CookiePath":     Config.CSRF.CookiePath,
			"CookieDomain":   Config.CSRF.CookieDomain,
			"CookieSecure":   Config.CSRF.CookieSecure,
			"CookieHttpOnly": Config.CSRF.CookieHttpOnly,
			"CookieSameSite": Config.CSRF.CookieSameSite,
		},
		"AutoMigrate": Config.AutoMigrate,
		"Plugins": map[string]interface{}{
			"Dir":            Config.Plugins.Dir,
			"WatcherEnabled": Config.Plugins.WatcherEnabled,
			"ScanInterval":   Config.Plugins.ScanInterval,
			"HotReload":      Config.Plugins.HotReload,
		},
	}
	
	return sanitized
}

// LoadConfigFile 从配置文件加载配置
func LoadConfigFile() error {
	// 检查配置文件是否存在
	if _, err := os.Stat(Config.ConfigFiles.Path); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置
		return nil
	} else if err != nil {
		return err
	}

	// 读取配置文件内容
	content, err := ioutil.ReadFile(Config.ConfigFiles.Path)
	if err != nil {
		return err
	}

	// 创建一个与Config结构相同的临时结构体用于解析
	// 使用map来保持灵活性
	var configMap map[string]interface{}

	// 根据文件类型解析配置
	switch Config.ConfigFiles.Type {
	case "yaml", "yml":
		if err := yaml.Unmarshal(content, &configMap); err != nil {
			return err
		}
	case "json":
		if err := json.Unmarshal(content, &configMap); err != nil {
			return err
		}
	default:
		return nil
	}

	// 将解析后的值映射到Config结构体
	if serverMap, ok := configMap["server"].(map[string]interface{}); ok {
		mapToServerConfig(serverMap)
	}

	if databaseMap, ok := configMap["database"].(map[string]interface{}); ok {
		mapToDatabaseConfig(databaseMap)
	}

	if loggerMap, ok := configMap["logger"].(map[string]interface{}); ok {
		mapToLoggerConfig(loggerMap)
	}

	if jwtMap, ok := configMap["jwt"].(map[string]interface{}); ok {
		mapToJWTConfig(jwtMap)
	}

	if csrfMap, ok := configMap["csrf"].(map[string]interface{}); ok {
		mapToCSRFConfig(csrfMap)
	}

	if pluginsMap, ok := configMap["plugins"].(map[string]interface{}); ok {
		mapToPluginsConfig(pluginsMap)
	}

	return nil
}

// mapToServerConfig 将map映射到Server配置
func mapToServerConfig(configMap map[string]interface{}) {
	if port, ok := configMap["port"]; ok {
		Config.Server.Port = convertToInt(port)
	}
}

// mapToDatabaseConfig 将map映射到Database配置
func mapToDatabaseConfig(configMap map[string]interface{}) {
	if driver, ok := configMap["driver"].(string); ok {
		Config.Database.Driver = driver
	}
	if host, ok := configMap["host"].(string); ok {
		Config.Database.Host = host
	}
	if port, ok := configMap["port"]; ok {
		Config.Database.Port = convertToInt(port)
	}
	if username, ok := configMap["username"].(string); ok {
		Config.Database.Username = username
	}
	if password, ok := configMap["password"].(string); ok {
		Config.Database.Password = password
	}
	if dbname, ok := configMap["dbname"].(string); ok {
		Config.Database.DBName = dbname
	}
	if charset, ok := configMap["charset"].(string); ok {
		Config.Database.Charset = charset
	}
}

// mapToLoggerConfig 将map映射到Logger配置
func mapToLoggerConfig(configMap map[string]interface{}) {
	if level, ok := configMap["level"].(string); ok {
		Config.Logger.Level = level
	}
	if outputPath, ok := configMap["outputPath"].(string); ok {
		Config.Logger.OutputPath = outputPath
	}
	if errorPath, ok := configMap["errorPath"].(string); ok {
		Config.Logger.ErrorPath = errorPath
	}
	if development, ok := configMap["development"]; ok {
		Config.Logger.Development = convertToBool(development)
	}
}

// mapToJWTConfig 将map映射到JWT配置
func mapToJWTConfig(configMap map[string]interface{}) {
	if secret, ok := configMap["secret"].(string); ok {
		Config.JWT.Secret = secret
	}
	if accessTokenExpiry, ok := configMap["accessTokenExpiry"]; ok {
		Config.JWT.AccessTokenExpiry = convertToInt(accessTokenExpiry)
	}
	if refreshTokenExpiry, ok := configMap["refreshTokenExpiry"]; ok {
		Config.JWT.RefreshTokenExpiry = convertToInt(refreshTokenExpiry)
	}
}

// mapToCSRFConfig 将map映射到CSRF配置
func mapToCSRFConfig(configMap map[string]interface{}) {
	if enabled, ok := configMap["enabled"]; ok {
		Config.CSRF.Enabled = convertToBool(enabled)
	}
	if cookieName, ok := configMap["cookieName"].(string); ok {
		Config.CSRF.CookieName = cookieName
	}
	if headerName, ok := configMap["headerName"].(string); ok {
		Config.CSRF.HeaderName = headerName
	}
	if tokenLength, ok := configMap["tokenLength"]; ok {
		Config.CSRF.TokenLength = convertToInt(tokenLength)
	}
	if cookieMaxAge, ok := configMap["cookieMaxAge"]; ok {
		Config.CSRF.CookieMaxAge = convertToInt(cookieMaxAge)
	}
	if cookiePath, ok := configMap["cookiePath"].(string); ok {
		Config.CSRF.CookiePath = cookiePath
	}
	if cookieDomain, ok := configMap["cookieDomain"].(string); ok {
		Config.CSRF.CookieDomain = cookieDomain
	}
	if cookieSecure, ok := configMap["cookieSecure"]; ok {
		Config.CSRF.CookieSecure = convertToBool(cookieSecure)
	}
	if cookieHttpOnly, ok := configMap["cookieHttpOnly"]; ok {
		Config.CSRF.CookieHttpOnly = convertToBool(cookieHttpOnly)
	}
	if cookieSameSite, ok := configMap["cookieSameSite"].(string); ok {
		Config.CSRF.CookieSameSite = cookieSameSite
	}
}

// mapToPluginsConfig 将map映射到Plugins配置
func mapToPluginsConfig(configMap map[string]interface{}) {
	if dir, ok := configMap["dir"].(string); ok {
		Config.Plugins.Dir = dir
	}
	if watcherEnabled, ok := configMap["watcherEnabled"]; ok {
		Config.Plugins.WatcherEnabled = convertToBool(watcherEnabled)
	}
	if scanInterval, ok := configMap["scanInterval"]; ok {
		Config.Plugins.ScanInterval = convertToInt(scanInterval)
	}
	if hotReload, ok := configMap["hotReload"]; ok {
		Config.Plugins.HotReload = convertToBool(hotReload)
	}
}

// convertToInt 将interface{}转换为int
func convertToInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// convertToBool 将interface{}转换为bool
func convertToBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	case int:
		return v != 0
	case float64:
		return v != 0
	}
	return false
}

// GetAbsConfigFilePath 获取配置文件的绝对路径
func GetAbsConfigFilePath() (string, error) {
	if Config.ConfigFiles.Path == "" {
		return "", nil
	}

	absPath, err := filepath.Abs(Config.ConfigFiles.Path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// LoadConfig 从配置文件和环境变量加载配置
func LoadConfig() error {
	// 首先从配置文件加载配置
	if err := LoadConfigFile(); err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 然后从环境变量加载配置，环境变量会覆盖配置文件的值
	// 配置文件路径
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		Config.ConfigFiles.Path = configPath
		// 获取文件扩展名以确定类型
		ext := filepath.Ext(configPath)
		if ext == ".json" {
			Config.ConfigFiles.Type = "json"
		} else if ext == ".yaml" || ext == ".yml" {
			Config.ConfigFiles.Type = "yaml"
		}
		// 重新从新的配置文件加载
		if err := LoadConfigFile(); err != nil {
			return fmt.Errorf("加载指定配置文件失败: %w", err)
		}
	}

	// 服务器端口
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			Config.Server.Port = p
		}
	}

	// 插件配置
	if dir := os.Getenv("PLUGINS_DIR"); dir != "" {
		Config.Plugins.Dir = dir
	}

	if watcherEnabled := os.Getenv("PLUGINS_WATCHER_ENABLED"); watcherEnabled != "" {
		if enabled, err := strconv.ParseBool(watcherEnabled); err == nil {
			Config.Plugins.WatcherEnabled = enabled
		}
	}

	if scanInterval := os.Getenv("PLUGINS_SCAN_INTERVAL"); scanInterval != "" {
		if interval, err := strconv.Atoi(scanInterval); err == nil {
			Config.Plugins.ScanInterval = interval
		}
	}

	if hotReload := os.Getenv("PLUGINS_HOT_RELOAD"); hotReload != "" {
		if reload, err := strconv.ParseBool(hotReload); err == nil {
			Config.Plugins.HotReload = reload
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
	
	// 数据库用户名 - 敏感信息，优先使用环境变量
	if Config.Database.Username == "" {
		Config.Database.Username = os.Getenv("DB_USERNAME")
	}
	
	// 数据库密码 - 敏感信息，优先使用环境变量
	if Config.Database.Password == "" {
		Config.Database.Password = os.Getenv("DB_PASSWORD")
	}
	
	// JWT密钥 - 敏感信息，优先使用环境变量
	if Config.JWT.Secret == "" {
		Config.JWT.Secret = os.Getenv("JWT_SECRET")
	}

	// CSRF配置
	if csrfEnabled := os.Getenv("CSRF_ENABLED"); csrfEnabled != "" {
		if enabled, err := strconv.ParseBool(csrfEnabled); err == nil {
			Config.CSRF.Enabled = enabled
		}
	}

	if cookieName := os.Getenv("CSRF_COOKIE_NAME"); cookieName != "" {
		Config.CSRF.CookieName = cookieName
	}

	if headerName := os.Getenv("CSRF_HEADER_NAME"); headerName != "" {
		Config.CSRF.HeaderName = headerName
	}

	if tokenLength := os.Getenv("CSRF_TOKEN_LENGTH"); tokenLength != "" {
		if length, err := strconv.Atoi(tokenLength); err == nil {
			Config.CSRF.TokenLength = length
		}
	}

	if cookieMaxAge := os.Getenv("CSRF_COOKIE_MAX_AGE"); cookieMaxAge != "" {
		if maxAge, err := strconv.Atoi(cookieMaxAge); err == nil {
			Config.CSRF.CookieMaxAge = maxAge
		}
	}

	if cookiePath := os.Getenv("CSRF_COOKIE_PATH"); cookiePath != "" {
		Config.CSRF.CookiePath = cookiePath
	}

	if cookieDomain := os.Getenv("CSRF_COOKIE_DOMAIN"); cookieDomain != "" {
		Config.CSRF.CookieDomain = cookieDomain
	}

	if cookieSecure := os.Getenv("CSRF_COOKIE_SECURE"); cookieSecure != "" {
		if secure, err := strconv.ParseBool(cookieSecure); err == nil {
			Config.CSRF.CookieSecure = secure
		}
	}

	if cookieHttpOnly := os.Getenv("CSRF_COOKIE_HTTP_ONLY"); cookieHttpOnly != "" {
		if httpOnly, err := strconv.ParseBool(cookieHttpOnly); err == nil {
			Config.CSRF.CookieHttpOnly = httpOnly
		}
	}

	if cookieSameSite := os.Getenv("CSRF_COOKIE_SAME_SITE"); cookieSameSite != "" {
		Config.CSRF.CookieSameSite = cookieSameSite
	}
	
	// 验证配置有效性
	return ValidateConfig()
}
