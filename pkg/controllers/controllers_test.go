package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"ETM/pkg/app"
	"ETM/pkg/controllers"
	"ETM/pkg/models"
	"ETM/pkg/types"
	"ETM/pkg/utils"
	"ETM/pkg/vars"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
	vars.SecretKey = "test-secret-key-for-controllers"
}

// ── Test helpers ──────────────────────────────────────────────────────────────

// sqliteSchema mirrors the production schema without PostgreSQL-specific
// DEFAULT expressions (gen_random_uuid()) that SQLite cannot parse.
// BeforeCreate hooks handle UUID generation in Go, so no DB default is needed.
var sqliteSchema = []string{
	`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		uuid TEXT, username TEXT, password TEXT, gid INTEGER,
		is_admin TEXT, telegram TEXT, browser TEXT, email TEXT,
		oidc_subject TEXT, oidc_provider TEXT, active_category_id INTEGER
	)`,
	`CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		name TEXT, owner_id INTEGER
	)`,
	`CREATE TABLE IF NOT EXISTS user_groups (
		user_id INTEGER NOT NULL, group_id INTEGER NOT NULL, role TEXT,
		PRIMARY KEY (user_id, group_id)
	)`,
	`CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		name TEXT, color TEXT, user_id INTEGER, active INTEGER
	)`,
	`CREATE TABLE IF NOT EXISTS category_groups (
		category_id INTEGER NOT NULL, group_id INTEGER NOT NULL,
		PRIMARY KEY (category_id, group_id)
	)`,
	`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		name TEXT, comment TEXT, link TEXT, is_completed INTEGER, is_started INTEGER,
		is_back_log INTEGER, category_id INTEGER, priority INTEGER, urgency INTEGER,
		due_date DATETIME, user_id INTEGER
	)`,
	`CREATE TABLE IF NOT EXISTS keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		pubkey TEXT, privkey TEXT
	)`,
	`CREATE TABLE IF NOT EXISTS devices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		uuid TEXT, user_id INTEGER, device_id TEXT, device_name TEXT,
		trusted INTEGER, last_used_at DATETIME
	)`,
	`CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME, updated_at DATETIME, deleted_at DATETIME,
		uuid TEXT, user_id INTEGER, device_id TEXT,
		refresh_token TEXT, country TEXT
	)`,
}

func newTestApp(t *testing.T) *app.App {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	for _, sql := range sqliteSchema {
		if err := gdb.Exec(sql).Error; err != nil {
			t.Fatalf("create table: %v\nSQL: %s", err, sql)
		}
	}
	db := models.DB{DB: *gdb}
	return &app.App{DB: db, Logger: zerolog.Nop()}
}

func newRouter(a *app.App, handlers ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.Handle("GET", "/test", handlers...)
	r.Handle("POST", "/test", handlers...)
	r.Handle("DELETE", "/test", handlers...)
	return r
}

func bearerFor(userID uint, userUUID uuid.UUID) string {
	claims := jwt.MapClaims{
		"sub":  float64(userID),
		"uuid": userUUID.String(),
		"exp":  time.Now().Add(30 * time.Minute).Unix(),
		"iss":  "etm",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(vars.SecretKey))
	return "Bearer " + signed
}

func createUser(t *testing.T, a *app.App, username, password string) models.Users {
	t.Helper()
	hash, err := utils.GenerateHashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if _, err := a.DB.CreateUser(types.UserBody{Username: username, Password: hash, Email: username + "@test.com"}, false); err != nil {
		t.Fatalf("create user: %v", err)
	}
	var user models.Users
	a.DB.Where("username = ?", username).First(&user)
	return user
}

func jsonBody(t *testing.T, v interface{}) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewReader(b)
}

// ── IsAuthorized middleware ───────────────────────────────────────────────────

func TestIsAuthorized_NoHeader(t *testing.T) {
	a := newTestApp(t)
	r := newRouter(a, controllers.IsAuthorized(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))
	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestIsAuthorized_InvalidFormat(t *testing.T) {
	a := newTestApp(t)
	r := newRouter(a, controllers.IsAuthorized(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "NotBearer abc")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestIsAuthorized_InvalidSignature(t *testing.T) {
	a := newTestApp(t)
	r := newRouter(a, controllers.IsAuthorized(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	// Sign with a different secret
	claims := jwt.MapClaims{"sub": float64(1), "exp": time.Now().Add(time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("wrong-secret"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestIsAuthorized_AlgNone(t *testing.T) {
	a := newTestApp(t)
	r := newRouter(a, controllers.IsAuthorized(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	claims := jwt.MapClaims{"sub": float64(1), "exp": time.Now().Add(time.Hour).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		t.Error("alg:none token must be rejected")
	}
}

func TestIsAuthorized_Valid(t *testing.T) {
	a := newTestApp(t)
	user := createUser(t, a, "alice", "pass")

	r := newRouter(a, controllers.IsAuthorized(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", bearerFor(user.ID, user.UUID))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("got %d, want %d", w.Code, http.StatusOK)
	}
}

// ── Login ─────────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	a := newTestApp(t)
	createUser(t, a, "alice", "password123")

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.POST("/login", controllers.Login)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/login",
		jsonBody(t, map[string]string{"username": "alice", "password": "password123"})))

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want 200; body: %s", w.Code, w.Body)
	}
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Error("expected non-empty token in response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	a := newTestApp(t)
	createUser(t, a, "alice", "password123")

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.POST("/login", controllers.Login)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/login",
		jsonBody(t, map[string]string{"username": "alice", "password": "wrongpass"})))

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want 400", w.Code)
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	a := newTestApp(t)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.POST("/login", controllers.Login)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/login",
		jsonBody(t, map[string]string{"username": "nobody", "password": "x"})))

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want 400", w.Code)
	}
}

// ── Register ──────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	a := newTestApp(t)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.POST("/register", controllers.Register)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/register",
		jsonBody(t, map[string]string{"username": "newuser", "password": "secret", "email": "n@test.com"})))

	if w.Code != http.StatusCreated {
		t.Errorf("got %d, want 201; body: %s", w.Code, w.Body)
	}
}

func TestRegister_Duplicate(t *testing.T) {
	a := newTestApp(t)
	createUser(t, a, "alice", "pass")

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.POST("/register", controllers.Register)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/register",
		jsonBody(t, map[string]string{"username": "alice", "password": "other"})))

	if w.Code != http.StatusConflict {
		t.Errorf("got %d, want 409", w.Code)
	}
}

// ── Tasks ─────────────────────────────────────────────────────────────────────

func setupTaskRouter(a *app.App) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.POST("/task", controllers.IsAuthorized(), controllers.CreateTask)
	r.POST("/task/:taskId", controllers.IsAuthorized(), controllers.UpdateTask)
	r.GET("/task/:taskId", controllers.IsAuthorized(), controllers.GetTask)
	r.DELETE("/task/:taskId", controllers.IsAuthorized(), controllers.DeleteTask)
	r.GET("/tasks/:categoryId", controllers.IsAuthorized(), controllers.GetTasks)
	return r
}

func TestCreateTask_Success(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	cat, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Cat", Color: "#0f0", Active: true}, owner)

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/task", jsonBody(t, map[string]interface{}{
		"name":       "My task",
		"comment":    "details",
		"duedate":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"categoryid": cat.ID,
	}))
	req.Header.Set("Authorization", bearerFor(owner.ID, owner.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("got %d, want 201; body: %s", w.Code, w.Body)
	}
}

func TestCreateTask_ForbiddenCategory(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	attacker := createUser(t, a, "attacker", "pass")
	cat, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Cat", Color: "#0f0", Active: true}, owner)

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/task", jsonBody(t, map[string]interface{}{
		"name":       "Evil task",
		"duedate":    time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"categoryid": cat.ID,
	}))
	req.Header.Set("Authorization", bearerFor(attacker.ID, attacker.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want 403", w.Code)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	cat, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Cat", Color: "#0f0", Active: true}, owner)
	task := &models.Tasks{Name: "Old", CategoryID: cat.ID, UserID: owner.ID, DueDate: time.Now().Add(time.Hour)}
	_ = a.DB.CreateTask(task)

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/task/"+strconv.Itoa(int(task.ID)),
		jsonBody(t, map[string]interface{}{
			"name":       "Updated",
			"duedate":    time.Now().Add(48 * time.Hour).Format(time.RFC3339),
			"categoryid": cat.ID,
		}))
	req.Header.Set("Authorization", bearerFor(owner.ID, owner.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200; body: %s", w.Code, w.Body)
	}
}

func TestUpdateTask_ForbiddenOtherUser(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	other := createUser(t, a, "other", "pass")
	cat, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Cat", Color: "#0f0", Active: true}, owner)
	task := &models.Tasks{Name: "Secret", CategoryID: cat.ID, UserID: owner.ID, DueDate: time.Now().Add(time.Hour)}
	_ = a.DB.CreateTask(task)

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/task/"+strconv.Itoa(int(task.ID)),
		jsonBody(t, map[string]interface{}{
			"name":       "Hijack",
			"duedate":    time.Now().Add(48 * time.Hour).Format(time.RFC3339),
			"categoryid": cat.ID,
		}))
	req.Header.Set("Authorization", bearerFor(other.ID, other.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want 403", w.Code)
	}
}

func TestDeleteTask_NotOwner(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	other := createUser(t, a, "other", "pass")
	cat, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Cat", Color: "#0f0", Active: true}, owner)
	task := &models.Tasks{Name: "Owned", CategoryID: cat.ID, UserID: owner.ID, DueDate: time.Now().Add(time.Hour)}
	_ = a.DB.CreateTask(task)

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/task/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", bearerFor(other.ID, other.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want 403", w.Code)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/task/9999", nil)
	req.Header.Set("Authorization", bearerFor(owner.ID, owner.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("got %d, want 404", w.Code)
	}
}

// ── Delete category ───────────────────────────────────────────────────────────

func setupCategoryRouter(a *app.App) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.DELETE("/category/:categoryId", controllers.IsAuthorized(), controllers.DeleteCategoryHandler)
	return r
}

func TestDeleteCategory_Owner(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	cat, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Mine", Color: "#f00", Active: true}, owner)

	r := setupCategoryRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/category/"+strconv.Itoa(int(cat.ID)), nil)
	req.Header.Set("Authorization", bearerFor(owner.ID, owner.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200; body: %s", w.Code, w.Body)
	}
}

func TestDeleteCategory_NotOwner(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	other := createUser(t, a, "other", "pass")
	cat, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Mine", Color: "#f00", Active: true}, owner)

	r := setupCategoryRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/category/"+strconv.Itoa(int(cat.ID)), nil)
	req.Header.Set("Authorization", bearerFor(other.ID, other.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want 403", w.Code)
	}
}

// ── Cross-category task move ──────────────────────────────────────────────────

func TestUpdateTask_MoveToOwnCategory(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	src, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Src", Color: "#0f0", Active: true}, owner)
	dst, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Dst", Color: "#00f", Active: true}, owner)
	task := &models.Tasks{Name: "Move me", CategoryID: src.ID, UserID: owner.ID, DueDate: time.Now().Add(time.Hour)}
	_ = a.DB.CreateTask(task)

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/task/"+strconv.Itoa(int(task.ID)),
		jsonBody(t, map[string]interface{}{
			"name":       "Move me",
			"duedate":    time.Now().Add(48 * time.Hour).Format(time.RFC3339),
			"categoryid": dst.ID,
		}))
	req.Header.Set("Authorization", bearerFor(owner.ID, owner.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200; body: %s", w.Code, w.Body)
	}
}

func TestUpdateTask_MoveToForbiddenCategory(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	attacker := createUser(t, a, "attacker", "pass")
	src, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Src", Color: "#0f0", Active: true}, owner)
	dst, _ := a.DB.CreateCategory(types.CategoryBody{Name: "Dst", Color: "#00f", Active: true}, owner)
	task := &models.Tasks{Name: "Owned", CategoryID: src.ID, UserID: owner.ID, DueDate: time.Now().Add(time.Hour)}
	_ = a.DB.CreateTask(task)

	r := setupTaskRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/task/"+strconv.Itoa(int(task.ID)),
		jsonBody(t, map[string]interface{}{
			"name":       "Owned",
			"duedate":    time.Now().Add(48 * time.Hour).Format(time.RFC3339),
			"categoryid": dst.ID,
		}))
	req.Header.Set("Authorization", bearerFor(attacker.ID, attacker.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want 403", w.Code)
	}
}

// ── Token ownership ───────────────────────────────────────────────────────────

func setupTokenRouter(a *app.App) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("App", a); c.Next() })
	r.DELETE("/token/:id", controllers.IsAuthorized(), controllers.DeleteToken)
	r.GET("/token/uuid/:uuid", controllers.IsAuthorized(), controllers.GetTokenByUUID)
	return r
}

func TestDeleteToken_WrongUser(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")
	attacker := createUser(t, a, "attacker", "pass")

	tok := &models.Tokens{UserID: owner.ID, RefreshToken: "secret"}
	_ = a.DB.CreateToken(tok)

	r := setupTokenRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/token/"+strconv.Itoa(int(tok.ID)), nil)
	req.Header.Set("Authorization", bearerFor(attacker.ID, attacker.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("got %d, want 403", w.Code)
	}
	// Token must still exist
	if _, err := a.DB.GetTokenByID(tok.ID); err != nil {
		t.Error("token was deleted by wrong user — ownership check failed")
	}
}

func TestDeleteToken_Owner(t *testing.T) {
	a := newTestApp(t)
	owner := createUser(t, a, "owner", "pass")

	tok := &models.Tokens{UserID: owner.ID, RefreshToken: "mine"}
	_ = a.DB.CreateToken(tok)

	r := setupTokenRouter(a)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/token/"+strconv.Itoa(int(tok.ID)), nil)
	req.Header.Set("Authorization", bearerFor(owner.ID, owner.UUID))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200; body: %s", w.Code, w.Body)
	}
}
