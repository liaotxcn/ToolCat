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

	"github.com/gin-gonic/gin"
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
	registerPlugins(router)

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
func registerPlugins(router *gin.Engine) {
	// 注册示例插件
	helloPlugin := &plugins.HelloPlugin{}
	if err := plugins.PluginManager.Register(helloPlugin); err != nil {
		log.Printf("Failed to register plugin %s: %v", helloPlugin.Name(), err)
	} else {
		// 注册插件路由
		helloPlugin.RegisterRoutes(router)
		log.Printf("Successfully registered plugin: %s", helloPlugin.Name())
	}

	// 注册记事本插件
	notePlugin := &plugins.NotePlugin{}
	if err := plugins.PluginManager.Register(notePlugin); err != nil {
		log.Printf("Failed to register plugin %s: %v", notePlugin.Name(), err)
	} else {
		// 注册插件路由
		notePlugin.RegisterRoutes(router)
		log.Printf("Successfully registered plugin: %s", notePlugin.Name())
	}
}
