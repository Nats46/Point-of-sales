package service

import (
	"context"
	"errors"
	"time"

	"point-of-sales/internal/model"
	"point-of-sales/internal/repository"
)

type InventoryService interface {
	CreateInventory(ctx context.Context, itemCode, itemName string, price float64, unit string, createdBy string) (*model.Inventory, error)
	GetInventory(ctx context.Context, id int64) (*model.Inventory, error)
	UpdateInventory(ctx context.Context, id int64, itemCode, itemName string, price float64, unit string, updatedBy string) (*model.Inventory, error)
	DeleteInventory(ctx context.Context, id int64) error
	ListInventories(ctx context.Context, page, pageSize int64, filter string) ([]model.Inventory, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) CreateInventory(ctx context.Context, itemCode, itemName string, price float64, unit string, createdBy string) (*model.Inventory, error) {
	if itemCode == "" || itemName == "" || unit == "" {
		return nil, errors.New("item code, item name and unit are required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	inventory := &model.Inventory{
		ItemCode:  itemCode,
		ItemName:  itemName,
		Price:     price,
		Unit:      unit,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}
	if err := s.repo.InsertInventory(ctx, inventory); err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *inventoryService) GetInventory(ctx context.Context, id int64) (*model.Inventory, error) {
	return s.repo.GetInventory(ctx, id)
}

func (s *inventoryService) UpdateInventory(ctx context.Context, id int64, itemCode, itemName string, price float64, unit string, updatedBy string) (*model.Inventory, error) {
	inventory, err := s.repo.GetInventory(ctx, id)
	if err != nil {
		return nil, err
	}
	if itemCode != "" {
		inventory.ItemCode = itemCode
	}
	if itemName != "" {
		inventory.ItemName = itemName
	}
	if price > 0 {
		inventory.Price = price
	}
	if unit != "" {
		inventory.Unit = unit
	}
	inventory.UpdatedAt = time.Now()
	inventory.UpdatedBy = updatedBy
	if err := s.repo.UpdateInventory(ctx, inventory); err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *inventoryService) DeleteInventory(ctx context.Context, id int64) error {
	return s.repo.DeleteInventory(ctx, id)
}

func (s *inventoryService) ListInventories(ctx context.Context, page, pageSize int64, filter string) ([]model.Inventory, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.ListInventories(ctx, page, pageSize, filter)
}