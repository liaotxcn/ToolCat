package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"weave/config"
	"weave/middleware"
	"weave/utils"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareSetsContextKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Config.JWT.Secret = "testsecret"

	userID := uint(42)
	tenantID := uint(1001)

	token, err := utils.GenerateToken(userID, tenantID)
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id":   c.GetUint("user_id"),
			"tenant_id": c.GetUint("tenant_id"),
			"userID":    c.GetUint("userID"),
			"tenantID":  c.GetUint("tenantID"),
		})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var body map[string]uint
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}

	if body["user_id"] != userID || body["userID"] != userID {
		t.Errorf("user id mismatch, expected %d", userID)
	}
	if body["tenant_id"] != tenantID || body["tenantID"] != tenantID {
		t.Errorf("tenant id mismatch, expected %d", tenantID)
	}
}
