package routers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"toolcat/config"
	"toolcat/pkg"
	"toolcat/routers"
	"toolcat/utils"
)

func TestRootRouteOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	if body["api_base"] != "/api/v1" {
		t.Fatalf("expected api_base=/api/v1, got %#v", body["api_base"])
	}
	msg, _ := body["message"].(string)
	if !strings.Contains(msg, "ToolCat") {
		t.Fatalf("expected message to contain 'ToolCat', got %q", msg)
	}
}

func TestHealthRouteOK_WithSQLite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Setup in-memory SQLite DB to satisfy health check
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		t.Fatalf("gorm open error: %v", err)
	}
	pkg.DB = db

	router := routers.SetupRouter()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status=ok, got %#v", body["status"])
	}
	dbInfo, _ := body["database"].(map[string]interface{})
	if dbInfo == nil {
		t.Fatalf("expected database info in response")
	}
	if healthy, _ := dbInfo["healthy"].(bool); !healthy {
		t.Fatalf("expected database.healthy=true, got false")
	}
}

func TestPluginHealth_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/health/plugins/unknown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if body["message"] != "插件不存在" {
		t.Fatalf("expected message '插件不存在', got %#v", body["message"])
	}
}

func TestPluginsList_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/plugins/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if body["error"] != "Authorization header is required" {
		t.Fatalf("expected error 'Authorization header is required', got %#v", body["error"])
	}
}

func TestPluginsList_Authorized_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Ensure config defaults loaded (JWT secret may be empty; both generator and verifier will use the same key)
	_ = config.LoadConfig()
	accessToken, err := utils.GenerateToken(1, 1)
	if err != nil {
		t.Fatalf("generate token error: %v", err)
	}

	router := routers.SetupRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/plugins/", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body []interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if len(body) != 0 {
		t.Fatalf("expected empty plugins list, got %d items", len(body))
	}
}

func TestMetricsEndpoint_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
