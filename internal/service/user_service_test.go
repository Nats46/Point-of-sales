package service_test

import (
    "context"
    "log"
    "os"
    "testing"
    "point-of-sales/config"
    "point-of-sales/db"
    "point-of-sales/internal/repository"
    "point-of-sales/internal/service"

    "github.com/jmoiron/sqlx"
    "github.com/joho/godotenv"
)

var testDB *sqlx.DB
var jwtSecret = "test_jwt_secret_key"

func TestMain(m *testing.M) {
    // Load .env relative to internal/service package
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

func TestUserService(t *testing.T) {
    cleanTables(t)
    repo := repository.NewUserRepository(testDB)
    svc := service.NewUserService(repo, jwtSecret)
    ctx := context.Background()

    // 1. Test Register
    username := "newuser1"
    password := "my_plain_password_123"
    group := "admin"

    user, err := svc.Register(ctx, username, password, group, nil)
    if err != nil {
        t.Fatalf("Failed to register user: %v", err)
    }

    if user.Username != username {
        t.Errorf("Expected username %s, got %s", username, user.Username)
    }
    if user.Password == password {
        t.Errorf("Password should have been encrypted/hashed in the database")
    }

    // 2. Test Register Duplicate Username
    _, err = svc.Register(ctx, username, "other_pass", group, nil)
    if err == nil {
        t.Errorf("Expected register duplicate username to fail, but succeeded")
    }

    // 3. Test Login Success
    token, err := svc.Login(ctx, username, password)
    if err != nil {
        t.Fatalf("Failed to login: %v", err)
    }
    if token == "" {
        t.Errorf("Expected a valid JWT token string, got empty string")
    }

    // 4. Test Login Fail - Wrong Password
    _, err = svc.Login(ctx, username, "wrong_password")
    if err == nil {
        t.Errorf("Expected login with wrong password to fail, but succeeded")
    }

    // 5. Test Login Fail - Non-existent User
    _, err = svc.Login(ctx, "nonexistent", password)
    if err == nil {
        t.Errorf("Expected login with non-existent user to fail, but succeeded")
    }
}
