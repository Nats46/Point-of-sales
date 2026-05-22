package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"point-of-sales/internal/model"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	CreateSalesHeader(ctx context.Context, header *model.SalesHeader) error
	GetSalesHeaderByID(ctx context.Context, id int64) (*model.SalesHeader, error)
	GetSalesHeaderByNo(ctx context.Context, no string) (*model.SalesHeader, error)
	UpdateSalesHeader(ctx context.Context, header *model.SalesHeader) error
	DeleteSalesHeader(ctx context.Context, id int64) error
	ListSalesHeaders(ctx context.Context, page, pageSize int64, filter string) ([]*model.SalesHeader, error)
}

type SalesDetailRepository interface {
	CreateSalesDetail(ctx context.Context, detail *model.SalesDetail) error
	GetSalesDetailByID(ctx context.Context, id int64) (*model.SalesDetail, error)
	ListSalesDetailsByTransaction(ctx context.Context, transactionId int64) ([]*model.SalesDetail, error)
	UpdateSalesDetail(ctx context.Context, detail *model.SalesDetail) error
	DeleteSalesDetail(ctx context.Context, id int64) error
}

type transactionRepository struct {
	db *sqlx.DB
}

type salesDetailRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func NewSalesDetailRepository(db *sqlx.DB) SalesDetailRepository {
	return &salesDetailRepository{db: db}
}

func (r *transactionRepository) CreateSalesHeader(ctx context.Context, header *model.SalesHeader) error {
	query := `
		INSERT INTO sales_headers (transaction_no, customer_id, cashier_id, subtotal, discount, tax, grand_total, payment_method, status, created_at, updated_at, created_by, updated_by)
		VALUES (:transaction_no, :customer_id, :cashier_id, :subtotal, :discount, :tax, :grand_total, :payment_method, :status, :created_at, :updated_at, :created_by, :updated_by)
		RETURNING transaction_id
	`
	rows, err := r.db.NamedQueryContext(ctx, query, header)
	if err != nil {
		return fmt.Errorf("failed to create sales header: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&header.TransactionId)
	}
	return nil
}

func (r *transactionRepository) GetSalesHeaderByID(ctx context.Context, id int64) (*model.SalesHeader, error) {
	var header model.SalesHeader
	query := `SELECT * FROM sales_headers WHERE transaction_id = $1`
	err := r.db.GetContext(ctx, &header, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sales header not found")
		}
		return nil, fmt.Errorf("failed to get sales header by id: %w", err)
	}
	return &header, nil
}

func (r *transactionRepository) GetSalesHeaderByNo(ctx context.Context, no string) (*model.SalesHeader, error) {
	var header model.SalesHeader
	query := `SELECT * FROM sales_headers WHERE transaction_no = $1`
	err := r.db.GetContext(ctx, &header, query, no)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sales header not found")
		}
		return nil, fmt.Errorf("failed to get sales header by no: %w", err)
	}
	return &header, nil
}

func (r *transactionRepository) UpdateSalesHeader(ctx context.Context, header *model.SalesHeader) error {
	query := `
		UPDATE sales_headers
		SET customer_id = :customer_id, cashier_id = :cashier_id, subtotal = :subtotal,
			discount = :discount, tax = :tax, grand_total = :grand_total, 
			payment_method = :payment_method, status = :status,
			updated_at = :updated_at, updated_by = :updated_by
		WHERE transaction_id = :transaction_id
	`
	_, err := r.db.NamedExecContext(ctx, query, header)
	if err != nil {
		return fmt.Errorf("failed to update sales header: %w", err)
	}
	return nil
}

func (r *transactionRepository) DeleteSalesHeader(ctx context.Context, id int64) error {
	query := `DELETE FROM sales_headers WHERE transaction_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sales header: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("sales header not found for deletion")
	}
	return nil
}

func (r *transactionRepository) ListSalesHeaders(ctx context.Context, page, pageSize int64, filter string) ([]*model.SalesHeader, error) {
	var headers []*model.SalesHeader
	query := `SELECT * FROM sales_headers WHERE transaction_no ILIKE $1 ORDER BY transaction_id LIMIT $2 OFFSET $3`
	filterParam := "%" + filter + "%"
	offset := (page - 1) * pageSize
	err := r.db.SelectContext(ctx, &headers, query, filterParam, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sales headers: %w", err)
	}
	return headers, nil
}

func (r *salesDetailRepository) CreateSalesDetail(ctx context.Context, detail *model.SalesDetail) error {
	query := `
		INSERT INTO sales_details (transaction_id, item_id, item_name, item_code, quantity, price, subtotal, discount, tax, grand_total, created_at, updated_at, created_by, updated_by)
		VALUES (:transaction_id, :item_id, :item_name, :item_code, :quantity, :price, :subtotal, :discount, :tax, :grand_total, :created_at, :updated_at, :created_by, :updated_by)
		RETURNING sales_detail_id
	`
	rows, err := r.db.NamedQueryContext(ctx, query, detail)
	if err != nil {
		return fmt.Errorf("failed to create sales detail: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&detail.SalesDetailId)
	}
	return nil
}

func (r *salesDetailRepository) GetSalesDetailByID(ctx context.Context, id int64) (*model.SalesDetail, error) {
	var detail model.SalesDetail
	query := `SELECT * FROM sales_details WHERE sales_detail_id = $1`
	err := r.db.GetContext(ctx, &detail, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sales detail not found")
		}
		return nil, fmt.Errorf("failed to get sales detail by id: %w", err)
	}
	return &detail, nil
}

func (r *salesDetailRepository) ListSalesDetailsByTransaction(ctx context.Context, transactionId int64) ([]*model.SalesDetail, error) {
	var details []*model.SalesDetail
	query := `SELECT * FROM sales_details WHERE transaction_id = $1 ORDER BY sales_detail_id`
	err := r.db.SelectContext(ctx, &details, query, transactionId)
	if err != nil {
		return nil, fmt.Errorf("failed to list sales details by transaction: %w", err)
	}
	return details, nil
}

func (r *salesDetailRepository) UpdateSalesDetail(ctx context.Context, detail *model.SalesDetail) error {
	query := `
		UPDATE sales_details
		SET item_id = :item_id, item_name = :item_name, item_code = :item_code,
			quantity = :quantity, price = :price, subtotal = :subtotal, discount = :discount,
			tax = :tax, grand_total = :grand_total, updated_at = :updated_at, updated_by = :updated_by
		WHERE sales_detail_id = :sales_detail_id
	`
	_, err := r.db.NamedExecContext(ctx, query, detail)
	if err != nil {
		return fmt.Errorf("failed to update sales detail: %w", err)
	}
	return nil
}

func (r *salesDetailRepository) DeleteSalesDetail(ctx context.Context, id int64) error {
	query := `DELETE FROM sales_details WHERE sales_detail_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sales detail: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("sales detail not found for deletion")
	}
	return nil
}