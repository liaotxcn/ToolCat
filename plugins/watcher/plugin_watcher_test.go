package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"toolcat/config"
	"toolcat/pkg"
)

// stubManager implements PluginManager for watcher unit tests
type stubManager struct {
	reloaded     []string
	unregistered []string
	registered   []string
	plugins      map[string]bool
}

func newStubManager() *stubManager { return &stubManager{plugins: make(map[string]bool)} }
func (sm *stubManager) ReloadPlugin(name string) error {
	sm.reloaded = append(sm.reloaded, name)
	return nil
}
func (sm *stubManager) GetPlugin(name string) (Plugin, bool) {
	if sm.plugins[name] {
		return stubWatchedPlugin{name: name}, true
	}
	return nil, false
}
func (sm *stubManager) Unregister(name string) error {
	sm.unregistered = append(sm.unregistered, name)
	delete(sm.plugins, name)
	return nil
}
func (sm *stubManager) Register(p Plugin) error {
	sm.registered = append(sm.registered, p.Name())
	sm.plugins[p.Name()] = true
	return nil
}

// stubWatchedPlugin minimal watcher.Plugin
type stubWatchedPlugin struct{ name string }

func (p stubWatchedPlugin) Name() string { return p.name }

func TestWatcherHotReloadDisabledSkipsActions(t *testing.T) {
	config.Config.Plugins.HotReload = false
	// Prepare temp directory
	dir, err := os.MkdirTemp("", "watcher_unit")
	if err != nil {
		t.Fatalf("temp dir error: %v", err)
	}
	defer os.RemoveAll(dir)
	// Create a .go and a .so to simulate presence
	goPath := filepath.Join(dir, "hot.go")
	soPath := filepath.Join(dir, "hot.so")
	if err := os.WriteFile(goPath, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("write go error: %v", err)
	}
	if err := os.WriteFile(soPath, []byte(""), 0644); err != nil {
		t.Fatalf("write so error: %v", err)
	}

	sm := newStubManager()
	logger := pkg.GetLogger()
	pw, err := NewPluginWatcher(dir, sm, logger)
	if err != nil {
		t.Fatalf("new watcher error: %v", err)
	}
	if err := pw.Start(); err != nil {
		t.Fatalf("start watcher error: %v", err)
	}
	defer pw.Stop()
	// Wait for initial scan debounce window
	time.Sleep(1200 * time.Millisecond)
	if len(sm.registered) != 0 {
		t.Fatalf("expected no registered when HotReload disabled, got %#v", sm.registered)
	}
	if len(sm.reloaded) != 0 {
		t.Fatalf("expected no reloaded when HotReload disabled, got %#v", sm.reloaded)
	}
}

func TestWatcherReloadOnChangeWhenEnabled(t *testing.T) {
	config.Config.Plugins.HotReload = true
	// Prepare temp directory
	dir, err := os.MkdirTemp("", "watcher_unit2")
	if err != nil {
		t.Fatalf("temp dir error: %v", err)
	}
	defer os.RemoveAll(dir)
	// Create plugin files
	goPath := filepath.Join(dir, "hot.go")
	soPath := filepath.Join(dir, "hot.so")
	if err := os.WriteFile(goPath, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("write go error: %v", err)
	}
	if err := os.WriteFile(soPath, []byte(""), 0644); err != nil {
		t.Fatalf("write so error: %v", err)
	}
	// Manager knows plugin exists, so change should trigger reload
	sm := newStubManager()
	sm.plugins["hot"] = true
	logger := pkg.GetLogger()
	pw, err := NewPluginWatcher(dir, sm, logger)
	if err != nil {
		t.Fatalf("new watcher error: %v", err)
	}
	if err := pw.Start(); err != nil {
		t.Fatalf("start watcher error: %v", err)
	}
	defer pw.Stop()
	// Modify the file to trigger change event
	if err := os.WriteFile(goPath, []byte("package main\n// change"), 0644); err != nil {
		t.Fatalf("rewrite file error: %v", err)
	}
	time.Sleep(1500 * time.Millisecond)
	if len(sm.reloaded) == 0 || sm.reloaded[0] != "hot" {
		t.Fatalf("expected 'hot' to be reloaded, got %#v", sm.reloaded)
	}
}

// Additional tests to improve coverage
func TestGetPluginNameFromPath(t *testing.T) {
	p := filepath.Join("/opt", "plugins", "abc.go")
	if got := getPluginNameFromPath(p); got != "abc" {
		t.Fatalf("getPluginNameFromPath: got %s, want %s", got, "abc")
	}
}

func TestIsTempFileVariants(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"file.swp", true},
		{"file.swo", true},
		{"file~", false},
		{".hidden.go", true},
		{"normal.go", false},
	}
	for _, c := range cases {
		if got := isTempFile(c.name); got != c.want {
			t.Fatalf("isTempFile(%s): got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestIsDirectoryTrueFalse(t *testing.T) {
	d := t.TempDir()
	// true case
	if !isDirectory(d) {
		t.Fatalf("expected isDirectory(%s) true", d)
	}
	// false case with file
	f := filepath.Join(d, "x.go")
	if err := os.WriteFile(f, []byte("package main"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if isDirectory(f) {
		t.Fatalf("expected isDirectory(%s) false for file", f)
	}
	// false case with non-existent path
	np := filepath.Join(d, "nope")
	if isDirectory(np) {
		t.Fatalf("expected isDirectory(%s) false for non-existent", np)
	}
}

func TestHandlePluginRemovalUnregisters(t *testing.T) {
	d := t.TempDir()
	sm := newStubManager()
	sm.plugins["gone"] = true
	pw, err := NewPluginWatcher(d, sm, pkg.GetLogger())
	if err != nil {
		t.Fatalf("new watcher error: %v", err)
	}
	path := filepath.Join(d, "gone.go")
	pw.handlePluginRemoval(path)
	if len(sm.unregistered) != 1 || sm.unregistered[0] != "gone" {
		t.Fatalf("expected 'gone' to be unregistered, got %#v", sm.unregistered)
	}
}

func TestTryLoadNewPlugin_FileMissing(t *testing.T) {
	d := t.TempDir()
	sm := newStubManager()
	pw, err := NewPluginWatcher(d, sm, pkg.GetLogger())
	if err != nil {
		t.Fatalf("new watcher error: %v", err)
	}
	config.Config.Plugins.HotReload = true
	// No .so file present for 'hot'
	pw.tryLoadNewPlugin("hot")
	// Should not register anything
	if len(sm.registered) != 0 {
		t.Fatalf("expected no registration when .so missing, got %#v", sm.registered)
	}
}

func TestTryLoadNewPlugin_LoadError(t *testing.T) {
	d := t.TempDir()
	sm := newStubManager()
	pw, err := NewPluginWatcher(d, sm, pkg.GetLogger())
	if err != nil {
		t.Fatalf("new watcher error: %v", err)
	}
	config.Config.Plugins.HotReload = true
	// Create an empty .so to satisfy existence check; loading should fail
	soPath := filepath.Join(d, "hot.so")
	if err := os.WriteFile(soPath, []byte(""), 0644); err != nil {
		t.Fatalf("write so: %v", err)
	}
	pw.tryLoadNewPlugin("hot")
	// On failure, no registration should occur
	if len(sm.registered) != 0 {
		t.Fatalf("expected no registration on load error, got %#v", sm.registered)
	}
}

func TestLoadPluginManifest_Success(t *testing.T) {
	d := t.TempDir()
	m := filepath.Join(d, "manifest.json")
	content := `{
	  "name": "p1",
	  "description": "d",
	  "version": "v1",
	  "author": "a",
	  "dependencies": ["x"],
	  "conflicts": ["y"],
	  "entry_point": "main",
	  "build_tags": ["t"],
	  "required_go_version": "1.22"
	}`
	if err := os.WriteFile(m, []byte(content), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	mf, err := LoadPluginManifest(m)
	if err != nil {
		t.Fatalf("LoadPluginManifest error: %v", err)
	}
	if mf == nil || mf.Name != "p1" || mf.Version != "v1" {
		t.Fatalf("unexpected manifest: %#v", mf)
	}
}

func TestLoadPluginManifest_InvalidJSON(t *testing.T) {
	d := t.TempDir()
	m := filepath.Join(d, "bad.json")
	if err := os.WriteFile(m, []byte("{"), 0644); err != nil {
		t.Fatalf("write bad manifest: %v", err)
	}
	if _, err := LoadPluginManifest(m); err == nil {
		t.Fatalf("expected error for invalid json")
	}
}
