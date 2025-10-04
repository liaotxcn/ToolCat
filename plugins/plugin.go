package plugins

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

// Route 定义路由结构
type Route struct {
	Path         string            // 路由路径（不包含插件前缀）
	Method       string            // HTTP方法
	Handler      gin.HandlerFunc   // 处理函数
	Middlewares  []gin.HandlerFunc // 路由特定中间件
	Description  string            // 路由描述
	AuthRequired bool              // 是否需要认证
	Tags         []string          // 路由标签
	Params       map[string]string // 参数说明
} // 路由结构定义

// Plugin 插件接口定义
type Plugin interface {
	// 基础信息接口
	Name() string        // 返回插件名称
	Description() string // 返回插件描述
	Version() string     // 返回插件版本

	// 生命周期接口
	Init() error     // 初始化插件
	Shutdown() error // 关闭插件

	// 路由注册接口
	// 新版接口：提供路由定义，由PluginManager统一注册
	GetRoutes() []Route // 获取插件路由定义
	// 旧版接口：为了兼容现有插件保留
	RegisterRoutes(router *gin.Engine) // 注册插件路由

	// 执行功能接口
	Execute(params map[string]interface{}) (interface{}, error) // 执行插件功能

	// 插件配置接口（可选）
	GetDefaultMiddlewares() []gin.HandlerFunc // 获取插件默认中间件
}

// PluginInfo 存储插件信息和路由元数据
type PluginInfo struct {
	Plugin       Plugin  // 插件实例
	Routes       []Route // 插件路由
	IsRegistered bool    // 路由是否已注册
}

// PluginManager 插件管理器
var PluginManager = &pluginManager{
	plugins: make(map[string]PluginInfo),
	router:  nil,
	mutex:   &sync.RWMutex{},
}

type pluginManager struct {
	plugins map[string]PluginInfo // 存储插件信息和路由
	router  *gin.Engine           // 路由引擎引用
	mutex   *sync.RWMutex         // 读写锁，保证线程安全
}

// SetRouter 设置路由引擎
func (pm *pluginManager) SetRouter(router *gin.Engine) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.router = router
}

// Register 注册插件
func (pm *pluginManager) Register(plugin Plugin) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	name := plugin.Name()
	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("插件 '%s' 已存在", name)
	}

	// 初始化插件
	if err := plugin.Init(); err != nil {
		return fmt.Errorf("插件 '%s' 初始化失败: %w", name, err)
	}

	// 创建插件信息
	info := PluginInfo{
		Plugin:       plugin,
		Routes:       plugin.GetRoutes(),
		IsRegistered: false,
	}

	pm.plugins[name] = info

	// 如果路由引擎已设置，自动注册路由
	if pm.router != nil {
		if err := pm.registerPluginRoutes(name); err != nil {
			return fmt.Errorf("插件 '%s' 路由注册失败: %w", name, err)
		}
	}

	return nil
}

// registerPluginRoutes 注册单个插件的路由
func (pm *pluginManager) registerPluginRoutes(name string) error {
	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("插件 '%s' 不存在", name)
	}

	if pm.router == nil {
		return fmt.Errorf("路由引擎未初始化")
	}

	// 如果路由已经注册，先清理
	if info.IsRegistered {
		// 注意：Gin不支持动态删除路由，这里只能标记为未注册
		// 实际生产环境中可能需要重启服务或使用其他路由方案
		info.IsRegistered = false
	}

	plugin := info.Plugin
	pluginName := plugin.Name()

	// 创建插件路由组
	pluginGroup := pm.router.Group(fmt.Sprintf("/plugins/%s", pluginName))

	// 添加插件默认中间件
	if defaultMiddlewares := plugin.GetDefaultMiddlewares(); len(defaultMiddlewares) > 0 {
		pluginGroup.Use(defaultMiddlewares...)
	}

	// 获取插件路由
	routes := plugin.GetRoutes()

	// 如果没有通过GetRoutes提供路由，则回退到旧版的RegisterRoutes方法
	if len(routes) == 0 {
		plugin.RegisterRoutes(pm.router)
		info.IsRegistered = true
		pm.plugins[name] = info
		return nil
	}

	// 注册每个路由
	for _, route := range routes {
		// 创建路由处理函数链
		handlers := append(route.Middlewares, route.Handler)

		// 根据HTTP方法注册路由
		switch route.Method {
		case "GET":
			pluginGroup.GET(route.Path, handlers...)
		case "POST":
			pluginGroup.POST(route.Path, handlers...)
		case "PUT":
			pluginGroup.PUT(route.Path, handlers...)
		case "DELETE":
			pluginGroup.DELETE(route.Path, handlers...)
		case "PATCH":
			pluginGroup.PATCH(route.Path, handlers...)
		case "OPTIONS":
			pluginGroup.OPTIONS(route.Path, handlers...)
		default:
			return fmt.Errorf("不支持的HTTP方法: %s", route.Method)
		}

		// 更新路由信息
		info.Routes = routes
	}

	info.IsRegistered = true
	pm.plugins[name] = info
	return nil
}

// Unregister 注销插件
func (pm *pluginManager) Unregister(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("插件 '%s' 不存在", name)
	}

	// 关闭插件
	plugin := info.Plugin
	if err := plugin.Shutdown(); err != nil {
		return fmt.Errorf("插件 '%s' 关闭失败: %w", name, err)
	}

	// 标记路由为未注册
	info.IsRegistered = false

	// 从管理器中删除插件
	delete(pm.plugins, name)
	return nil
}

// GetPlugin 获取插件
func (pm *pluginManager) GetPlugin(name string) (Plugin, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	info, exists := pm.plugins[name]
	if !exists {
		return nil, false
	}
	return info.Plugin, true
}

// ListPlugins 列出所有插件
func (pm *pluginManager) ListPlugins() []string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	names := make([]string, 0, len(pm.plugins))
	for name := range pm.plugins {
		names = append(names, name)
	}
	return names
}

// GetPluginInfo 获取插件详细信息
func (pm *pluginManager) GetPluginInfo(name string) (*PluginInfo, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	info, exists := pm.plugins[name]
	if !exists {
		return nil, false
	}
	return &info, true
}

// RegisterAllRoutes 注册所有插件的路由
func (pm *pluginManager) RegisterAllRoutes() error {
	pm.mutex.RLock()
	// 复制插件名称列表，避免在注册过程中锁定太久
	pluginNames := make([]string, 0, len(pm.plugins))
	for name := range pm.plugins {
		pluginNames = append(pluginNames, name)
	}
	pm.mutex.RUnlock()

	// 逐个注册插件路由
	for _, name := range pluginNames {
		pm.mutex.Lock()
		err := pm.registerPluginRoutes(name)
		pm.mutex.Unlock()

		if err != nil {
			return fmt.Errorf("注册插件 '%s' 路由失败: %w", name, err)
		}
	}

	return nil
}

// GetAllRoutes 获取所有插件的路由信息
func (pm *pluginManager) GetAllRoutes() map[string][]Route {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	allRoutes := make(map[string][]Route)
	for name, info := range pm.plugins {
		allRoutes[name] = info.Routes
	}
	return allRoutes
}

// ExecutePlugin 执行插件功能
func (pm *pluginManager) ExecutePlugin(name string, params map[string]interface{}) (interface{}, error) {
	pm.mutex.RLock()
	info, exists := pm.plugins[name]
	pm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("插件 '%s' 不存在", name)
	}

	return info.Plugin.Execute(params)
}
