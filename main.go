package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"toolcat/config"
	"toolcat/middleware"
	"toolcat/models"
	"toolcat/pkg"
	"toolcat/pkg/migrate/migration"
	"toolcat/plugins"
	"toolcat/plugins/examples"
	"toolcat/plugins/features"
	"toolcat/routers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志系统
	if err := pkg.InitLogger(pkg.Options{
		Level:       config.Config.Logger.Level,
		OutputPath:  config.Config.Logger.OutputPath,
		ErrorPath:   config.Config.Logger.ErrorPath,
		Development: config.Config.Logger.Development,
	}); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer pkg.Sync()

	// 设置PluginManager的日志记录器
	plugins.PluginManager.SetLogger(pkg.GetLogger())

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		pkg.Fatal("Failed to load configuration", zap.Error(err))
	}

	// 输出清理后的配置信息（隐藏敏感数据）
	pkg.Info("Configuration loaded successfully", zap.Any("config", config.SanitizeConfig()))

	// 验证配置完整性（确保所有配置项都经过验证）
	if err := config.ValidateConfig(); err != nil {
		pkg.Fatal("Configuration validation failed", zap.Error(err))
	}
	pkg.Info("Configuration validation passed successfully")

	// 监控指标将在路由设置中初始化

	// 初始化数据库
	if err := pkg.InitDatabase(); err != nil {
		pkg.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer pkg.CloseDatabase()

	// 执行数据库迁移
	// 如果禁用了自动迁移，使用SQL迁移文件
	if !config.Config.AutoMigrate {
		log.Println("Starting SQL migrations...")
		mm := migration.NewMigrationManager()
		if err := mm.Init(); err != nil {
			log.Printf("Warning: Failed to initialize migration manager: %v", err)
		} else {
			if err := mm.Up(); err != nil {
				log.Printf("Warning: Migration errors: %v", err)
			} else {
				log.Println("SQL migrations completed successfully")
			}
		}
	} else {
		// 仅当启用自动迁移时才使用GORM自动迁移
		log.Println("Starting GORM auto-migration...")
		if err := models.MigrateTables(pkg.DB); err != nil {
			pkg.Warn("Failed to migrate database tables", zap.Error(err))
		} else {
			log.Println("GORM auto-migration completed successfully")
		}
	}

	// 初始化路由
	router := routers.SetupRouter()

	// 添加错误处理中间件
	errHandler := middleware.NewErrorHandler()
	router.Use(errHandler.HandlerFunc())

	// 监控指标和中间件已在路由设置中配置

	// 注册插件
	registerPlugins(router)

	// Prometheus指标导出路由已在路由设置中注册

	// 监控系统已在路由设置中初始化

	// 初始化插件系统
	if err := plugins.InitPluginSystem(); err != nil {
		pkg.Error("Failed to initialize plugin system", zap.Error(err))
	}

	// 启动服务器
	port := config.Config.Server.Port
	// 创建HTTP服务器并配置连接复用参数
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        router,
		ReadTimeout:    15 * time.Second, // 请求读取超时时间
		WriteTimeout:   15 * time.Second, // 响应写入超时时间
		IdleTimeout:    60 * time.Second, // 空闲连接超时时间（影响Keep-Alive）
		MaxHeaderBytes: 1 << 20,          // 最大请求头大小（1MB）
	}

	go func() {
		pkg.Info("ToolCat 服务启动成功", zap.String("address", fmt.Sprintf("http://localhost:%d", port)))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkg.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	pkg.Info("Shutting down server...")

	// 停止插件监控器
	plugins.PluginManager.StopPluginWatcher()

	// 创建超时上下文，用于优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		pkg.Fatal("Server forced to shutdown", zap.Error(err))
	}

	pkg.Info("Server exiting")
}

// 注册插件
func registerPlugins(router *gin.Engine) {
	// 设置路由引擎到PluginManager
	plugins.PluginManager.SetRouter(router)

	// 注册Hello插件
	helloPlugin := &examples.HelloPlugin{}
	if err := plugins.PluginManager.Register(helloPlugin); err != nil {
		pkg.Error("Failed to register plugin", zap.String("plugin", helloPlugin.Name()), zap.Error(err))
	} else {
		pkg.Info("Successfully registered plugin", zap.String("plugin", helloPlugin.Name()))
	}

	// 注册Note插件
	notePlugin := &features.NotePlugin{}
	if err := plugins.PluginManager.Register(notePlugin); err != nil {
		pkg.Error("Failed to register plugin", zap.String("plugin", notePlugin.Name()), zap.Error(err))
	} else {
		pkg.Info("Successfully registered plugin", zap.String("plugin", notePlugin.Name()))
	}

	// 统一注册所有插件路由
	// 注意：由于我们使用了新的注册机制，插件在注册时已经自动注册了路由
	// 这里可以省略，或者保留作为额外的确认步骤
	// if err := plugins.PluginManager.RegisterAllRoutes(); err != nil {
	//	log.Printf("Failed to register all plugin routes: %v", err)
	// }

	// 注册优化插件
	sampleOptimizedPlugin := examples.NewSampleOptimizedPlugin()
	if err := plugins.PluginManager.Register(sampleOptimizedPlugin); err != nil {
		pkg.Error("Failed to register plugin", zap.String("plugin", sampleOptimizedPlugin.Name()), zap.Error(err))
	} else {
		pkg.Info("Successfully registered plugin", zap.String("plugin", sampleOptimizedPlugin.Name()))
	}

	// 注册依赖插件
	sampleDependentPlugin := examples.NewSampleDependentPlugin()
	if err := plugins.PluginManager.Register(sampleDependentPlugin); err != nil {
		pkg.Error("Failed to register plugin", zap.String("plugin", sampleDependentPlugin.Name()), zap.Error(err))
	} else {
		pkg.Info("Successfully registered plugin", zap.String("plugin", sampleDependentPlugin.Name()))
	}

	// 所有插件注册完成，输出确认日志
	pkg.Info("插件已全部注册运行成功")
}
