package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"toolcat/middleware"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterBurstExceeded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// 高速率，突发容量为2，第三个请求应被限制
	r.Use(middleware.RateLimiter(1000, 2))
	r.GET("/rl", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	mkReq := func() *http.Request {
		req, _ := http.NewRequest("GET", "/rl", nil)
		req.Header.Set("X-Forwarded-For", "1.2.3.4")
		return req
	}

	// 1st
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, mkReq())
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200 for first request, got %d", w1.Code)
	}

	// 2nd
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, mkReq())
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 for second request, got %d", w2.Code)
	}

	// 3rd should be 429
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, mkReq())
	if w3.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 for third request, got %d", w3.Code)
	}
}
