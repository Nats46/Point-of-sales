package service

import (
	"context"
	"fmt"
	"point-of-sales/internal/model"
	"point-of-sales/internal/repository"
	"time"
)

type DiscountService interface {
	CreateDiscount(ctx context.Context, code, discountType string, value float64, startDate, endDate time.Time, createdBy string) (*model.Discount, error)
	GetDiscountByID(ctx context.Context, id int64) (*model.Discount, error)
	GetDiscountByCode(ctx context.Context, code string) (*model.Discount, error)
	UpdateDiscount(ctx context.Context, id int64, code, discountType string, value float64, startDate, endDate time.Time, isActive bool, updatedBy string) (*model.Discount, error)
	DeleteDiscount(ctx context.Context, id int64) error
	ListDiscounts(ctx context.Context, page, pageSize int64, filter string) ([]*model.Discount, error)
}

type discountService struct {
	repo repository.DiscountRepository
}

func NewDiscountService(repo repository.DiscountRepository) DiscountService {
	return &discountService{repo: repo}
}

func (s *discountService) CreateDiscount(ctx context.Context, code, discountType string, value float64, startDate, endDate time.Time, createdBy string) (*model.Discount, error) {
	if code == "" || discountType == "" || createdBy == "" {
		return nil, fmt.Errorf("code, discount type, and created by are required")
	}

	discount := &model.Discount{
		DiscountCode: code,
		Value:        value,
		DiscountType: discountType,
		StartDate:    startDate,
		EndDate:      endDate,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
	}

	if err := s.repo.Create(ctx, discount); err != nil {
		return nil, fmt.Errorf("failed to create discount: %w", err)
	}

	return discount, nil
}

func (s *discountService) GetDiscountByID(ctx context.Context, id int64) (*model.Discount, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *discountService) GetDiscountByCode(ctx context.Context, code string) (*model.Discount, error) {
	return s.repo.GetByCode(ctx, code)
}

func (s *discountService) UpdateDiscount(ctx context.Context, id int64, code, discountType string, value float64, startDate, endDate time.Time, isActive bool, updatedBy string) (*model.Discount, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.DiscountCode = code
	existing.DiscountType = discountType
	existing.Value = value
	existing.StartDate = startDate
	existing.EndDate = endDate
	existing.IsActive = isActive
	existing.UpdatedAt = time.Now()
	existing.UpdatedBy = updatedBy

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update discount: %w", err)
	}

	return existing, nil
}

func (s *discountService) DeleteDiscount(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *discountService) ListDiscounts(ctx context.Context, page, pageSize int64, filter string) ([]*model.Discount, error) {
	return s.repo.List(ctx, page, pageSize, filter)
}