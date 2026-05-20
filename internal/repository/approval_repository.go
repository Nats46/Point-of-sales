package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "point-of-sales/internal/model"

    "github.com/jmoiron/sqlx"
)

type ApprovalRepository interface {
    Create(ctx context.Context, approval *model.Approval) error
    GetByID(ctx context.Context, id int64) (*model.Approval, error)
    List(ctx context.Context) ([]*model.Approval, error)
}

type approvalRepository struct {
    db *sqlx.DB
}

func NewApprovalRepository(db *sqlx.DB) ApprovalRepository {
    return &approvalRepository{db: db}
}

func (r *approvalRepository) Create(ctx context.Context, approval *model.Approval) error {
    query := `
        INSERT INTO approvals (transaction, requester, approver)
        VALUES ($1, $2, $3)
        RETURNING id
    `
    err := r.db.QueryRowContext(ctx, query, approval.Transaction, approval.Requester, approval.Approver).Scan(&approval.ID)
    if err != nil {
        return fmt.Errorf("failed to create approval: %w", err)
    }
    return nil
}

func (r *approvalRepository) GetByID(ctx context.Context, id int64) (*model.Approval, error) {
    var approval model.Approval
    query := `
        SELECT id, transaction, requester, approver
        FROM approvals
        WHERE id = $1
    `
    err := r.db.GetContext(ctx, &approval, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("approval not found: %w", err)
        }
        return nil, fmt.Errorf("failed to get approval by id: %w", err)
    }
    return &approval, nil
}

func (r *approvalRepository) List(ctx context.Context) ([]*model.Approval, error) {
    var approvals []*model.Approval
    query := `
        SELECT id, transaction, requester, approver
        FROM approvals
        ORDER BY id ASC
    `
    err := r.db.SelectContext(ctx, &approvals, query)
    if err != nil {
        return nil, fmt.Errorf("failed to list approvals: %w", err)
    }
    return approvals, nil
}
