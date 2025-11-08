package plugins

import (
	"fmt"
	"weave/config"
	"weave/pkg"
	"weave/plugins/core"
	"weave/plugins/watcher"

	"go.uber.org/zap"
)

// 全局插件管理器实例
var PluginManager = core.GlobalPluginManager

// pluginManagerAdapter 适配器，将core.PluginManager适配到watcher.PluginManager接口
type pluginManagerAdapter struct {
	manager *core.PluginManager
}

// ReloadPlugin 实现watcher.PluginManager接口
func (adapter *pluginManagerAdapter) ReloadPlugin(name string) error {
	return adapter.manager.ReloadPlugin(name)
}

// GetPlugin 实现watcher.PluginManager接口
func (adapter *pluginManagerAdapter) GetPlugin(name string) (watcher.Plugin, bool) {
	plugin, exists := adapter.manager.GetPlugin(name)
	return plugin, exists
}

// Unregister 实现watcher.PluginManager接口
func (adapter *pluginManagerAdapter) Unregister(name string) error {
	return adapter.manager.Unregister(name)
}

// Register 实现watcher.PluginManager接口
func (adapter *pluginManagerAdapter) Register(plugin watcher.Plugin) error {
	// 这里需要断言，因为core包的Register方法需要core.Plugin类型
	// 实际使用中，传入的应该是实现了core.Plugin接口的实例
	if corePlugin, ok := plugin.(core.Plugin); ok {
		return adapter.manager.Register(corePlugin)
	}
	return fmt.Errorf("plugin does not implement core.Plugin interface")
}

// InitPluginSystem 初始化插件系统
// 包括创建和设置PluginWatcher实例
func InitPluginSystem() error {
	// 获取插件目录
	pluginsDir := config.Config.Plugins.Dir
	if pluginsDir == "" {
		pluginsDir = "./plugins"
	}
	PluginManager.SetPluginDir(pluginsDir)

	// 如果配置启用了插件监控器，则创建并设置监控器
	if config.Config.Plugins.WatcherEnabled {
		// 创建适配器
		adapter := &pluginManagerAdapter{manager: PluginManager}

		// 创建插件监控器实例
		pw, err := watcher.NewPluginWatcher(pluginsDir, adapter, pkg.GetLogger())
		if err != nil {
			return err
		}

		// 设置插件监控器
		PluginManager.SetPluginWatcher(pw)

		// 启动插件监控器
		if err := pw.Start(); err != nil {
			return err
		}

		pkg.Info("插件监控器已启动", zap.String("pluginDir", pluginsDir))
	} else {
		pkg.Info("插件监控器已被配置禁用")
	}

	return nil
}
