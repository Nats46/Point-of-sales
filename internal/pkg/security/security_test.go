package security_test

import (
    "testing"
    "point-of-sales/internal/pkg/security"
)

func TestPasswordHashing(t *testing.T) {
    password := "my_secure_password_123"

    // 1. Hash password
    hash, err := security.HashPassword(password)
    if err != nil {
        t.Fatalf("Failed to hash password: %v", err)
    }

    if hash == password {
        t.Errorf("Hashed password should not be equal to plain text password")
    }

    // 2. Validate correct password
    if !security.CheckPasswordHash(password, hash) {
        t.Errorf("CheckPasswordHash failed for correct password")
    }

    // 3. Validate incorrect password
    if security.CheckPasswordHash("wrong_password", hash) {
        t.Errorf("CheckPasswordHash should have failed for incorrect password")
    }
}

func TestJWT(t *testing.T) {
    secret := "test_secret_key_987654321"
    userID := int64(42)
    username := "testuser"
    group := "admin"

    // 1. Generate Token
    token, err := security.GenerateToken(userID, username, group, secret)
    if err != nil {
        t.Fatalf("Failed to generate token: %v", err)
    }

    if token == "" {
        t.Errorf("Expected token string to be non-empty")
    }

    // 2. Validate Token
    claims, err := security.ValidateToken(token, secret)
    if err != nil {
        t.Fatalf("Failed to validate token: %v", err)
    }

    if claims.UserID != userID {
        t.Errorf("Expected user ID %d, got %d", userID, claims.UserID)
    }
    if claims.Username != username {
        t.Errorf("Expected username %s, got %s", username, claims.Username)
    }
    if claims.Group != group {
        t.Errorf("Expected group %s, got %s", group, claims.Group)
    }

    // 3. Validate Token with wrong secret
    _, err = security.ValidateToken(token, "wrong_secret_key")
    if err == nil {
        t.Errorf("Expected validation to fail with wrong secret, but succeeded")
    }
}
