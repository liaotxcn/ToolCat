package pkg_test

import (
	"os"
	"testing"
	"toolcat/pkg"
)

func TestMain(m *testing.M) {
	// 设置测试环境变量
	testDBHost := os.Getenv("TEST_DB_HOST")
	testDBPort := os.Getenv("TEST_DB_PORT")
	testDBUsername := os.Getenv("TEST_DB_USERNAME")
	testDBPassword := os.Getenv("TEST_DB_PASSWORD")
	testDBName := os.Getenv("TEST_DB_NAME")

	// 如果没有设置测试环境变量，则使用默认的测试数据库配置
	if testDBHost == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if testDBPort == "" {
		os.Setenv("DB_PORT", "3306")
	}
	if testDBUsername == "" {
		os.Setenv("DB_USERNAME", "root")
	}
	if testDBPassword == "" {
		os.Setenv("DB_PASSWORD", "123456")
	}
	if testDBName == "" {
		os.Setenv("DB_NAME", "toolcat_test")
	}

	// 运行测试
	code := m.Run()

	os.Exit(code)
}

// TestInitDatabase 测试数据库初始化功能
func TestInitDatabase(t *testing.T) {
	// 清理之前的数据库连接
	if pkg.DB != nil {
		pkg.CloseDatabase()
	}

	// 测试数据库初始化
	err := pkg.InitDatabase()
	if err != nil {
		t.Skipf("Skipping database test: %v\n请确保MySQL服务已启动，并且配置了正确的连接参数", err)
	}

	// 验证数据库连接是否成功
	if pkg.DB == nil {
		t.Fatal("Database connection is nil after InitDatabase")
	}

	// 获取底层sql.DB对象
	sqlDB, err := pkg.DB.DB()
	if err != nil {
		t.Fatalf("Get DB instance failed: %v", err)
	}

	// 测试连接健康状态
	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Database ping failed: %v", err)
	}
}

// TestCloseDatabase 测试数据库关闭功能
func TestCloseDatabase(t *testing.T) {
	// 确保数据库已初始化
	err := pkg.InitDatabase()
	if err != nil {
		t.Skipf("Skipping database test: %v\n请确保MySQL服务已启动，并且配置了正确的连接参数", err)
	}

	// 测试关闭数据库连接
	err = pkg.CloseDatabase()
	if err != nil {
		t.Fatalf("CloseDatabase failed: %v", err)
	}

	// 验证数据库连接是否已关闭
	if pkg.DB != nil {
		sqlDB, _ := pkg.DB.DB()
		if sqlDB != nil {
			err = sqlDB.Ping()
			if err == nil {
				t.Fatal("Database connection is still active after CloseDatabase")
			}
		}
	}
}

// TestDatabaseReconnect 测试数据库重连功能
func TestDatabaseReconnect(t *testing.T) {
	// 初始化数据库
	err := pkg.InitDatabase()
	if err != nil {
		t.Skipf("Skipping database test: %v\n请确保MySQL服务已启动，并且配置了正确的连接参数", err)
	}

	// 关闭连接
	err = pkg.CloseDatabase()
	if err != nil {
		t.Fatalf("CloseDatabase failed: %v", err)
	}

	// 重新连接
	err = pkg.InitDatabase()
	if err != nil {
		t.Skipf("Skipping database reconnect test: %v", err)
	}

	// 验证重连是否成功
	if pkg.DB == nil {
		t.Fatal("Database connection is nil after reconnect")
	}
}