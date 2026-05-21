package repository

import (
	"context"
	"point-of-sales/internal/model"

	"github.com/jmoiron/sqlx"
)

type BatchRepository interface {
	InsertBatch(ctx context.Context, b *model.Batch) error
	GetBatch(ctx context.Context, id int64) (*model.Batch, error)
	UpdateBatch(ctx context.Context, b *model.Batch) error
	DeleteBatch(ctx context.Context, id int64) error
	ListBatches(ctx context.Context, page, pageSize int64, filter string) ([]model.Batch, error)
}

type batchRepository struct {
	DB *sqlx.DB
}

func NewBatchRepository(db *sqlx.DB) BatchRepository {
	return &batchRepository{DB: db}
}

func (br *batchRepository) InsertBatch(ctx context.Context, b *model.Batch) error {
	query := `
		INSERT INTO batches (item_code, item_name, created_at, updated_at, created_by, updated_by, stock, status, batch_number, expired_date, stock_date)
		VALUES (:item_code, :item_name, :created_at, :updated_at, :created_by, :updated_by, :stock, :status, :batch_number, :expired_date, :stock_date)
		RETURNING id
	`
	rows, err := br.DB.NamedQueryContext(ctx, query, b)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&b.Id)
	}
	return nil
}

func (br *batchRepository) GetBatch(ctx context.Context, id int64) (*model.Batch, error) {
	var batch model.Batch
	query := `SELECT * FROM batches WHERE id = $1`
	err := br.DB.GetContext(ctx, &batch, query, id)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (br *batchRepository) UpdateBatch(ctx context.Context, b *model.Batch) error {
	query := `
		UPDATE batches
		SET item_code = :item_code, item_name = :item_name, updated_at = :updated_at, updated_by = :updated_by,
			stock = :stock, status = :status, batch_number = :batch_number, expired_date = :expired_date, stock_date = :stock_date
		WHERE id = :id
	`
	_, err := br.DB.NamedExecContext(ctx, query, b)
	return err
}

func (br *batchRepository) DeleteBatch(ctx context.Context, id int64) error {
	query := `DELETE FROM batches WHERE id = $1`
	_, err := br.DB.ExecContext(ctx, query, id)
	return err
}

func (br *batchRepository) ListBatches(ctx context.Context, page, pageSize int64, filter string) ([]model.Batch, error) {
	var batches []model.Batch
	query := `SELECT * FROM batches WHERE item_code ILIKE $1 OR item_name ILIKE $1 OR batch_number ILIKE $1 LIMIT $2 OFFSET $3`
	filterParam := "%" + filter + "%"
	offset := (page - 1) * pageSize
	err := br.DB.SelectContext(ctx, &batches, query, filterParam, pageSize, offset)
	return batches, err
}