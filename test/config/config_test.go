package config_test

import (
	"os"
	"testing"
	"toolcat/config"
)

// 重置环境变量为测试前状态
func resetEnvVars() {
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USERNAME")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("CSRF_COOKIE_NAME")
	os.Unsetenv("CSRF_COOKIE_HTTP_ONLY")
	os.Unsetenv("CSRF_COOKIE_SAME_SITE")
	os.Unsetenv("CSRF_ENABLED")
}

// TestDefaultConfig 测试默认配置值
func TestDefaultConfig(t *testing.T) {
	resetEnvVars()

	// 验证默认配置是否正确
	if config.Config.Server.Port != 8081 {
		t.Errorf("Expected default server port 8081, got %d", config.Config.Server.Port)
	}

	if config.Config.Database.Host != "localhost" {
		t.Errorf("Expected default DB host 'localhost', got %s", config.Config.Database.Host)
	}

	if config.Config.Database.Port != 3306 {
		t.Errorf("Expected default DB port 3306, got %d", config.Config.Database.Port)
	}

	// 敏感字段默认值应为空，由环境变量或配置文件提供
	if config.Config.Database.Username != "" {
		t.Errorf("Expected default DB username '', got %s", config.Config.Database.Username)
	}

	if config.Config.JWT.Secret != "" {
		t.Errorf("Expected default JWT secret '', got %s", config.Config.JWT.Secret)
	}

	if config.Config.CSRF.Enabled != true {
		t.Errorf("Expected default CSRF.Enabled to be true, got %v", config.Config.CSRF.Enabled)
	}
}

// TestLoadConfigFromEnv 测试从环境变量加载配置
func TestLoadConfigFromEnv(t *testing.T) {
	resetEnvVars()

	// 设置测试环境变量
	testVars := map[string]string{
		"SERVER_PORT":        "8082",
		"DB_HOST":            "test-db-host",
		"DB_PORT":            "3307",
		"DB_USERNAME":        "test-user",
		"DB_PASSWORD":        "test-pass",
		"DB_NAME":            "test-db",
		"JWT_SECRET":         "test-jwt-secret",
		"CSRF_COOKIE_NAME":   "test-csrf",
		"CSRF_COOKIE_HTTP_ONLY": "true",
		"CSRF_COOKIE_SAME_SITE": "Lax",
	}

	// 设置环境变量
	for key, value := range testVars {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}

	// 加载配置
	config.LoadConfig()

	// 验证配置是否从环境变量正确加载
	if config.Config.Server.Port != 8082 {
		t.Errorf("Expected server port 8082 from env, got %d", config.Config.Server.Port)
	}

	if config.Config.Database.Host != "test-db-host" {
		t.Errorf("Expected DB host 'test-db-host' from env, got %s", config.Config.Database.Host)
	}

	if config.Config.Database.Port != 3307 {
		t.Errorf("Expected DB port 3307 from env, got %d", config.Config.Database.Port)
	}

	if config.Config.Database.Username != "test-user" {
		t.Errorf("Expected DB username 'test-user' from env, got %s", config.Config.Database.Username)
	}

	if config.Config.Database.Password != "test-pass" {
		t.Errorf("Expected DB password 'test-pass' from env, got %s", config.Config.Database.Password)
	}

	if config.Config.Database.DBName != "test-db" {
		t.Errorf("Expected DB name 'test-db' from env, got %s", config.Config.Database.DBName)
	}

	if config.Config.JWT.Secret != "test-jwt-secret" {
		t.Errorf("Expected JWT secret 'test-jwt-secret' from env, got %s", config.Config.JWT.Secret)
	}

	if config.Config.CSRF.CookieName != "test-csrf" {
		t.Errorf("Expected CSRF cookie name 'test-csrf' from env, got %s", config.Config.CSRF.CookieName)
	}

	if !config.Config.CSRF.CookieHttpOnly {
		t.Errorf("Expected CSRF cookie HttpOnly true from env, got false")
	}

	if config.Config.CSRF.CookieSameSite != "Lax" {
		t.Errorf("Expected CSRF cookie SameSite 'Lax' from env, got %s", config.Config.CSRF.CookieSameSite)
	}
}

// TestCSRFConfigLoading 测试CSRF配置的特殊加载逻辑
func TestCSRFConfigLoading(t *testing.T) {
	resetEnvVars()

	// 测试CSRF CookieHttpOnly的各种值
	testCases := []struct {
		name        string
		envValue    string
		expected    bool
	}{{
		name:        "true string",
		envValue:    "true",
		expected:    true,
	}, {
		name:        "false string",
		envValue:    "false",
		expected:    false,
	}, {
		name:        "empty string",
		envValue:    "",
		expected:    false,
	}, {
		name:        "invalid string",
		envValue:    "invalid",
		expected:    false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetEnvVars()
			os.Setenv("CSRF_COOKIE_HTTP_ONLY", tc.envValue)
			config.LoadConfig()
			if config.Config.CSRF.CookieHttpOnly != tc.expected {
				t.Errorf("For env value '%s', expected %v, got %v", tc.envValue, tc.expected, config.Config.CSRF.CookieHttpOnly)
			}
		})
	}

	// 测试CSRF CookieSameSite
	testCasesSameSite := []struct {
		name        string
		envValue    string
		expected    string
	}{{
		name:        "Lax",
		envValue:    "Lax",
		expected:    "Lax",
	}, {
		name:        "Strict",
		envValue:    "Strict",
		expected:    "Strict",
	}, {
		name:        "None",
		envValue:    "None",
		expected:    "None",
	}, {
		name:        "empty string",
		envValue:    "",
		expected:    "Lax",
	}, {
		name:        "invalid string",
		envValue:    "invalid",
		expected:    "invalid",
	}}

	for _, tc := range testCasesSameSite {
		t.Run(tc.name, func(t *testing.T) {
			resetEnvVars()
			os.Setenv("CSRF_COOKIE_SAME_SITE", tc.envValue)
			config.LoadConfig()
			if config.Config.CSRF.CookieSameSite != tc.expected {
				t.Errorf("For env value '%s', expected %v, got %v", tc.envValue, tc.expected, config.Config.CSRF.CookieSameSite)
			}
		})
	}

	// 测试CSRF Enabled配置
	testCasesEnabled := []struct {
		name        string
		envValue    string
		expected    bool
	}{{
		name:        "true string",
		envValue:    "true",
		expected:    true,
	}, {
		name:        "false string",
		envValue:    "false",
		expected:    false,
	}}

	for _, tc := range testCasesEnabled {
		t.Run(tc.name, func(t *testing.T) {
			resetEnvVars()
			os.Setenv("CSRF_ENABLED", tc.envValue)
			config.LoadConfig()
			if config.Config.CSRF.Enabled != tc.expected {
				t.Errorf("For env value '%s', expected %v, got %v", tc.envValue, tc.expected, config.Config.CSRF.Enabled)
			}
		})
	}
}
