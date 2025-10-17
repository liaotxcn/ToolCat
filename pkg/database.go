package pkg

import (
	"fmt"
	"log"
	"os"
	"time"

	"toolcat/config"
	"toolcat/pkg/metrics"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	// 加载配置
	config.LoadConfig()

	var dsn string
	var dialector gorm.Dialector

	// 根据数据库驱动类型构建连接字符串
	switch config.Config.Database.Driver {
	case "postgres":
		// PostgreSQL连接字符串
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.Config.Database.Host,
			config.Config.Database.Port,
			config.Config.Database.Username,
			config.Config.Database.Password,
			config.Config.Database.DBName,
		)
		dialector = postgres.Open(dsn)
	case "mysql":
		fallthrough
	default:
		// MySQL连接字符串
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=30s&readTimeout=30s&writeTimeout=30s&collation=utf8mb4_unicode_ci&tls=false",
			config.Config.Database.Username,
			config.Config.Database.Password,
			config.Config.Database.Host,
			config.Config.Database.Port,
			config.Config.Database.DBName,
			config.Config.Database.Charset,
		)
		dialector = mysql.Open(dsn)
	}

	// 根据环境设置日志级别
	logLevel := logger.Info
	if !config.Config.Logger.Development {
		logLevel = logger.Warn // 生产环境使用Warn级别
	}

	// 配置自定义日志器
	customLogger := logger.New(
		log.New(os.Stdout, "[gorm] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  config.Config.Logger.Development,
		},
	)

	// 连接重试机制
	maxRetries := 3
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		// 创建带有性能监控的GORM配置
		gormConfig := &gorm.Config{
			Logger: customLogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		}

		// 添加GORM性能监控插件
		DB, lastErr = gorm.Open(dialector, gormConfig)
		if lastErr == nil {
			// 记录连接建立指标
			metrics.RecordDatabaseQuery("connect", "system", 0)
		}
		if lastErr == nil {
			break
		}
		log.Printf("Database connection attempt failed, retrying... attempt=%d/%d, error=%v",
			i+1, maxRetries, lastErr)
		time.Sleep(1 * time.Second) // 等待一秒后重试
	}
	if lastErr != nil {
		return fmt.Errorf("failed to connect database after %d retries: %w", maxRetries, lastErr)
	}

	// 获取底层数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置优化的连接池参数
	sqlDB.SetMaxIdleConns(20)                  // 增加空闲连接数，适应高峰期
	sqlDB.SetMaxOpenConns(100)                 // 保持最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(15 * time.Minute) // 添加连接最大空闲时间

	// 连接池预热
	for i := 0; i < 5; i++ {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 连接健康检查
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// 启动数据库连接监控
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // 监测时间间隔
		for range ticker.C {
			stats := sqlDB.Stats()
			idle := stats.Idle
			open := stats.OpenConnections
			metrics.UpdateDatabaseConnections(open)
			log.Printf("Database connection stats: idle=%d, open=%d", idle, open)
		}
	}()

	// 输出数据库连接成功日志
	dbType := "MySQL"
	if config.Config.Database.Driver == "postgres" {
		dbType = "PostgreSQL"
	}
	Info("数据库连接成功",
		zap.String("database_type", dbType),
		zap.String("driver", config.Config.Database.Driver),
		zap.String("host", config.Config.Database.Host),
		zap.Int("port", config.Config.Database.Port),
		zap.String("database", config.Config.Database.DBName),
	)
	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
