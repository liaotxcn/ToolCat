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
	GetDependencies() []string // 返回插件依赖的其他插件名称
	GetConflicts() []string // 返回与当前插件冲突的插件名称

	// 生命周期接口
	Init() error     // 初始化插件
	Shutdown() error // 关闭插件
	OnEnable() error // 插件启用时调用（热重载相关）
	OnDisable() error // 插件禁用时调用（热重载相关）

	// 路由注册接口
	// 新版接口：提供路由定义，由PluginManager统一注册
	GetRoutes() []Route // 获取插件路由定义
	// 旧版接口：为了兼容现有插件保留
	RegisterRoutes(router *gin.Engine) // 注册插件路由

	// 执行功能接口
	Execute(params map[string]interface{}) (interface{}, error) // 执行插件功能

	// 插件配置接口（可选）
	GetDefaultMiddlewares() []gin.HandlerFunc // 获取插件默认中间件
	SetPluginManager(manager *pluginManager) // 设置插件管理器引用
}

// PluginInfo 存储插件信息和路由元数据
type PluginInfo struct {
	Plugin       Plugin  // 插件实例
	Routes       []Route // 插件路由
	Dependencies []string // 依赖的插件名称列表
	Conflicts    []string // 冲突的插件名称列表
	IsRegistered bool    // 路由是否已注册
	IsEnabled    bool    // 插件是否启用
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

	// 设置插件管理器引用
	plugin.SetPluginManager(pm)

	// 检查冲突插件
	conflicts := plugin.GetConflicts()
	for _, conflictName := range conflicts {
		if _, exists := pm.plugins[conflictName]; exists {
			return fmt.Errorf("插件 '%s' 与已注册的插件 '%s' 冲突", name, conflictName)
		}
	}

	// 检查依赖插件
	dependencies := plugin.GetDependencies()
	for _, depName := range dependencies {
		if _, exists := pm.plugins[depName]; !exists {
			return fmt.Errorf("依赖的插件未注册: %s", depName)
		}
	}

	// 初始化插件
	if err := plugin.Init(); err != nil {
		return fmt.Errorf("插件 '%s' 初始化失败: %w", name, err)
	}

	// 创建插件信息
	info := PluginInfo{
		Plugin:       plugin,
		Routes:       plugin.GetRoutes(),
		Dependencies: dependencies,
		Conflicts:    conflicts,
		IsRegistered: false,
		IsEnabled:    true, // 默认为启用状态
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

// EnablePlugin 启用插件
func (pm *pluginManager) EnablePlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("插件 '%s' 不存在", name)
	}

	if info.IsEnabled {
		return nil // 已经是启用状态
	}

	// 检查依赖是否可用
	for _, depName := range info.Dependencies {
		depInfo, exists := pm.plugins[depName]
		if !exists || !depInfo.IsEnabled {
			return fmt.Errorf("依赖的插件 '%s' 未启用", depName)
		}
	}

	// 调用插件的OnEnable方法
	if err := info.Plugin.OnEnable(); err != nil {
		return fmt.Errorf("插件 '%s' 启用回调失败: %w", name, err)
	}

	// 启用插件
	info.IsEnabled = true
	pm.plugins[name] = info

	// 如果路由引擎已设置，注册路由
	if pm.router != nil && !info.IsRegistered {
		if err := pm.registerPluginRoutes(name); err != nil {
			return fmt.Errorf("插件 '%s' 路由注册失败: %w", name, err)
		}
	}

	return nil
}

// DisablePlugin 禁用插件
func (pm *pluginManager) DisablePlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("插件 '%s' 不存在", name)
	}

	if !info.IsEnabled {
		return nil // 已经是禁用状态
	}

	// 检查是否有其他插件依赖当前插件
	for pluginName, pluginInfo := range pm.plugins {
		if pluginName != name && pluginInfo.IsEnabled {
			for _, depName := range pluginInfo.Dependencies {
				if depName == name {
					return fmt.Errorf("插件 '%s' 被插件 '%s' 依赖，无法禁用", name, pluginName)
				}
			}
		}
	}

	// 调用插件的OnDisable方法
	if err := info.Plugin.OnDisable(); err != nil {
		return fmt.Errorf("插件 '%s' 禁用回调失败: %w", name, err)
	}

	// 禁用插件
	info.IsEnabled = false
	pm.plugins[name] = info

	// 注意：Gin不支持动态删除路由，这里只能标记为禁用
	// 在ExecutePlugin等方法中会检查IsEnabled状态

	return nil
}

// ReloadPlugin 重新加载插件
// 注意：这是一个简化实现，在实际生产环境中可能需要结合插件文件监控等功能
func (pm *pluginManager) ReloadPlugin(name string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// 获取当前插件信息
	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("插件 '%s' 不存在", name)
	}

	plugin := info.Plugin
	isEnabled := info.IsEnabled

	// 先禁用插件
	if isEnabled {
		info.IsEnabled = false
		pm.plugins[name] = info
	}

	// 关闭当前插件
	if err := plugin.Shutdown(); err != nil {
		return fmt.Errorf("插件 '%s' 关闭失败: %w", name, err)
	}

	// 从管理器中移除插件
	delete(pm.plugins, name)

	// 重新初始化插件
	if err := plugin.Init(); err != nil {
		return fmt.Errorf("插件 '%s' 重新初始化失败: %w", name, err)
	}

	// 重新创建插件信息
	newInfo := PluginInfo{
		Plugin:       plugin,
		Routes:       plugin.GetRoutes(),
		Dependencies: plugin.GetDependencies(),
		Conflicts:    plugin.GetConflicts(),
		IsRegistered: false,
		IsEnabled:    isEnabled,
	}

	pm.plugins[name] = newInfo

	// 如果路由引擎已设置且插件被启用，重新注册路由
	if pm.router != nil && isEnabled {
		if err := pm.registerPluginRoutes(name); err != nil {
			return fmt.Errorf("插件 '%s' 路由重新注册失败: %w", name, err)
		}
	}

	return nil
}

// GetPluginStatus 获取插件状态
func (pm *pluginManager) GetPluginStatus(name string) (string, bool) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	info, exists := pm.plugins[name]
	if !exists {
		return "not_registered", false
	}

	if info.IsEnabled {
		return "enabled", true
	}
	return "disabled", true
}

// GetAllPluginsInfo 获取所有插件信息
func (pm *pluginManager) GetAllPluginsInfo() []PluginInfo {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	infos := make([]PluginInfo, 0, len(pm.plugins))
	for _, info := range pm.plugins {
		infos = append(infos, info)
	}
	return infos
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

	// 检查插件是否启用
	if !info.IsEnabled {
		return nil, fmt.Errorf("插件 '%s' 已被禁用", name)
	}

	return info.Plugin.Execute(params)
}

// RegisterPlugins 批量注册插件，自动处理依赖顺序
func (pm *pluginManager) RegisterPlugins(plugins []Plugin) error {
	// 1. 构建依赖图
	dependencyGraph := make(map[string][]string)
	pluginMap := make(map[string]Plugin)

	for _, plugin := range plugins {
		name := plugin.Name()
		pluginMap[name] = plugin
		dependencyGraph[name] = plugin.GetDependencies()
	}

	// 2. 拓扑排序
	sortedNames, err := topologicalSort(dependencyGraph)
	if err != nil {
		return err
	}

	// 3. 按排序结果注册插件
	for _, name := range sortedNames {
		if err := pm.Register(pluginMap[name]); err != nil {
			return err
		}
	}

	return nil
}

// topologicalSort 执行拓扑排序
func topologicalSort(graph map[string][]string) ([]string, error) {
	// 计算每个节点的入度
	inDegree := make(map[string]int)
	for node := range graph {
		inDegree[node] = 0
	}

	for _, dependencies := range graph {
		for _, dep := range dependencies {
			inDegree[dep]++
		}
	}

	// 将入度为0的节点加入队列
	queue := []string{}
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	// 执行拓扑排序
	sorted := []string{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		// 减少相邻节点的入度
		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// 检查是否存在环
	if len(sorted) != len(inDegree) {
		return nil, fmt.Errorf("插件依赖关系存在循环依赖")
	}

	return sorted, nil
}

// CheckDependencies 检查所有插件的依赖关系
func (pm *pluginManager) CheckDependencies() []error {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	var errors []error

	for name, info := range pm.plugins {
		for _, depName := range info.Dependencies {
			if _, exists := pm.plugins[depName]; !exists {
				errors = append(errors, fmt.Errorf("插件 '%s' 依赖的插件 '%s' 未注册", name, depName))
			}
		}

		for _, conflictName := range info.Conflicts {
			if _, exists := pm.plugins[conflictName]; exists {
				errors = append(errors, fmt.Errorf("插件 '%s' 与插件 '%s' 冲突", name, conflictName))
			}
		}
	}

	return errors
}

// GetDependencyGraph 获取插件依赖图
func (pm *pluginManager) GetDependencyGraph() map[string]map[string]bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	graph := make(map[string]map[string]bool)

	for name, info := range pm.plugins {
		graph[name] = make(map[string]bool)
		for _, depName := range info.Dependencies {
			graph[name][depName] = true
		}
	}

	return graph
}
