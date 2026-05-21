package service

import (
	"context"
	"errors"
	"time"

	"point-of-sales/internal/model"
	"point-of-sales/internal/repository"
)

type BatchService interface {
	CreateBatch(ctx context.Context, itemCode, itemName, batchNumber string, stock float64, status string, expiredDate, stockDate time.Time, createdBy string) (*model.Batch, error)
	GetBatch(ctx context.Context, id int64) (*model.Batch, error)
	UpdateBatch(ctx context.Context, id int64, itemCode, itemName, batchNumber string, stock float64, status string, expiredDate, stockDate time.Time, updatedBy string) (*model.Batch, error)
	DeleteBatch(ctx context.Context, id int64) error
	ListBatches(ctx context.Context, page, pageSize int64, filter string) ([]model.Batch, error)
}

type batchService struct {
	repo repository.BatchRepository
}

func NewBatchService(repo repository.BatchRepository) BatchService {
	return &batchService{repo: repo}
}

func (s *batchService) CreateBatch(ctx context.Context, itemCode, itemName, batchNumber string, stock float64, status string, expiredDate, stockDate time.Time, createdBy string) (*model.Batch, error) {
	if itemCode == "" || itemName == "" || batchNumber == "" {
		return nil, errors.New("item code, item name and batch number are required")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	if expiredDate.IsZero() {
		return nil, errors.New("expired date is required")
	}
	if stockDate.IsZero() {
		return nil, errors.New("stock date is required")
	}
	batch := &model.Batch{
		ItemCode:    itemCode,
		ItemName:    itemName,
		BatchNumber: batchNumber,
		Stock:       stock,
		Status:      status,
		ExpiredDate: expiredDate,
		StockDate:   stockDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.InsertBatch(ctx, batch); err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *batchService) GetBatch(ctx context.Context, id int64) (*model.Batch, error) {
	return s.repo.GetBatch(ctx, id)
}

func (s *batchService) UpdateBatch(ctx context.Context, id int64, itemCode, itemName, batchNumber string, stock float64, status string, expiredDate, stockDate time.Time, updatedBy string) (*model.Batch, error) {
	batch, err := s.repo.GetBatch(ctx, id)
	if err != nil {
		return nil, err
	}
	if itemCode != "" {
		batch.ItemCode = itemCode
	}
	if itemName != "" {
		batch.ItemName = itemName
	}
	if batchNumber != "" {
		batch.BatchNumber = batchNumber
	}
	if stock >= 0 {
		batch.Stock = stock
	}
	if status != "" {
		batch.Status = status
	}
	if !expiredDate.IsZero() {
		batch.ExpiredDate = expiredDate
	}
	if !stockDate.IsZero() {
		batch.StockDate = stockDate
	}
	batch.UpdatedAt = time.Now()
	batch.UpdatedBy = updatedBy
	if err := s.repo.UpdateBatch(ctx, batch); err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *batchService) DeleteBatch(ctx context.Context, id int64) error {
	return s.repo.DeleteBatch(ctx, id)
}

func (s *batchService) ListBatches(ctx context.Context, page, pageSize int64, filter string) ([]model.Batch, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListBatches(ctx, page, pageSize, filter)
}