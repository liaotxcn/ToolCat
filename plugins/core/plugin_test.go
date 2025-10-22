package core

import (
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
		if n == "B" { foundB = true }
		if n == "C" { foundC = true }
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
	name           string
	deps           []string
	conflicts      []string
	routes         []Route
	pm             *PluginManager
	enableCalled   int
	disableCalled  int
	initCalled     int
	shutdownCalled int
	executeCalled  int
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

func (p *testPlugin) Name() string                             { return p.name }
func (p *testPlugin) Description() string                      { return "test plugin" }
func (p *testPlugin) Version() string                          { return "1.0.0" }
func (p *testPlugin) GetDependencies() []string                { return p.deps }
func (p *testPlugin) GetConflicts() []string                   { return p.conflicts }
func (p *testPlugin) Init() error                              { p.initCalled++; return nil }
func (p *testPlugin) Shutdown() error                          { p.shutdownCalled++; return nil }
func (p *testPlugin) OnEnable() error                          { p.enableCalled++; return nil }
func (p *testPlugin) OnDisable() error                         { p.disableCalled++; return nil }
func (p *testPlugin) GetRoutes() []Route                       { return p.routes }
func (p *testPlugin) GetDefaultMiddlewares() []gin.HandlerFunc { return nil }
func (p *testPlugin) SetPluginManager(manager *PluginManager)  { p.pm = manager }
func (p *testPlugin) Execute(params map[string]interface{}) (interface{}, error) {
	p.executeCalled++
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
