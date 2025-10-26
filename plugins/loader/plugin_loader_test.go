package loader

import (
	"os"
	"os/exec"
	"path/filepath"
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
