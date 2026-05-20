package handler_test

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "point-of-sales/config"
    "point-of-sales/db"
    "point-of-sales/internal/router"

    "github.com/jmoiron/sqlx"
    "github.com/joho/godotenv"
)

var testDB *sqlx.DB
var jwtSecret = "test_jwt_secret_key"

func TestMain(m *testing.M) {
    // Load .env relative to internal/handler package
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Printf("Warning: error loading .env file from root: %v", err)
    }

    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Config load failed for test: %v", err)
    }

    testDB, err = db.Connect(cfg.DSN)
    if err != nil {
        log.Fatalf("DB connect failed for test: %v", err)
    }
    defer testDB.Close()

    // Run migrations to ensure schema is present
    if err := db.RunMigrations(testDB); err != nil {
        log.Fatalf("Migration failed for test: %v", err)
    }

    os.Exit(m.Run())
}

func cleanTables(t *testing.T) {
    t.Helper()
    _, err := testDB.Exec("TRUNCATE TABLE approvals, users RESTART IDENTITY CASCADE")
    if err != nil {
        t.Fatalf("Failed to clean tables: %v", err)
    }
}

func TestUserHandler_RegisterAndLogin(t *testing.T) {
    cleanTables(t)
    r := router.SetupRouter(testDB, jwtSecret)

    // 1. Test Register Endpoint
    registerPayload := map[string]interface{}{
        "username": "handler_user",
        "password": "securepassword",
        "group":    "cashier",
    }
    jsonPayload, _ := json.Marshal(registerPayload)

    req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonPayload))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusCreated, w.Code, w.Body.String())
    }

    // 2. Test Register Validation Failure (missing group)
    invalidRegisterPayload := map[string]interface{}{
        "username": "invalid_user",
        "password": "securepassword",
    }
    jsonPayload, _ = json.Marshal(invalidRegisterPayload)
    req, _ = http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonPayload))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusBadRequest {
        t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusBadRequest, w.Code, w.Body.String())
    }

    // 3. Test Login Endpoint Success
    loginPayload := map[string]interface{}{
        "username": "handler_user",
        "password": "securepassword",
    }
    jsonPayload, _ = json.Marshal(loginPayload)
    req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("Expected status code %d, got %d. Response: %s", http.StatusOK, w.Code, w.Body.String())
    }

    var loginResponse map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &loginResponse); err != nil {
        t.Fatalf("Failed to parse login response: %v", err)
    }

    token, exists := loginResponse["token"]
    if !exists || token == "" {
        t.Errorf("Login response should contain a valid token")
    }

    // 4. Test Login Endpoint Unauthorized
    wrongLoginPayload := map[string]interface{}{
        "username": "handler_user",
        "password": "wrongpassword",
    }
    jsonPayload, _ = json.Marshal(wrongLoginPayload)
    req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonPayload))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Errorf("Expected status code %d, got %d. Response: %s", http.StatusUnauthorized, w.Code, w.Body.String())
    }
}
