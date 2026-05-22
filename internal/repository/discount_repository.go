package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"point-of-sales/internal/model"

	"github.com/jmoiron/sqlx"
)

type DiscountRepository interface {
	Create(ctx context.Context, discount *model.Discount) error
	GetByID(ctx context.Context, id int64) (*model.Discount, error)
	GetByCode(ctx context.Context, code string) (*model.Discount, error)
	Update(ctx context.Context, discount *model.Discount) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, pageSize int64, filter string) ([]*model.Discount, error)
}

type discountRepository struct {
	db *sqlx.DB
}

func NewDiscountRepository(db *sqlx.DB) DiscountRepository {
	return &discountRepository{db: db}
}

func (r *discountRepository) Create(ctx context.Context, discount *model.Discount) error {
	query := `
		INSERT INTO discounts (discount_code, value, discount_type, start_date, end_date, is_active, created_at, updated_at, created_by, updated_by)
		VALUES (:discount_code, :value, :discount_type, :start_date, :end_date, :is_active, :created_at, :updated_at, :created_by, :updated_by)
		RETURNING id
	`
	rows, err := r.db.NamedQueryContext(ctx, query, discount)
	if err != nil {
		return fmt.Errorf("failed to create discount: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&discount.Id)
	}
	return nil
}

func (r *discountRepository) GetByID(ctx context.Context, id int64) (*model.Discount, error) {
	var discount model.Discount
	query := `SELECT * FROM discounts WHERE id = $1`
	err := r.db.GetContext(ctx, &discount, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("discount not found")
		}
		return nil, fmt.Errorf("failed to get discount by id: %w", err)
	}
	return &discount, nil
}

func (r *discountRepository) GetByCode(ctx context.Context, code string) (*model.Discount, error) {
	var discount model.Discount
	query := `SELECT * FROM discounts WHERE discount_code = $1`
	err := r.db.GetContext(ctx, &discount, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("discount not found")
		}
		return nil, fmt.Errorf("failed to get discount by code: %w", err)
	}
	return &discount, nil
}

func (r *discountRepository) Update(ctx context.Context, discount *model.Discount) error {
	query := `
		UPDATE discounts
		SET discount_code = :discount_code, value = :value, discount_type = :discount_type, 
			start_date = :start_date, end_date = :end_date, is_active = :is_active,
			updated_at = :updated_at, updated_by = :updated_by
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, discount)
	if err != nil {
		return fmt.Errorf("failed to update discount: %w", err)
	}
	return nil
}

func (r *discountRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM discounts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete discount: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("discount not found for deletion")
	}
	return nil
}

func (r *discountRepository) List(ctx context.Context, page, pageSize int64, filter string) ([]*model.Discount, error) {
	var discounts []*model.Discount
	query := `SELECT * FROM discounts WHERE discount_code ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3`
	filterParam := "%" + filter + "%"
	offset := (page - 1) * pageSize
	err := r.db.SelectContext(ctx, &discounts, query, filterParam, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list discounts: %w", err)
	}
	return discounts, nil
}