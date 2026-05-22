package service

import (
	"context"
	"fmt"
	"point-of-sales/internal/model"
	"point-of-sales/internal/repository"
	"time"
)

type TransactionService interface {
	CreateSales(ctx context.Context, header *model.SalesHeader, details []model.SalesDetail, createdBy string) (*model.SalesHeader, error)
	GetSalesByID(ctx context.Context, id int64) (*model.SalesHeader, error)
	GetSalesByNo(ctx context.Context, no string) (*model.SalesHeader, error)
	UpdateSales(ctx context.Context, header *model.SalesHeader, details []model.SalesDetail, updatedBy string) (*model.SalesHeader, error)
	DeleteSales(ctx context.Context, id int64) error
	ListSales(ctx context.Context, page, pageSize int64, filter string) ([]*model.SalesHeader, error)
}

type transactionService struct {
	repo         repository.TransactionRepository
	detailRepo   repository.SalesDetailRepository
}

func NewTransactionService(repo repository.TransactionRepository, detailRepo repository.SalesDetailRepository) TransactionService {
	return &transactionService{repo: repo, detailRepo: detailRepo}
}

func (s *transactionService) CreateSales(ctx context.Context, header *model.SalesHeader, details []model.SalesDetail, createdBy string) (*model.SalesHeader, error) {
	header.CreatedAt = time.Now()
	header.UpdatedAt = time.Now()
	header.CreatedBy = createdBy
	header.UpdatedBy = createdBy

	if err := s.repo.CreateSalesHeader(ctx, header); err != nil {
		return nil, fmt.Errorf("failed to create sales header: %w", err)
	}

	var detailPtrs []*model.SalesDetail
	for i := range details {
		details[i].TransactionId = header.TransactionId
		details[i].CreatedAt = header.CreatedAt
		details[i].UpdatedAt = header.UpdatedAt
		details[i].CreatedBy = createdBy
		details[i].UpdatedBy = createdBy
		detailPtrs = append(detailPtrs, &details[i])
		if err := s.detailRepo.CreateSalesDetail(ctx, &details[i]); err != nil {
			return nil, fmt.Errorf("failed to create sales detail: %w", err)
		}
	}

	header.SalesDetails = detailPtrs
	return header, nil
}

func (s *transactionService) GetSalesByID(ctx context.Context, id int64) (*model.SalesHeader, error) {
	header, err := s.repo.GetSalesHeaderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	details, err := s.detailRepo.ListSalesDetailsByTransaction(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales details: %w", err)
	}

	header.SalesDetails = details
	return header, nil
}

func (s *transactionService) GetSalesByNo(ctx context.Context, no string) (*model.SalesHeader, error) {
	header, err := s.repo.GetSalesHeaderByNo(ctx, no)
	if err != nil {
		return nil, err
	}

	details, err := s.detailRepo.ListSalesDetailsByTransaction(ctx, header.TransactionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales details: %w", err)
	}

	header.SalesDetails = details
	return header, nil
}

func (s *transactionService) UpdateSales(ctx context.Context, header *model.SalesHeader, details []model.SalesDetail, updatedBy string) (*model.SalesHeader, error) {
	header.UpdatedAt = time.Now()
	header.UpdatedBy = updatedBy

	if err := s.repo.UpdateSalesHeader(ctx, header); err != nil {
		return nil, fmt.Errorf("failed to update sales header: %w", err)
	}

	var detailPtrs []*model.SalesDetail
	for i := range details {
		details[i].TransactionId = header.TransactionId
		details[i].UpdatedAt = header.UpdatedAt
		details[i].UpdatedBy = updatedBy
		detailPtrs = append(detailPtrs, &details[i])
		if err := s.detailRepo.UpdateSalesDetail(ctx, &details[i]); err != nil {
			return nil, fmt.Errorf("failed to update sales detail: %w", err)
		}
	}

	header.SalesDetails = detailPtrs
	return header, nil
}

func (s *transactionService) DeleteSales(ctx context.Context, id int64) error {
	return s.repo.DeleteSalesHeader(ctx, id)
}

func (s *transactionService) ListSales(ctx context.Context, page, pageSize int64, filter string) ([]*model.SalesHeader, error) {
	return s.repo.ListSalesHeaders(ctx, page, pageSize, filter)
}