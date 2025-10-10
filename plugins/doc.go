// Package plugins 提供插件系统的核心功能和各类插件实现
// 子目录说明：
// - core: 插件系统核心定义和接口
// - watcher: 插件监控相关功能
// - examples: 示例插件
// - features: 功能性插件
// - templates: 插件开发模板
package plugins

import "toolcat/plugins/core"

// PluginManager 提供对插件系统的统一管理接口
// 对core包中GlobalPluginManager的公开引用
var PluginManager = core.GlobalPluginManager
