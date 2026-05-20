package security

import (
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Group    string `json:"group"`
    jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a given user
func GenerateToken(userID int64, username, group string, secret string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Username: username,
        Group:    group,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }

    return tokenStr, nil
}

// ValidateToken validates a JWT token and returns claims if valid
func ValidateToken(tokenStr string, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }

    return claims, nil
}
