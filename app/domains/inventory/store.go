package inventory

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type IStore interface {
	Get(ctx context.Context, productID uint) (*Inventory, error)
	GetMultiple(ctx context.Context, productIDs []uint) (map[uint]Inventory, error)
	IncreaseStockBulk(ctx context.Context, updates map[uint]int) error
	DecreaseStockBulk(ctx context.Context, updates map[uint]int) error
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) IStore {
	return &store{db: db}
}

func (s *store) GetMultiple(ctx context.Context, productIDs []uint) (map[uint]Inventory, error) {
	var inventories []Inventory
	result := s.db.WithContext(ctx).
		Preload("Product").
		Where("product_id IN (?)", productIDs).Find(&inventories)
	if result.Error != nil {
		return nil, result.Error
	}

	inventoryMap := make(map[uint]Inventory)
	for i := range inventories {
		inventoryMap[inventories[i].ProductID] = inventories[i]
	}

	return inventoryMap, nil
}

func (s *store) Get(ctx context.Context, productID uint) (*Inventory, error) {
	var result Inventory
	query := s.db.
		WithContext(ctx).
		Preload("Product").
		First(&result, "product_id = ?", productID)

	if query.Error != nil {
		return nil, query.Error
	}
	return &result, nil
}

func (s *store) IncreaseStockBulk(ctx context.Context, updates map[uint]int) error {
	tx := s.db.WithContext(ctx).Begin()

	for productID, quantity := range updates {
		result := tx.Model(&Inventory{}).
			Where("product_id = ?", productID).
			Update("stock", gorm.Expr("stock + ?", quantity))

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	return tx.Commit().Error
}

func (s *store) DecreaseStockBulk(ctx context.Context, updates map[uint]int) error {
	tx := s.db.WithContext(ctx).Begin()

	for productID, quantity := range updates {
		result := tx.Model(&Inventory{}).
			Where("product_id = ? AND stock >= ?", productID, quantity).
			Update("stock", gorm.Expr("stock - ?", quantity))

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("not enough stock for product %d", productID)
		}
	}

	return tx.Commit().Error
}
