package routers

import (
	"time"
	"toolcat/controllers"
	"toolcat/middleware"
	"toolcat/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// 创建路由引擎，但不使用默认中间件，而是手动添加需要的中间件
	router := gin.New()

	// 初始化指标管理器
	mm := metrics.NewMetricsManager()

	// 添加基本中间件
	router.Use(gin.Recovery()) // 恢复中间件，处理panic
	router.Use(gin.Logger())   // 使用gin内置的日志中间件
	router.Use(middleware.CORSMiddleware())

	// 注册Prometheus指标导出路由
	mm.RegisterMetricsRouter(router)

	// 启动指标更新器，每30秒更新一次系统指标
	mm.StartMetricsUpdater(30 * time.Second)

	// 创建一个应用组，为所有其他路由应用完整的中间件链
	appGroup := router.Group("")
	{
		// 添加其他必要的中间件，但仅应用于appGroup而不是全局
		appGroup.Use(middleware.RequestBufferMiddleware())
		appGroup.Use(middleware.CSRFMiddleware())
		appGroup.Use(mm.HTTPMonitoringMiddleware()) // 添加HTTP请求监控中间件

		// 认证相关路由
		auth := appGroup.Group("/auth")
		{
			//  限流保护，为认证接口添加限流：每秒允许10个请求，突发容量20
			auth.Use(middleware.RateLimiter(10, 20))
			userCtrl := &controllers.UserController{}
			auth.POST("/register", userCtrl.Register)
			auth.POST("/login", userCtrl.Login)
			auth.POST("/refresh-token", userCtrl.RefreshToken)
		}

		// API分组
		api := appGroup.Group("/api/v1")
		{
			// 使用认证中间件
			api.Use(middleware.AuthMiddleware())
			// 为API接口添加限流：每秒允许20个请求，突发容量50
			api.Use(middleware.RateLimiter(20, 50))

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
				pluginCtrl := &controllers.PluginController{}
				// 获取所有插件信息
				plugins.GET("/", pluginCtrl.GetAllPlugins)
				// 获取插件状态
				plugins.GET("/:name/status", pluginCtrl.GetPluginStatus)
				// 启用插件
				plugins.POST("/:name/enable", pluginCtrl.EnablePlugin)
				// 禁用插件
				plugins.POST("/:name/disable", pluginCtrl.DisablePlugin)
				// 重载插件
				plugins.POST("/:name/reload", pluginCtrl.ReloadPlugin)
				// 获取插件依赖图
				plugins.GET("/dependency-graph", pluginCtrl.GetDependencyGraph)
			}
		}
	}

	// 根路径和健康检查路由放在appGroup内，确保一致的中间件处理
	healthCtrl := &controllers.HealthController{}
	appGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":             "欢迎使用 ToolCat 服务！",
			"version":             "1.0.0",
			"api_base":            "/api/v1",
			"health_check":        "/health",
			"available_endpoints": []string{"/api/v1/users", "/api/v1/tools", "/api/v1/plugins", "/health"},
		})
	})

	// 健康检查 - 健康检查控制器提供更全面的健康状态信息
	appGroup.GET("/health", healthCtrl.GetHealth)

	return router
}
