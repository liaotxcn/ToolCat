package routers

import (
	"toolcat/controllers"
	"toolcat/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 创建路由引擎
	router := gin.Default()

	// 全局中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LogMiddleware())
	router.Use(middleware.RequestBufferMiddleware())

	// 认证相关路由
	auth := router.Group("/auth")
	{
		userCtrl := &controllers.UserController{}
		auth.POST("/register", userCtrl.Register)
		auth.POST("/login", userCtrl.Login)
	}

	// API分组
	api := router.Group("/api/v1")
	{
		// 使用认证中间件
		api.Use(middleware.AuthMiddleware())

		// 用户相关路由
		users := api.Group("/users")
		{
			userCtrl := &controllers.UserController{}
			users.GET("/", userCtrl.GetUsers)
			users.GET("/:id", userCtrl.GetUser)
			users.POST("/", userCtrl.CreateUser)
			users.PUT("/:id", userCtrl.UpdateUser)
			users.DELETE("/:id", userCtrl.DeleteUser)
		}

		// 工具相关路由
		tools := api.Group("/tools")
		{
			toolCtrl := &controllers.ToolController{}
			tools.GET("/", toolCtrl.GetTools)
			tools.GET("/:id", toolCtrl.GetTool)
			tools.POST("/", toolCtrl.CreateTool)
			tools.PUT("/:id", toolCtrl.UpdateTool)
			tools.DELETE("/:id", toolCtrl.DeleteTool)
			tools.POST("/:id/execute", toolCtrl.ExecuteTool)
		}

		// 插件相关路由
		plugins := api.Group("/plugins")
		{
			// 插件管理接口
			plugins.GET("/", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Get all plugins"})
			})
			plugins.POST("/load", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Load plugin"})
			})
			plugins.POST("/unload/:name", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Unload plugin"})
			})
		}
	}

	// 根路径处理
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":             "欢迎使用工具猫(ToolCat)服务！",
			"version":             "1.0.0",
			"api_base":            "/api/v1",
			"health_check":        "/health",
			"available_endpoints": []string{"/api/v1/users", "/api/v1/tools", "/api/v1/plugins", "/health"},
		})
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
