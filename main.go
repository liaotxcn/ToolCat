package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"toolcat/config"
	"toolcat/middleware"
	"toolcat/models"
	"toolcat/pkg"
	"toolcat/plugins"
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

	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	if err := pkg.InitDatabase(); err != nil {
		pkg.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer pkg.CloseDatabase()

	// 执行数据库迁移
	if err := models.MigrateTables(); err != nil {
		pkg.Warn("Failed to migrate database tables", zap.Error(err))
	}

	// 初始化路由
	router := routers.SetupRouter()

	// 添加错误处理中间件
	errHandler := middleware.NewErrorHandler()
	router.Use(errHandler.HandlerFunc())

	// 注册插件
	registerPlugins(router)

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
		pkg.Info("工具猫服务启动成功", zap.String("address", fmt.Sprintf("http://localhost:%d", port)))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkg.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	pkg.Info("Shutting down server...")

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

	// 注册示例插件
	helloPlugin := &plugins.HelloPlugin{}
	if err := plugins.PluginManager.Register(helloPlugin); err != nil {
		pkg.Error("Failed to register plugin", zap.String("plugin", helloPlugin.Name()), zap.Error(err))
	} else {
		pkg.Info("Successfully registered plugin", zap.String("plugin", helloPlugin.Name()))
	}

	// 注册记事本插件
	notePlugin := &plugins.NotePlugin{}
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
}
