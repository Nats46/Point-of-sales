package repository_test

import (
    "context"
    "log"
    "os"
    "testing"
    "point-of-sales/config"
    "point-of-sales/db"
    "point-of-sales/internal/model"
    "point-of-sales/internal/repository"

    "github.com/jmoiron/sqlx"
    "github.com/joho/godotenv"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
    // Load .env relative to the test location (which is in internal/repository)
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

func TestUserRepository(t *testing.T) {
    cleanTables(t)
    repo := repository.NewUserRepository(testDB)
    ctx := context.Background()

    // 1. Create a Supervisor User
    spvUser := &model.User{
        Username: "supervisor1",
        Password: "hashedpassword123",
        Group:    "supervisor",
    }
    err := repo.Create(ctx, spvUser)
    if err != nil {
        t.Fatalf("Failed to create supervisor user: %v", err)
    }
    if spvUser.ID == 0 {
        t.Errorf("Expected non-zero ID for supervisor user")
    }

    // 2. Create a Regular User with Supervisor
    regUser := &model.User{
        Username: "cashier1",
        Password: "hashedpassword456",
        Spv:      &spvUser.ID,
        Group:    "cashier",
    }
    err = repo.Create(ctx, regUser)
    if err != nil {
        t.Fatalf("Failed to create regular user: %v", err)
    }

    // 3. Get User By ID
    fetchedUser, err := repo.GetByID(ctx, regUser.ID)
    if err != nil {
        t.Fatalf("Failed to get user by ID: %v", err)
    }
    if fetchedUser.Username != regUser.Username {
        t.Errorf("Expected username %s, got %s", regUser.Username, fetchedUser.Username)
    }
    if fetchedUser.Spv == nil || *fetchedUser.Spv != spvUser.ID {
        t.Errorf("Expected supervisor ID %d, got %v", spvUser.ID, fetchedUser.Spv)
    }

    // 4. Get User By Username
    fetchedByUsername, err := repo.GetByUsername(ctx, "cashier1")
    if err != nil {
        t.Fatalf("Failed to get user by username: %v", err)
    }
    if fetchedByUsername.ID != regUser.ID {
        t.Errorf("Expected user ID %d, got %d", regUser.ID, fetchedByUsername.ID)
    }

    // 5. Update User
    fetchedUser.Group = "admin"
    err = repo.Update(ctx, fetchedUser)
    if err != nil {
        t.Fatalf("Failed to update user: %v", err)
    }

    updatedUser, err := repo.GetByID(ctx, fetchedUser.ID)
    if err != nil {
        t.Fatalf("Failed to get updated user: %v", err)
    }
    if updatedUser.Group != "admin" {
        t.Errorf("Expected group to be updated to admin, got %s", updatedUser.Group)
    }

    // 6. List Users
    users, err := repo.List(ctx)
    if err != nil {
        t.Fatalf("Failed to list users: %v", err)
    }
    if len(users) != 2 {
        t.Errorf("Expected 2 users, got %d", len(users))
    }

    // 7. Delete User
    err = repo.Delete(ctx, regUser.ID)
    if err != nil {
        t.Fatalf("Failed to delete user: %v", err)
    }

    _, err = repo.GetByID(ctx, regUser.ID)
    if err == nil {
        t.Errorf("Expected error getting deleted user, got nil")
    }
}

func TestApprovalRepository(t *testing.T) {
    cleanTables(t)
    userRepo := repository.NewUserRepository(testDB)
    approvalRepo := repository.NewApprovalRepository(testDB)
    ctx := context.Background()

    // Create requester and approver users
    requester := &model.User{
        Username: "requester1",
        Password: "pass",
        Group:    "cashier",
    }
    approver := &model.User{
        Username: "approver1",
        Password: "pass",
        Group:    "supervisor",
    }
    if err := userRepo.Create(ctx, requester); err != nil {
        t.Fatalf("Failed to create requester user: %v", err)
    }
    if err := userRepo.Create(ctx, approver); err != nil {
        t.Fatalf("Failed to create approver user: %v", err)
    }

    // Create Approval
    approval := &model.Approval{
        Transaction: "TX-100293",
        Requester:   requester.ID,
        Approver:    approver.ID,
    }
    err := approvalRepo.Create(ctx, approval)
    if err != nil {
        t.Fatalf("Failed to create approval record: %v", err)
    }
    if approval.ID == 0 {
        t.Errorf("Expected non-zero approval ID")
    }

    // Get Approval By ID
    fetched, err := approvalRepo.GetByID(ctx, approval.ID)
    if err != nil {
        t.Fatalf("Failed to get approval by ID: %v", err)
    }
    if fetched.Transaction != approval.Transaction {
        t.Errorf("Expected transaction %s, got %s", approval.Transaction, fetched.Transaction)
    }
    if fetched.Requester != approval.Requester || fetched.Approver != approval.Approver {
        t.Errorf("Expected requester %d and approver %d, got %d and %d", approval.Requester, approval.Approver, fetched.Requester, fetched.Approver)
    }

    // List Approvals
    list, err := approvalRepo.List(ctx)
    if err != nil {
        t.Fatalf("Failed to list approvals: %v", err)
    }
    if len(list) != 1 {
        t.Errorf("Expected 1 approval record, got %d", len(list))
    }
}
