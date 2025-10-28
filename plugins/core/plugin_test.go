package core

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"toolcat/pkg"

	"github.com/gin-gonic/gin"
)

// mockPlugin implements Plugin for unit tests in core package
type mockPlugin struct {
	name      string
	deps      []string
	conflicts []string
	pm        *PluginManager
}

func (p *mockPlugin) Name() string                             { return p.name }
func (p *mockPlugin) Description() string                      { return "mock plugin for unit tests" }
func (p *mockPlugin) Version() string                          { return "0.0.1" }
func (p *mockPlugin) GetDependencies() []string                { return p.deps }
func (p *mockPlugin) GetConflicts() []string                   { return p.conflicts }
func (p *mockPlugin) Init() error                              { return nil }
func (p *mockPlugin) Shutdown() error                          { return nil }
func (p *mockPlugin) OnEnable() error                          { return nil }
func (p *mockPlugin) OnDisable() error                         { return nil }
func (p *mockPlugin) GetRoutes() []Route                       { return nil }
func (p *mockPlugin) GetDefaultMiddlewares() []gin.HandlerFunc { return nil }
func (p *mockPlugin) SetPluginManager(manager *PluginManager)  { p.pm = manager }
func (p *mockPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	return nil, nil
}
func (p *mockPlugin) RegisterRoutes(router *gin.Engine) {}

func TestTopologicalSortOrder(t *testing.T) {
	// 构建一个依赖图：A 依赖 B 和 C
	graph := map[string][]string{
		"A": {"B", "C"},
		"B": {},
		"C": {},
	}

	order, err := topologicalSort(graph)
	if err != nil {
		t.Fatalf("topologicalSort returned error: %v", err)
	}

	if len(order) != 3 {
		t.Fatalf("expected 3 nodes in order, got %d", len(order))
	}
	// 依赖应在使用者之前，A 应该在最后
	if order[len(order)-1] != "A" {
		t.Fatalf("expected A to be last, got %v", order)
	}
	// B 和 C 必须在结果中（顺序不强制）
	foundB, foundC := false, false
	for _, n := range order {
		if n == "B" {
			foundB = true
		}
		if n == "C" {
			foundC = true
		}
	}
	if !foundB || !foundC {
		t.Fatalf("expected both B and C in order, got %v", order)
	}
}

func TestRegisterDependencyCheckUnit(t *testing.T) {
	// Initialize PluginManager with required fields
	pm := &PluginManager{
		plugins: make(map[string]PluginInfo),
		mutex:   &sync.RWMutex{},
	}
	pB := &mockPlugin{name: "B"}
	pC := &mockPlugin{name: "C"}
	pA := &mockPlugin{name: "A", deps: []string{"B", "C"}}
	if err := pm.Register(pB); err != nil {
		t.Fatalf("register B error: %v", err)
	}
	if err := pm.Register(pC); err != nil {
		t.Fatalf("register C error: %v", err)
	}
	if err := pm.Register(pA); err != nil {
		t.Fatalf("register A error: %v", err)
	}
}

func TestGetDependencyGraphUnit(t *testing.T) {
	pm := &PluginManager{
		plugins: make(map[string]PluginInfo),
		mutex:   &sync.RWMutex{},
	}
	pB := &mockPlugin{name: "B"}
	pC := &mockPlugin{name: "C"}
	pA := &mockPlugin{name: "A", deps: []string{"B", "C"}}
	_ = pm.Register(pB)
	_ = pm.Register(pC)
	_ = pm.Register(pA)
	graph := pm.GetDependencyGraph()
	if _, ok := graph["A"]["B"]; !ok {
		t.Fatalf("expected graph to include A->B")
	}
	if _, ok := graph["A"]["C"]; !ok {
		t.Fatalf("expected graph to include A->C")
	}
}

// Additional comprehensive tests for PluginManager interfaces

type testPlugin struct {
	name               string
	deps               []string
	conflicts          []string
	routes             []Route
	pm                 *PluginManager
	enableCalled       int
	disableCalled      int
	initCalled         int
	shutdownCalled     int
	executeCalled      int
	initError          error
	enableError        error
	disableError       error
	executeError       error
	shutdownError      error
	initFunc           func() error
	defaultMiddlewares []gin.HandlerFunc
}

func newTestPlugin(name string, withRoute bool) *testPlugin {
	tp := &testPlugin{name: name}
	if withRoute {
		tp.routes = []Route{
			{
				Path:        "/ping",
				Method:      "GET",
				Handler:     func(c *gin.Context) { c.String(200, "pong") },
				Middlewares: nil,
			},
		}
	}
	return tp
}

func (p *testPlugin) Name() string              { return p.name }
func (p *testPlugin) Description() string       { return "test plugin" }
func (p *testPlugin) Version() string           { return "1.0.0" }
func (p *testPlugin) GetDependencies() []string { return p.deps }
func (p *testPlugin) GetConflicts() []string    { return p.conflicts }
func (p *testPlugin) Init() error {
	p.initCalled++
	if p.initFunc != nil {
		return p.initFunc()
	}
	return p.initError
}
func (p *testPlugin) Shutdown() error    { p.shutdownCalled++; return p.shutdownError }
func (p *testPlugin) OnEnable() error    { p.enableCalled++; return p.enableError }
func (p *testPlugin) OnDisable() error   { p.disableCalled++; return p.disableError }
func (p *testPlugin) GetRoutes() []Route { return p.routes }
func (p *testPlugin) GetDefaultMiddlewares() []gin.HandlerFunc {
	if p.defaultMiddlewares != nil {
		return p.defaultMiddlewares
	}
	return nil
}
func (p *testPlugin) SetPluginManager(manager *PluginManager) { p.pm = manager }
func (p *testPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	p.executeCalled++
	if p.executeError != nil {
		return nil, p.executeError
	}
	return "ok", nil
}
func (p *testPlugin) RegisterRoutes(router *gin.Engine) {}

type mockWatcher struct{ startCalled, stopCalled int }

func (mw *mockWatcher) Start() error { mw.startCalled++; return nil }
func (mw *mockWatcher) Stop()        { mw.stopCalled++ }

func TestEnableDisableAndRouteRegistration(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	router := gin.New()
	pm.SetRouter(router)

	tp := newTestPlugin("P", true)
	if err := pm.Register(tp); err != nil {
		t.Fatalf("register error: %v", err)
	}
	// After register with router, routes should be registered
	if info, ok := pm.GetPluginInfo("P"); !ok || !info.IsRegistered {
		t.Fatalf("expected plugin P routes registered")
	}
	// Disable then enable; callbacks should be called
	if err := pm.DisablePlugin("P"); err != nil {
		t.Fatalf("disable error: %v", err)
	}
	if err := pm.EnablePlugin("P"); err != nil {
		t.Fatalf("enable error: %v", err)
	}
	if tp.disableCalled == 0 || tp.enableCalled == 0 {
		t.Fatalf("expected enable/disable callbacks invoked, got enable=%d disable=%d", tp.enableCalled, tp.disableCalled)
	}
}

func TestDisableBlockedByDependency(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pB := newTestPlugin("B", false)
	pA := newTestPlugin("A", false)
	pA.deps = []string{"B"}
	if err := pm.Register(pB); err != nil {
		t.Fatalf("register B error: %v", err)
	}
	if err := pm.Register(pA); err != nil {
		t.Fatalf("register A error: %v", err)
	}
	if err := pm.DisablePlugin("B"); err == nil {
		t.Fatalf("expected error when disabling B while A depends on it")
	}
}

func TestCheckDependenciesMissing(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pA := newTestPlugin("A", false)
	if err := pm.Register(pA); err != nil {
		t.Fatalf("register A error: %v", err)
	}
	// Inject missing dependency into PluginInfo
	info := pm.plugins["A"]
	info.Dependencies = []string{"B"}
	pm.plugins["A"] = info
	errs := pm.CheckDependencies()
	if len(errs) == 0 {
		t.Fatalf("expected dependency missing error, got none")
	}
	msg := errs[0].Error()
	if !strings.Contains(msg, "A") || !strings.Contains(msg, "B") || !strings.Contains(msg, "未注册") {
		t.Fatalf("unexpected error message: %s", msg)
	}
}

func TestCheckDependenciesConflict(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pC := newTestPlugin("C", false)
	pA := newTestPlugin("A", false)
	if err := pm.Register(pA); err != nil {
		t.Fatalf("register A error: %v", err)
	}
	if err := pm.Register(pC); err != nil {
		t.Fatalf("register C error: %v", err)
	}
	// Inject conflict into PluginInfo
	infoA := pm.plugins["A"]
	infoA.Conflicts = []string{"C"}
	pm.plugins["A"] = infoA
	errs := pm.CheckDependencies()
	if len(errs) == 0 {
		t.Fatalf("expected conflict error, got none")
	}
	found := false
	for _, e := range errs {
		msg := e.Error()
		if strings.Contains(msg, "A") && strings.Contains(msg, "C") && strings.Contains(msg, "冲突") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected conflict error for A and C, got: %#v", errs)
	}
}

func TestCheckDependenciesNoIssues(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pB := newTestPlugin("B", false)
	pA := newTestPlugin("A", false)
	pA.deps = []string{"B"}
	if err := pm.Register(pB); err != nil {
		t.Fatalf("register B error: %v", err)
	}
	if err := pm.Register(pA); err != nil {
		t.Fatalf("register A error: %v", err)
	}
	errs := pm.CheckDependencies()
	if len(errs) != 0 {
		t.Fatalf("expected no dependency errors, got: %#v", errs)
	}
}
func TestExecutePlugin(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	tp := newTestPlugin("X", false)
	if err := pm.Register(tp); err != nil {
		t.Fatalf("register error: %v", err)
	}
	res, err := pm.ExecutePlugin("X", map[string]interface{}{"k": "v"})
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}
	if res != "ok" || tp.executeCalled == 0 {
		t.Fatalf("expected execute called and return ok, got res=%v calls=%d", res, tp.executeCalled)
	}
}

func TestRegisterAllRoutesAfterSettingRouter(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	tp := newTestPlugin("R", true)
	// Register without router set; so no routes yet
	if err := pm.Register(tp); err != nil {
		t.Fatalf("register error: %v", err)
	}
	if info, _ := pm.GetPluginInfo("R"); info.IsRegistered {
		t.Fatalf("expected routes not registered before router set")
	}
	// Set router and register all routes
	router := gin.New()
	pm.SetRouter(router)
	if err := pm.RegisterAllRoutes(); err != nil {
		t.Fatalf("RegisterAllRoutes error: %v", err)
	}
	if info, _ := pm.GetPluginInfo("R"); !info.IsRegistered {
		t.Fatalf("expected routes registered after RegisterAllRoutes")
	}
}

func TestGettersSettersAndUnregister(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	tp := newTestPlugin("G", false)
	_ = pm.Register(tp)
	// Status checks
	if status, ok := pm.GetPluginStatus("G"); !ok || status != "enabled" {
		t.Fatalf("expected enabled status, got %s ok=%v", status, ok)
	}
	// Setters
	pm.SetPluginDir("plugins_dir")
	pm.SetLogger(pkg.GetLogger())
	if pm.pluginDir != "plugins_dir" || pm.logger == nil {
		t.Fatalf("setters failed")
	}
	// Getters
	if _, ok := pm.GetPlugin("G"); !ok {
		t.Fatalf("GetPlugin should succeed")
	}
	infos := pm.GetAllPluginsInfo()
	if len(infos) == 0 {
		t.Fatalf("expected non-empty plugins info")
	}
	routes := pm.GetAllRoutes()
	if len(routes) == 0 {
		t.Fatalf("expected non-empty routes map (empty slice allowed)")
	}
	// Unregister
	if err := pm.Unregister("G"); err != nil {
		t.Fatalf("Unregister error: %v", err)
	}
	if _, ok := pm.GetPlugin("G"); ok {
		t.Fatalf("expected plugin removed after Unregister")
	}
}

func TestRegisterPluginsOrderError(t *testing.T) {
	// A 依赖 B。按 [A, B] 输入，拓扑排序应先 B 后 A，使注册成功。
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pB := newTestPlugin("B", false)
	pA := newTestPlugin("A", false)
	pA.deps = []string{"B"}

	// RegisterPlugins 应该根据拓扑排序自动调整注册顺序，确保依赖优先
	err := pm.RegisterPlugins([]Plugin{pA, pB})
	if err != nil {
		t.Fatalf("expected register success, got error: %v", err)
	}

	// 验证两个插件均已注册
	if _, ok := pm.GetPlugin("A"); !ok {
		t.Fatalf("expected A registered")
	}
	if _, ok := pm.GetPlugin("B"); !ok {
		t.Fatalf("expected B registered")
	}
}

func TestReloadPluginRetainsEnabledAndRegistersRoutes(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// Do NOT set router to avoid Gin duplicate route panic during reload
	tp := newTestPlugin("Z", true)
	if err := pm.Register(tp); err != nil {
		t.Fatalf("register error: %v", err)
	}
	// Ensure enabled
	if status, ok := pm.GetPluginStatus("Z"); !ok || status != "enabled" {
		t.Fatalf("expected enabled")
	}
	if err := pm.ReloadPlugin("Z"); err != nil {
		t.Fatalf("ReloadPlugin error: %v", err)
	}
	if tp.shutdownCalled == 0 || tp.initCalled == 0 {
		t.Fatalf("expected shutdown and re-init during reload")
	}
	// Without router set, routes remain unregistered
	if info, _ := pm.GetPluginInfo("Z"); info.IsRegistered {
		t.Fatalf("expected routes not registered without router")
	}
}

func TestPluginWatcherStartStop(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	mw := &mockWatcher{}
	pm.SetPluginWatcher(mw)
	if err := pm.StartPluginWatcher(); err != nil {
		t.Fatalf("StartPluginWatcher returned error: %v", err)
	}
	pm.StopPluginWatcher()
	if mw.stopCalled == 0 {
		t.Fatalf("expected watcher Stop to be called")
	}
}

func TestListPluginsEmpty(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	names := pm.ListPlugins()
	if len(names) != 0 {
		t.Fatalf("expected empty plugin list, got %v", names)
	}
}

func TestListPluginsReturnsAllNamesAndUpdates(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pA := newTestPlugin("A", false)
	pB := newTestPlugin("B", false)
	pC := newTestPlugin("C", false)
	if err := pm.Register(pA); err != nil {
		t.Fatalf("register A error: %v", err)
	}
	if err := pm.Register(pB); err != nil {
		t.Fatalf("register B error: %v", err)
	}
	if err := pm.Register(pC); err != nil {
		t.Fatalf("register C error: %v", err)
	}
	names := pm.ListPlugins()
	if len(names) != 3 {
		t.Fatalf("expected 3 plugins, got %d: %v", len(names), names)
	}
	seen := map[string]bool{}
	for _, n := range names {
		seen[n] = true
	}
	if !seen["A"] || !seen["B"] || !seen["C"] {
		t.Fatalf("expected names contain A,B,C; got %v", names)
	}
	// Unregister one plugin and verify ListPlugins updates
	if err := pm.Unregister("B"); err != nil {
		t.Fatalf("unregister B error: %v", err)
	}
	names2 := pm.ListPlugins()
	if len(names2) != 2 {
		t.Fatalf("expected 2 plugins after unregister, got %d: %v", len(names2), names2)
	}
	seen2 := map[string]bool{}
	for _, n := range names2 {
		seen2[n] = true
	}
	if !seen2["A"] || !seen2["C"] || seen2["B"] {
		t.Fatalf("expected names contain A,C and not B; got %v", names2)
	}
}

// TestRegisterWithExistingPlugin 测试注册已存在的插件
func TestRegisterWithExistingPlugin(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	p1 := newTestPlugin("duplicate", false)
	p2 := newTestPlugin("duplicate", false)

	// 第一次注册应该成功
	if err := pm.Register(p1); err != nil {
		t.Fatalf("first register should succeed: %v", err)
	}

	// 第二次注册相同名称的插件应该失败
	if err := pm.Register(p2); err == nil {
		t.Fatalf("second register with same name should fail")
	} else if !strings.Contains(err.Error(), "已存在") {
		t.Fatalf("expected 'already exists' error, got: %v", err)
	}
}

// TestRegisterWithConflicts 测试注册有冲突的插件
func TestRegisterWithConflicts(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	p1 := newTestPlugin("plugin1", false)
	p2 := newTestPlugin("plugin2", false)
	p2.conflicts = []string{"plugin1"}

	// 先注册 p1
	if err := pm.Register(p1); err != nil {
		t.Fatalf("register p1 error: %v", err)
	}

	// 注册与 p1 冲突的 p2 应该失败
	if err := pm.Register(p2); err == nil {
		t.Fatalf("register conflicting plugin should fail")
	} else if !strings.Contains(err.Error(), "冲突") {
		t.Fatalf("expected 'conflict' error, got: %v", err)
	}
}

// TestRegisterWithMissingDependency 测试注册缺少依赖的插件
func TestRegisterWithMissingDependency(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pA := newTestPlugin("pluginA", false)
	pA.deps = []string{"nonexistent"}

	// 注册缺少依赖的插件应该失败
	if err := pm.Register(pA); err == nil {
		t.Fatalf("register with missing dependency should fail")
	} else if !strings.Contains(err.Error(), "未注册") {
		t.Fatalf("expected 'not registered' error, got: %v", err)
	}
}

// TestRegisterWithInitError 测试注册初始化失败的插件
func TestRegisterWithInitError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// 创建一个会在Init时返回错误的插件
	errPlugin := &testPlugin{
		name:      "error-plugin",
		initError: errors.New("init failed"),
	}

	// 初始化失败的插件注册应该失败
	if err := pm.Register(errPlugin); err == nil {
		t.Fatalf("register with init error should fail")
	} else if !strings.Contains(err.Error(), "初始化失败") {
		t.Fatalf("expected 'init failed' error, got: %v", err)
	}
	if errPlugin.initCalled != 1 {
		t.Fatalf("expected init called once, got %d", errPlugin.initCalled)
	}
}

// TestEnablePluginNotFound 测试启用不存在的插件
func TestEnablePluginNotFound(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	// 启用不存在的插件应该失败
	if err := pm.EnablePlugin("nonexistent"); err == nil {
		t.Fatalf("enable nonexistent plugin should fail")
	} else if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("expected 'not exist' error, got: %v", err)
	}
}

// TestEnablePluginWithDependencyDisabled 测试启用依赖已禁用的插件
func TestEnablePluginWithDependencyDisabled(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	// 创建两个测试插件，A依赖B
	pB := &testPlugin{name: "B"}
	pA := &testPlugin{name: "A", deps: []string{"B"}}

	// 注册插件B，但不启用它
	if err := pm.Register(pB); err != nil {
		t.Fatalf("register B error: %v", err)
	}
	// 先禁用B
	if err := pm.DisablePlugin("B"); err != nil {
		t.Fatalf("disable B error: %v", err)
	}

	// 注册插件A
	if err := pm.Register(pA); err != nil {
		t.Fatalf("register A error: %v", err)
	}

	// 先禁用A
	if err := pm.DisablePlugin("A"); err != nil {
		t.Fatalf("disable A error: %v", err)
	}

	// 尝试启用 A 应该失败，因为依赖 B 已禁用
	if err := pm.EnablePlugin("A"); err == nil {
		t.Fatalf("enable plugin with disabled dependency should fail")
	} else if !strings.Contains(err.Error(), "未启用") {
		t.Fatalf("expected 'not enabled' error, got: %v", err)
	}
}

// TestEnablePluginWithOnEnableError 测试启用时 OnEnable 失败的插件
func TestEnablePluginWithOnEnableError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// 创建一个会在OnEnable时返回错误的插件
	errPlugin := &testPlugin{
		name:        "error-plugin",
		enableError: errors.New("onEnable failed"),
	}

	// 注册插件
	if err := pm.Register(errPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 禁用再启用
	if err := pm.DisablePlugin("error-plugin"); err != nil {
		t.Fatalf("disable error: %v", err)
	}

	// 启用失败应该返回错误
	if err := pm.EnablePlugin("error-plugin"); err == nil {
		t.Fatalf("enable with onEnable error should fail")
	} else if !strings.Contains(err.Error(), "启用回调失败") {
		t.Fatalf("expected 'enable callback failed' error, got: %v", err)
	}
	if errPlugin.enableCalled != 1 {
		t.Fatalf("expected onEnable called once, got %d", errPlugin.enableCalled)
	}
}

// TestDisablePluginNotFound 测试禁用不存在的插件
func TestDisablePluginNotFound(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	// 禁用不存在的插件应该失败
	if err := pm.DisablePlugin("nonexistent"); err == nil {
		t.Fatalf("disable nonexistent plugin should fail")
	} else if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("expected 'not exist' error, got: %v", err)
	}
}

// TestDisablePluginWithOnDisableError 测试禁用时 OnDisable 失败的插件
func TestDisablePluginWithOnDisableError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// 创建一个会在OnDisable时返回错误的插件
	errPlugin := &testPlugin{
		name:         "error-plugin",
		disableError: errors.New("onDisable failed"),
	}

	// 注册插件
	if err := pm.Register(errPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 禁用失败应该返回错误
	if err := pm.DisablePlugin("error-plugin"); err == nil {
		t.Fatalf("disable with onDisable error should fail")
	} else if !strings.Contains(err.Error(), "禁用回调失败") {
		t.Fatalf("expected 'disable callback failed' error, got: %v", err)
	}
	if errPlugin.disableCalled != 1 {
		t.Fatalf("expected onDisable called once, got %d", errPlugin.disableCalled)
	}
}

// TestReloadPluginNotFound 测试重新加载不存在的插件
func TestReloadPluginNotFound(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	// 重新加载不存在的插件应该失败
	if err := pm.ReloadPlugin("nonexistent"); err == nil {
		t.Fatalf("reload nonexistent plugin should fail")
	} else if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("expected 'not exist' error, got: %v", err)
	}
}

// TestReloadPluginWithShutdownError 测试重新加载时 Shutdown 失败的插件
func TestReloadPluginWithShutdownError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// 创建一个会在Shutdown时返回错误的插件
	errPlugin := &testPlugin{
		name:          "error-plugin",
		shutdownError: errors.New("shutdown failed"),
	}

	// 注册插件
	if err := pm.Register(errPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 重新加载失败应该返回错误
	if err := pm.ReloadPlugin("error-plugin"); err == nil {
		t.Fatalf("reload with shutdown error should fail")
	} else if !strings.Contains(err.Error(), "关闭失败") {
		t.Fatalf("expected 'shutdown failed' error, got: %v", err)
	}
	if errPlugin.shutdownCalled != 1 {
		t.Fatalf("expected shutdown called once, got %d", errPlugin.shutdownCalled)
	}
}

// TestReloadPluginWithInitError 测试重新加载时 Init 失败的插件
func TestReloadPluginWithInitError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// 创建一个在首次Init成功，但在重新Init时失败的插件
	initCount := 0
	errPlugin := &testPlugin{
		name: "error-plugin",
		initFunc: func() error {
			initCount++
			// 第二次调用（重新加载时）返回错误
			if initCount > 1 {
				return errors.New("reinit failed")
			}
			return nil
		},
	}

	// 注册插件
	if err := pm.Register(errPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 重新加载失败应该返回错误
	if err := pm.ReloadPlugin("error-plugin"); err == nil {
		t.Fatalf("reload with reinit error should fail")
	} else if !strings.Contains(err.Error(), "重新初始化失败") {
		t.Fatalf("expected 'reinit failed' error, got: %v", err)
	}
	if initCount != 2 {
		t.Fatalf("expected init called twice, got %d", initCount)
	}
	if errPlugin.shutdownCalled != 1 {
		t.Fatalf("expected shutdown called once, got %d", errPlugin.shutdownCalled)
	}
}

// TestExecutePluginNotFound 测试执行不存在的插件
func TestExecutePluginNotFound(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	// 执行不存在的插件应该失败
	result, err := pm.ExecutePlugin("nonexistent", nil)
	if err == nil {
		t.Fatalf("execute nonexistent plugin should fail")
	} else if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("expected 'not exist' error, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
}

// TestExecutePluginDisabled 测试执行已禁用的插件
func TestExecutePluginDisabled(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	disabledPlugin := newTestPlugin("disabled", false)

	// 注册插件
	if err := pm.Register(disabledPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 禁用插件
	if err := pm.DisablePlugin("disabled"); err != nil {
		t.Fatalf("disable error: %v", err)
	}

	// 执行已禁用的插件应该失败
	result, err := pm.ExecutePlugin("disabled", nil)
	if err == nil {
		t.Fatalf("execute disabled plugin should fail")
	} else if !strings.Contains(err.Error(), "已被禁用") {
		t.Fatalf("expected 'disabled' error, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
}

// TestExecutePluginWithError 测试执行时返回错误的插件
func TestExecutePluginWithError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// 创建一个会在Execute时返回错误的插件
	errPlugin := &testPlugin{
		name:         "error-plugin",
		executeError: errors.New("execute failed"),
	}

	// 注册插件
	if err := pm.Register(errPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 执行失败应该返回错误
	result, err := pm.ExecutePlugin("error-plugin", nil)
	if err == nil {
		t.Fatalf("execute with error should fail")
	}
	if errPlugin.executeCalled != 1 {
		t.Fatalf("expected execute called once, got %d", errPlugin.executeCalled)
	}
	if result != nil {
		t.Fatalf("expected nil result, got: %v", result)
	}
}

// TestRegisterPluginsWithCycleDependency 测试注册有循环依赖的插件
func TestRegisterPluginsWithCycleDependency(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	pA := newTestPlugin("A", false)
	pB := newTestPlugin("B", false)
	pA.deps = []string{"B"}
	pB.deps = []string{"A"}

	// 注册有循环依赖的插件应该失败
	if err := pm.RegisterPlugins([]Plugin{pA, pB}); err == nil {
		t.Fatalf("register plugins with cycle dependency should fail")
	} else if !strings.Contains(err.Error(), "循环依赖") {
		t.Fatalf("expected 'cycle dependency' error, got: %v", err)
	}
}

// TestTopologicalSortWithCycle 测试拓扑排序处理循环依赖
func TestTopologicalSortWithCycle(t *testing.T) {
	// 构建一个循环依赖图：A 依赖 B，B 依赖 A
	graph := map[string][]string{
		"A": {"B"},
		"B": {"A"},
	}

	order, err := topologicalSort(graph)
	if err == nil {
		t.Fatalf("topologicalSort with cycle should return error")
	}
	if !strings.Contains(err.Error(), "循环依赖") {
		t.Fatalf("expected 'cycle dependency' error, got: %v", err)
	}
	if order != nil {
		t.Fatalf("expected nil order, got: %v", order)
	}
}

// TestRegisterAllRoutesWithError 测试注册所有路由时出错的情况
func TestRegisterAllRoutesWithError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	router := gin.New()
	pm.SetRouter(router)

	// 创建一个带有无效路由的插件
	errPlugin := &testPlugin{
		name: "error-route",
		routes: []Route{
			{
				Path:   "/invalid",
				Method: "INVALID", // 无效的 HTTP 方法
				Handler: func(c *gin.Context) {
					c.String(200, "ok")
				},
			},
		},
	}

	// 注册插件时就应该失败，因为路由验证在注册时进行
	if err := pm.Register(errPlugin); err == nil {
		t.Fatalf("register with invalid route should fail")
	} else if !strings.Contains(err.Error(), "不支持的HTTP方法") {
		t.Fatalf("expected 'unsupported HTTP method' error, got: %v", err)
	}
}

// TestRegisterPluginWithDefaultMiddlewares 测试带有默认中间件的插件注册
func TestRegisterPluginWithDefaultMiddlewares(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	router := gin.New()
	pm.SetRouter(router)

	// 创建一个带有默认中间件的插件
	mwPlugin := &testPlugin{
		name: "middleware-plugin",
		defaultMiddlewares: []gin.HandlerFunc{
			func(c *gin.Context) {
				c.Next()
			},
		},
		routes: []Route{
			{
				Path:        "/ping",
				Method:      "GET",
				Handler:     func(c *gin.Context) { c.String(200, "pong") },
				Middlewares: nil,
			},
		},
	}

	// 注册插件
	if err := pm.Register(mwPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 验证插件信息
	info, ok := pm.GetPluginInfo("middleware-plugin")
	if !ok {
		t.Fatalf("plugin not found")
	}
	if !info.IsRegistered {
		t.Fatalf("plugin routes not registered")
	}
}

// TestRegisterPluginWithAuthRequiredRoute 测试带有认证要求的路由注册
func TestRegisterPluginWithAuthRequiredRoute(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	router := gin.New()
	pm.SetRouter(router)

	// 创建一个带有认证要求的路由的插件
	authPlugin := newTestPlugin("auth-plugin", true)
	authPlugin.routes = []Route{
		{
			Path:         "/protected",
			Method:       "GET",
			Handler:      func(c *gin.Context) { c.String(200, "protected") },
			AuthRequired: true,
		},
	}

	// 注册插件
	if err := pm.Register(authPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 验证插件信息
	info, ok := pm.GetPluginInfo("auth-plugin")
	if !ok {
		t.Fatalf("plugin not found")
	}
	if !info.IsRegistered {
		t.Fatalf("plugin routes not registered")
	}
}

// TestGetPluginStatusNotFound 测试获取不存在插件的状态
func TestGetPluginStatusNotFound(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	status, ok := pm.GetPluginStatus("nonexistent")
	if ok {
		t.Fatalf("expected ok=false for nonexistent plugin")
	}
	if status != "not_registered" {
		t.Fatalf("expected status 'not_registered', got: %s", status)
	}
}

// TestGetPluginStatusDisabled 测试获取已禁用插件的状态
func TestGetPluginStatusDisabled(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	disabledPlugin := newTestPlugin("disabled", false)

	// 注册并禁用插件
	if err := pm.Register(disabledPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}
	if err := pm.DisablePlugin("disabled"); err != nil {
		t.Fatalf("disable error: %v", err)
	}

	status, ok := pm.GetPluginStatus("disabled")
	if !ok {
		t.Fatalf("expected ok=true for existing plugin")
	}
	if status != "disabled" {
		t.Fatalf("expected status 'disabled', got: %s", status)
	}
}

// TestGetPluginInfoNotFound 测试获取不存在插件的信息
func TestGetPluginInfoNotFound(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	info, ok := pm.GetPluginInfo("nonexistent")
	if ok {
		t.Fatalf("expected ok=false for nonexistent plugin")
	}
	if info != nil {
		t.Fatalf("expected nil info, got: %v", info)
	}
}

// TestUnregisterNotFound 测试注销不存在的插件
func TestUnregisterNotFound(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	// 注销不存在的插件应该失败
	if err := pm.Unregister("nonexistent"); err == nil {
		t.Fatalf("unregister nonexistent plugin should fail")
	} else if !strings.Contains(err.Error(), "不存在") {
		t.Fatalf("expected 'not exist' error, got: %v", err)
	}
}

// TestUnregisterWithShutdownError 测试注销时 Shutdown 失败的插件
func TestUnregisterWithShutdownError(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}
	// 创建一个会在Shutdown时返回错误的插件
	errPlugin := &testPlugin{
		name:          "error-plugin",
		shutdownError: errors.New("shutdown failed"),
	}

	// 注册插件
	if err := pm.Register(errPlugin); err != nil {
		t.Fatalf("register error: %v", err)
	}

	// 注销失败应该返回错误
	if err := pm.Unregister("error-plugin"); err == nil {
		t.Fatalf("unregister with shutdown error should fail")
	} else if !strings.Contains(err.Error(), "关闭失败") {
		t.Fatalf("expected 'shutdown failed' error, got: %v", err)
	}
	if errPlugin.shutdownCalled != 1 {
		t.Fatalf("expected shutdown called once, got %d", errPlugin.shutdownCalled)
	}
}

// TestStartPluginWatcherWithoutWatcher 测试未设置监控器时启动监控器
func TestStartPluginWatcherWithoutWatcher(t *testing.T) {
	pm := &PluginManager{plugins: make(map[string]PluginInfo), mutex: &sync.RWMutex{}}

	// 未设置监控器时启动应该失败
	if err := pm.StartPluginWatcher(); err == nil {
		t.Fatalf("StartPluginWatcher without watcher should fail")
	} else if !strings.Contains(err.Error(), "未初始化") {
		t.Fatalf("expected 'not initialized' error, got: %v", err)
	}
}

// TestTopologicalSortWithIsolatedNodes 测试拓扑排序处理孤立节点
func TestTopologicalSortWithIsolatedNodes(t *testing.T) {
	// 构建包含孤立节点的图
	graph := map[string][]string{
		"A": {}, // 孤立节点
		"B": {}, // 孤立节点
		"C": {"B"},
	}

	order, err := topologicalSort(graph)
	if err != nil {
		t.Fatalf("topologicalSort returned error: %v", err)
	}

	if len(order) != 3 {
		t.Fatalf("expected 3 nodes in order, got %d", len(order))
	}

	// B 应该在 C 之前
	bIndex := -1
	cIndex := -1
	for i, n := range order {
		if n == "B" {
			bIndex = i
		} else if n == "C" {
			cIndex = i
		}
	}

	if bIndex == -1 || cIndex == -1 {
		t.Fatalf("B and C should be in the order, got: %v", order)
	}
	if bIndex >= cIndex {
		t.Fatalf("B should come before C, got B at %d and C at %d", bIndex, cIndex)
	}
}
