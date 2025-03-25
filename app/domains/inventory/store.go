package inventory

import (
	"gorm.io/gorm"
)

type IStore interface {
	GetStock(productID uint) (int, error)
	UpdateStock(productID uint, quantityChange int) error
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) IStore {
	return &store{db: db}
}

func (s *store) GetStock(productID uint) (int, error) {
	var stock int
	err := s.db.
		Model(&Inventory{}).
		Select("stock").
		Where("product_id = ?", productID).
		Scan(&stock).
		Error

	return stock, err
}

func (s *store) UpdateStock(productID uint, quantityChange int) error {
	return s.db.
		Model(&Inventory{}).
		Where("product_id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantityChange)).
		Error
}
