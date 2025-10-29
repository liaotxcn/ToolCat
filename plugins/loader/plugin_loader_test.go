package loader

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"runtime"
	"sync"
	"testing"
	"time"

	"toolcat/pkg"
)

func TestNewPluginLoaderInitialState(t *testing.T) {
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)
	if pl == nil {
		t.Fatalf("expected non-nil PluginLoader")
	}
	if pl.logger == nil {
		t.Fatalf("expected logger to be set")
	}
	if len(pl.loadedPlugins) != 0 {
		t.Fatalf("expected no loaded plugins initially; got %d", len(pl.loadedPlugins))
	}
}

func TestGetPluginPath(t *testing.T) {
	dir := "plugins"
	name := "hello"
	got := GetPluginPath(dir, name)
	want := filepath.Join(dir, "hello.so")
	if got != want {
		t.Fatalf("GetPluginPath mismatch: got %s, want %s", got, want)
	}
}

func TestGetPluginPathAbsolute(t *testing.T) {
	dir := filepath.FromSlash("/opt/toolcat/plugins")
	name := "world"
	got := GetPluginPath(dir, name)
	want := filepath.Join(dir, "world.so")
	if got != want {
		t.Fatalf("GetPluginPath absolute mismatch: got %s, want %s", got, want)
	}
}

func TestUnloadPluginWhenNotLoaded(t *testing.T) {
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)

	if pl.GetLoadedPlugin("unknown") {
		t.Fatalf("expected plugin 'unknown' to be not loaded")
	}
	if err := pl.UnloadPlugin("unknown"); err != nil {
		t.Fatalf("UnloadPlugin returned error for not loaded plugin: %v", err)
	}
}

func TestGetLoadedPluginToggle(t *testing.T) {
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)

	// Simulate a loaded plugin by inserting a nil entry (presence of key is enough)
	pl.loadedPlugins["demo"] = nil
	if !pl.GetLoadedPlugin("demo") {
		t.Fatalf("expected plugin 'demo' to be reported as loaded")
	}
	if err := pl.UnloadPlugin("demo"); err != nil {
		t.Fatalf("UnloadPlugin returned error: %v", err)
	}
	if pl.GetLoadedPlugin("demo") {
		t.Fatalf("expected plugin 'demo' to be not loaded after unload")
	}
}

func TestUnloadMultiplePlugins(t *testing.T) {
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)
	names := []string{"A", "B", "C", "D"}
	for _, n := range names {
		pl.loadedPlugins[n] = nil
	}
	for _, n := range names {
		if !pl.GetLoadedPlugin(n) {
			t.Fatalf("expected %s to be loaded before unload", n)
		}
		if err := pl.UnloadPlugin(n); err != nil {
			t.Fatalf("UnloadPlugin(%s) error: %v", n, err)
		}
		if pl.GetLoadedPlugin(n) {
			t.Fatalf("expected %s to be not loaded after unload", n)
		}
	}
}

func TestConcurrentGetAndUnload(t *testing.T) {
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)
	names := []string{"x", "y", "z"}
	for _, n := range names {
		pl.loadedPlugins[n] = nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// alternate between checking and unloading
			n := names[idx%len(names)]
			_ = pl.GetLoadedPlugin(n)
			_ = pl.UnloadPlugin(n)
			_ = pl.GetLoadedPlugin(n)
		}(i)
	}
	wg.Wait()

	// Eventually all should be unloaded
	time.Sleep(10 * time.Millisecond)
	for _, n := range names {
		if pl.GetLoadedPlugin(n) {
			t.Fatalf("expected %s to be not loaded after concurrent operations", n)
		}
	}
}

func TestLoadPluginUnloadsExistingOnError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Go plugin not supported on Windows; skipping LoadPlugin tests")
	}
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)

	// Preload an entry to force the unload path
	pl.loadedPlugins["A"] = nil
	// Use a clearly non-existent path so Open fails on supported OS
	_, err := pl.LoadPlugin(filepath.Join("nonexistent_dir", "A.so"), "A")
	if err == nil {
		t.Fatalf("expected LoadPlugin to fail for nonexistent path")
	}
	// After failure, existing entry should be removed by UnloadPlugin branch
	if pl.GetLoadedPlugin("A") {
		t.Fatalf("expected plugin 'A' to be removed after failed reload")
	}
}

func TestLoadRealPluginOnLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("real plugin build/load test requires Linux")
	}
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)

	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "covplugin.go")
	so := filepath.Join(tmpDir, "covplugin.so")

	source := `package main
import (
    "toolcat/plugins/core"
    "github.com/gin-gonic/gin"
)

type covPlugin struct{ pm *core.PluginManager }
func (p *covPlugin) Name() string { return "covplugin" }
func (p *covPlugin) Description() string { return "" }
func (p *covPlugin) Version() string { return "v0" }
func (p *covPlugin) GetDependencies() []string { return nil }
func (p *covPlugin) GetConflicts() []string { return nil }
func (p *covPlugin) Init() error { return nil }
func (p *covPlugin) Shutdown() error { return nil }
func (p *covPlugin) OnEnable() error { return nil }
func (p *covPlugin) OnDisable() error { return nil }
func (p *covPlugin) GetRoutes() []core.Route { return nil }
func (p *covPlugin) RegisterRoutes(_ *gin.Engine) {}
func (p *covPlugin) Execute(_ map[string]interface{}) (interface{}, error) { return nil, nil }
func (p *covPlugin) GetDefaultMiddlewares() []gin.HandlerFunc { return nil }
func (p *covPlugin) SetPluginManager(m *core.PluginManager) { p.pm = m }

func NewPlugin() core.Plugin { return &covPlugin{} }
`
	if err := os.WriteFile(src, []byte(source), 0644); err != nil {
		t.Fatalf("failed to write plugin source: %v", err)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", so, src)
	cmd.Env = os.Environ()
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build plugin: %v\n%s", err, string(out))
	}

	pluginPath := so
	inst, err := pl.LoadPlugin(pluginPath, "covplugin")
	if err != nil {
		t.Fatalf("LoadPlugin returned error: %v", err)
	}
	if inst == nil || inst.Name() != "covplugin" {
		t.Fatalf("unexpected plugin instance: %#v", inst)
	}
	if !pl.GetLoadedPlugin("covplugin") {
		t.Fatalf("expected covplugin to be marked as loaded")
	}
	if err := pl.UnloadPlugin("covplugin"); err != nil {
		t.Fatalf("UnloadPlugin error: %v", err)
	}
	if pl.GetLoadedPlugin("covplugin") {
		t.Fatalf("expected covplugin to be not loaded after unload")
	}
}

func TestLoadPluginUnsupportedOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific behavior: plugin package unsupported")
	}
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)
	_, err := pl.LoadPlugin(filepath.Join("nonexistent", "x.so"), "x")
	if err == nil {
		t.Fatalf("expected error from LoadPlugin on Windows (plugin unsupported)")
	}
}

// 模拟插件包，用于测试

type mockPluginLoader struct {
	*PluginLoader
	originalOpenPlugin func(path string) (*plugin.Plugin, error)
}

// TestLoadPluginNameMismatch 测试插件名称不匹配的情况
func TestLoadPluginNameMismatch(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Go plugin not supported on Windows")
	}

	// 创建一个模拟插件并构建
	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)

	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "namemismatch.go")
	so := filepath.Join(tmpDir, "namemismatch.so")

	source := `package main
import (
    "toolcat/plugins/core"
    "github.com/gin-gonic/gin"
)

type mockPlugin struct{ pm *core.PluginManager }
func (p *mockPlugin) Name() string { return "wrongname" } // 与加载时名称不匹配
func (p *mockPlugin) Description() string { return "" }
func (p *mockPlugin) Version() string { return "v0" }
func (p *mockPlugin) GetDependencies() []string { return nil }
func (p *mockPlugin) GetConflicts() []string { return nil }
func (p *mockPlugin) Init() error { return nil }
func (p *mockPlugin) Shutdown() error { return nil }
func (p *mockPlugin) OnEnable() error { return nil }
func (p *mockPlugin) OnDisable() error { return nil }
func (p *mockPlugin) GetRoutes() []core.Route { return nil }
func (p *mockPlugin) RegisterRoutes(_ *gin.Engine) {}
func (p *mockPlugin) Execute(_ map[string]interface{}) (interface{}, error) { return nil, nil }
func (p *mockPlugin) GetDefaultMiddlewares() []gin.HandlerFunc { return nil }
func (p *mockPlugin) SetPluginManager(m *core.PluginManager) { p.pm = m }

func NewPlugin() core.Plugin { return &mockPlugin{} }
`

	if err := os.WriteFile(src, []byte(source), 0644); err != nil {
		t.Fatalf("failed to write plugin source: %v", err)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", so, src)
	cmd.Env = os.Environ()
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build plugin: %v\n%s", err, string(out))
	}

	// 尝试加载插件，期望名称不匹配错误
	_, err = pl.LoadPlugin(so, "expectedname")
	if err == nil {
		t.Fatalf("expected name mismatch error but got nil")
	} else if !errors.Is(err, errors.New("插件名称不匹配")) && !contains(err.Error(), "名称不匹配") {
		t.Fatalf("expected name mismatch error, got: %v", err)
	}
}

// TestLoadPluginSuccessReload 测试成功重新加载插件
func TestLoadPluginSuccessReload(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Go plugin not supported on Windows")
	}

	logger := pkg.GetLogger()
	pl := NewPluginLoader(logger)

	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "reloadplugin.go")
	so := filepath.Join(tmpDir, "reloadplugin.so")

	source := `package main
import (
    "toolcat/plugins/core"
    "github.com/gin-gonic/gin"
)

type reloadPlugin struct{ pm *core.PluginManager }
func (p *reloadPlugin) Name() string { return "reloadplugin" }
func (p *reloadPlugin) Description() string { return "" }
func (p *reloadPlugin) Version() string { return "v0" }
func (p *reloadPlugin) GetDependencies() []string { return nil }
func (p *reloadPlugin) GetConflicts() []string { return nil }
func (p *reloadPlugin) Init() error { return nil }
func (p *reloadPlugin) Shutdown() error { return nil }
func (p *reloadPlugin) OnEnable() error { return nil }
func (p *reloadPlugin) OnDisable() error { return nil }
func (p *reloadPlugin) GetRoutes() []core.Route { return nil }
func (p *reloadPlugin) RegisterRoutes(_ *gin.Engine) {}
func (p *reloadPlugin) Execute(_ map[string]interface{}) (interface{}, error) { return nil, nil }
func (p *reloadPlugin) GetDefaultMiddlewares() []gin.HandlerFunc { return nil }
func (p *reloadPlugin) SetPluginManager(m *core.PluginManager) { p.pm = m }

func NewPlugin() core.Plugin { return &reloadPlugin{} }
`

	if err := os.WriteFile(src, []byte(source), 0644); err != nil {
		t.Fatalf("failed to write plugin source: %v", err)
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", so, src)
	cmd.Env = os.Environ()
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build plugin: %v\n%s", err, string(out))
	}

	// 第一次加载
	inst1, err := pl.LoadPlugin(so, "reloadplugin")
	if err != nil {
		t.Fatalf("first LoadPlugin error: %v", err)
	}

	// 第二次加载（重新加载）
	inst2, err := pl.LoadPlugin(so, "reloadplugin")
	if err != nil {
		t.Fatalf("second LoadPlugin error: %v", err)
	}

	// 验证两个实例都不是nil且名称正确
	if inst1 == nil || inst1.Name() != "reloadplugin" {
		t.Fatalf("invalid first plugin instance")
	}
	if inst2 == nil || inst2.Name() != "reloadplugin" {
		t.Fatalf("invalid second plugin instance")
	}

	// 卸载插件
	if err := pl.UnloadPlugin("reloadplugin"); err != nil {
		t.Fatalf("UnloadPlugin error: %v", err)
	}
}

// 辅助函数：检查字符串是否包含子串
func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestLoadPluginWrongEntryPointType 测试插件入口点类型错误的情况
func TestLoadPluginWrongEntryPointType(t *testing.T) {
	// 只有在Linux上才能运行插件测试
	if runtime.GOOS != "linux" {
		t.Skip("插件功能仅在Linux上支持")
	}

	// 创建一个临时目录
	dir := t.TempDir()
	logger := pkg.GetLogger()
	loader := NewPluginLoader(logger)

	// 创建一个错误的插件Go文件，入口点是一个变量而不是函数
	pluginGoPath := filepath.Join(dir, "wrong_entry.go")
	pluginContent := `package main

// 错误的入口点类型
var Plugin = "这不是一个函数"
`

	if err := os.WriteFile(pluginGoPath, []byte(pluginContent), 0644); err != nil {
		t.Fatalf("写入插件文件失败: %v", err)
	}

	// 编译插件
	soPath := filepath.Join(dir, "wrong_entry.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, pluginGoPath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("编译插件失败: %v", err)
	}

	// 尝试加载插件
	_, err := loader.LoadPlugin(soPath, "wrong_entry")
	if err == nil {
		t.Fatal("应该返回错误，但没有")
	}

	// 验证错误消息包含预期内容
	if !contains(err.Error(), "不是一个有效的插件入口点") {
		t.Fatalf("错误消息不包含预期内容: %v", err)
	}
}

// TestLoadPluginEntryPointNotFound 测试插件入口点不存在的情况
func TestLoadPluginEntryPointNotFound(t *testing.T) {
	// 只有在Linux上才能运行插件测试
	if runtime.GOOS != "linux" {
		t.Skip("插件功能仅在Linux上支持")
	}

	// 创建一个临时目录
	dir := t.TempDir()
	logger := pkg.GetLogger()
	loader := NewPluginLoader(logger)

	// 创建一个没有Plugin入口点的插件Go文件
	pluginGoPath := filepath.Join(dir, "no_entry.go")
	pluginContent := `package main

// 没有定义Plugin变量
func SomeOtherFunction() {}
`

	if err := os.WriteFile(pluginGoPath, []byte(pluginContent), 0644); err != nil {
		t.Fatalf("写入插件文件失败: %v", err)
	}

	// 编译插件
	soPath := filepath.Join(dir, "no_entry.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, pluginGoPath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("编译插件失败: %v", err)
	}

	// 尝试加载插件
	_, err := loader.LoadPlugin(soPath, "no_entry")
	if err == nil {
		t.Fatal("应该返回错误，但没有")
	}

	// 验证错误消息包含预期内容
	if !contains(err.Error(), "未找到插件入口点") {
		t.Fatalf("错误消息不包含预期内容: %v", err)
	}
}

// TestUnloadNonExistentPlugin 测试卸载不存在的插件
func TestUnloadNonExistentPlugin(t *testing.T) {
	logger := pkg.GetLogger()
	loader := NewPluginLoader(logger)

	// 尝试卸载一个不存在的插件
	loader.UnloadPlugin("non_existent_plugin")

	// 这个操作应该是安全的，不会产生错误
	t.Log("Unloaded non-existent plugin without errors")
}

// TestMultipleLoadAndUnload 测试多次加载和卸载插件的情况
func TestMultipleLoadAndUnload(t *testing.T) {
	// 只有在Linux上才能运行插件测试
	if runtime.GOOS != "linux" {
		t.Skip("插件功能仅在Linux上支持")
	}

	// 创建一个临时目录
	dir := t.TempDir()
	logger := pkg.GetLogger()
	loader := NewPluginLoader(logger)

	// 创建一个简单的插件Go文件
	pluginGoPath := filepath.Join(dir, "test_plugin.go")
	pluginContent := `package main

import "fmt"

// 简单的插件接口实现
type TestPlugin struct{}

func (p *TestPlugin) Name() string {
	return "test_plugin"
}

func (p *TestPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	return "success", nil
}

// 插件入口点
var Plugin = &TestPlugin{}
`

	if err := os.WriteFile(pluginGoPath, []byte(pluginContent), 0644); err != nil {
		t.Fatalf("写入插件文件失败: %v", err)
	}

	// 编译插件
	soPath := filepath.Join(dir, "test_plugin.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, pluginGoPath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("编译插件失败: %v", err)
	}

	// 多次加载和卸载插件
	for i := 0; i < 3; i++ {
		// 加载插件
		plugin, err := loader.LoadPlugin(soPath, "test_plugin")
		if err != nil {
			t.Fatalf("第%d次加载插件失败: %v", i+1, err)
		}

		// 验证插件名称
		if p, ok := plugin.(interface{ Name() string }); ok {
			if p.Name() != "test_plugin" {
				t.Fatalf("插件名称不匹配: %s", p.Name())
			}
		} else {
			t.Fatal("插件不实现Name()方法")
		}

		// 卸载插件
		loader.UnloadPlugin("test_plugin")
	}

	t.Log("Multiple load/unload test completed successfully")
}
