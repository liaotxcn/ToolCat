package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"weave/controllers"
	"weave/plugins"
	"weave/plugins/core"
)

// clearPlugins removes all registered plugins for test isolation
func clearPlugins(t *testing.T) {
	names := plugins.PluginManager.ListPlugins()
	for _, n := range names {
		_ = plugins.PluginManager.Unregister(n)
	}
}

type pcTestPlugin struct{ pm *core.PluginManager }

func (p *pcTestPlugin) Name() string                                               { return "pc_demo" }
func (p *pcTestPlugin) Description() string                                        { return "plugin controller test plugin" }
func (p *pcTestPlugin) Version() string                                            { return "1.0.0" }
func (p *pcTestPlugin) GetDependencies() []string                                  { return nil }
func (p *pcTestPlugin) GetConflicts() []string                                     { return nil }
func (p *pcTestPlugin) Init() error                                                { return nil }
func (p *pcTestPlugin) Shutdown() error                                            { return nil }
func (p *pcTestPlugin) OnEnable() error                                            { return nil }
func (p *pcTestPlugin) OnDisable() error                                           { return nil }
func (p *pcTestPlugin) GetRoutes() []core.Route                                    { return nil }
func (p *pcTestPlugin) GetDefaultMiddlewares() []gin.HandlerFunc                   { return nil }
func (p *pcTestPlugin) SetPluginManager(manager *core.PluginManager)               { p.pm = manager }
func (p *pcTestPlugin) RegisterRoutes(router *gin.Engine)                          {}
func (p *pcTestPlugin) Execute(params map[string]interface{}) (interface{}, error) { return nil, nil }

func TestGetAllPlugins_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	clearPlugins(t)

	pc := controllers.PluginController{}
	r := gin.New()
	// Note trailing slash to avoid redirect
	r.GET("/api/v1/plugins/", pc.GetAllPlugins)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/plugins/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty plugins list, got %d", len(list))
	}
}

func TestGetPluginStatus_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	clearPlugins(t)

	pc := controllers.PluginController{}
	r := gin.New()
	r.GET("/api/v1/plugins/:name/status", pc.GetPluginStatus)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/plugins/ghost/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if body["error"] != "插件不存在" || body["plugin"] != "ghost" {
		t.Fatalf("unexpected response: %#v", body)
	}
}

func TestGetDependencyGraph_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	clearPlugins(t)

	pc := controllers.PluginController{}
	r := gin.New()
	r.GET("/api/v1/plugins/dependency-graph", pc.GetDependencyGraph)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/plugins/dependency-graph", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var graph map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &graph); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if len(graph) != 0 {
		t.Fatalf("expected empty graph, got %#v", graph)
	}
}

func TestGetPluginStatus_Enabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	clearPlugins(t)
	if err := plugins.PluginManager.Register(&pcTestPlugin{}); err != nil {
		t.Fatalf("register plugin error: %v", err)
	}
	defer func() { _ = plugins.PluginManager.Unregister("pc_demo") }()

	pc := controllers.PluginController{}
	r := gin.New()
	r.GET("/api/v1/plugins/:name/status", pc.GetPluginStatus)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/plugins/pc_demo/status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if body["status"] != "enabled" || body["plugin"] != "pc_demo" {
		t.Fatalf("unexpected response: %#v", body)
	}
}
