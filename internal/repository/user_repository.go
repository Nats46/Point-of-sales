package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "point-of-sales/internal/model"

    "github.com/jmoiron/sqlx"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id int64) (*model.User, error)
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context) ([]*model.User, error)
}

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users (username, password, spv, "group")
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Spv, user.Group).Scan(&user.ID)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    query := `
        SELECT id, username, password, spv, "group"
        FROM users
        WHERE id = $1
    `
    err := r.db.GetContext(ctx, &user, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }
    return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    query := `
        SELECT id, username, password, spv, "group"
        FROM users
        WHERE username = $1
    `
    err := r.db.GetContext(ctx, &user, query, username)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get user by username: %w", err)
    }
    return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
    query := `
        UPDATE users
        SET username = $1, password = $2, spv = $3, "group" = $4
        WHERE id = $5
    `
    res, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Spv, user.Group, user.ID)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("user not found for update")
    }
    return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
    query := `
        DELETE FROM users
        WHERE id = $1
    `
    res, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    if rowsAffected == 0 {
        return fmt.Errorf("user not found for deletion")
    }
    return nil
}

func (r *userRepository) List(ctx context.Context) ([]*model.User, error) {
    var users []*model.User
    query := `
        SELECT id, username, password, spv, "group"
        FROM users
        ORDER BY id ASC
    `
    err := r.db.SelectContext(ctx, &users, query)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }
    return users, nil
}
