package order

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type IStore interface {
	Create(ctx context.Context, order *Order) error
	Get(ctx context.Context, id uint) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uint) error
	DeleteOrderItems(ctx context.Context, orderItemIDs []uint) error
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) IStore {
	return &store{db: db}
}

func (s *store) Create(ctx context.Context, order *Order) error {
	return s.db.WithContext(ctx).Create(order).Error
}

func (s *store) Get(ctx context.Context, id uint) (*Order, error) {
	var order Order
	err := s.db.
		WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Where("id = ?", id).
		First(&order).Error

	return &order, err
}

func (s *store) Update(ctx context.Context, order *Order) error {
	return s.db.WithContext(ctx).Save(order).Error
}

func (s *store) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&Order{}).Error
}

func (s *store) DeleteOrderItems(ctx context.Context, orderItemIDs []uint) error {
	if err := s.db.WithContext(ctx).Where("id IN ?", orderItemIDs).Delete(&OrderItem{}).Error; err != nil {
		return fmt.Errorf("failed to bulk delete order items: %w", err)
	}
	return nil
}
