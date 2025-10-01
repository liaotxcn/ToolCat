package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"toolcat/config"
	"toolcat/models"
	"toolcat/pkg"
	"toolcat/plugins"
	"toolcat/routers"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	if err := pkg.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer pkg.CloseDatabase()

	// 执行数据库迁移
	if err := models.MigrateTables(); err != nil {
		log.Printf("Warning: Failed to migrate database tables: %v", err)
	}

	// 初始化路由
	router := routers.SetupRouter()

	// 注册插件
	registerPlugins()

	// 启动服务器
	port := config.Config.Server.Port
	go func() {
		fmt.Printf("工具猫服务启动成功，访问地址：http://localhost:%d\n", port)
		if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

// 注册插件
func registerPlugins() {
	// 注册示例插件
	helloPlugin := &plugins.HelloPlugin{}
	if err := plugins.PluginManager.Register(helloPlugin); err != nil {
		log.Printf("Failed to register plugin %s: %v", helloPlugin.Name(), err)
	} else {
		// 注册插件路由
		router := routers.SetupRouter()
		helloPlugin.RegisterRoutes(router)
		log.Printf("Successfully registered plugin: %s", helloPlugin.Name())
	}

	// 注册计算器插件
	calcPlugin := &plugins.CalcPlugin{}
	if err := plugins.PluginManager.Register(calcPlugin); err != nil {
		log.Printf("Failed to register plugin %s: %v", calcPlugin.Name(), err)
	} else {
		// 注册插件路由
		router := routers.SetupRouter()
		calcPlugin.RegisterRoutes(router)
		log.Printf("Successfully registered plugin: %s", calcPlugin.Name())
	}

	// 可以在这里注册更多插件
}
