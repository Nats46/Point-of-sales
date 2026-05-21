package repository

import (
	"context"
	"point-of-sales/internal/model"

	"github.com/jmoiron/sqlx"
)

type InventoryRepository interface {
	InsertInventory(ctx context.Context, i *model.Inventory) error
	GetInventory(ctx context.Context, id int64) (*model.Inventory, error)
	UpdateInventory(ctx context.Context, i *model.Inventory) error
	DeleteInventory(ctx context.Context, id int64) error
	ListInventories(ctx context.Context, page, pageSize int64, filter string) ([]model.Inventory, error)
}

type inventoryRepository struct {
	DB *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) InventoryRepository {
	return &inventoryRepository{DB: db}
}

func (ir *inventoryRepository) InsertInventory(ctx context.Context, i *model.Inventory) error {
	query := `
		INSERT INTO inventories (item_code, item_name, price, unit, created_at, updated_at, created_by, updated_by)
		VALUES (:item_code, :item_name, :price, :unit, :created_at, :updated_at, :created_by, :updated_by)
		RETURNING id
	`
	rows, err := ir.DB.NamedQueryContext(ctx, query, i)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&i.Id)
	}
	return nil
}

func (ir *inventoryRepository) GetInventory(ctx context.Context, id int64) (*model.Inventory, error) {
	var inventory model.Inventory
	query := `SELECT * FROM inventories WHERE id = $1`
	err := ir.DB.GetContext(ctx, &inventory, query, id)
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (ir *inventoryRepository) UpdateInventory(ctx context.Context, i *model.Inventory) error {
	query := `
		UPDATE inventories
		SET item_code = :item_code, item_name = :item_name, price = :price, unit = :unit,
			updated_at = :updated_at, updated_by = :updated_by
		WHERE id = :id
	`
	_, err := ir.DB.NamedExecContext(ctx, query, i)
	return err
}

func (ir *inventoryRepository) DeleteInventory(ctx context.Context, id int64) error {
	query := `DELETE FROM inventories WHERE id = $1`
	_, err := ir.DB.ExecContext(ctx, query, id)
	return err
}

func (ir *inventoryRepository) ListInventories(ctx context.Context, page, pageSize int64, filter string) ([]model.Inventory, error) {
	var inventories []model.Inventory
	query := `SELECT * FROM inventories WHERE item_code ILIKE $1 OR item_name ILIKE $1 LIMIT $2 OFFSET $3`
	filterParam := "%" + filter + "%"
	offset := (page - 1) * pageSize
	err := ir.DB.SelectContext(ctx, &inventories, query, filterParam, pageSize, offset)
	return inventories, err
}