package service

import (
    "context"
    "errors"
    "fmt"
    "point-of-sales/internal/model"
    "point-of-sales/internal/pkg/security"
    "point-of-sales/internal/repository"
)

type UserService interface {
    Register(ctx context.Context, username, password, group string, spv *int64) (*model.User, error)
    Login(ctx context.Context, username, password string) (string, error)
}

type userService struct {
    repo      repository.UserRepository
    jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
    return &userService{
        repo:      repo,
        jwtSecret: jwtSecret,
    }
}

func (s *userService) Register(ctx context.Context, username, password, group string, spv *int64) (*model.User, error) {
    if username == "" || password == "" || group == "" {
        return nil, errors.New("username, password and group are required")
    }

    // Check if user already exists
    existingUser, err := s.repo.GetByUsername(ctx, username)
    if err == nil && existingUser != nil {
        return nil, errors.New("username is already taken")
    }

    // Hash the password
    hashedPassword, err := security.HashPassword(password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    user := &model.User{
        Username: username,
        Password: hashedPassword,
        Group:    group,
        Spv:      spv,
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
    if username == "" || password == "" {
        return "", errors.New("username and password are required")
    }

    user, err := s.repo.GetByUsername(ctx, username)
    if err != nil {
        return "", errors.New("invalid username or password")
    }

    // Verify password
    if !security.CheckPasswordHash(password, user.Password) {
        return "", errors.New("invalid username or password")
    }

    // Generate JWT token
    token, err := security.GenerateToken(user.ID, user.Username, user.Group, s.jwtSecret)
    if err != nil {
        return "", fmt.Errorf("failed to generate token: %w", err)
    }

    return token, nil
}
