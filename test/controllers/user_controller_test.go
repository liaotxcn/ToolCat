package controllers_test

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

	"weave/controllers"
	"weave/models"
	"weave/pkg"
	"weave/utils"
)

func setupMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		t.Fatalf("gorm open error: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.LoginHistory{}, &models.AuditLog{}); err != nil {
		t.Fatalf("auto migrate user/audit tables error: %v", err)
	}
	pkg.DB = db
	return db
}

func TestUserRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMemoryDB(t)

	uc := controllers.UserController{}
	r := gin.New()
	r.POST("/register", func(c *gin.Context) { uc.Register(c) })

	payload := `{"username":"alice","password":"secret123","confirm_password":"secret123","email":"alice@example.com"}`
	req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var body struct {
		Message string      `json:"message"`
		User    models.User `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if body.Message != "注册成功" {
		t.Fatalf("expected message=注册成功, got %q", body.Message)
	}
	if body.User.ID == 0 || body.User.Username != "alice" || body.User.Email != "alice@example.com" {
		t.Fatalf("unexpected user in response: %#v", body.User)
	}
	if body.User.Password != "" {
		t.Fatalf("password should not be returned")
	}
}

func TestUserLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupMemoryDB(t)

	// Seed user with hashed password
	hash, err := utils.HashPassword("secret123")
	if err != nil {
		t.Fatalf("hash password error: %v", err)
	}
	user := models.User{Username: "alice", Password: hash, Email: "alice@example.com", TenantID: 1}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("seed user error: %v", err)
	}

	uc := controllers.UserController{}
	r := gin.New()
	r.POST("/login", func(c *gin.Context) { uc.Login(c) })

	payload := `{"username":"alice","password":"secret123"}`
	req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body struct {
		Message      string      `json:"message"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		User         models.User `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if body.Message != "登录成功" {
		t.Fatalf("expected message=登录成功, got %q", body.Message)
	}
	if body.AccessToken == "" || body.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens")
	}
	if body.User.ID == 0 || body.User.Password != "" {
		t.Fatalf("unexpected user in response: %#v", body.User)
	}
}

func TestGetUsers_TenantIsolation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupMemoryDB(t)

	u1 := models.User{Username: "alice", Password: "x", Email: "a@example.com", TenantID: 1}
	u2 := models.User{Username: "bob", Password: "y", Email: "b@example.com", TenantID: 2}
	if err := db.Create(&u1).Error; err != nil {
		t.Fatalf("seed u1 error: %v", err)
	}
	if err := db.Create(&u2).Error; err != nil {
		t.Fatalf("seed u2 error: %v", err)
	}

	uc := controllers.UserController{}
	r := gin.New()
	// Set tenant_id=1 via middleware
	r.Use(func(c *gin.Context) { c.Set("tenant_id", uint(1)); c.Next() })
	r.GET("/users", func(c *gin.Context) { uc.GetUsers(c) })

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var users []models.User
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if len(users) != 1 || users[0].TenantID != 1 || users[0].Username != "alice" {
		t.Fatalf("expected 1 user for tenant 1, got %#v", users)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupMemoryDB(t)

	uc := controllers.UserController{}
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("tenant_id", uint(1)); c.Next() })
	r.GET("/users/:id", func(c *gin.Context) { uc.GetUser(c) })

	req, _ := http.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if body["message"] != "User not found" {
		t.Fatalf("expected message 'User not found', got %#v", body["message"])
	}
}
